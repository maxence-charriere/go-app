package core

import "github.com/murlokswarm/app"

// Controller is a base struct to embed in app.Controller implementations.
type Controller struct {
	Elem
}

// Close satisfies the app.Controller interface.
func (c *Controller) Close() {
	c.SetErr(app.ErrNotSupported)
}
