package app

import "encoding/json"

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
	NewContextMenu(c MenuConfig) (Menu, error)

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
	MenuBar() (Menu, error)

	// NewStatusMenu creates a status menu.
	NewStatusMenu(c StatusMenuConfig) (StatusMenu, error)

	// Dock returns the dock tile.
	Dock() (DockTile, error)

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

type driverWithLogs struct {
	Driver
}

func (d *driverWithLogs) Run(f Factory) error {
	WhenDebug(func() {
		Debug("running %T driver", d)
	})

	err := d.Driver.Run(f)
	if err != nil {
		Log("driver stopped running: %s", err)
	}
	return err
}

func (d *driverWithLogs) Render(c Compo) error {
	WhenDebug(func() {
		Debug("rendering %T", c)
	})

	err := d.Driver.Render(c)
	if err != nil {
		Log("rendering %T failed: %s", err)
	}
	return err
}

func (d *driverWithLogs) ElemByCompo(c Compo) Elem {
	WhenDebug(func() {
		Debug("getting element from %T", c)
	})

	switch e := d.Driver.ElemByCompo(c).(type) {
	case Window:
		return &windowWithLogs{Window: e}

	default:
		return e
	}
}

func (d *driverWithLogs) NewWindow(c WindowConfig) Window {
	WhenDebug(func() {
		config, _ := json.MarshalIndent(c, "", "    ")
		Debug("creating window: %s", config)
	})

	w := d.Driver.NewWindow(c)
	if w.Err() != nil {
		Log("creating window failed: %s", w.Err())
	}

	return &windowWithLogs{Window: w}
}

func (d *driverWithLogs) NewContextMenu(c MenuConfig) (Menu, error) {
	c.Type = "context menu"

	WhenDebug(func() {
		config, _ := json.MarshalIndent(c, "", "  ")
		Debug("creating context menu: %s", config)
	})

	menu, err := d.Driver.NewContextMenu(c)
	if err != nil {
		Log("creating context menu failed: %s", err)
		return nil, err
	}

	menu = &menuWithLogs{
		Menu: menu,
	}
	return menu, nil
}

func (d *driverWithLogs) NewPage(c PageConfig) error {
	WhenDebug(func() {
		config, _ := json.MarshalIndent(c, "", "  ")
		Debug("creating page: %s", config)
	})

	err := d.Driver.NewPage(c)
	if err != nil {
		Log("creating page failed: %s", err)
	}
	return err
}

func (d *driverWithLogs) NewFilePanel(c FilePanelConfig) error {
	WhenDebug(func() {
		config, _ := json.MarshalIndent(c, "", "  ")
		Debug("creating file panel: %s", config)
	})

	err := d.Driver.NewFilePanel(c)
	if err != nil {
		Log("creating file panel failed: %s", err)
	}
	return err
}

func (d *driverWithLogs) NewSaveFilePanel(c SaveFilePanelConfig) error {
	WhenDebug(func() {
		config, _ := json.MarshalIndent(c, "", "  ")
		Debug("creating save file panel: %s", config)
	})

	err := d.Driver.NewSaveFilePanel(c)
	if err != nil {
		Log("creating save file panel failed: %s", err)
	}
	return err
}

func (d *driverWithLogs) NewShare(v interface{}) error {
	WhenDebug(func() {
		Debug("creating share: %v", v)
	})

	err := d.Driver.NewShare(v)
	if err != nil {
		Log("creating share failed: %s", err)
	}
	return err
}

func (d *driverWithLogs) NewNotification(c NotificationConfig) error {
	WhenDebug(func() {
		config, _ := json.MarshalIndent(c, "", "  ")
		Debug("creating notification: %s", config)
	})

	err := d.Driver.NewNotification(c)
	if err != nil {
		Log("creating notification failed: %s", err)
	}
	return err
}

func (d *driverWithLogs) MenuBar() (Menu, error) {
	WhenDebug(func() {
		Debug("getting menubar")
	})

	menubar, err := d.Driver.MenuBar()
	if err != nil {
		Log("getting menubar failed: %s", err)
		return nil, err
	}

	menubar = &menuWithLogs{
		Menu: menubar,
	}
	return menubar, nil
}

func (d *driverWithLogs) NewStatusMenu(c StatusMenuConfig) (StatusMenu, error) {
	WhenDebug(func() {
		config, _ := json.MarshalIndent(c, "", "  ")
		Debug("creating status menu: %s", config)
	})

	menu, err := d.Driver.NewStatusMenu(c)
	if err != nil {
		Log("getting status menu failed: %s", err)
		return nil, err
	}

	menu = &statusMenuWithLogs{
		StatusMenu: menu,
	}
	return menu, nil
}

func (d *driverWithLogs) Dock() (DockTile, error) {
	WhenDebug(func() {
		Debug("getting dock tile")
	})

	dockTile, err := d.Driver.Dock()
	if err != nil {
		Log("getting dock tile failed: %s", err)
		return nil, err
	}

	dockTile = &dockWithLogs{
		DockTile: dockTile,
	}
	return dockTile, nil
}
