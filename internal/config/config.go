package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/BurntSushi/toml"
)

// Config holds all AutoLang configuration.
type Config struct {
	Translation TranslationConfig `toml:"translation"`
	Proxy       ProxyConfig       `toml:"proxy"`
	Output      OutputConfig      `toml:"output"`
	Rules       RulesConfig       `toml:"rules"`
}

type TranslationConfig struct {
	Provider   string `toml:"provider"`
	APIKey     string `toml:"api_key"`
	SourceLang string `toml:"source_lang"`
	TargetLang string `toml:"target_lang"`
}

type ProxyConfig struct {
	Port     int    `toml:"port"`
	LogLevel string `toml:"log_level"`
}

type OutputConfig struct {
	ShowTokenSavings         bool   `toml:"show_token_savings"`
	ShowTranslationIndicator bool   `toml:"show_translation_indicator"`
	ResponseLang             string `toml:"response_lang"`
}

type RulesConfig struct {
	SkipCodeBlocks bool `toml:"skip_code_blocks"`
	SkipURLs       bool `toml:"skip_urls"`
	SkipFilePaths  bool `toml:"skip_file_paths"`
}

func Default() *Config {
	return &Config{
		Translation: TranslationConfig{
			Provider:   "libretranslate",
			SourceLang: "ko",
			TargetLang: "en",
		},
		Proxy: ProxyConfig{
			Port:     7878,
			LogLevel: "info",
		},
		Output: OutputConfig{
			ShowTokenSavings:         true,
			ShowTranslationIndicator: true,
			ResponseLang:             "ko",
		},
		Rules: RulesConfig{
			SkipCodeBlocks: true,
			SkipURLs:       true,
			SkipFilePaths:  true,
		},
	}
}

const defaultConfigTOML = `# AutoLang configuration
# https://github.com/Kyeong6/autolang

[translation]
provider = "libretranslate"   # deepl | openai | google | libretranslate
# api_key = ""                # set via AUTOLANG_API_KEY env var
source_lang = "ko"
target_lang = "en"

[proxy]
port = 7878
log_level = "info"            # info | debug | off

[output]
show_token_savings = true
show_translation_indicator = true
response_lang = "ko"          # language for Claude's responses

[rules]
skip_code_blocks = true
skip_urls = true
skip_file_paths = true
`

// Load reads ~/.autolang/config.toml, creates it with defaults if absent,
// then applies AUTOLANG_* environment variable overrides.
func Load() (*Config, error) {
	cfg := Default()

	path, err := configFilePath()
	if err != nil {
		// Cannot determine home dir; fall back to defaults + env.
		applyEnvOverrides(cfg)
		return cfg, nil
	}

	if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
		if writeErr := writeDefaultConfig(path); writeErr != nil {
			// Non-fatal: proceed with in-memory defaults.
			fmt.Fprintf(os.Stderr, "autolang: could not write default config: %v\n", writeErr)
		}
	} else if statErr == nil {
		if _, decodeErr := parseFile(path, cfg); decodeErr != nil {
			return nil, fmt.Errorf("invalid config file %s: %w", path, decodeErr)
		}
	}

	applyEnvOverrides(cfg)
	return cfg, nil
}

func configFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".autolang", "config.toml"), nil
}

func parseFile(path string, cfg *Config) (toml.MetaData, error) {
	return toml.DecodeFile(path, cfg)
}

func writeDefaultConfig(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(defaultConfigTOML), 0o644)
}

// applyEnvOverrides applies AUTOLANG_* environment variables over the loaded config.
// Environment variables always take precedence over the config file.
func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("AUTOLANG_PROVIDER"); v != "" {
		cfg.Translation.Provider = v
	}
	if v := os.Getenv("AUTOLANG_API_KEY"); v != "" {
		cfg.Translation.APIKey = v
	}
	if v := os.Getenv("AUTOLANG_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			cfg.Proxy.Port = p
		}
	}
	if v := os.Getenv("AUTOLANG_LOG"); v != "" {
		cfg.Proxy.LogLevel = v
	}
}
