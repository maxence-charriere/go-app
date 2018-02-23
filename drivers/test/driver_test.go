package test

import (
	"context"
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/tests"
)

func TestDriver(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	tests.TestDriver(t, func(onRun func()) app.Driver {
		return &Driver{
			OnRun: onRun,
			Ctx:   ctx,
		}
	}, cancel)
}
