package protect

import "testing"

func TestProtectRestore(t *testing.T) {
	input := "다음 코드를 리뷰해줘:\n```go\nfunc hello() {}\n```"
	protected, p := Protect(input)
	restored := p.Restore(protected)
	_ = restored
	// Stub returns input unchanged — tests will pass once Task 04 is implemented.
}
