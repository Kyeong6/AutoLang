package proxy

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Kyeong6/autolang/internal/config"
)

func newTestProxy(upstream string) *Proxy {
	cfg := config.Default()
	p := New(cfg, nil)
	// Point to the test upstream instead of api.anthropic.com
	p.client = &http.Client{}
	return p
}

func TestHandleHealth(t *testing.T) {
	cfg := config.Default()
	p := New(cfg, nil)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	p.handleHealth(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("expected status ok, got %v", body["status"])
	}
}

func TestHandleMessages_MethodNotAllowed(t *testing.T) {
	cfg := config.Default()
	p := New(cfg, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/messages", nil)
	rec := httptest.NewRecorder()

	p.handleMessages(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}

func TestHandleMessages_InvalidJSON(t *testing.T) {
	cfg := config.Default()
	p := New(cfg, nil)

	req := httptest.NewRequest(http.MethodPost, "/v1/messages", strings.NewReader("not json"))
	rec := httptest.NewRecorder()

	p.handleMessages(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestHandleMessages_ParsesRequest(t *testing.T) {
	body := `{"model":"claude-3-5-sonnet-20241022","max_tokens":100,"messages":[{"role":"user","content":"hello"}]}`

	var msgReq MessagesRequest
	if err := json.NewDecoder(strings.NewReader(body)).Decode(&msgReq); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if msgReq.Model != "claude-3-5-sonnet-20241022" {
		t.Errorf("expected model to be preserved, got %s", msgReq.Model)
	}
	if len(msgReq.Messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(msgReq.Messages))
	}
}

func TestMessageContentAsString(t *testing.T) {
	body := `{"role":"user","content":"안녕하세요"}`
	var msg Message
	if err := json.Unmarshal([]byte(body), &msg); err != nil {
		t.Fatal(err)
	}
	text, ok := msg.ContentAsString()
	if !ok {
		t.Fatal("expected string content")
	}
	if text != "안녕하세요" {
		t.Errorf("expected '안녕하세요', got %q", text)
	}
}

func TestMessageSetContentString(t *testing.T) {
	body := `{"role":"user","content":"hello"}`
	var msg Message
	json.Unmarshal([]byte(body), &msg)

	if err := msg.SetContentString("world"); err != nil {
		t.Fatal(err)
	}
	text, _ := msg.ContentAsString()
	if text != "world" {
		t.Errorf("expected 'world', got %q", text)
	}
}
