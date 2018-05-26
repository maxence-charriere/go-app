package app_test

import (
	"context"
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/test"
	"github.com/murlokswarm/app/tests"
)

func TestBaseDriver(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	tests.TestDriver(t, func(onRun func()) app.Driver {
		return &test.Driver{
			OnRun:         onRun,
			Ctx:           ctx,
			UseBaseDriver: true,
		}
	}, cancel)
}
