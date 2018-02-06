package app

import (
	"encoding/json"
)

// Driver is the interface that describes a backend for app rendering.
type Driver interface {
	// Name returns the driver name.
	Name() string

	// Run runs the application with the components resistered in the given
	// factory.
	Run(f Factory) error

	// AppName returns the appliction name.
	AppName() string

	// Resources returns the given path prefixed by the resources directory
	// location.
	Resources(path ...string) string

	// Storage returns the given path prefixed by the storage directory
	// location.
	Storage(path ...string) string

	// NewWindow creates and displays the window described by the given
	// configuration.
	NewWindow(c WindowConfig) (Window, error)

	// NewContextMenu creates and displays the context menu described by the
	// given configuration.
	NewContextMenu(c MenuConfig) (Menu, error)

	// Render renders the given component.
	Render(c Component) error

	// ElementByComponent returns the element where the given component is mounted.
	ElementByComponent(c Component) (ElementWithComponent, error)

	// NewFilePanel creates and displays the file panel described by the given
	// configuration.
	NewFilePanel(c FilePanelConfig) error

	// NewShare creates and display the share pannel to share the given value.
	NewShare(v interface{}) error

	// NewNotification creates and displays the notification described in the
	// given configuration.
	NewNotification(c NotificationConfig) error

	// MenuBar returns the menu bar.
	MenuBar() Menu

	// Dock returns the dock tile.
	Dock() DockTile

	// CallOnUIGoroutine calls a function on the UI goroutine.
	CallOnUIGoroutine(f func())
}

// NewDriverWithLogs returns a decorated version of the given driver that logs
// all the operations.
// It uses the default logger.
func NewDriverWithLogs(driver Driver) Driver {
	return &driverWithLogs{
		base: driver,
	}
}

type driverWithLogs struct {
	base Driver
}

func (d *driverWithLogs) Name() string {
	name := d.base.Name()
	Log("driver name:", name)
	return name
}

func (d *driverWithLogs) Run(f Factory) error {
	Log("running driver", d.base.Name())

	err := d.base.Run(f)
	if err != nil {
		Error("running driver returned an error:", err)
	}
	return err
}

func (d *driverWithLogs) AppName() string {
	appName := d.base.AppName()
	Log("app name:", appName)
	return appName
}

func (d *driverWithLogs) Resources(path ...string) string {
	resources := d.base.Resources(path...)
	Log("resources path:", resources)
	return resources
}

func (d *driverWithLogs) Storage(path ...string) string {
	storage := d.base.Storage(path...)
	Log("storage path:", storage)
	return storage
}

func (d *driverWithLogs) NewWindow(c WindowConfig) (Window, error) {
	Log("creating window:", indentedJSON(c))

	win, err := d.base.NewWindow(c)
	if err != nil {
		Error("creating window failed:", err)
	}
	return win, err
}

func (d *driverWithLogs) NewContextMenu(c MenuConfig) (Menu, error) {
	Log("creating context menu:", indentedJSON(c))

	menu, err := d.base.NewContextMenu(c)
	if err != nil {
		Error("creating context menu failed:", err)
	}
	return menu, err
}

func (d *driverWithLogs) Render(c Component) error {
	Logf("rendering %T", c)

	err := d.base.Render(c)
	if err != nil {
		Errorf("rendering %T failed: %s", c, err)
	}
	return err
}

func (d *driverWithLogs) ElementByComponent(c Component) (ElementWithComponent, error) {
	Logf("returning element that hosts %T", c)

	elem, err := d.base.ElementByComponent(c)
	if err != nil {
		Errorf("returning element that hosts %T failed: %s", c, err)
	}
	return elem, err
}

func (d *driverWithLogs) NewFilePanel(c FilePanelConfig) error {
	Log("creating file panel:", indentedJSON(c))

	err := d.base.NewFilePanel(c)
	if err != nil {
		Error("creating file panel failed:", err)
	}
	return err
}

func (d *driverWithLogs) NewShare(v interface{}) error {
	Log("sharing", v)

	err := d.base.NewShare(v)
	if err != nil {
		Error("sharing failed:", err)
	}
	return err
}

func (d *driverWithLogs) NewNotification(c NotificationConfig) error {
	Log("creating notification:", indentedJSON(c))

	err := d.base.NewNotification(c)
	if err != nil {
		Error("creating notification failed:", err)
	}
	return err
}

func (d *driverWithLogs) MenuBar() Menu {
	Log("returning menu bar")
	return d.base.MenuBar()
}

func (d *driverWithLogs) Dock() DockTile {
	Log("returning dock tile")
	return d.base.Dock()
}

func (d *driverWithLogs) CallOnUIGoroutine(f func()) {
	Log("calling a function on the UI goroutine")
	d.base.CallOnUIGoroutine(f)
}

func indentedJSON(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}
