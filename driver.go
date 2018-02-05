package app

// Driver is the interface that describes a backend for app rendering.
type Driver interface {
	// Run runs the application with the components resistered in the factory.
	Run(factory Factory) error

	// Render renders the component.
	Render(c Component) error

	// Context returns the element where the component is mounted.
	// It returns an error if c is not mounted.
	Context(c Component) (ElementWithComponent, error)

	// NewContextMenu creates and displays the context menu described in the
	// configuration.
	NewContextMenu(c MenuConfig) Menu

	// AppName returns the appliction name.
	AppName() string

	// Resources returns the location of the resources directory.
	Resources() string

	// CallOnUIGoroutine calls a function on the UI goroutine.
	CallOnUIGoroutine(f func())
}

// DriverWithStorage is the interface that describes a driver which supports
// storage.
// A storage is the local directory where an application can save files.
// This is the ideal location to save dynamic contents and local database files.
type DriverWithStorage interface {
	Driver

	// Storage returns the location of the storage directory.
	Storage() string
}

// DriverWithWindows is the interface that describes a driver able to create
// windows.
type DriverWithWindows interface {
	Driver

	// NewWindow creates and displays the window described in the configuration.
	NewWindow(c WindowConfig) Window
}

// DriverWithMenuBar is the interface that describes a driver with a menu bar.
type DriverWithMenuBar interface {
	Driver

	// MenuBar returns the menu bar.
	MenuBar() Menu
}

// DriverWithDock is the interface that describes a driver with a dock.
type DriverWithDock interface {
	Driver

	// Dock returns the dock tile.
	Dock() DockTile
}

// DriverWithShare is the interface that describes a driver with sharing
// support.
type DriverWithShare interface {
	Driver

	// Share shares the value.
	Share(v interface{}) error
}

// DriverWithFilePanels is the interface that describes a driver able to open a
// file panel.
type DriverWithFilePanels interface {
	Driver

	// NewFilePanel creates and displays the file panel described in the
	// configuration.
	NewFilePanel(c FilePanelConfig) Element
}

// DriverWithNotifications is the interface that describes a driver able to
// display notifications.
type DriverWithNotifications interface {
	// NewNotification creates and displays the notification described in the
	// given configuration.
	NewNotification(c NotificationConfig) error
}
