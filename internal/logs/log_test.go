package logs

import (
	"bytes"
	"errors"
	"testing"
)

func TestFromWriter(t *testing.T) {
	b := &bytes.Buffer{}
	l := ToWriter(b)

	testLogger(t, l)
	t.Log(b.String())
}

func TestWithPromt(t *testing.T) {
	b := &bytes.Buffer{}
	l := ToWriter(b)
	l = WithPrompt(l)

	testLogger(t, l)
	t.Log(b.String())
}

func TestWithColoredPromt(t *testing.T) {
	b := &bytes.Buffer{}
	l := ToWriter(b)
	l = WithColoredPrompt(l)

	testLogger(t, l)
	t.Log(b.String())
}

func testLogger(t *testing.T, l Logger) {
	l("logging a message")
	l("logging a message with arg: %s", "hello")
	l("logging an error: %s", errors.New("dumb error"))
}
