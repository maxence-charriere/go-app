package core

import (
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/tests"
	"github.com/tidwall/gjson"
)

func TestDriverMinimal(t *testing.T) {
	ui := make(chan func(), 64)

	c := app.DriverConfig{
		Events:  app.NewEventRegistry(ui),
		Factory: app.NewFactory(),
		UI:      ui,
	}

	d := &Driver{
		Elems:    NewElemDB(),
		Events:   c.Events,
		Factory:  c.Factory,
		Platform: &Platform{},
		UIChan:   ui,
	}

	handler := func(call string) error {
		returnID := gjson.Get(call, "ReturnID").Str
		d.Platform.Return(returnID, "", "not implemented")
		return nil
	}

	d.Platform.Handler = handler

	tests.TestDriver(t, d, c)
}
