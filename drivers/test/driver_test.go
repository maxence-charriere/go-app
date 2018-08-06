package test

import (
	"testing"

	"github.com/murlokswarm/app"

	"github.com/murlokswarm/app/internal/tests"
)

func TestDriver(t *testing.T) {
	setup := func(onRun func()) app.Driver {
		return &Driver{
			OnRun: onRun,
		}
	}

	tests.TestDriver(t, setup)
}
