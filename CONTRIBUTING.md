# Contributing to AutoLang

Thank you for your interest in contributing. This document covers how to get started.

## Ways to Contribute

- **Bug reports** — open an issue with the bug report template
- **Feature requests** — open an issue with the feature request template
- **Code** — fix a bug, implement a feature from the roadmap, improve tests
- **Documentation** — fix typos, improve examples, add translations
- **Translation providers** — add support for a new translation API

## Development Setup

### Prerequisites

- Go 1.22+
- macOS (primary development platform)

### Build

```bash
git clone https://github.com/YOUR_USERNAME/autolang
cd autolang
go build ./...
```

### Run tests

```bash
go test ./...
```

### Run locally

```bash
# Set a translation API key
export AUTOLANG_PROVIDER=deepl
export AUTOLANG_API_KEY=your-key

# Start the proxy
go run . start

# In another terminal, test it
export ANTHROPIC_BASE_URL=http://localhost:7878
claude "테스트 메시지"
```

## Project Structure

```
autolang/
├── cmd/
│   └── autolang/
│       └── main.go          # CLI entry point (cobra)
├── internal/
│   ├── proxy/
│   │   ├── proxy.go         # Proxy server (net/http)
│   │   └── anthropic.go     # Claude API format handling
│   ├── translate/
│   │   ├── translator.go    # Translator interface
│   │   ├── deepl.go         # DeepL backend
│   │   ├── openai.go        # OpenAI backend
│   │   ├── google.go        # Google Translate backend
│   │   └── libre.go         # LibreTranslate backend
│   ├── detect/
│   │   └── korean.go        # Korean language detection
│   ├── protect/
│   │   └── codeblock.go     # Code block protection
│   ├── config/
│   │   └── config.go        # Config file management
│   └── stats/
│       └── stats.go         # Token savings tracking
├── tests/
├── go.mod
├── go.sum
└── README.md
```

## Pull Request Guidelines

1. **One PR per change** — keep PRs small and focused
2. **Tests** — add tests for new behavior; existing tests must pass
3. **No breaking changes** without discussion in an issue first
4. **Commit messages** — use conventional commits (`feat:`, `fix:`, `docs:`, etc.)

### PR checklist

- [ ] `go test ./...` passes
- [ ] `go vet ./...` passes with no warnings
- [ ] `gofmt` applied
- [ ] Updated `CHANGELOG.md` if applicable
- [ ] Added/updated tests for the change

## Adding a Translation Provider

1. Create `internal/translate/your_provider.go`
2. Implement the `Translator` interface:

```go
type Translator interface {
    Translate(ctx context.Context, text, from, to string) (string, error)
}
```

3. Register it in `internal/translate/translator.go`
4. Add the provider name to `Config` and README

## Reporting Bugs

Use the [bug report template](.github/ISSUE_TEMPLATE/bug_report.md). Include:
- AutoLang version (`autolang --version`)
- macOS version
- Which AI CLI you're using
- Steps to reproduce

## Code of Conduct

By participating, you agree to the [Code of Conduct](CODE_OF_CONDUCT.md).
