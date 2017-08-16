package app

import (
	"github.com/murlokswarm/app/markup"
	"github.com/pkg/errors"
)

var (
	driver       Driver
	compoBuilder = markup.NewCompoBuilder()
)

// Import imports component c into the app.
// Components must be imported in order the be used by the app package.
// This mechanism allows components to be created dynamically when they are
// found into HTML code.
// Import should be called during app initialization.
func Import(c markup.Component) {
	if err := compoBuilder.Register(c); err != nil {
		err = errors.Wrap(err, "invalid component import")
		panic(err)
	}
}

// Run runs the app with driver d as backend.
func Run(d Driver) error {
	driver = d
	return d.Run(compoBuilder)
}

// CurrentDriver returns the used driver.
// It panics if called before Run.
func CurrentDriver() Driver {
	if driver == nil {
		panic("no current driver")
	}
	return driver
}

// Resources returns the location of the resources directory.
// Resources should be used only for read only operations.
func Resources() string {
	return driver.Resources()
}

// Storage returns the location of the storage directory.
// It panics if the running driver is not a DriverWithStorage.
func Storage() string {
	d := driver.(DriverWithStorage)
	return d.Storage()
}

// NewWindow creates and displays a window described by configuration c.
// It panics if the running driver is not a DriverWithWindows.
func NewWindow(c WindowConfig) Window {
	d := driver.(DriverWithWindows)
	return d.NewWindow(c)
}
