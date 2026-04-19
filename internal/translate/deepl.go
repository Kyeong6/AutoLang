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

const deeplEndpoint = "https://api-free.deepl.com/v2/translate"

type DeepL struct {
	apiKey string
	client *http.Client
}

func NewDeepL(apiKey string) *DeepL {
	return &DeepL{
		apiKey: apiKey,
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

func (d *DeepL) Name() string { return "deepl" }

func (d *DeepL) Translate(ctx context.Context, text, from, to string) (string, error) {
	if strings.TrimSpace(text) == "" {
		return text, nil
	}

	payload := map[string]interface{}{
		"text":        []string{text},
		"source_lang": strings.ToUpper(from),
		"target_lang": strings.ToUpper(to),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, deeplEndpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "DeepL-Auth-Key "+d.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("deepl: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("deepl: unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Translations []struct {
			Text string `json:"text"`
		} `json:"translations"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("deepl: failed to decode response: %w", err)
	}
	if len(result.Translations) == 0 {
		return "", fmt.Errorf("deepl: empty translation response")
	}
	return result.Translations[0].Text, nil
}
