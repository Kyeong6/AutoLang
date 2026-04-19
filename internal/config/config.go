package config

// Config holds all AutoLang configuration.
// Implementation: Task 07
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

// Load reads config from ~/.autolang/config.toml.
// Falls back to Default() if the file does not exist.
// Implementation: Task 07
func Load() (*Config, error) {
	return Default(), nil
}
