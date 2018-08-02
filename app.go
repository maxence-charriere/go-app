package app

import (
	"github.com/pkg/errors"
)

var (
	// ErrNotSupported describes an error that occurs when an unsupported
	// feature is used.
	ErrNotSupported = errors.New("not supported")

	// ErrElemNotSet describes an error that reports if an element is set.
	ErrElemNotSet = errors.New("element not set")

	driver  Driver
	factory = NewFactory()
)

func init() {
	events := NewEventRegistry(CallOnUIGoroutine)
	events = ConcurrentEventRegistry(events)
	DefaultEventRegistry = events

	actions := NewActionRegistry(events)
	DefaultActionRegistry = actions
}

// Import imports the component into the app.
// Components must be imported in order the be used by the app package.
// This allows components to be created dynamically when they are found into
// markup.
func Import(c Compo) {
	if _, err := factory.RegisterCompo(c); err != nil {
		err = errors.Wrap(err, "import component failed")
		panic(err)
	}
}

// Run runs the app with the driver as backend.
func Run(d Driver, addons ...Addon) error {
	for _, addon := range addons {
		d = addon(d)
	}
	driver = d
	return driver.Run(factory)
}

// RunningDriver returns the running driver.
func RunningDriver() Driver {
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

// Render renders the given component.
// It should be called when the display of component c have to be updated.
//
// It panics if called before Run.
func Render(c Compo) {
	driver.CallOnUIGoroutine(func() {
		driver.Render(c)
	})
}

// ElemByCompo returns the element where the given component is mounted.
//
// It panics if called before Run.
func ElemByCompo(c Compo) Elem {
	return driver.ElemByCompo(c)
}

// NewWindow creates and displays the window described by the given
// configuration.
//
// It panics if called before Run.
func NewWindow(c WindowConfig) Window {
	return driver.NewWindow(c)
}

// NewPage creates the page described by the given configuration.
//
// It panics if called before Run.
func NewPage(c PageConfig) Elem {
	return driver.NewPage(c)
}

// NewContextMenu creates and displays the context menu described by the
// given configuration.
//
// It panics if called before Run.
func NewContextMenu(c MenuConfig) Menu {
	return driver.NewContextMenu(c)
}

// NewFilePanel creates and displays the file panel described by the given
// configuration.
//
// It panics if called before Run.
func NewFilePanel(c FilePanelConfig) Elem {
	return driver.NewFilePanel(c)
}

// NewSaveFilePanel creates and displays the save file panel described by the
// given configuration.
//
// It panics if called before Run.
func NewSaveFilePanel(c SaveFilePanelConfig) Elem {
	return driver.NewSaveFilePanel(c)
}

// NewShare creates and display the share pannel to share the given value.
//
// It panics if called before Run.
func NewShare(v interface{}) Elem {
	return driver.NewShare(v)
}

// NewNotification creates and displays the notification described in the
// given configuration.
//
// It panics if called before Run.
func NewNotification(c NotificationConfig) Elem {
	return driver.NewNotification(c)
}

// MenuBar returns the menu bar.
//
// It panics if called before Run.
func MenuBar() Menu {
	return driver.MenuBar()
}

// NewStatusMenu creates and displays the status menu described in the given
// configuration.
//
// It panics if called before Run.
func NewStatusMenu(c StatusMenuConfig) StatusMenu {
	return driver.NewStatusMenu(c)
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
