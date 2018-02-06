package app

import (
	"github.com/pkg/errors"
)

var (
	driver     Driver
	components Factory = make(factory)
)

// Import imports the component into the app.
// Components must be imported in order the be used by the app package.
// This allows components to be created dynamically when they are found into
// markup.
func Import(c Component) {
	if _, err := components.RegisterComponent(c); err != nil {
		err = errors.Wrap(err, "import component failed")
		panic(err)
	}
}

// Run runs the app with the driver as backend.
func Run(d Driver) error {
	if driver != nil {
		return errors.Errorf("driver %T is already running", driver)
	}

	driver = NewDriverWithLogs(d)
	return driver.Run(components)

}

// RunningDriver returns the running driver.
//
// It panics if called before Run.
func RunningDriver() Driver {
	if driver == nil {
		panic("no current driver")
	}
	return driver
}

// Name returns the application name.
//
// It panics if called before Run.
func Name() string {
	return driver.AppName()
}

// Resources returns the given path prefixed by the resources directory
// location.
// Resources should be used only for read only operations.
//
// It panics if called before Run.
func Resources(path ...string) string {
	return driver.Resources(path...)
}

// Storage returns the given path prefixed by the storage directory
// location.
//
// It panics if called before Run.
func Storage(path ...string) string {
	return driver.Storage(path...)
}

// NewWindow creates and displays the window described by the given
// configuration.
//
// It panics if called before Run.
func NewWindow(c WindowConfig) (Window, error) {
	return driver.NewWindow(c)
}

// NewContextMenu creates and displays the context menu described by the
// given configuration.
//
// It panics if called before Run.
func NewContextMenu(c MenuConfig) (Menu, error) {
	return driver.NewContextMenu(c)
}

// Render renders the given component.
// It should be called when the display of component c have to be updated.
//
// It panics if called before Run.
func Render(c Component) {
	driver.Render(c)
}

// ElementByComponent returns the element where the given component is mounted.
//
// It panics if called before Run.
func ElementByComponent(c Component) (ElementWithComponent, error) {
	return driver.ElementByComponent(c)
}

// WindowByComponent returns the window where the given component is mounted.
//
// It panics if called before Run.
func WindowByComponent(c Component) (Window, error) {
	elem, err := driver.ElementByComponent(c)
	if err != nil {
		return nil, err
	}

	win, ok := elem.(Window)
	if !ok {
		return nil, errors.New("component is not mounted in a window")
	}
	return win, nil
}

// NewFilePanel creates and displays the file panel described by the given
// configuration.
//
// It panics if called before Run.
func NewFilePanel(c FilePanelConfig) error {
	return driver.NewFilePanel(c)
}

// NewShare creates and display the share pannel to share the given value.
//
// It panics if called before Run.
func NewShare(v interface{}) error {
	return driver.NewShare(v)
}

// NewNotification creates and displays the notification described in the
// given configuration.
//
// It panics if called before Run.
func NewNotification(c NotificationConfig) error {
	return driver.NewNotification(c)
}

// MenuBar returns the menu bar.
//
// It panics if called before Run.
func MenuBar() Menu {
	return driver.MenuBar()
}

// Dock returns the dock tile.
//
// It panics if called before Run.
func Dock() DockTile {
	return driver.Dock()
}

// CallOnUIGoroutine calls a function on the UI goroutine.
// UI goroutine is the running application main thread.
func CallOnUIGoroutine(f func()) {
	driver.CallOnUIGoroutine(f)
}
