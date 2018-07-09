package app

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImport(t *testing.T) {
	Import(&Bar{})

	defer func() { recover() }()
	Import(&EmptyCompo{})
	t.Error("no panic")
}

func TestRun(t *testing.T) {
	buff := &bytes.Buffer{}
	Loggers = []Logger{
		NewLogger(buff, buff, true, true),
	}
	defer t.Log(buff.String())

	AddBackend(&backend{})
	err := Run()
	assert.NoError(t, err)

	AddBackend(&backend{
		simulateError: true,
	})
	err = Run()
	assert.Error(t, err)
}

func TestRender(t *testing.T) {
	buff := &bytes.Buffer{}
	Loggers = []Logger{
		NewLogger(buff, buff, true, true),
	}
	defer t.Log(buff.String())

	AddBackend(&backend{})
	Render(nil)

	AddBackend(&backend{
		simulateError: true,
	})
	Render(nil)
}

func TestCallOnUIGoroutine(t *testing.T) {
	CallOnUIGoroutine(func() {
		t.Log("hello from ui goroutine")
	})

	f := <-uiChan
	f()
}

func TestActions(t *testing.T) {
	HandleAction("test", func(e EventDispatcher, a Action) {})

	PostAction("test", 42)
	PostActions(
		Action{Name: "test", Arg: 21},
		Action{Name: "test", Arg: 84},
	)
}

func TestLogs(t *testing.T) {
	buff := &bytes.Buffer{}
	Loggers = []Logger{
		NewLogger(buff, buff, true, true),
	}
	defer t.Log(buff.String())

	Log("hello world")
	WhenDebug(func() {
		Debug("goodbye world")
	})

	t.Log(buff.String())
}
