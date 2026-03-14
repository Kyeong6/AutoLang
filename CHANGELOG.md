# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial Claude Code support via local proxy
- Korean ↔ English automatic translation
- Code block protection (code is never translated)
- BYOK translation API support (DeepL, OpenAI, Google, LibreTranslate)
- Token savings report per session
- `eval "$(autolang init)"` shell integration
- `~/.autolang/config.toml` configuration
- `autolang stats` — session token savings summary
- `autolang status` — proxy health check
