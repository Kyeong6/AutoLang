<div align="center">

<h1>AutoLang</h1>

<p>한국어로 작성하면, 영어로 쓴 것과 같은 AI 성능을.</p>

[![CI](https://github.com/Kyeong6/autolang/actions/workflows/ci.yml/badge.svg)](https://github.com/Kyeong6/autolang/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/Kyeong6/autolang)](https://github.com/Kyeong6/autolang/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://go.dev)
[![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20Linux-lightgrey)](#installation)

</div>

---

LLM은 영어로 학습됩니다. 동일한 내용을 한국어로 작성하면 **2~3배 더 많은 토큰**을 소모하고, 응답 품질도 떨어집니다.

AutoLang은 CLI와 AI API 사이에 투명하게 동작하는 **로컬 번역 프록시**입니다. 한국어 입력을 자동으로 영어로 변환해 AI에 전달하고, 응답을 한국어로 돌려줍니다. **기존 워크플로우는 그대로 유지됩니다.**

```
$ claude "이 함수의 시간 복잡도를 O(n log n)으로 줄여줘"

  ✦ AutoLang  KO → EN  (18 tokens → 9 tokens, 50% saved)

  이 함수를 최적화하려면 이중 루프를 정렬 기반 접근으로 교체하세요.

  def optimized(arr: list[int]) -> list[int]:
      arr.sort()
      left, right = 0, len(arr) - 1
      ...

  이렇게 하면 O(n²) → O(n log n)으로 개선됩니다.
```

<br>

## 왜 AutoLang인가

| | AutoLang 없이 | AutoLang 사용 |
|---|---|---|
| 입력 | 영어로 직접 작성하거나 번역 후 붙여넣기 | 한국어로 그대로 작성 |
| 토큰 사용량 | 한국어: 영어 대비 2~3배 소모 | 영어 수준의 토큰 효율 |
| 응답 품질 | 한국어 입력 시 저하됨 | 영어 입력과 동일한 품질 |
| 워크플로우 | 번역 앱과 터미널 사이 컨텍스트 스위칭 | 터미널을 벗어나지 않음 |

**토큰 절약 효과** (Claude Sonnet, 1,000단어 기준)

| | 한국어 | 영어 | 절감 |
|---|---|---|---|
| 토큰 수 | ~2,800 | ~1,300 | **53%** |
| 비용 ($3/1M 기준) | $0.0084 | $0.0039 | **$0.0045 / 메시지** |

<br>

## 기능

- **투명한 프록시** — `claude` 사용 방식을 전혀 바꾸지 않아도 됨
- **자동 언어 감지** — 한국어 입력만 번역, 영어는 그대로 통과
- **코드 블록 보호** — 코드, URL, 파일 경로는 번역하지 않음
- **BYOK** — 원하는 번역 API 키를 직접 연결 (DeepL, OpenAI, Google)
- **토큰 절약 리포트** — 세션별 절약량 확인
- **단일 바이너리** — 런타임 의존성 없음

<br>

## 설치

### Homebrew (권장)

```bash
brew tap Kyeong6/tap
brew install autolang
```

### curl

```bash
curl -fsSL https://raw.githubusercontent.com/Kyeong6/autolang/main/install.sh | sh
```

### go install

```bash
go install github.com/Kyeong6/autolang@latest
```

바이너리 직접 다운로드: [Releases](https://github.com/Kyeong6/autolang/releases)

<br>

## 시작하기

```bash
# 1. 번역 API 키 설정
export AUTOLANG_PROVIDER=deepl       # deepl | openai | google
export AUTOLANG_API_KEY=your-api-key

# 2. Shell 통합 (최초 1회)
echo 'eval "$(autolang init)"' >> ~/.zshrc
source ~/.zshrc

# 3. 평소처럼 사용하면 됩니다
claude "인증 미들웨어 JWT 기반으로 구현해줘"
```

`autolang init`은 터미널을 열 때마다 프록시를 자동으로 백그라운드에서 시작하고, `ANTHROPIC_BASE_URL`을 자동으로 설정합니다.

<br>

## 사용법

```bash
# 한국어 입력 → 영어로 번역 → Claude에 전달 → 한국어로 응답
claude "리팩토링할 때 고려해야 할 점이 뭐야?"

# 영어 입력 → 번역 없이 그대로 전달
claude "what should I consider when refactoring?"

# 프록시 상태 확인
autolang status

# 세션 토큰 절약 통계
autolang stats

# 현재 세션에서 번역 끄기
autolang off
```

<br>

## 설정

`~/.autolang/config.toml`로 동작을 커스터마이즈할 수 있습니다.

```toml
[translation]
provider = "deepl"         # deepl | openai | google | libretranslate
source_lang = "ko"
target_lang = "en"

[proxy]
port = 7878
log_level = "info"         # info | debug | off

[output]
show_token_savings = true
response_lang = "ko"       # ko | en

[rules]
skip_code_blocks = true
skip_urls = true
skip_file_paths = true
```

<br>

## 번역 Provider

| Provider | 키 필요 | 무료 한도 | 품질 | 비고 |
|---|---|---|---|---|
| **DeepL** | 필요 | 500K chars/월 | ★★★★★ | 권장 |
| **OpenAI** | 필요 | 없음 | ★★★★☆ | OpenAI 키 보유 시 |
| **Google** | 필요 | 500K chars/월 | ★★★★☆ | |
| LibreTranslate | 불필요 | 제한 있음 | ★★★☆☆ | 테스트용 |

```bash
export AUTOLANG_PROVIDER=deepl
export AUTOLANG_API_KEY=your-deepl-free-key
```

> DeepL Free 키는 [deepl.com/pro-api](https://www.deepl.com/pro-api)에서 무료로 발급받을 수 있습니다.

<br>

## 지원 CLI

| 도구 | 환경변수 | 상태 |
|---|---|---|
| [Claude Code](https://github.com/anthropics/claude-code) | `ANTHROPIC_BASE_URL` | ✅ v0.1 |
| OpenAI CLI / Aider | `OPENAI_BASE_URL` | 🔜 v0.2 |
| Gemini CLI | `GOOGLE_AI_BASE_URL` | 🔜 v0.2 |
| OpenAI 호환 API | Custom base URL | 🔜 v0.2 |

<br>

## 기여

기여는 언제나 환영입니다. 시작하기 전에 [CONTRIBUTING.md](CONTRIBUTING.md)를 읽어주세요.

```bash
git clone https://github.com/Kyeong6/autolang
cd autolang
go build ./...
go test ./...
```

<br>

## 라이선스

MIT — [LICENSE](LICENSE)

---

<div align="center">
AutoLang이 도움이 됐다면 ⭐을 눌러주세요. 더 많은 개발자에게 닿을 수 있습니다.
</div>
