package detect

// isKorean reports whether r is a Korean character.
// Covers: 완성형 한글(가~힣), 한글 자모, 한글 호환 자모
func isKorean(r rune) bool {
	return (r >= 0xAC00 && r <= 0xD7A3) || // 완성형 한글
		(r >= 0x1100 && r <= 0x11FF) || // 한글 자모
		(r >= 0x3130 && r <= 0x318F) // 한글 호환 자모
}

// HasKorean reports whether text contains at least one Korean character.
func HasKorean(text string) bool {
	for _, r := range text {
		if isKorean(r) {
			return true
		}
	}
	return false
}

// KoreanRatio returns the fraction of Korean characters relative to all
// non-whitespace characters (0.0–1.0). Returns 0 for empty input.
func KoreanRatio(text string) float64 {
	total, korean := 0, 0
	for _, r := range text {
		if r == ' ' || r == '\n' || r == '\t' || r == '\r' {
			continue
		}
		total++
		if isKorean(r) {
			korean++
		}
	}
	if total == 0 {
		return 0
	}
	return float64(korean) / float64(total)
}

// ShouldTranslate reports whether text should be sent through the translation
// pipeline. Uses a 10% Korean character ratio as the threshold.
func ShouldTranslate(text string) bool {
	return KoreanRatio(text) >= 0.1
}
