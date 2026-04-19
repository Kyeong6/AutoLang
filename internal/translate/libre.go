package translate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const libreEndpoint = "https://libretranslate.com/translate"

// LibreTranslate uses the public LibreTranslate instance.
// No API key required, but subject to rate limits.
type LibreTranslate struct {
	client *http.Client
}

func NewLibreTranslate() *LibreTranslate {
	return &LibreTranslate{
		client: &http.Client{Timeout: 20 * time.Second},
	}
}

func (l *LibreTranslate) Name() string { return "libretranslate" }

func (l *LibreTranslate) Translate(ctx context.Context, text, from, to string) (string, error) {
	if strings.TrimSpace(text) == "" {
		return text, nil
	}

	payload := map[string]string{
		"q":      text,
		"source": strings.ToLower(from),
		"target": strings.ToLower(to),
		"format": "text",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, libreEndpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := l.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("libretranslate: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("libretranslate: unexpected status %d (rate limit?)", resp.StatusCode)
	}

	var result struct {
		TranslatedText string `json:"translatedText"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("libretranslate: failed to decode response: %w", err)
	}
	return result.TranslatedText, nil
}
