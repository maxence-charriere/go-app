package app_test

import (
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/test"
	"github.com/murlokswarm/app/tests"
)

func TestDriverWithLogs(t *testing.T) {
	// Specific tests.
	var d app.Driver = &test.Driver{
		SimulateErr: true,
	}
	d = app.NewDriverWithLogs(d)

	factory := app.NewFactory()
	factory = app.NewConcurrentFactory(factory)

	d.Run(factory)
	d.Render(&tests.Foo{})
	d.NewWindow(app.WindowConfig{})
	d.NewContextMenu(app.MenuConfig{})

	// Test suite.
	setup := func(onRun func()) app.Driver {
		var d app.Driver = &test.Driver{
			OnRun: onRun,
		}
		d = app.NewDriverWithLogs(d)
		return d
	}
	tests.TestDriver(t, setup)
}
