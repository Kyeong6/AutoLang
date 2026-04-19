package stats

import (
	"fmt"
	"sync"
	"time"
)

// Stats tracks token savings per session.
// Implementation: Task 09
type Stats struct {
	mu           sync.Mutex
	SessionStart time.Time
	InputKO      int
	InputEN      int
	OutputEN     int
	Requests     int
}

func New() *Stats {
	return &Stats{SessionStart: time.Now()}
}

func (s *Stats) RecordInput(koText, enText string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.InputKO += estimateTokens(koText, true)
	s.InputEN += estimateTokens(enText, false)
	s.Requests++
}

func (s *Stats) Summary() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	saved := s.InputKO - s.InputEN
	if s.InputKO == 0 {
		return "No translations recorded yet."
	}
	pct := float64(saved) / float64(s.InputKO) * 100
	return fmt.Sprintf(
		"Requests: %d | Input: %d KO → %d EN tokens (saved %.0f%%)",
		s.Requests, s.InputKO, s.InputEN, pct,
	)
}

func estimateTokens(text string, isKorean bool) int {
	if isKorean {
		return len([]rune(text)) * 25 / 10
	}
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
		return 1
	}
	return words * 13 / 10
}
