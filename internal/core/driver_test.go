package core

import (
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/tests"
)

func TestDriver(t *testing.T) {
	setup := func() app.Driver {
		return &Driver{
			Elems:  NewElemDB(),
			UIChan: make(chan func(), 32),
		}
	}

	tests.TestDriver(t, setup)
}
