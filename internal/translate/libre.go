package translate

import "context"

// LibreTranslate backend (no API key required).
// Implementation: Task 05
type LibreTranslate struct{}

func NewLibreTranslate() *LibreTranslate {
	return &LibreTranslate{}
}

func (l *LibreTranslate) Name() string { return "libretranslate" }

func (l *LibreTranslate) Translate(_ context.Context, text, _, _ string) (string, error) {
	return text, nil
}
