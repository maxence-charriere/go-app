package app_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/test"
	"github.com/murlokswarm/app/tests"
)

func TestDriverWithLogs(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	buff := &bytes.Buffer{}
	app.Loggers = []app.Logger{
		app.NewLogger(buff, buff, true),
	}

	tests.TestDriver(t, func(onRun func()) app.Driver {
		d := &test.Driver{
			OnRun: onRun,
			Ctx:   ctx,
		}
		return app.Logs()(d)
	}, cancel)

	t.Log(buff.String())
}
