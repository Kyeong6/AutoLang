package translate

import "context"

// DeepL translation backend.
// Implementation: Task 05
type DeepL struct {
	apiKey string
}

func NewDeepL(apiKey string) *DeepL {
	return &DeepL{apiKey: apiKey}
}

func (d *DeepL) Name() string { return "deepl" }

func (d *DeepL) Translate(_ context.Context, text, _, _ string) (string, error) {
	return text, nil
}
