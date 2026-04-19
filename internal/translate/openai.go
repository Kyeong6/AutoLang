package translate

import "context"

// OpenAI translation backend (uses gpt-4o-mini).
// Implementation: Task 05
type OpenAI struct {
	apiKey string
}

func NewOpenAI(apiKey string) *OpenAI {
	return &OpenAI{apiKey: apiKey}
}

func (o *OpenAI) Name() string { return "openai" }

func (o *OpenAI) Translate(_ context.Context, text, _, _ string) (string, error) {
	return text, nil
}
