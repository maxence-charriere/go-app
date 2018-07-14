package app_test

import (
	"context"
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/test"
	"github.com/murlokswarm/app/internal/tests"
)

func TestBaseDriver(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	shutdown := func() error {
		cancel()
		return nil
	}

	tests.TestDriver(t, func(onRun func()) app.Driver {
		return &test.Driver{
			OnRun:         onRun,
			Ctx:           ctx,
			UseBaseDriver: true,
		}
	}, shutdown)
}
