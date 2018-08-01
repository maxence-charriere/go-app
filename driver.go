package app

// Driver is the interface that describes a backend for app rendering.
type Driver interface {
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

	// Render renders the given component.
	Render(c Compo) error

	// ElemByCompo returns the element where the given component is mounted.
	ElemByCompo(c Compo) Elem

	// NewWindow creates and displays the window described by the given
	// configuration.
	NewWindow(c WindowConfig) Window

	// NewContextMenu creates and displays the context menu described by the
	// given configuration.
	NewContextMenu(c MenuConfig) Menu

	// NewPage creates the webpage described in the given configuration.
	NewPage(c PageConfig) error

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
	MenuBar() Menu

	// NewStatusMenu creates a status menu.
	NewStatusMenu(c StatusMenuConfig) (StatusMenu, error)

	// Dock returns the dock tile.
	Dock() DockTile

	// CallOnUIGoroutine calls a function on the UI goroutine.
	CallOnUIGoroutine(f func())
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
