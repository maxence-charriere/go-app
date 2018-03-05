// Package web is the driver to be used for web applications.
package web

import (
	"net/http"

	"github.com/murlokswarm/app"
)

// Driver is an app.Driver implementation for web.
type Driver struct {
	app.BaseDriver

	// The URL of the component to load when a navigating on the website root.
	DefaultURL string

	// The URL of the component to load when a 404 errors occurs.
	// Default is /web.NotFound
	NotFoundURL string

	// The server used to save request.
	// Default is a server that listens on port 7042.
	Server *http.Server

	factory     app.Factory
	elements    app.ElemDB
	page        app.Page
	uichan      chan func()
	cancel      func()
	fileHandler http.Handler
}

// Name satisfies the app.Driver interface.
func (d *Driver) Name() string {
	return "Web"
}

// Base satisfies the app.Driver interface.
func (d *Driver) Base() app.Driver {
	return d
}
