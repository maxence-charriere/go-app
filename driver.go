package app

import "encoding/json"

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
	NewPage(c PageConfig) error

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
func (d *BaseDriver) NewPage(c PageConfig) error {
	return NewErrNotSupported("page")
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

// Addon represents a driver addon.
type Addon func(Driver) Driver

// Logs returns an addons that logs all the driver operations.
// It uses the loggers defined in app.Loggers.
func Logs() func(Driver) Driver {
	return func(d Driver) Driver {
		return &driverWithLogs{
			Driver: d,
		}
	}
}

type driverWithLogs struct {
	Driver
}

func (d *driverWithLogs) Run(f Factory) error {
	WhenDebug(func() {
		Debug("running %s driver", d.Name())
	})

	err := d.Driver.Run(f)
	if err != nil {
		Log("driver stopped running: %s", err)
	}
	return err
}

func (d *driverWithLogs) Resources(path ...string) string {
	resources := d.Driver.Resources(path...)
	Log("resources path: %s", resources)
	return resources
}

func (d *driverWithLogs) Storage(path ...string) string {
	storage := d.Driver.Resources(path...)
	Log("storage path: %s", storage)
	return storage
}

func (d *driverWithLogs) NewWindow(c WindowConfig) (Window, error) {
	WhenDebug(func() {
		config, _ := json.Marshal(c)
		Log("creating window: %s", config)
	})

	win, err := d.Driver.NewWindow(c)
	if err != nil {
		Log("creating window failed: %s", err)
		return nil, err
	}

	win = &windowWithLogs{
		Window: win,
	}
	return win, err
}
