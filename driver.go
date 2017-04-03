package app

import "github.com/murlokswarm/log"

var (
	driver Driver
)

// Driver is the interface that describes the implementation to handle platform
// specific rendering.
type Driver interface {
	// Run runs the application. Should start the app loop.
	Run()

	// NewElement create an app element. e should be a struct describing the
	// element (e.g. Window, ContextMenu).
	// Implementation:
	// - Should check the type of e and then create the native element
	//   described.
	NewElement(e interface{}) Elementer

	// MenuBar returns the element that represents the menu bar.
	// Implementation:
	// - Should be created in a driver.
	// - Calling Close should make the program panic.
	// - If there is no menu bar in the native platform, methods should do
	//   nothing (except Close).
	MenuBar() Contexter

	// Dock returns the element that represents the dock.
	// Implementation:
	// - Should be created in the driver implementation.
	// - Calling Close should make the program panic.
	// - If there is no dock in the native platform, methods should do nothing
	//   (except Close).
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
