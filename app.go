package app

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

var (
	driver     Driver
	components Factory
)

func init() {
	components = NewFactory()
	components = ConcurrentFactory(components)

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
func Import(c Component) {
	if _, err := components.Register(c); err != nil {
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
	return driver.Run(components)
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

// CSSResources returns a list that contains the path of the css files located
// in the resource/css directory.
func CSSResources() []string {
	var css []string

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if ext := filepath.Ext(path); ext != ".css" {
			return nil
		}

		css = append(css, path)
		return nil
	}

	filepath.Walk(Resources("css"), walker)
	return css
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

// NewPage creates the page described by the given configuration.
//
// It panics if called before Run.
func NewPage(c PageConfig) error {
	return driver.NewPage(c)
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
	driver.CallOnUIGoroutine(func() {
		driver.Render(c)
	})
}

// ElementByComponent returns the element where the given component is mounted.
//
// It panics if called before Run.
func ElementByComponent(c Component) (ElementWithComponent, error) {
	return driver.ElementByComponent(c)
}

// NavigatorByComponent returns the navigator where the given component is
// mounted.
func NavigatorByComponent(c Component) (Navigator, error) {
	elem, err := driver.ElementByComponent(c)
	if err != nil {
		return nil, err
	}

	nav, ok := elem.(Navigator)
	if !ok {
		return nil, errors.New("component is not mounted into a navigator")
	}
	return nav, nil
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

// PageByComponent returns the page where the given component is mounted.
func PageByComponent(c Component) (Page, error) {
	elem, err := driver.ElementByComponent(c)
	if err != nil {
		return nil, err
	}

	page, ok := elem.(Page)
	if !ok {
		return nil, errors.New("component is not mounted in a page")
	}
	return page, nil
}

// NewFilePanel creates and displays the file panel described by the given
// configuration.
//
// It panics if called before Run.
func NewFilePanel(c FilePanelConfig) error {
	return driver.NewFilePanel(c)
}

// NewSaveFilePanel creates and displays the save file panel described by the
// given configuration.
//
// It panics if called before Run.
func NewSaveFilePanel(c SaveFilePanelConfig) error {
	return driver.NewSaveFilePanel(c)
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
func MenuBar() (Menu, error) {
	return driver.MenuBar()
}

// StatusBar returns the status bar menu.
//
// It panics if called before Run.
func StatusBar() (StatusBarMenu, error) {
	return driver.StatusBar()
}

// Dock returns the dock tile.
//
// It panics if called before Run.
func Dock() (DockTile, error) {
	return driver.Dock()
}

// CallOnUIGoroutine calls a function on the UI goroutine.
// UI goroutine is the running application main thread.
func CallOnUIGoroutine(f func()) {
	driver.CallOnUIGoroutine(f)
}
