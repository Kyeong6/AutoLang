package proxy

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Kyeong6/autolang/internal/config"
	"github.com/Kyeong6/autolang/internal/stats"
	"github.com/Kyeong6/autolang/internal/translate"
)

type sseEvent struct {
	Type  string   `json:"type"`
	Index int      `json:"index"`
	Delta sseDelta `json:"delta"`
}

type sseDelta struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// StreamTranslator buffers SSE text_delta chunks and translates at sentence boundaries.
type StreamTranslator struct {
	ctx         context.Context
	translator  translate.Translator
	cfg         config.TranslationConfig
	stats       *stats.Stats
	buffer      strings.Builder
	rawBuffer   strings.Builder // original English text for stats
	inCodeBlock bool
	writer      http.ResponseWriter
	flusher     http.Flusher
	logger      *log.Logger
}

func NewStreamTranslator(
	ctx context.Context,
	w http.ResponseWriter,
	tr translate.Translator,
	cfg config.TranslationConfig,
	st *stats.Stats,
	logger *log.Logger,
) *StreamTranslator {
	s := &StreamTranslator{
		ctx:        ctx,
		translator: tr,
		cfg:        cfg,
		stats:      st,
		writer:     w,
		logger:     logger,
	}
	if f, ok := w.(http.Flusher); ok {
		s.flusher = f
	}
	return s
}

// Relay reads the SSE response body, translating text at sentence boundaries.
func (s *StreamTranslator) Relay(body io.Reader) {
	scanner := bufio.NewScanner(body)
	for scanner.Scan() {
		if err := s.processLine(scanner.Text()); err != nil {
			s.logger.Printf("stream translate error: %v", err)
		}
	}
	// Flush any text remaining in the buffer after stream ends.
	if err := s.Flush(); err != nil {
		s.logger.Printf("stream flush error: %v", err)
	}
}

func (s *StreamTranslator) processLine(line string) error {
	if !strings.HasPrefix(line, "data: ") {
		s.writeLine(line)
		return nil
	}

	payload := line[6:]
	if payload == "[DONE]" {
		s.writeLine(line)
		return nil
	}

	var ev sseEvent
	if err := json.Unmarshal([]byte(payload), &ev); err != nil {
		s.writeLine(line)
		return nil
	}

	switch ev.Type {
	case "content_block_delta":
		if ev.Delta.Type == "text_delta" {
			return s.processText(ev.Delta.Text)
		}
		s.writeLine(line)
	case "message_stop":
		if err := s.Flush(); err != nil {
			return err
		}
		s.writeLine(line)
	default:
		s.writeLine(line)
	}
	return nil
}

// processText routes each text chunk: code-block content passes through,
// natural-language text is buffered until a sentence boundary is detected.
func (s *StreamTranslator) processText(text string) error {
	if !strings.Contains(text, "```") {
		if s.inCodeBlock {
			s.writeTextDelta(text)
		} else {
			s.buffer.WriteString(text)
			s.rawBuffer.WriteString(text)
			return s.flushIfSentence()
		}
		return nil
	}

	// The chunk contains one or more ``` fences — split and toggle.
	parts := strings.Split(text, "```")
	for i, part := range parts {
		if s.inCodeBlock {
			s.writeTextDelta(part)
			if i < len(parts)-1 {
				s.writeTextDelta("```")
				s.inCodeBlock = false
			}
		} else {
			s.buffer.WriteString(part)
			s.rawBuffer.WriteString(part)
			if err := s.flushIfSentence(); err != nil {
				return err
			}
			if i < len(parts)-1 {
				if err := s.Flush(); err != nil {
					return err
				}
				s.writeTextDelta("```")
				s.inCodeBlock = true
			}
		}
	}
	return nil
}

// flushIfSentence translates and emits buffered text when a sentence boundary is found.
func (s *StreamTranslator) flushIfSentence() error {
	text := s.buffer.String()
	if strings.Contains(text, "\n\n") ||
		strings.HasSuffix(text, ". ") ||
		strings.HasSuffix(text, "? ") ||
		strings.HasSuffix(text, "! ") {
		return s.Flush()
	}
	return nil
}

// Flush translates any buffered text (en→ko) and writes it to the client.
func (s *StreamTranslator) Flush() error {
	text := s.buffer.String()
	if text == "" {
		return nil
	}
	s.buffer.Reset()
	raw := s.rawBuffer.String()
	s.rawBuffer.Reset()

	// Response direction: TargetLang (en) → SourceLang (ko)
	translated, err := s.translator.Translate(s.ctx, text, s.cfg.TargetLang, s.cfg.SourceLang)
	if err != nil {
		s.logger.Printf("response translate error: %v", err)
		s.writeTextDelta(text) // fallback to original on error
		return nil
	}
	s.writeTextDelta(translated)
	if s.stats != nil {
		s.stats.RecordOutput(raw, translated)
	}
	return nil
}

func (s *StreamTranslator) writeTextDelta(text string) {
	if text == "" {
		return
	}
	event := map[string]interface{}{
		"type":  "content_block_delta",
		"index": 0,
		"delta": map[string]string{
			"type": "text_delta",
			"text": text,
		},
	}
	data, err := json.Marshal(event)
	if err != nil {
		return
	}
	fmt.Fprintf(s.writer, "data: %s\n\n", data)
	s.doFlush()
}

func (s *StreamTranslator) writeLine(line string) {
	fmt.Fprintf(s.writer, "%s\n", line)
	s.doFlush()
}

func (s *StreamTranslator) doFlush() {
	if s.flusher != nil {
		s.flusher.Flush()
	}
}
