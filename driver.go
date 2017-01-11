package app

import "github.com/murlokswarm/log"

var (
	driver Driver
)

// Driver is the interface that describes the implementation to handle platform
// specific rendering.
type Driver interface {
	Run()

	NewContext(ctx interface{}) Contexter

	MenuBar() Contexter

	Dock() Docker

	Storage() Storer

	JavascriptBridge() string

	Share() Sharer
}

// RegisterDriver registers the driver to be used when using the app package.
// Should be used only into a driver implementation, in an init function.
func RegisterDriver(d Driver) {
	driver = d
	log.Infof("driver %T is loaded", d)
}
