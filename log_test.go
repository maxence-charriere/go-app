package app

import "testing"

func TestConsole(t *testing.T) {
	cons := Console{}

	cons.Log("hello", "world")
	cons.Logf("%s %s", "hello", "world")

	cons.Debug = true

	cons.Log("hello", "world")
	cons.Logf("%s %s", "hello", "world")

	cons.Error("hello", "world")
	cons.Errorf("%s %s", "hello", "world")
}

func TestMultiLogger(t *testing.T) {
	l := NewMultiLogger(&Console{Debug: true})

	l.Log("multi", "hello", "world")
	l.Logf("%s %s %s", "multi", "hello", "world")

	l.Error("multi", "hello", "world")
	l.Errorf("%s %s %s", "multi", "hello", "world")
}
