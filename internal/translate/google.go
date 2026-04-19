package translate

import "context"

// Google Cloud Translation backend.
// Implementation: Task 05
type Google struct {
	apiKey string
}

func NewGoogle(apiKey string) *Google {
	return &Google{apiKey: apiKey}
}

func (g *Google) Name() string { return "google" }

func (g *Google) Translate(_ context.Context, text, _, _ string) (string, error) {
	return text, nil
}
