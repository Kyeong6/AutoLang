package detect

import (
	"testing"
)

func TestHasKorean(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		// 한글 포함
		{"이 함수를 최적화해줘", true},
		{"JWT 인증 미들웨어 구현해줘", true}, // 영어 + 한글 혼용
		{"함수", true},
		// 영어만
		{"optimize this function", false},
		{"func main() { return nil }", false},
		// 경계 케이스
		{"", false},
		{"123 !@# $%^", false},
		{"https://github.com/Kyeong6/autolang", false},
	}

	for _, tt := range tests {
		got := HasKorean(tt.input)
		if got != tt.want {
			t.Errorf("HasKorean(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestKoreanRatio(t *testing.T) {
	tests := []struct {
		input    string
		wantGte  float64 // ratio >= wantGte
		wantLte  float64 // ratio <= wantLte
	}{
		{"이 함수를 최적화해줘", 0.8, 1.0},  // 거의 전부 한글
		{"JWT 인증", 0.3, 0.8},            // 한글 + 영어 혼용
		{"optimize this", 0.0, 0.0},      // 영어만 → 0
		{"", 0.0, 0.0},                   // 빈 문자열 → 0
	}

	for _, tt := range tests {
		got := KoreanRatio(tt.input)
		if got < tt.wantGte || got > tt.wantLte {
			t.Errorf("KoreanRatio(%q) = %.2f, want [%.2f, %.2f]", tt.input, got, tt.wantGte, tt.wantLte)
		}
	}
}

func TestShouldTranslate(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"이 함수를 최적화해줘", true},
		{"JWT 인증 미들웨어 구현해줘", true},
		{"optimize this function", false},
		{"func main() {}", false},
		{"", false},
		// 한글 비율 10% 미만 — 영어 코드에 한글 주석 한 글자 수준
		{"hello world 가", false},
	}

	for _, tt := range tests {
		got := ShouldTranslate(tt.input)
		if got != tt.want {
			t.Errorf("ShouldTranslate(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}
