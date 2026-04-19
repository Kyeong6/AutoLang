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

const googleEndpoint = "https://translation.googleapis.com/language/translate/v2"

type Google struct {
	apiKey string
	client *http.Client
}

func NewGoogle(apiKey string) *Google {
	return &Google{
		apiKey: apiKey,
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

func (g *Google) Name() string { return "google" }

func (g *Google) Translate(ctx context.Context, text, from, to string) (string, error) {
	if strings.TrimSpace(text) == "" {
		return text, nil
	}

	payload := map[string]interface{}{
		"q":      text,
		"source": strings.ToLower(from),
		"target": strings.ToLower(to),
		"format": "text",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s?key=%s", googleEndpoint, g.apiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("google: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("google: unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Data struct {
			Translations []struct {
				TranslatedText string `json:"translatedText"`
			} `json:"translations"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("google: failed to decode response: %w", err)
	}
	if len(result.Data.Translations) == 0 {
		return "", fmt.Errorf("google: empty translation response")
	}
	return result.Data.Translations[0].TranslatedText, nil
}
