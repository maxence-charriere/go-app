package app

import "github.com/murlokswarm/log"

var (
	driver Driver
)

// Driver is the interface that describes the implementation to handle platform
// specific rendering.
type Driver interface {
	// Run runs the application.
	//
	// Driver implementation:
	// - Should start the app loop.
	Run()

	// NewElement create an app element. e should be a struct describing the
	// element (e.g. Window, ContextMenu).
	//
	// Driver implementation:
	// - Should check the type of e and then create the native element
	//   described.
	NewElement(e interface{}) Elementer

	// MenuBar returns the element that represents the menu bar.
	// ok will be false if there is no menubar available.
	//
	// Driver implementation:
	// - Should be created in a driver.
	MenuBar() (menu Contexter, ok bool)

	// Dock returns the element that represents the dock.
	// ok will be false if there is no dock available.
	//
	// Driver implementation:
	// - Should be created in the driver implementation.
	Dock() (d Docker, ok bool)

	Storage() Storer

	JavascriptBridge() string
}

// RegisterDriver registers the driver to be used when using the app package.
//
// Driver implementation:
// - Should be called once in an init() func.
func RegisterDriver(d Driver) {
	driver = d
	log.Infof("driver %T is loaded", d)
}
