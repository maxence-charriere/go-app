package test

import (
	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/core"
)

// Controller implements the core.Elem struct
type Controller struct {
	core.Controller

	id     string
	driver *Driver

	// event handlers
	onDpadChange   func(app.ControllerInput, float64, float64)
	onButtonChange func(app.ControllerInput, float64, bool)
	onConnected    func()
	onDisconnected func()
	onPause        func()
	onClose        func()
}

func newController(d *Driver, c app.ControllerConfig) *Controller {
	controller := &Controller{
		id:             uuid.New().String(),
		driver:         d,
		onDpadChange:   c.OnDpadChange,
		onButtonChange: c.OnButtonChange,
		onConnected:    c.OnConnected,
		onDisconnected: c.OnDisconnected,
		onPause:        c.OnPause,
		onClose:        c.OnClose,
	}

	d.elems.Put(controller)
	return controller
}

// ID satistfies the app.Elem interface.
func (c *Controller) ID() string {
	return c.id
}

func onControllerDpadChange(c *Controller, in map[string]interface{}) interface{} {
	if c.onDpadChange != nil {
		input := app.ControllerInput(in["Input"].(float64))
		x := in["X"].(float64)
		y := in["Y"].(float64)
		c.onDpadChange(input, x, y)
	}

	return nil
}

func onControllerButtonChange(c *Controller, in map[string]interface{}) interface{} {
	if c.onButtonChange != nil {
		input := app.ControllerInput(in["Input"].(float64))
		value := in["Value"].(float64)
		pressed := in["Pressed"].(bool)
		c.onButtonChange(input, value, pressed)
	}

	return nil
}

func onControllerConnected(c *Controller, in map[string]interface{}) interface{} {
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
