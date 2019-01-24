package core_test

import (
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/tests"
)

func TestDriverMinimal(t *testing.T) {
	ui := make(chan func(), 64)

	c := app.DriverConfig{
		Events:  app.NewEventRegistry(ui),
		Factory: app.NewFactory(),
		UI:      ui,
	}

	d := tests.NewMinimalDriver(c)
	tests.TestDriver(t, d, c)
}

func TestDriver(t *testing.T) {
	ui := make(chan func(), 64)

	c := app.DriverConfig{
		Events:  app.NewEventRegistry(ui),
		Factory: app.NewFactory(),
		UI:      ui,
	}

	d := tests.NewDriver(c)
	tests.TestDriver(t, d, c)
}
