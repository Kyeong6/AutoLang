package stats

import "testing"

func TestEstimateKorean(t *testing.T) {
	// "이 함수를 최적화해줘" = 10 chars → 10 × 2.5 = 25
	got := estimateKorean("이 함수를 최적화해줘")
	if got <= 0 {
		t.Errorf("expected positive Korean token count, got %d", got)
	}
}

func TestEstimateEnglish(t *testing.T) {
	// "optimize this function" = 3 words → 3 × 1.3 = 3 (int)
	got := estimateEnglish("optimize this function")
	if got <= 0 {
		t.Errorf("expected positive English token count, got %d", got)
	}
}

func TestEstimateEmpty(t *testing.T) {
	if estimateKorean("") != 0 {
		t.Error("empty Korean should be 0")
	}
	if estimateEnglish("") != 0 {
		t.Error("empty English should be 0")
	}
}

func TestRecordInput(t *testing.T) {
	s := New()
	s.RecordInput("안녕하세요 세계", "Hello world")
	if s.Requests != 1 {
		t.Errorf("expected 1 request, got %d", s.Requests)
	}
	if s.InputKO <= 0 {
		t.Error("expected InputKO > 0")
	}
	if s.InputEN <= 0 {
		t.Error("expected InputEN > 0")
	}
}

func TestRecordOutput(t *testing.T) {
	s := New()
	s.RecordOutput("Hello world", "안녕하세요 세계")
	if s.OutputEN <= 0 {
		t.Error("expected OutputEN > 0")
	}
	if s.OutputKO <= 0 {
		t.Error("expected OutputKO > 0")
	}
}

func TestSummaryNoRequests(t *testing.T) {
	s := New()
	out := s.Summary()
	if out == "" {
		t.Error("Summary should return non-empty string")
	}
}

func TestSummaryWithRequests(t *testing.T) {
	s := New()
	s.RecordInput("이 코드를 리팩토링해줘", "refactor this code")
	s.RecordOutput("Here is the refactored code.", "리팩토링된 코드입니다.")
	out := s.Summary()
	if out == "" {
		t.Error("Summary should not be empty")
	}
	// Should contain token numbers
	if len(out) < 50 {
		t.Errorf("Summary too short: %q", out)
	}
}

func TestIndicator(t *testing.T) {
	s := New()
	line := s.Indicator(20, 8)
	if line == "" {
		t.Error("Indicator should return non-empty string")
	}
}

func TestFormatInt(t *testing.T) {
	cases := map[int]string{
		0:    "0",
		999:  "999",
		1000: "1,000",
		1234: "1,234",
	}
	for n, want := range cases {
		if got := formatInt(n); got != want {
			t.Errorf("formatInt(%d) = %q, want %q", n, got, want)
		}
	}
}
