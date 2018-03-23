package app_test

import (
	"context"
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/test"
	"github.com/murlokswarm/app/tests"
)

func TestDriverWithLogs(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	tests.TestDriver(t, func(onRun func()) app.Driver {
		d := &test.Driver{
			OnRun: onRun,
			Ctx:   ctx,
		}
		return app.DriverWithLogs(d)
	}, cancel)
}

func TestDriverWithLogsBaseDriver(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	tests.TestDriver(t, func(onRun func()) app.Driver {
		d := &test.Driver{
			Ctx:           ctx,
			UseBaseDriver: true,

			OnRun: onRun,
		}
		return app.DriverWithLogs(d)
	}, cancel)
}

func TestDriverWithLogsError(t *testing.T) {
	var d app.Driver = &test.Driver{
		SimulateErr: true,
	}
	d = app.DriverWithLogs(d)

	err := d.Run(app.NewFactory())
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)

	if _, err = d.NewWindow(app.WindowConfig{}); err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)

	if _, err = d.NewContextMenu(app.MenuConfig{}); err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)

	if err = d.Render(&tests.Hello{}); err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}
