<div align="center">

# AutoLang

**Write in Korean. Get English-quality AI responses.**

[![CI](https://github.com/Kyeong6/autolang/actions/workflows/ci.yml/badge.svg)](https://github.com/Kyeong6/autolang/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/Kyeong6/autolang)](https://github.com/Kyeong6/autolang/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20Linux-lightgrey)](#installation)

</div>

---

LLMs are trained on English. Korean uses **2–3× more tokens** for the same content — meaning higher costs, shorter context windows, and lower response quality.

**AutoLang** is a transparent local proxy that sits between your CLI and AI APIs. It automatically detects Korean input, translates it to English before sending, and returns responses in Korean — **with zero changes to your workflow**.

```
$ claude "이 함수의 시간 복잡도를 O(n log n)으로 줄여줘"

  ✦ [AutoLang] KO → EN  (tokens: 18 → 9, saved 50%)

  이 함수를 최적화하려면 이중 루프를 제거하고 정렬 기반 접근으로 바꾸세요.

  def optimized(arr: list[int]) -> list[int]:
      arr.sort()
      ...

  이렇게 하면 O(n²) → O(n log n)으로 개선됩니다.
```

> Demo GIF — coming soon

---

## Why AutoLang?

| | Without AutoLang | With AutoLang |
|-|-----------------|---------------|
| Input language | Must write in English | Write in Korean |
| Token usage | 2–3× more tokens for Korean | English-level token efficiency |
| Response quality | Degraded for non-English | Full English-model quality |
| Workflow | Context switch to translator | Stay in the terminal |

**Token savings example** (Claude Sonnet, 1,000 words):

| | Korean | English | Savings |
|-|--------|---------|---------|
| Token count | ~2,800 | ~1,300 | **53%** |
| Cost (at $3/1M) | $0.0084 | $0.0039 | **$0.0045/msg** |

---

## Features

- **Transparent proxy** — no changes to how you use `claude` or other AI CLIs
- **Korean detection** — automatically detects Korean input; passes English through unchanged
- **Code block protection** — code, file paths, and variable names are never translated
- **BYOK** — bring your own translation API key (DeepL, OpenAI, etc.)
- **Token savings report** — see how many tokens you saved per session
- **Single binary** — no runtime dependencies, fast startup

---

## Installation

### Homebrew (recommended)

```bash
brew tap Kyeong6/tap
brew install autolang
```

### curl (Linux / quick install)

```bash
curl -fsSL https://raw.githubusercontent.com/Kyeong6/autolang/main/install.sh | sh
```

### go install

```bash
go install github.com/Kyeong6/autolang@latest
```

### Manual

Download the binary for your platform from [Releases](https://github.com/Kyeong6/autolang/releases).

---

## Setup

```bash
# 1. Set your translation API key (required)
export AUTOLANG_API_KEY="your-deepl-or-openai-key"
export AUTOLANG_PROVIDER="deepl"   # deepl | openai | google (default: deepl)

# 2. Add AutoLang to your shell (one-time setup)
echo 'eval "$(autolang init)"' >> ~/.zshrc
source ~/.zshrc

# 3. That's it — use Claude Code as usual
claude "이 코드 리뷰해줘"
```

**Shell init** automatically:
- Starts the AutoLang proxy in the background
- Sets `ANTHROPIC_BASE_URL` to point to the local proxy
- Restarts the proxy if it's not running

---

## Usage

No new commands. Use your AI CLI tools exactly as before.

```bash
# Korean input → translated to English → response returned in Korean
claude "인증 미들웨어 구현해줘, JWT 기반으로"

# English input → passed through unchanged
claude "implement a JWT authentication middleware"

# View session token savings
autolang stats

# Toggle Korean output off for current session
autolang off

# Check proxy status
autolang status
```

---

## Configuration

AutoLang is configured via `~/.autolang/config.toml`:

```toml
[translation]
provider = "deepl"        # deepl | openai | google | libretranslate
source_lang = "ko"
target_lang = "en"

[proxy]
port = 7878
log_level = "info"        # info | debug | off

[output]
show_token_savings = true
show_translation_indicator = true
response_lang = "ko"      # ko | en (language for AI responses)

[rules]
skip_code_blocks = true
skip_urls = true
skip_file_paths = true
```

Override any option with CLI flags:

```bash
claude --autolang-off "no translation for this message"
```

---

## Translation Providers

AutoLang supports multiple translation backends via `AUTOLANG_PROVIDER`:

| Provider | Key Required | Free Tier | Quality | Best For |
|----------|-------------|-----------|---------|----------|
| **DeepL** | Yes | 500K chars/month | ★★★★★ | Best quality (recommended) |
| **OpenAI** | Yes | No | ★★★★☆ | Already have OpenAI key |
| **Google** | Yes | 500K chars/month | ★★★★☆ | Multi-language |
| LibreTranslate | No | Unlimited* | ★★★☆☆ | No key, testing |

\* Public instance has rate limits. Self-host for unlimited usage.

**Setting your provider:**

```bash
# DeepL (recommended)
export AUTOLANG_PROVIDER=deepl
export AUTOLANG_API_KEY=your-deepl-free-key

# OpenAI (if you already have a key)
export AUTOLANG_PROVIDER=openai
export AUTOLANG_API_KEY=your-openai-key

# No key (free, lower quality)
export AUTOLANG_PROVIDER=libretranslate
```

---

## Supported AI CLIs

| Tool | Environment Variable | Status |
|------|---------------------|--------|
| [Claude Code](https://github.com/anthropics/claude-code) | `ANTHROPIC_BASE_URL` | ✅ Supported |
| OpenAI CLI | `OPENAI_BASE_URL` | 🔜 Coming in v0.2 |
| [Aider](https://github.com/paul-gauthier/aider) | `OPENAI_API_BASE` | 🔜 Coming in v0.2 |
| Any OpenAI-compatible API | Custom base URL | 🔜 Coming in v0.2 |

---

## How It Works

```
┌─────────────────────────────────────────────────────┐
│                    Your Terminal                      │
│                                                       │
│  claude "한글 입력"                                   │
│       │                                               │
│       ▼                                               │
│  AutoLang Proxy (localhost:7878)                      │
│       │                                               │
│   [1] Korean detected                                 │
│   [2] Translate KO → EN via API                       │
│   [3] Forward English request to Claude API           │
│       │                                               │
│       ▼                                               │
│  Claude API (api.anthropic.com)                       │
│       │                                               │
│   [4] English response received                       │
│   [5] Translate natural language EN → KO              │
│   [6] Code blocks preserved as-is                     │
│       │                                               │
│       ▼                                               │
│  Korean output in your terminal                       │
└─────────────────────────────────────────────────────┘
```

The proxy is set via `ANTHROPIC_BASE_URL=http://localhost:7878` — Claude Code sends requests there, AutoLang handles the rest transparently.

---

## Roadmap

- [x] Claude Code support (v0.1)
- [x] DeepL / OpenAI / Google translation backends
- [x] Code block protection
- [x] Token savings report
- [ ] OpenAI / Aider / Gemini CLI support (v0.2)
- [ ] Japanese and Chinese support (v0.3)
- [ ] Streaming translation (word-by-word output)
- [ ] Translation history and search
- [ ] Homebrew core submission (at 75+ stars)

---

## Contributing

Contributions are welcome. Please read [CONTRIBUTING.md](CONTRIBUTING.md) first.

```bash
git clone https://github.com/Kyeong6/autolang
cd autolang
cargo build
cargo test
```

---

## License

MIT — see [LICENSE](LICENSE).

---

<div align="center">

If AutoLang saves you tokens, a ⭐ helps others find it.

</div>
