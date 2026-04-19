package stats

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// Stats tracks per-session translation token estimates.
type Stats struct {
	mu           sync.Mutex
	SessionStart time.Time
	InputKO      int // estimated tokens before translation (Korean)
	InputEN      int // estimated tokens after translation (English)
	OutputEN     int // estimated tokens of Claude's English response
	OutputKO     int // estimated tokens of translated Korean response
	Requests     int // number of translated requests
}

func New() *Stats {
	return &Stats{SessionStart: time.Now()}
}

// RecordInput records one translated user message (ko→en).
func (s *Stats) RecordInput(koText, enText string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.InputKO += estimateKorean(koText)
	s.InputEN += estimateEnglish(enText)
	s.Requests++
}

// RecordOutput records one translated Claude response (en→ko).
func (s *Stats) RecordOutput(enText, koText string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.OutputEN += estimateEnglish(enText)
	s.OutputKO += estimateKorean(koText)
}

// Summary returns a formatted multi-line session report.
func (s *Stats) Summary() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Requests == 0 {
		return "No translations recorded yet.\nStart a session with: autolang start"
	}

	dur := time.Since(s.SessionStart).Round(time.Second)
	saved := s.InputKO - s.InputEN
	var pct float64
	if s.InputKO > 0 {
		pct = float64(saved) / float64(s.InputKO) * 100
	}

	// Rough cost estimate: ~$0.000003 per token saved (gpt-4o-mini input pricing)
	costSaved := float64(saved) * 0.000003

	sep := strings.Repeat("─", 38)
	return fmt.Sprintf(`
  AutoLang Session Stats
  %s
  Session duration    : %s
  Translated requests : %d

  Input  (your messages)
    Korean  : ~%s tokens
    English : ~%s tokens
    Saved   : %s tokens (%.0f%%)

  Output  (Claude responses)
    English : ~%s tokens
    Korean  : ~%s tokens  (읽기 전용, API 비용 없음)

  Estimated cost saved : ~$%.4f
  %s`,
		sep,
		dur,
		s.Requests,
		formatInt(s.InputKO),
		formatInt(s.InputEN),
		formatInt(saved),
		pct,
		formatInt(s.OutputEN),
		formatInt(s.OutputKO),
		costSaved,
		sep,
	)
}

// Indicator returns a one-line translation indicator for stderr.
func (s *Stats) Indicator(koTokens, enTokens int) string {
	saved := koTokens - enTokens
	var pct float64
	if koTokens > 0 {
		pct = float64(saved) / float64(koTokens) * 100
	}
	return fmt.Sprintf("✦ [AutoLang] KO → EN  (tokens: %d → %d, saved %.0f%%)", koTokens, enTokens, pct)
}

// EstimateKorean returns the estimated token count for Korean text.
func EstimateKorean(text string) int { return estimateKorean(text) }

// EstimateEnglish returns the estimated token count for English text.
func EstimateEnglish(text string) int { return estimateEnglish(text) }

func estimateKorean(text string) int {
	chars := len([]rune(text))
	if chars == 0 {
		return 0
	}
	return chars * 25 / 10 // × 2.5
}

func estimateEnglish(text string) int {
	words := 0
	inWord := false
	for _, r := range text {
		if r == ' ' || r == '\n' || r == '\t' {
			inWord = false
		} else if !inWord {
			words++
			inWord = true
		}
	}
	if words == 0 {
		return 0
	}
	return words * 13 / 10 // × 1.3
}

func formatInt(n int) string {
	if n < 0 {
		return "0"
	}
	s := fmt.Sprintf("%d", n)
	// Insert thousands separators
	out := []byte{}
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			out = append(out, ',')
		}
		out = append(out, byte(c))
	}
	return string(out)
}
