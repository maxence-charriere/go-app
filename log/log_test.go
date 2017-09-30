package log

import "testing"

func TestLogger(t *testing.T) {
	l := Logger{}

	l.Log("hello", "world")
	l.Logf("%s %s\n", "hello", "world")

	l.Debug = true

	l.Log("hello", "world")
	l.Logf("%s %s\n", "hello", "world")

	l.Error("hello", "world")
	l.Errorf("%s %s\n", "hello", "world")
}
