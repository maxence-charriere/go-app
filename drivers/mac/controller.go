// +build darwin,amd64

package mac

import (
	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/core"
)

// Controller implements the app.Controller interface.
//
// OnButtonPressed for LeftThumbstick and RightThumbstick is not available
// before MacOS 10.14.1.
type Controller struct {
	core.Controller

	id string

	onDirectionChange func(app.ControllerInput, float64, float64)
	onButtonPressed   func(app.ControllerInput, float64, bool)
	onConnected       func()
	onDisconnected    func()
	onPause           func()
	onClose           func()
}

func newController(c app.ControllerConfig) *Controller {
	controller := &Controller{
		id: uuid.New().String(),

		onDirectionChange: c.OnDirectionChange,
		onButtonPressed:   c.OnButtonPressed,
		onConnected:       c.OnConnected,
		onDisconnected:    c.OnDisconnected,
		onPause:           c.OnPause,
		onClose:           c.OnClose,
	}

	if err := driver.Platform.Call("controller.New", nil, struct {
		ID string
	}{
		ID: controller.id,
	}); err != nil {
		controller.SetErr(err)
		return controller
	}

	driver.Elems.Put(controller)
	return controller
}

// ID satistfies the app.Controller interface.
func (c *Controller) ID() string {
	return c.id
}

// Close satisfies the app.Controller interface.
func (c *Controller) Close() {
	err := driver.Platform.Call("controller.Close", nil, struct {
		ID string
	}{
		ID: c.id,
	})

	c.SetErr(err)
	driver.Elems.Delete(c)
}

func onControllerDirectionChange(c *Controller, in map[string]interface{}) {
	if c.onDirectionChange != nil {
		c.onDirectionChange(
			app.ControllerInput(in["Input"].(float64)),
			in["X"].(float64),
			in["Y"].(float64),
		)
	}
}

func onControllerButtonPressed(c *Controller, in map[string]interface{}) {
	if c.onButtonPressed != nil {
		c.onButtonPressed(
			app.ControllerInput(in["Input"].(float64)),
			in["Value"].(float64),
			in["Pressed"].(bool),
		)
	}
}

func onControllerPause(c *Controller, in map[string]interface{}) {
	if c.onPause != nil {
		c.onPause()
	}
}

func onControllerClose(c *Controller, in map[string]interface{}) {
	if c.onClose != nil {
		c.onClose()
	}
}

func onControllerConnected(c *Controller, in map[string]interface{}) {
	if err := driver.Platform.Call("controller.Listen", nil, struct {
		ID string
	}{
		ID: c.id,
	}); err != nil {
		c.SetErr(err)
		return
	}

	if c.onConnected != nil {
		c.onConnected()
	}
}

func onControllerDisconnected(c *Controller, in map[string]interface{}) {
	if c.onDisconnected != nil {
		c.onDisconnected()
	}
}

func handleController(h func(c *Controller, in map[string]interface{})) core.GoHandler {
	return func(in map[string]interface{}) {
		id, _ := in["ID"].(string)

		e := driver.Elems.GetByID(id)
		if e.Err() == app.ErrElemNotSet {
			return
		}

		h(e.(*Controller), in)
	}
}
