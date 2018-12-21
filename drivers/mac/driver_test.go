// +build darwin,amd64

package mac

import (
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/tests"
)

func TestDriver(t *testing.T) {
	setup := func() app.Driver {
		return &Driver{}
	}

	tests.TestDriver(t, setup)
}
