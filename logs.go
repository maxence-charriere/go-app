package app

import (
	"encoding/json"
	"fmt"
)

// Driver logs.
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

	case DockTile:
		return &dockWithLogs{DockTile: e}

	case Menu:
		return &menuWithLogs{Menu: e}

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

func (d *driverWithLogs) NewContextMenu(c MenuConfig) Menu {
	WhenDebug(func() {
		config, _ := json.MarshalIndent(c, "", "  ")
		Debug("creating context menu: %s", config)
	})

	m := d.Driver.NewContextMenu(c)
	if m.Err() != nil {
		Log("creating context menu failed: %s", m.Err())
	}

	return &menuWithLogs{Menu: m}
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

func (d *driverWithLogs) MenuBar() Menu {
	WhenDebug(func() {
		Debug("getting menubar")
	})

	m := d.Driver.MenuBar()
	if m.Err() != nil {
		Log("getting menubar failed: %s", m.Err())
	}

	return &menuWithLogs{Menu: m}
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

func (d *driverWithLogs) Dock() DockTile {
	WhenDebug(func() {
		Debug("getting dock tile")
	})

	dt := d.Driver.Dock()
	if dt.Err() != nil {
		Log("getting dock tile failed: %s", dt.Err())
	}

	return &dockWithLogs{DockTile: dt}
}

// Window logs.
type windowWithLogs struct {
	Window
}

func (w *windowWithLogs) WhenWindow(f func(Window)) {
	f(w)
}

func (w *windowWithLogs) WhenNavigator(f func(Navigator)) {
	f(w)
}

func (w *windowWithLogs) Load(url string, v ...interface{}) {
	parsedURL := fmt.Sprintf(url, v...)

	WhenDebug(func() {
		Debug("window %s is loading %s",
			w.ID(),
			parsedURL,
		)
	})

	w.Window.Load(url, v...)
	if w.Err() != nil {
		Log("window %s failed to load %s: %s",
			w.ID(),
			parsedURL,
			w.Err(),
		)
	}
}

func (w *windowWithLogs) Render(c Compo) {
	WhenDebug(func() {
		Debug("window %s is rendering %T",
			w.ID(),
			c,
		)
	})

	w.Window.Render(c)
	if w.Err() != nil {
		Log("window %s failed to render %T: %s",
			w.ID(),
			c,
			w.Err(),
		)
	}
}

func (w *windowWithLogs) Reload() {
	WhenDebug(func() {
		Debug("window %s is reloading", w.ID())
	})

	w.Window.Reload()
	if w.Err() != nil {
		Log("window %s failed to reload: %s",
			w.ID(),
			w.Err(),
		)
	}
}

func (w *windowWithLogs) Previous() {
	WhenDebug(func() {
		Debug("window %s is loading previous", w.ID())
	})

	w.Window.Previous()
	if w.Err() != nil {
		Log("window %s failed to load previous: %s",
			w.ID(),
			w.Err(),
		)
	}
}

func (w *windowWithLogs) Next() {
	WhenDebug(func() {
		Debug("window %s is loading next", w.ID())
	})

	w.Window.Next()
	if w.Err() != nil {
		Log("window %s failed to load next: %s",
			w.ID(),
			w.Err(),
		)
	}
}

func (w *windowWithLogs) Close() {
	WhenDebug(func() {
		Debug("window %s is closing", w.ID())
	})

	w.Window.Close()
	if w.Err() != nil {
		Log("window %s failed to close: %s",
			w.ID(),
			w.Err(),
		)
	}
}

func (w *windowWithLogs) Move(x, y float64) {
	WhenDebug(func() {
		Debug("window %s is moving to x:%.2f y:%.2f",
			w.ID(),
			x,
			y,
		)
	})

	w.Window.Move(x, y)
}

func (w *windowWithLogs) Center() {
	WhenDebug(func() {
		Debug("window %s is moving to center", w.ID())
	})

	w.Window.Center()
}

func (w *windowWithLogs) Resize(width, height float64) {
	WhenDebug(func() {
		Debug("window %s is resizing to width:%.2f height:%.2f",
			w.ID(),
			width,
			height,
		)
	})

	w.Window.Resize(width, height)
}

func (w *windowWithLogs) Focus() {
	WhenDebug(func() {
		Debug("window %s is getting focus", w.ID())
	})

	w.Window.Focus()
}

func (w *windowWithLogs) FullScreen() {
	WhenDebug(func() {
		Debug("window %s is entering full screen", w.ID())
	})

	w.Window.FullScreen()
}

func (w *windowWithLogs) ExitFullScreen() {
	WhenDebug(func() {
		Debug("window %s is exiting full screen", w.ID())
	})

	w.Window.ExitFullScreen()
}

func (w *windowWithLogs) Minimize() {
	WhenDebug(func() {
		Debug("window %s is minimizing", w.ID())
	})

	w.Window.Minimize()
}

func (w *windowWithLogs) Deminimize() {
	WhenDebug(func() {
		Debug("window %s is deminimizing", w.ID())
	})

	w.Window.Deminimize()
}

// Menu logs.
type menuWithLogs struct {
	Menu
}

func (m *menuWithLogs) Load(url string, v ...interface{}) {
	parsedURL := fmt.Sprintf(url, v...)

	WhenDebug(func() {
		Debug("%s %s is loading %s",
			m.Type(),
			m.ID(),
			parsedURL,
		)
	})

	m.Menu.Load(url, v...)
	if m.Err() != nil {
		Log("%s %s failed to load %s: %s",
			m.Type(),
			m.ID(),
			parsedURL,
			m.Err(),
		)
	}
}

func (m *menuWithLogs) Render(c Compo) {
	WhenDebug(func() {
		Debug("%s %s is rendering %T",
			m.Type(),
			m.ID(),
			c,
		)
	})

	m.Menu.Render(c)
	if m.Err() != nil {
		Log("%s %s failed to render %T: %s",
			m.Type(),
			m.ID(),
			c,
			m.Err(),
		)
	}
}

// Dock tile logs.
type dockWithLogs struct {
	DockTile
}

func (d *dockWithLogs) Load(url string, v ...interface{}) {
	parsedURL := fmt.Sprintf(url, v...)

	WhenDebug(func() {
		Debug("dock tile is loading %s", parsedURL)
	})

	d.DockTile.Load(url, v...)
	if d.Err() != nil {
		Log("dock tile failed to load %s: %s",
			parsedURL,
			d.Err(),
		)
	}
}

func (d *dockWithLogs) Render(c Compo) {
	WhenDebug(func() {
		Debug("dock tile is rendering %T", c)
	})

	d.DockTile.Render(c)
	if d.Err() != nil {
		Log("dock tile failed to render %T: %s",
			c,
			d.Err(),
		)
	}
}

func (d *dockWithLogs) SetIcon(name string) {
	WhenDebug(func() {
		Debug("dock tile is setting its icon to %s", name)
	})

	d.DockTile.SetIcon(name)
	if d.Err() != nil {
		Log("dock tile failed to set its icon: %s", d.Err())
	}
}

func (d *dockWithLogs) SetBadge(v interface{}) {
	WhenDebug(func() {
		Debug("dock tile is setting its badge to %d", v)
	})

	d.DockTile.SetBadge(v)
	if d.Err() != nil {
		Log("dock tile failed to set its badge: %s", d.Err())
	}
}

// Status menu logs.
type statusMenuWithLogs struct {
	StatusMenu
}

func (s *statusMenuWithLogs) Load(url string, v ...interface{}) {
	parsedURL := fmt.Sprintf(url, v...)

	WhenDebug(func() {
		Debug("status menu %s is loading %s",
			s.ID(),
			parsedURL,
		)
	})

	s.StatusMenu.Load(url, v...)
	if s.Err() != nil {
		Log("status menu %T failed to load %s: %s",
			s.ID(),
			parsedURL,
			s.Err(),
		)
	}
}

func (s *statusMenuWithLogs) Render(c Compo) {
	WhenDebug(func() {
		Debug("status menu %s is rendering %T",
			s.ID(),
			c,
		)
	})

	s.StatusMenu.Render(c)
	if s.Err() != nil {
		Log("status menu %s failed to render %T: %s",
			s.ID(),
			c,
			s.Err(),
		)
	}
}

func (s *statusMenuWithLogs) SetIcon(name string) error {
	WhenDebug(func() {
		Debug("status menu %s is setting icon to %s",
			s.ID(),
			name,
		)
	})

	err := s.StatusMenu.SetIcon(name)
	if err != nil {
		Log("status menu %s failed to set icon: %s",
			s.ID(),
			err,
		)
	}
	return err
}

func (s *statusMenuWithLogs) SetText(text string) error {
	WhenDebug(func() {
		Debug("status menu %s is setting text to %s",
			s.ID(),
			text,
		)
	})

	err := s.StatusMenu.SetText(text)
	if err != nil {
		Log("status menu %s failed to set text: %s",
			s.ID(),
			err,
		)
	}
	return err
}

func (s *statusMenuWithLogs) Close() {
	WhenDebug(func() {
		Debug("status menu %s is closing", s.ID())
	})

	s.StatusMenu.Close()
	if s.Err() != nil {
		Log("status menu %s failed to close: %s",
			s.ID(),
			s.Err(),
		)
	}
}
