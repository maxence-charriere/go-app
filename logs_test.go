package app_test

import (
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/tests"
)

func TestDriverWithLogsMinimal(t *testing.T) {
	app.EnableDebug(true)
	ui := make(chan func(), 64)

	c := app.DriverConfig{
		Events:  app.NewEventRegistry(ui),
		Factory: app.NewFactory(),
		UI:      ui,
	}

	d := tests.NewMinimalDriver(c)
	tests.TestDriver(t, app.Logs()(d), c)
}

func TestDriverWithLogs(t *testing.T) {
	app.EnableDebug(true)
	ui := make(chan func(), 64)

	c := app.DriverConfig{
		Events:  app.NewEventRegistry(ui),
		Factory: app.NewFactory(),
		UI:      ui,
	}

	d := tests.NewDriver(c)
	tests.TestDriver(t, app.Logs()(d), c)
}
