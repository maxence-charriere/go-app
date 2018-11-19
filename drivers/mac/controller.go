// +build darwin,amd64

package mac

import (
	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/bridge"
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

	if err := driver.macRPC.Call("controller.New", nil, struct {
		ID string
	}{
		ID: controller.id,
	}); err != nil {
		controller.SetErr(err)
		return controller
	}

	driver.elems.Put(controller)
	return controller
}

// ID satistfies the app.Controller interface.
func (c *Controller) ID() string {
	return c.id
}

// Close satisfies the app.Controller interface.
func (c *Controller) Close() {
	err := driver.macRPC.Call("controller.Close", nil, struct {
		ID string
	}{
		ID: c.id,
	})

	c.SetErr(err)
	driver.elems.Delete(c)
}

func onControllerDirectionChange(c *Controller, in map[string]interface{}) interface{} {
	if c.onDirectionChange != nil {
		c.onDirectionChange(
			app.ControllerInput(in["Input"].(float64)),
			in["X"].(float64),
			in["Y"].(float64),
		)
	}

	return nil
}

func onControllerButtonPressed(c *Controller, in map[string]interface{}) interface{} {
	if c.onButtonPressed != nil {
		c.onButtonPressed(
			app.ControllerInput(in["Input"].(float64)),
			in["Value"].(float64),
			in["Pressed"].(bool),
		)
	}

	return nil
}

func onControllerPause(c *Controller, in map[string]interface{}) interface{} {
	if c.onPause != nil {
		c.onPause()
	}

	return nil
}

func onControllerClose(c *Controller, in map[string]interface{}) interface{} {
	if c.onClose != nil {
		c.onClose()
	}

	return nil
}

func onControllerConnected(c *Controller, in map[string]interface{}) interface{} {
	if err := driver.macRPC.Call("controller.Listen", nil, struct {
		ID string
	}{
		ID: c.id,
	}); err != nil {
		c.SetErr(err)
		return c
	}

	if c.onConnected != nil {
		c.onConnected()
	}

	return nil
}

func onControllerDisconnected(c *Controller, in map[string]interface{}) interface{} {
	if c.onDisconnected != nil {
		c.onDisconnected()
	}

	return nil
}

func handleController(h func(c *Controller, in map[string]interface{}) interface{}) bridge.GoRPCHandler {
	return func(in map[string]interface{}) interface{} {
		id, _ := in["ID"].(string)

		e := driver.elems.GetByID(id)
		if e.Err() == app.ErrElemNotSet {
			return nil
		}

		return h(e.(*Controller), in)
	}
}
