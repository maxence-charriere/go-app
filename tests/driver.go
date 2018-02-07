package tests

import (
	"testing"

	"github.com/murlokswarm/app"
)

func TestDriver(t *testing.T, setup func(onRun func()) app.Driver) {
	var driver app.Driver

	factory := app.NewFactory()
	factory = app.NewConcurrentFactory(factory)

	onRun := func() {
		t.Log("test driver", driver.Name())
	}

	driver = setup(onRun)
	if err := driver.Run(factory); err != nil {
		t.Error(err)
	}
}
