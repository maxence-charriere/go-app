package log

import "testing"

func TestLogger(t *testing.T) {
	l := Logger{}

	l.Log("hello", "world")
	l.Logf("%s %s", "hello", "world")

	l.Debug = true

	l.Log("hello", "world")
	l.Logf("%s %s", "hello", "world")

	l.Error("hello", "world")
	l.Errorf("%s %s", "hello", "world")
}
