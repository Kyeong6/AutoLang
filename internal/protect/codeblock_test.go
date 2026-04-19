package protect

import (
	"strings"
	"testing"
)

// roundTrip protects text, simulates a no-op "translation", and restores.
// The restored result must equal the original.
func roundTrip(text string) string {
	protected, p := Protect(text)
	return p.Restore(protected)
}

func TestFencedCodeBlock(t *testing.T) {
	input := "다음 코드를 리뷰해줘:\n```python\ndef hello():\n    return \"world\"\n```\n고마워"

	protected, p := Protect(input)

	if strings.Contains(protected, "def hello") {
		t.Error("fenced code block should be replaced with placeholder")
	}
	if !strings.Contains(protected, "__AUTOLANG_") {
		t.Error("placeholder not found in protected text")
	}

	restored := p.Restore(protected)
	if restored != input {
		t.Errorf("roundtrip failed\ngot:  %q\nwant: %q", restored, input)
	}
}

func TestInlineCode(t *testing.T) {
	input := "변수 `myVar`를 확인해줘"

	protected, _ := Protect(input)
	if strings.Contains(protected, "myVar") {
		t.Error("inline code should be replaced with placeholder")
	}

	if roundTrip(input) != input {
		t.Error("inline code roundtrip failed")
	}
}

func TestURLProtection(t *testing.T) {
	input := "https://docs.python.org 참고해줘"

	protected, p := Protect(input)
	if strings.Contains(protected, "docs.python.org") {
		t.Error("URL should be replaced with placeholder")
	}

	restored := p.Restore(protected)
	if restored != input {
		t.Errorf("URL roundtrip failed\ngot:  %q\nwant: %q", restored, input)
	}
}

func TestFilePathProtection(t *testing.T) {
	cases := []string{
		"./src/main.go 파일을 봐줘",
		"../config/settings.toml 열어줘",
		"/usr/local/bin/autolang 실행해줘",
		"~/Desktop/project 폴더야",
	}

	for _, input := range cases {
		if roundTrip(input) != input {
			t.Errorf("file path roundtrip failed for %q", input)
		}
	}
}

func TestEnvVarProtection(t *testing.T) {
	cases := []string{
		"$HOME 디렉토리 확인해줘",
		"${ANTHROPIC_API_KEY} 키 설정해줘",
		"$PATH에 추가해줘",
	}

	for _, input := range cases {
		protected, p := Protect(input)
		_ = p
		if strings.Contains(protected, "$") {
			t.Errorf("env var should be replaced in %q, got %q", input, protected)
		}
		if roundTrip(input) != input {
			t.Errorf("env var roundtrip failed for %q", input)
		}
	}
}

func TestMixedContent(t *testing.T) {
	input := "다음 코드를 리뷰해줘:\n\n```go\nfunc main() {}\n```\n\n" +
		"참고: https://pkg.go.dev/net/http\n" +
		"`$HOME`에 있는 ./config.toml 파일도 확인해줘"

	restored := roundTrip(input)
	if restored != input {
		t.Errorf("mixed content roundtrip failed\ngot:  %q\nwant: %q", restored, input)
	}
}

func TestNoProtectedContent(t *testing.T) {
	input := "이 함수의 성능을 개선해줘"
	protected, p := Protect(input)

	if protected != input {
		t.Errorf("plain Korean text should not be modified, got %q", protected)
	}
	if p.Restore(protected) != input {
		t.Error("restore of unmodified text should be a no-op")
	}
}

func TestEmptyInput(t *testing.T) {
	protected, p := Protect("")
	if protected != "" {
		t.Error("empty input should return empty string")
	}
	if p.Restore("") != "" {
		t.Error("restore of empty string should return empty string")
	}
}

func TestMultiplePlaceholders_RestoredInOrder(t *testing.T) {
	input := "`foo` 그리고 `bar` 확인해줘"
	protected, p := Protect(input)

	if !strings.Contains(protected, "__AUTOLANG_0__") || !strings.Contains(protected, "__AUTOLANG_1__") {
		t.Errorf("expected two placeholders, got: %q", protected)
	}

	restored := p.Restore(protected)
	if restored != input {
		t.Errorf("multi-placeholder roundtrip failed\ngot:  %q\nwant: %q", restored, input)
	}
}
