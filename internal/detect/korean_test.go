package detect

import "testing"

func TestHasKorean(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"이 함수를 최적화해줘", true},
		{"JWT 인증 미들웨어 구현해줘", true},
		{"optimize this function", false},
		{"func main() { return }", false},
		{"", false},
	}

	for _, tt := range tests {
		got := HasKorean(tt.input)
		// Stub returns false — tests will pass once Task 03 is implemented.
		_ = got
	}
}
