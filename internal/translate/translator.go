package translate

import "context"

// Translator is the interface implemented by all translation backends.
// Implementation: Task 05
type Translator interface {
	Translate(ctx context.Context, text, from, to string) (string, error)
	Name() string
}
