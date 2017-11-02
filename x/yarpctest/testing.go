package yarpctest

import "testing"

// TestingT is an interface wrapper around *testing.T and *testing.B
type TestingT interface {
	testing.TB
}

// Run will cast the TestingT to it's sub and call the appropriate Run func.
func Run(name string, t TestingT, f func(TestingT)) {
	if tt, ok := t.(*testing.T); ok {
		tt.Run(name, func(ttt *testing.T) { f(ttt) })
		return
	}
	if tb, ok := t.(*testing.B); ok {
		tb.Run(name, func(ttb *testing.B) { f(ttb) })
		return
	}
	t.Error("invalid test harness")
	t.FailNow()
}
