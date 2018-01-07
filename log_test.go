package app_test

import (
	"bytes"
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/tests"
)

func TestLogger(t *testing.T) {
	buffer := &bytes.Buffer{}
	tests.TestLogger(t, app.NewLogger(buffer, true))
	tests.TestLogger(t, app.NewLogger(buffer, false))
	t.Log(buffer.String())
}

func TestConsole(t *testing.T) {
	tests.TestLogger(t, app.NewConsole(false))
	tests.TestLogger(t, app.NewConsole(true))
}

func TestConcurrentLogger(t *testing.T) {
	buffer := &bytes.Buffer{}
	tests.TestLogger(t, app.NewConcurrentLogger(app.NewLogger(buffer, true)))
}

func TestMultiLogger(t *testing.T) {
	buffer := &bytes.Buffer{}
	tests.TestLogger(t, app.NewMultiLogger(
		app.NewConsole(false),
		app.NewLogger(buffer, true),
	))
}
