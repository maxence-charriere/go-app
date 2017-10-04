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
	if driver != nil {
		panic(errors.Errorf("importing component %T failed: can't import when a driver is running", c))
	}

	if err := compoBuilder.Register(c); err != nil {
		err = errors.Wrapf(err, "importing component %T failed", c)
		panic(err)
	}
}

// Run runs the app with driver d as backend.
func Run(d Driver) {
	if driver != nil {
		panic(errors.Errorf("driver %T is already running", driver))
	}

	driver = d
	if err := d.Run(compoBuilder); err != nil {
		driver = nil
		panic(errors.Wrap(err, "running the app failed"))
	}
}

// RunningDriver returns the running driver.
// It panics if called before Run.
func RunningDriver() Driver {
	if driver == nil {
		panic("no current driver")
	}
	return driver
}

// Render renders component c.
// It should be called when the display of component c have to be updated.
// It panics if called before Run.
func Render(c markup.Component) {
	if err := driver.Render(c); err != nil {
		Logs().Error(err)
	}
}

// Context returns the element where component c is mounted.
// It returns an error if c is not mounted.
// It panics if called before Run.
func Context(c markup.Component) (ElementWithComponent, error) {
	return driver.Context(c)
}

// NewContextMenu creates and displays the context menu described in
// configuration c.
// Context menu are displayed in the window or page in use.
// It panics if called before Run.
func NewContextMenu(c MenuConfig) Menu {
	return driver.NewContextMenu(c)
}

// Resources returns the location of the resources directory.
// Resources should be used only for read only operations.
// It panics if called before Run.
func Resources() string {
	return driver.Resources()
}

// Logs returns the application logger.
// It panics if called before Run.
func Logs() Logger {
	return driver.Logs()
}

// CallOnUIGoroutine calls func f and ensure it's called from the UI goroutine.
// UI goroutine is the running application main thread.
func CallOnUIGoroutine(f func()) {
	driver.CallOnUIGoroutine(f)
}

// SupportsStorage reports whether storage is supported.
func SupportsStorage() bool {
	_, ok := driver.(DriverWithStorage)
	return ok
}

// Storage returns the location of the storage directory.
// It panics if storage is not supported.
func Storage() string {
	d := driver.(DriverWithStorage)
	return d.Storage()
}

// SupportsWindows reports whether windows are supported.
func SupportsWindows() bool {
	_, ok := driver.(DriverWithWindows)
	return ok
}

// NewWindow creates and displays the window described in configuration c.
// It panics if windows are not supported.
func NewWindow(c WindowConfig) Window {
	d := driver.(DriverWithWindows)
	return d.NewWindow(c)
}

// SupportsMenuBar reports whether menu bar is supported.
func SupportsMenuBar() bool {
	_, ok := driver.(DriverWithMenuBar)
	return ok
}

// MenuBar returns the menu bar.
// It panics if menu bar is not supported.
func MenuBar() Menu {
	d := driver.(DriverWithMenuBar)
	return d.MenuBar()
}

// SupportsDock reports whether dock is supported.
func SupportsDock() bool {
	_, ok := driver.(DriverWithDock)
	return ok
}

// Dock returns the dock tile.
// It panics if dock is not supported.
func Dock() DockTile {
	d := driver.(DriverWithDock)
	return d.Dock()
}

// SupportsShare reports whether share is supported.
func SupportsShare() bool {
	_, ok := driver.(DriverWithShare)
	return ok
}

// Share shares the value v.
// It panics if share is not supported.
func Share(v interface{}) {
	d := driver.(DriverWithShare)
	d.Share(v)
}

// SupportsFilePanels reports whether file panels are supported.
func SupportsFilePanels() bool {
	_, ok := driver.(DriverWithFilePanels)
	return ok
}

// NewFilePanel creates and displays the file panel described inconfiguration c.
// It panics if file panels are not supported.
func NewFilePanel(c FilePanelConfig) Element {
	d := driver.(DriverWithFilePanels)
	return d.NewFilePanel(c)
}

// SupportsPopupNotifications reports whether popup notifications are supported.
func SupportsPopupNotifications() bool {
	_, ok := driver.(DriverWithPopupNotifications)
	return ok
}

// NewPopupNotification creates and displays the popup notification described in
// configuration c.
// It panics if popup notifications are not supported.
func NewPopupNotification(c PopupNotificationConfig) Element {
	d := driver.(DriverWithPopupNotifications)
	return d.NewPopupNotification(c)
}
