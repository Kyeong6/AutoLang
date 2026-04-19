package stats

import "testing"

func TestEstimateTokens(t *testing.T) {
	koTokens := estimateTokens("이 함수를 최적화해줘", true)
	enTokens := estimateTokens("optimize this function", false)

	if koTokens <= 0 {
		t.Errorf("expected positive Korean token count, got %d", koTokens)
	}
	if enTokens <= 0 {
		t.Errorf("expected positive English token count, got %d", enTokens)
	}
}
