// +build darwin,amd64

package mac

import (
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/tests"
)

func TestDriver(t *testing.T) {
	ui := make(chan func(), 64)

	c := app.DriverConfig{
		Events:  app.NewEventRegistry(ui),
		Factory: app.NewFactory(),
		UI:      ui,
	}

	tests.TestDriver(t, &Driver{}, c)
}
