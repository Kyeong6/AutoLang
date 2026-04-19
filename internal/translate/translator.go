package translate

import (
	"context"
	"fmt"
	"os"

	"github.com/Kyeong6/autolang/internal/config"
)

// Translator is the interface implemented by all translation backends.
type Translator interface {
	Translate(ctx context.Context, text, from, to string) (string, error)
	Name() string
}

// New returns a Translator based on config and environment variables.
// Environment variables take precedence over config file values.
//
//	AUTOLANG_PROVIDER — deepl | openai | google | libretranslate
//	AUTOLANG_API_KEY  — API key for the chosen provider
func New(cfg *config.Config) (Translator, error) {
	provider := cfg.Translation.Provider
	apiKey := cfg.Translation.APIKey

	// Environment variables override config file
	if v := os.Getenv("AUTOLANG_PROVIDER"); v != "" {
		provider = v
	}
	if v := os.Getenv("AUTOLANG_API_KEY"); v != "" {
		apiKey = v
	}

	switch provider {
	case "deepl":
		if apiKey == "" {
			return nil, fmt.Errorf("AUTOLANG_API_KEY is required for DeepL")
		}
		return NewDeepL(apiKey), nil
	case "openai":
		if apiKey == "" {
			return nil, fmt.Errorf("AUTOLANG_API_KEY is required for OpenAI")
		}
		return NewOpenAI(apiKey), nil
	case "google":
		if apiKey == "" {
			return nil, fmt.Errorf("AUTOLANG_API_KEY is required for Google Translate")
		}
		return NewGoogle(apiKey), nil
	case "libretranslate", "":
		return NewLibreTranslate(), nil
	default:
		return nil, fmt.Errorf("unknown provider %q (supported: deepl, openai, google, libretranslate)", provider)
	}
}
