package app_test

import (
	"fmt"
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/test"
	"github.com/murlokswarm/app/internal/tests"
)

func TestLogs(t *testing.T) {
	app.Logger = func(format string, a ...interface{}) {
		log := fmt.Sprintf(format, a...)
		t.Log(log)
	}

	app.EnableDebug(true)

	setup := func() app.Driver {
		d := &test.Driver{}
		withLogs := app.Logs()
		return withLogs(d)
	}

	tests.TestDriver(t, setup)
}

func TestLogsErrors(t *testing.T) {
	app.Logger = func(format string, a ...interface{}) {
		log := fmt.Sprintf(format, a...)
		t.Log(log)
	}

	setup := func() app.Driver {
		d := &test.Driver{Err: true}
		withLogs := app.Logs()
		return withLogs(d)
	}

	tests.TestDriver(t, setup)
}
