package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()
	if cfg.Translation.Provider != "libretranslate" {
		t.Errorf("unexpected default provider: %s", cfg.Translation.Provider)
	}
	if cfg.Proxy.Port != 7878 {
		t.Errorf("unexpected default port: %d", cfg.Proxy.Port)
	}
	if !cfg.Rules.SkipCodeBlocks {
		t.Error("expected SkipCodeBlocks to be true by default")
	}
}

func TestLoad_CreatesDefaultConfigFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")

	if err := writeDefaultConfig(path); err != nil {
		t.Fatalf("writeDefaultConfig: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected config file to be created: %v", err)
	}
}

func TestLoad_ParsesTOMLFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")

	tomlContent := `
[translation]
provider = "deepl"
source_lang = "ko"
target_lang = "en"

[proxy]
port = 9090
log_level = "debug"
`
	if err := os.WriteFile(path, []byte(tomlContent), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := Default()
	if _, err := parseFile(path, cfg); err != nil {
		t.Fatalf("parseFile: %v", err)
	}
	if cfg.Translation.Provider != "deepl" {
		t.Errorf("expected deepl, got %s", cfg.Translation.Provider)
	}
	if cfg.Proxy.Port != 9090 {
		t.Errorf("expected port 9090, got %d", cfg.Proxy.Port)
	}
	if cfg.Proxy.LogLevel != "debug" {
		t.Errorf("expected debug, got %s", cfg.Proxy.LogLevel)
	}
}

func TestApplyEnvOverrides(t *testing.T) {
	t.Setenv("AUTOLANG_PROVIDER", "openai")
	t.Setenv("AUTOLANG_API_KEY", "sk-test")
	t.Setenv("AUTOLANG_PORT", "8080")
	t.Setenv("AUTOLANG_LOG", "debug")

	cfg := Default()
	applyEnvOverrides(cfg)

	if cfg.Translation.Provider != "openai" {
		t.Errorf("expected openai, got %s", cfg.Translation.Provider)
	}
	if cfg.Translation.APIKey != "sk-test" {
		t.Errorf("expected sk-test, got %s", cfg.Translation.APIKey)
	}
	if cfg.Proxy.Port != 8080 {
		t.Errorf("expected 8080, got %d", cfg.Proxy.Port)
	}
	if cfg.Proxy.LogLevel != "debug" {
		t.Errorf("expected debug, got %s", cfg.Proxy.LogLevel)
	}
}

func TestEnvOverridesFileConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	os.WriteFile(path, []byte("[translation]\nprovider = \"deepl\"\n"), 0o644) //nolint:errcheck

	t.Setenv("AUTOLANG_PROVIDER", "google")

	cfg := Default()
	if _, err := parseFile(path, cfg); err != nil {
		t.Fatal(err)
	}
	applyEnvOverrides(cfg)

	if cfg.Translation.Provider != "google" {
		t.Errorf("env should override file; got %s", cfg.Translation.Provider)
	}
}
