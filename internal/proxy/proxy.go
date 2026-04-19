package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/Kyeong6/autolang/internal/config"
	"github.com/Kyeong6/autolang/internal/detect"
	"github.com/Kyeong6/autolang/internal/protect"
	"github.com/Kyeong6/autolang/internal/stats"
	"github.com/Kyeong6/autolang/internal/translate"
)

const anthropicBase = "https://api.anthropic.com"

// Proxy is the local HTTP proxy server that intercepts Claude API requests,
// applies Korean↔English translation, and forwards to the real Anthropic API.
type Proxy struct {
	cfg        *config.Config
	translator translate.Translator
	stats      *stats.Stats
	server     *http.Server
	client     *http.Client
	logger     *log.Logger
}

// New creates a Proxy. translator may be nil (passthrough mode).
func New(cfg *config.Config, t translate.Translator) *Proxy {
	p := &Proxy{
		cfg:        cfg,
		translator: t,
		stats:      stats.New(),
		client:     &http.Client{Timeout: 0}, // no timeout — streaming responses can be long
		logger:     log.New(io.Discard, "[autolang] ", log.LstdFlags),
	}

	if cfg.Proxy.LogLevel != "off" {
		p.logger.SetOutput(log.Writer())
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", p.handleHealth)
	mux.HandleFunc("/stats", p.handleStats)
	mux.HandleFunc("/v1/messages", p.handleMessages)
	mux.HandleFunc("/", p.handlePassthrough)

	p.server = &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Proxy.Port),
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	return p
}

// Start starts the proxy and blocks until Stop is called or an error occurs.
func (p *Proxy) Start() error {
	ln, err := net.Listen("tcp", p.server.Addr)
	if err != nil {
		return fmt.Errorf("failed to bind %s: %w", p.server.Addr, err)
	}
	p.logger.Printf("proxy listening on %s", p.server.Addr)
	return p.server.Serve(ln)
}

// Stop gracefully shuts down the proxy server.
func (p *Proxy) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return p.server.Shutdown(ctx)
}

// handleHealth responds to liveness checks from `autolang status`.
func (p *Proxy) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"ok","port":%d}`, p.cfg.Proxy.Port)
}

// handleStats returns a plain-text session summary for `autolang stats`.
func (p *Proxy) handleStats(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, p.stats.Summary())
}

// handleMessages intercepts POST /v1/messages, applies translation, and forwards.
func (p *Proxy) handleMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "cannot read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req MessagesRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	// --- translation hook (Task 03 / 04 / 05 will populate this) ---
	if p.translator != nil {
		if err := p.translateRequest(r.Context(), &req); err != nil {
			p.logger.Printf("translation error: %v", err)
			// fall through — send original on translation failure
		}
	}
	// ----------------------------------------------------------------

	translated, err := json.Marshal(req)
	if err != nil {
		http.Error(w, "cannot marshal request", http.StatusInternalServerError)
		return
	}

	p.forward(w, r, translated, req.Stream)
}

// handlePassthrough transparently forwards any other Anthropic API endpoints.
func (p *Proxy) handlePassthrough(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	p.forward(w, r, body, false)
}

// forward sends body to the real Anthropic API and relays the response.
// When stream is true it flushes incrementally so SSE reaches the client in real time.
func (p *Proxy) forward(w http.ResponseWriter, r *http.Request, body []byte, stream bool) {
	upstream := anthropicBase + r.URL.Path
	if r.URL.RawQuery != "" {
		upstream += "?" + r.URL.RawQuery
	}

	req, err := http.NewRequestWithContext(r.Context(), r.Method, upstream, bytes.NewReader(body))
	if err != nil {
		http.Error(w, "cannot build upstream request", http.StatusInternalServerError)
		return
	}

	// Forward all original headers (auth key, anthropic-version, beta flags, etc.)
	for key, vals := range r.Header {
		for _, v := range vals {
			req.Header.Add(key, v)
		}
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		p.logger.Printf("upstream error: %v", err)
		http.Error(w, "upstream request failed: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Relay response headers
	for key, vals := range resp.Header {
		for _, v := range vals {
			w.Header().Add(key, v)
		}
	}
	w.WriteHeader(resp.StatusCode)

	if stream {
		p.relayStream(r.Context(), w, resp.Body)
	} else {
		io.Copy(w, resp.Body) //nolint:errcheck
	}
}

// relayStream copies an SSE response body to the client.
// When a translator is configured, text is translated sentence-by-sentence (en→ko).
func (p *Proxy) relayStream(ctx context.Context, w http.ResponseWriter, body io.Reader) {
	if p.translator != nil {
		st := NewStreamTranslator(ctx, w, p.translator, p.cfg.Translation, p.stats, p.logger)
		st.Relay(body)
		return
	}

	flusher, canFlush := w.(http.Flusher)
	buf := make([]byte, 4096)
	for {
		n, err := body.Read(buf)
		if n > 0 {
			w.Write(buf[:n]) //nolint:errcheck
			if canFlush {
				flusher.Flush()
			}
		}
		if err != nil {
			break
		}
	}
}

// translateRequest applies Korean detection, code-block protection, and translation
// to all user messages in the request.
func (p *Proxy) translateRequest(ctx context.Context, req *MessagesRequest) error {
	cfg := p.cfg.Translation
	for i := range req.Messages {
		msg := &req.Messages[i]
		if msg.Role != "user" {
			continue
		}
		text, ok := msg.ContentAsString()
		if !ok {
			continue
		}

		// Skip messages that don't need translation
		if !detect.ShouldTranslate(text) {
			continue
		}

		// Extract code blocks, URLs, file paths before translating
		sanitised, protector := protect.Protect(text)

		translated, err := p.translator.Translate(ctx, sanitised, cfg.SourceLang, cfg.TargetLang)
		if err != nil {
			return err
		}

		// Restore protected regions in the translated text
		restored := protector.Restore(translated)

		if err := msg.SetContentString(restored); err != nil {
			return err
		}

		koTokens := stats.EstimateKorean(text)
		enTokens := stats.EstimateEnglish(restored)
		p.stats.RecordInput(text, restored)

		if p.cfg.Output.ShowTranslationIndicator {
			fmt.Fprintln(os.Stderr, p.stats.Indicator(koTokens, enTokens))
		}
	}
	return nil
}
