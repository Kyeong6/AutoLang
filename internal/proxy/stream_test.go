package proxy

import (
	"context"
	"io"
	"log"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Kyeong6/autolang/internal/config"
)

// echoTranslator returns "[KO] " + text to simulate en→ko translation.
type echoTranslator struct{}

func (e *echoTranslator) Name() string { return "echo" }
func (e *echoTranslator) Translate(_ context.Context, text, _, _ string) (string, error) {
	if strings.TrimSpace(text) == "" {
		return text, nil
	}
	return "[KO] " + text, nil
}

func newTestStreamTranslator(w *httptest.ResponseRecorder) *StreamTranslator {
	cfg := config.TranslationConfig{SourceLang: "ko", TargetLang: "en"}
	return NewStreamTranslator(context.Background(), w, &echoTranslator{}, cfg, nil, log.New(io.Discard, "", 0))
}

func TestStreamTranslator_PassThroughNonDataLines(t *testing.T) {
	rec := httptest.NewRecorder()
	st := newTestStreamTranslator(rec)

	body := strings.NewReader("event: content_block_delta\n")
	st.Relay(body)

	if !strings.Contains(rec.Body.String(), "event: content_block_delta") {
		t.Errorf("expected event line to pass through, got: %q", rec.Body.String())
	}
}

func TestStreamTranslator_TranslatesOnMessageStop(t *testing.T) {
	rec := httptest.NewRecorder()
	st := newTestStreamTranslator(rec)

	sseBody := `data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Hello world"}}
data: {"type":"message_stop"}
`
	st.Relay(strings.NewReader(sseBody))

	out := rec.Body.String()
	if !strings.Contains(out, "[KO]") {
		t.Errorf("expected translated output, got: %q", out)
	}
	if !strings.Contains(out, "message_stop") {
		t.Errorf("expected message_stop to be forwarded, got: %q", out)
	}
}

func TestStreamTranslator_FlushOnSentenceBoundary(t *testing.T) {
	rec := httptest.NewRecorder()
	st := newTestStreamTranslator(rec)

	// Two chunks: "Hello. " triggers flush; "World" stays buffered until message_stop.
	st.processText("Hello. ") //nolint:errcheck
	if !strings.Contains(rec.Body.String(), "[KO]") {
		t.Errorf("expected flush after sentence boundary, got: %q", rec.Body.String())
	}
}

func TestStreamTranslator_CodeBlockPassthrough(t *testing.T) {
	rec := httptest.NewRecorder()
	st := newTestStreamTranslator(rec)

	// Opening fence
	st.processText("```") //nolint:errcheck
	// Content inside code block — should NOT be translated
	st.processText("print('hello')") //nolint:errcheck
	// Closing fence
	st.processText("```") //nolint:errcheck

	out := rec.Body.String()
	// Code content must appear verbatim (not wrapped in [KO])
	if strings.Contains(out, "[KO] print") {
		t.Errorf("code block content should not be translated, got: %q", out)
	}
	if !strings.Contains(out, "print('hello')") {
		t.Errorf("expected code block content in output, got: %q", out)
	}
}

func TestStreamTranslator_MixedCodeAndText(t *testing.T) {
	rec := httptest.NewRecorder()
	st := newTestStreamTranslator(rec)

	// Single chunk containing text + code block
	st.processText("Use this: ```go\nfmt.Println()\n``` to print. ") //nolint:errcheck

	out := rec.Body.String()
	// Text outside code block should be translated
	if !strings.Contains(out, "[KO]") {
		t.Errorf("expected text parts to be translated, got: %q", out)
	}
	// Code inside block must not be translated
	if strings.Contains(out, "[KO] fmt") {
		t.Errorf("code inside block must not be translated, got: %q", out)
	}
}

func TestStreamTranslator_EmptyBuffer_NoWrite(t *testing.T) {
	rec := httptest.NewRecorder()
	st := newTestStreamTranslator(rec)

	// Flush on empty buffer should produce no output
	st.Flush() //nolint:errcheck
	if rec.Body.Len() > 0 {
		t.Errorf("expected no output for empty flush, got: %q", rec.Body.String())
	}
}

func TestStreamTranslator_DoublNewlineFlushesParagraph(t *testing.T) {
	rec := httptest.NewRecorder()
	st := newTestStreamTranslator(rec)

	st.processText("First paragraph\n\n") //nolint:errcheck

	if !strings.Contains(rec.Body.String(), "[KO]") {
		t.Errorf("expected paragraph to be flushed on double newline, got: %q", rec.Body.String())
	}
}
