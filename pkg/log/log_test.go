package log

import (
	"bytes"
	"testing"
)

func TestLogs(t *testing.T) {
	b := bytes.Buffer{}
	outpout = &b

	Info("hello", "world").T("foo", 42)
	Infof("hello %s", "world").T("foo", 42)

	Error("hello", "world").T("foo", 42)
	Errorf("hello %s", "world").T("foo", 42)

	Warn("hello", "world").T("foo", 42)
	Warnf("hello %s", "world").T("foo", 42)

	Debug("hello", "world").T("foo", 42)

	CurrentLevel = DebugLevel

	Debug("hello", "world").T("foo", 42)
	Debugf("hello %s", "world").T("foo", 42)

	Log(Entry{
		Level:   InfoLevel,
		Message: "hello",
		Tags:    map[string]string{"foo": "bar"},
	})

	t.Log(b.String())
}

func TestEntryPanic(t *testing.T) {
	defer func() {
		recover()
	}()

	Error("error!").Panic()
	t.Error("did not panic")
}
