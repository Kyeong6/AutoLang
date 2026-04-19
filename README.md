<div align="center">

<h1>AutoLang</h1>

<p><strong>한국어로 작성하면, 영어로 쓴 것과 같은 AI 성능을.</strong><br>
<sub>A transparent local proxy that translates your Korean prompts to English — invisibly.</sub></p>

[![CI](https://github.com/Kyeong6/autolang/actions/workflows/ci.yml/badge.svg)](https://github.com/Kyeong6/autolang/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/Kyeong6/autolang?color=brightgreen)](https://github.com/Kyeong6/autolang/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://go.dev)
[![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20Linux-lightgrey)](#installation)

</div>

---

LLM은 영어 데이터로 학습됩니다. 동일한 내용을 한국어로 작성하면 **토큰을 2~3배 더 소모**하고, 응답 품질도 떨어집니다.

AutoLang은 CLI와 AI API 사이에서 투명하게 동작하는 **로컬 HTTP 프록시**입니다.  
한국어 입력을 감지해 자동으로 영어로 번역한 뒤 AI에 전달하고, 응답을 다시 한국어로 반환합니다.  
**기존 워크플로우는 전혀 바꾸지 않아도 됩니다.**

```
you          autolang proxy         Claude API
  │                │                    │
  │ "이 함수 O(n²)인데   │                    │
  │  O(n log n)으로   │                    │
  │  줄여줄 수 있어?" ──▶│ "Can you reduce the │
  │                │  complexity of this  │
  │                │  function from O(n²) │
  │                │  to O(n log n)?"   ──▶│
  │                │                    │ (English response)
  │ "네, 이중 루프를  ◀──│◀───────────────────│
  │  정렬 기반으로..."   │                    │
```

<br>

## 왜 AutoLang인가

**토큰 절약 효과** (Claude Sonnet, 1,000단어 기준)

| | 한국어 직접 입력 | AutoLang 사용 | 절감 |
|---|---|---|---|
| 토큰 수 | ~2,800 | ~1,300 | **53%** |
| 비용 ($3/1M 기준) | $0.0084 | $0.0039 | **$0.0045/메시지** |

| | AutoLang 없이 | AutoLang 사용 |
|---|---|---|
| 입력 방식 | 영어로 쓰거나, 번역 후 붙여넣기 | 한국어로 그대로 작성 |
| 응답 품질 | 한국어 입력 시 저하 | 영어 입력과 동일한 품질 |
| 워크플로우 | 번역 앱 ↔ 터미널 컨텍스트 스위칭 | 터미널을 벗어나지 않음 |

<br>

## 주요 기능

- **투명한 프록시** — `claude` 사용 방식을 그대로 유지
- **자동 언어 감지** — 한국어 비율이 10% 미만이면 번역 없이 통과
- **코드 블록 보호** — 코드, URL, 파일 경로, 환경변수는 번역하지 않음
- **스트리밍 지원** — Claude SSE 스트리밍 응답을 실시간으로 번역
- **BYOK** — DeepL · OpenAI · Google · LibreTranslate 중 선택
- **토큰 절약 통계** — `autolang stats`로 세션별 절약량 확인
- **단일 바이너리** — 런타임 의존성 없음, `brew install` 한 줄로 설치

<br>

## 설치

### Homebrew (권장)

```bash
brew tap Kyeong6/tap
brew install autolang
```

### curl (macOS / Linux)

```bash
curl -fsSL https://raw.githubusercontent.com/Kyeong6/autolang/main/install.sh | sh
```

### go install

```bash
go install github.com/Kyeong6/autolang@latest
```

바이너리 직접 다운로드 → [Releases](https://github.com/Kyeong6/autolang/releases)

<br>

## 시작하기

### 1. 번역 API 키 설정

```bash
export AUTOLANG_PROVIDER=deepl        # deepl | openai | google | libretranslate
export AUTOLANG_API_KEY=your-api-key  # libretranslate는 키 불필요
```

> DeepL Free 키는 [deepl.com/pro-api](https://www.deepl.com/pro-api)에서 무료로 발급받을 수 있습니다 (500K chars/월).

### 2. Shell 통합 (최초 1회)

```bash
echo 'eval "$(autolang init)"' >> ~/.zshrc
source ~/.zshrc
```

터미널을 열 때마다 프록시가 자동으로 백그라운드에서 시작되고 `ANTHROPIC_BASE_URL`이 설정됩니다.

### 3. 평소처럼 사용

```bash
claude "인증 미들웨어를 JWT 기반으로 구현해줘"
# ✦ [AutoLang] KO → EN  (tokens: 12 → 5, saved 58%)
# → Claude가 영어로 최고 품질 응답 → AutoLang이 한국어로 반환
```

<br>

## 사용법

```bash
# 프록시 시작 (포그라운드)
autolang start

# 프록시 백그라운드 시작
autolang start --daemon

# 프록시 상태 확인
autolang status

# 프록시 종료
autolang stop

# 세션 토큰 절약 통계
autolang stats

# 버전 확인
autolang --version
```

<br>

## 설정

`~/.autolang/config.toml` — 최초 실행 시 자동으로 생성됩니다.

```toml
[translation]
provider = "libretranslate"   # deepl | openai | google | libretranslate
# api_key = ""                # AUTOLANG_API_KEY 환경변수로도 설정 가능
source_lang = "ko"
target_lang = "en"

[proxy]
port = 7878
log_level = "info"            # info | debug | off

[output]
show_token_savings = true
show_translation_indicator = true
response_lang = "ko"

[rules]
skip_code_blocks = true
skip_urls = true
skip_file_paths = true
```

**환경변수 우선순위**: 환경변수 > config.toml > 기본값

| 환경변수 | 설명 |
|---|---|
| `AUTOLANG_PROVIDER` | 번역 provider |
| `AUTOLANG_API_KEY` | API 키 |
| `AUTOLANG_PORT` | 프록시 포트 (기본값: 7878) |
| `AUTOLANG_LOG` | 로그 레벨 |

<br>

## 번역 Provider

| Provider | 키 필요 | 무료 한도 | 품질 |
|---|---|---|---|
| **DeepL** ⭐ | 필요 | 500K chars/월 | 최고 |
| **OpenAI** | 필요 | — | 우수 |
| **Google** | 필요 | 500K chars/월 | 우수 |
| **LibreTranslate** | 불필요 | Rate limit 있음 | 기본 (테스트용) |

<br>

## 지원 CLI

| 도구 | 상태 |
|---|---|
| [Claude Code](https://github.com/anthropics/claude-code) | ✅ v0.1 지원 |
| OpenAI CLI · Aider | 🔜 v0.2 예정 |
| Gemini CLI | 🔜 v0.2 예정 |

<br>

## 기여

기여는 언제나 환영입니다. 시작하기 전에 [CONTRIBUTING.md](CONTRIBUTING.md)를 읽어주세요.

```bash
git clone https://github.com/Kyeong6/autolang
cd autolang
go build ./...
go test ./...
```

버그 리포트 · 기능 요청 → [Issues](https://github.com/Kyeong6/autolang/issues)

<br>

## 라이선스

MIT — [LICENSE](LICENSE)

---

<div align="center">
<sub>AutoLang이 도움이 됐다면 ⭐을 눌러주세요. 더 많은 한국 개발자에게 닿을 수 있습니다.</sub>
</div>
