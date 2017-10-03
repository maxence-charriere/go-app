package log

import "testing"

func TestMultiLogger(t *testing.T) {
	l := NewMultiLogger(&Logger{Debug: true})

	l.Log("multi", "hello", "world")
	l.Logf("%s %s %s", "multi", "hello", "world")

	l.Error("multi", "hello", "world")
	l.Errorf("%s %s %s", "multi", "hello", "world")
}
