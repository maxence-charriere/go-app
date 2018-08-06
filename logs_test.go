package app_test

import (
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/test"
	"github.com/murlokswarm/app/internal/tests"
)

func TestLogs(t *testing.T) {
	setup := func(onRun func()) app.Driver {
		d := &test.Driver{
			OnRun: onRun,
		}

		withLogs := app.Logs()
		return withLogs(d)
	}

	tests.TestDriver(t, setup)
}

func TestLogsErrors(t *testing.T) {
	setup := func(onRun func()) app.Driver {
		d := &test.Driver{
			Err:   true,
			OnRun: onRun,
		}

		withLogs := app.Logs()
		return withLogs(d)
	}

	tests.TestDriver(t, setup)
}
