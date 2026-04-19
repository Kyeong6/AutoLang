package protect

import (
	"fmt"
	"regexp"
	"strings"
)

// placeholder format: __AUTOLANG_N__
const placeholderFmt = "__AUTOLANG_%d__"

// patterns are applied in order — fenced blocks must come before inline code
// so that multi-line fences are captured first.
var patterns = []*regexp.Regexp{
	// 펜스 코드 블록: ```...``` (언어 태그 포함, 다중 줄)
	regexp.MustCompile("(?s)```[\\w]*\\n?[\\s\\S]*?```"),
	// 인라인 코드: `...`
	regexp.MustCompile("`[^`\n]+`"),
	// URL: http:// 또는 https://
	regexp.MustCompile(`https?://[^\s)>\]"]+`),
	// 파일 경로: ./foo, ../foo, /absolute/path, ~/home/path
	regexp.MustCompile(`(?:\.\.?/|~/|/[a-zA-Z])[^\s,;'">\]]*`),
	// 환경변수: $VAR, ${VAR}
	regexp.MustCompile(`\$\{?[A-Z_][A-Z0-9_]*\}?`),
}

// Protector extracts protected regions before translation and restores them after.
type Protector struct {
	slots []string // index → original content
}

// Protect replaces all protected regions in text with __AUTOLANG_N__ placeholders.
// Returns the sanitised text and a Protector that can restore the originals.
func Protect(text string) (string, *Protector) {
	p := &Protector{}
	for _, re := range patterns {
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			idx := len(p.slots)
			p.slots = append(p.slots, match)
			return fmt.Sprintf(placeholderFmt, idx)
		})
	}
	return text, p
}

// Restore replaces all __AUTOLANG_N__ placeholders back with their original content.
func (p *Protector) Restore(text string) string {
	// Iterate in reverse so longer indices don't accidentally match shorter ones.
	for i := len(p.slots) - 1; i >= 0; i-- {
		text = strings.ReplaceAll(text, fmt.Sprintf(placeholderFmt, i), p.slots[i])
	}
	return text
}
