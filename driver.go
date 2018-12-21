package app

// Driver is the interface that describes a backend for app rendering.
type Driver interface {
	// The operating system the driver is for.
	Target() string

	// Run runs the application with the components registered in the given
	// factory.
	Run(DriverConfig) error

	// AppName returns the appliction name.
	AppName() string

	// Resources returns the given path prefixed by the resources directory
	// location.
	Resources(path ...string) string

	// Storage returns the given path prefixed by the storage directory
	// location.
	Storage(path ...string) string

	// Render renders the given component.
	Render(Compo)

	// ElemByCompo returns the element where the given component is mounted.
	ElemByCompo(Compo) Elem

	// NewWindow creates and displays the window described by the given
	// configuration.
	NewWindow(WindowConfig) Window

	// NewPage creates the webpage described in the given configuration.
	NewPage(PageConfig) Page

	// NewContextMenu creates and displays the context menu described by the
	// given configuration.
	NewContextMenu(MenuConfig) Menu

	// NewController creates the controller described by the given
	// configuration.
	NewController(ControllerConfig) Controller

	// NewFilePanel creates and displays the file panel described by the given
	// configuration.
	NewFilePanel(FilePanelConfig) Elem

	// NewSaveFilePanel creates and displays the save file panel described in
	// the given configuration.
	NewSaveFilePanel(SaveFilePanelConfig) Elem

	// NewShare creates and display the share pannel to share the given value.
	NewShare(interface{}) Elem

	// NewNotification creates and displays the notification described in the
	// given configuration.
	NewNotification(NotificationConfig) Elem

	// MenuBar returns the menu bar.
	MenuBar() Menu

	// NewStatusMenu creates a status menu.
	NewStatusMenu(StatusMenuConfig) StatusMenu

	// Dock returns the dock tile.
	DockTile() DockTile

	// UI calls a function on the UI goroutine.
	UI(func())

	// Stop stops the driver.
	// Calling it make run return with an error.
	Stop()
}

// DriverConfig contains driver configuration.
type DriverConfig struct {
	// The channel to send function to execute on UI goroutine.
	UI chan func()

	// The factory used to create components.
	Factory *Factory

	// The event registery to emit events.
	Events *EventRegistry
}

// Addon represents a driver addon.
type Addon func(Driver) Driver
