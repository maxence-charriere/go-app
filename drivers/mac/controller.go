// +build darwin,amd64

package mac

import (
	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/bridge"
	"github.com/murlokswarm/app/internal/core"
)

// Controller implements the core.Elem struct
type Controller struct {
	core.Controller

	id string

	// Event handlers
	onDpadChange   func(app.ControllerInput, float64, float64)
	onButtonChange func(app.ControllerInput, float64, bool)
	onConnected    func()
	onDisconnected func()
	onPause        func()
	onClose        func()
}

func newController(c app.ControllerConfig) *Controller {
	controller := &Controller{
		id:             uuid.New().String(),
		onDpadChange:   c.OnDpadChange,
		onButtonChange: c.OnButtonChange,
		onConnected:    c.OnConnected,
		onDisconnected: c.OnDisconnected,
		onPause:        c.OnPause,
		onClose:        c.OnClose,
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

// ID satistfies the app.Elem interface.
func (c *Controller) ID() string {
	return c.id
}

// Close satisfies the app.Controller interface
func (c *Controller) Close() {
	err := driver.macRPC.Call("controller.Close", nil, struct {
		ID string
	}{
		ID: c.id,
	})

	c.SetErr(err)
	driver.elems.Delete(c)
}

// Called when either the directional pad or thumbsticks' values are changed (controller.onDpadChange)
// Available values for `button`:
//		dpad, leftThumbstick, rightThumbstick
func onControllerDpadChange(c *Controller, in map[string]interface{}) interface{} {
	if c.onDpadChange != nil {
		c.onDpadChange(
			app.ControllerInput(in["Input"].(float64)),
			in["X"].(float64),
			in["Y"].(float64),
		)
	}

	return nil
}

// Called when one of the following buttons' values are changed (controller.onButtonChange)
// Available values for `button`:
//		buttonA, buttonB, buttonX, buttonY,
//		leftShoulder, rightShoulder
//		leftTrigger, rightTrigger
func onControllerButtonChange(c *Controller, in map[string]interface{}) interface{} {
	if c.onButtonChange != nil {
		c.onButtonChange(
			app.ControllerInput(in["Input"].(float64)),
			in["Value"].(float64),
			in["Pressed"].(bool),
		)
	}

	return nil
}

// Called when the pause button is called.
func onControllerPause(c *Controller, in map[string]interface{}) interface{} {
	if c.onPause != nil {
		c.onPause()
	}

	return nil
}

// Called when the controller is intentionally removed from the application
func onControllerClose(c *Controller, in map[string]interface{}) interface{} {
	if c.onClose != nil {
		c.onClose()
	}

	return nil
}

// Called when a controller becomes connected.
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

// Called when a controller becomes disconnected.
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

		c := e.(*Controller)
		return h(c, in)
	}
}
