package app

// Driver is the interface that describes a backend for app rendering.
type Driver interface {

	// Returns the appliction name.
	AppName() string

	// Dock returns the dock tile.
	DockTile() DockTile

	// Returns the element where the given component is mounted.
	ElemByCompo(Compo) Elem

	// Returns the current menu bar.
	MenuBar() Menu

	// Creates and displays the context menu described by the given
	// configuration.
	NewContextMenu(MenuConfig) Menu

	// Creates the controller described by the given configuration.
	NewController(ControllerConfig) Controller

	// Creates and displays the file panel described by the given configuration.
	NewFilePanel(FilePanelConfig) Elem

	// Creates and displays the notification described in the given
	// configuration.
	NewNotification(NotificationConfig) Elem

	// Creates and displays the save file panel described in the given
	// configuration.
	NewSaveFilePanel(SaveFilePanelConfig) Elem

	// Creates a status menu.
	NewStatusMenu(StatusMenuConfig) StatusMenu

	// Creates and display the share pannel to share the given value.
	NewShare(interface{}) Elem

	// Creates and displays the window described by the given configuration.
	NewWindow(WindowConfig) Window

	// Opens the given URL on the operating system default browser.
	OpenDefaultBrowser(string) error

	// Runs the application.
	Run(DriverConfig) error

	// Renders the given component.
	Render(Compo)

	// Returns the given path prefixed by the resources directory location.
	Resources(path ...string) string

	// Returns the given path prefixed by the storage directory location.
	Storage(path ...string) string

	// Stops the driver. Calling it make Run to return with an error.
	Stop()

	// The operating system the driver is for.
	Target() string

	// Calls a function on the UI goroutine.
	UI(func())
}

// DriverConfig contains driver configuration.
type DriverConfig struct {
	// The event registery to emit events.
	Events *EventRegistry

	// The factory used to create components.
	Factory *Factory

	// The channel to send function to execute on UI goroutine.
	UI chan func()
}

// Addon represents a driver addon.
type Addon func(Driver) Driver
