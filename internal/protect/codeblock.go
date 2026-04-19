package protect

// Protector extracts protected regions (code blocks, URLs, paths) before translation
// and restores them afterward.
// Implementation: Task 04
type Protector struct {
	placeholders map[string]string
}

func Protect(text string) (string, *Protector) {
	return text, &Protector{placeholders: make(map[string]string)}
}

func (p *Protector) Restore(text string) string {
	return text
}
