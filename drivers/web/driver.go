// Package web is the driver to be used for web applications.
// It is build on the top of GopherJS.
package web

import (
	"net/http"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/core"
)

// Driver is an app.Driver implementation for web.
type Driver struct {
	core.Driver

	// The URL of the component to load when a navigating on the website root.
	URL string

	// The URL of the component to load when a 404 errors occurs.
	// Default is /web.NotFound
	NotFoundURL string

	// The server used to save request.
	// Default is a server that listens on port 7042.
	Server *http.Server

	// OnServerRun is called when the web server is running.
	// http.Handler overrides should be performed here.
	OnServerRun func()

	factory     *app.Factory
	elems       *core.ElemDB
	page        app.Page
	uichan      chan func()
	stop        func()
	fileHandler http.Handler
}

// Name satisfies the app.Driver interface.
func (d *Driver) Name() string {
	return "Web"
}
