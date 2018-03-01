package app

import (
	"encoding/json"
)

// Driver is the interface that describes a backend for app rendering.
type Driver interface {
	// Name returns the driver name.
	Name() string

	// Run runs the application with the components registered in the given
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

	// NewPage creates the webpage described in the given configuration.
	NewPage(c PageConfig) (Page, error)

	// Render renders the given component.
	Render(c Component) error

	// ElementByComponent returns the element where the given component is mounted.
	ElementByComponent(c Component) (ElementWithComponent, error)

	// NewFilePanel creates and displays the file panel described by the given
	// configuration.
	NewFilePanel(c FilePanelConfig) error

	// NewSaveFilePanel creates and displays the save file panel described in
	// the given configuration.
	NewSaveFilePanel(c SaveFilePanelConfig) error

	// NewShare creates and display the share pannel to share the given value.
	NewShare(v interface{}) error

	// NewNotification creates and displays the notification described in the
	// given configuration.
	NewNotification(c NotificationConfig) error

	// MenuBar returns the menu bar.
	MenuBar() (Menu, error)

	// Dock returns the dock tile.
	Dock() (DockTile, error)

	// CallOnUIGoroutine calls a function on the UI goroutine.
	CallOnUIGoroutine(f func())
}

// BaseDriver represents a base driver to be embedded in app.Driver
// implementations.
// It only contains methods related to features.
// All the methods return not supported error.
type BaseDriver struct{}

// NewWindow satisfies the app.Driver interface.
func (d *BaseDriver) NewWindow(c WindowConfig) (Window, error) {
	return nil, NewErrNotSupported("window")
}

// NewContextMenu satisfies the app.Driver interface.
func (d *BaseDriver) NewContextMenu(c MenuConfig) (Menu, error) {
	return nil, NewErrNotSupported("context menu")
}

// NewPage satisfies the app.Driver interface.
func (d *BaseDriver) NewPage(c PageConfig) (Page, error) {
	return nil, NewErrNotSupported("page")
}

// NewFilePanel satisfies the app.Driver interface.
func (d *BaseDriver) NewFilePanel(c FilePanelConfig) error {
	return NewErrNotSupported("file panel")
}

// NewSaveFilePanel satisfies the app.Driver interface.
func (d *BaseDriver) NewSaveFilePanel(c SaveFilePanelConfig) error {
	return NewErrNotSupported("save file panel")
}

// NewShare satisfies the app.Driver interface.
func (d *BaseDriver) NewShare(v interface{}) error {
	return NewErrNotSupported("share")
}

// NewNotification satisfies the app.Driver interface.
func (d *BaseDriver) NewNotification(c NotificationConfig) error {
	return NewErrNotSupported("notification")
}

// MenuBar satisfies the app.Driver interface.
func (d *BaseDriver) MenuBar() (Menu, error) {
	return nil, NewErrNotSupported("menubar")
}

// Dock satisfies the app.Driver interface.
func (d *BaseDriver) Dock() (DockTile, error) {
	return nil, NewErrNotSupported("dock")
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
	return d.base.AppName()
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

func (d *driverWithLogs) NewPage(c PageConfig) (Page, error) {
	Log("creating page:", indentedJSON(c))

	page, err := d.base.NewPage(c)
	if err != nil {
		Error("creating page failed:", err)
	}
	return page, err
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

func (d *driverWithLogs) NewSaveFilePanel(c SaveFilePanelConfig) error {
	Log("creating save file panel:", indentedJSON(c))

	err := d.base.NewSaveFilePanel(c)
	if err != nil {
		Error("creating save file panel failed:", err)
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

func (d *driverWithLogs) MenuBar() (Menu, error) {
	Log("returning menu bar")

	menu, err := d.base.MenuBar()
	if err != nil {
		Errorf("returning menubar failed: %s", err)
	}
	return menu, err
}

func (d *driverWithLogs) Dock() (DockTile, error) {
	Log("returning dock")

	dock, err := d.base.Dock()
	if err != nil {
		Errorf("returning dock failed: %s", err)
	}
	return dock, err
}

func (d *driverWithLogs) CallOnUIGoroutine(f func()) {
	Log("calling a function on the UI goroutine")
	d.base.CallOnUIGoroutine(f)
}

func indentedJSON(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}
