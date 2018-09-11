package app

import (
	"encoding/json"
	"fmt"
)

// Logs returns an addons that logs all the driver operations.
// It uses the loggers defined in app.Loggers.
func Logs() func(Driver) Driver {
	return func(d Driver) Driver {
		return &driverWithLogs{
			Driver: d,
		}
	}
}

// Driver logs.
type driverWithLogs struct {
	Driver
}

func (d *driverWithLogs) Run(f *Factory) error {
	Logf("running %T driver", d.Driver)

	err := d.Driver.Run(f)
	if err != nil {
		Logf("driver stopped running: %s", err)
	}
	return err
}

func (d *driverWithLogs) Render(c Compo) {
	e := d.ElemByCompo(c)
	if e.Err() == nil {
		e.(ElemWithCompo).Render(c)
	}
}

func (d *driverWithLogs) ElemByCompo(c Compo) Elem {
	switch e := d.Driver.ElemByCompo(c).(type) {
	case Window:
		return &windowWithLogs{Window: e}

	case Page:
		return &pageWithLogs{Page: e}

	case DockTile:
		return &dockWithLogs{DockTile: e}

	case StatusMenu:
		return &statusMenuWithLogs{StatusMenu: e}

	case Menu:
		return &menuWithLogs{Menu: e}

	default:
		return e
	}
}

func (d *driverWithLogs) NewWindow(c WindowConfig) Window {
	WhenDebug(func() {
		config, _ := json.MarshalIndent(c, "", "    ")
		Logf("creating window: %s", config)
	})

	w := d.Driver.NewWindow(c)
	if w.Err() != nil {
		Logf("creating window failed: %s", w.Err())
	}

	return &windowWithLogs{Window: w}
}

func (d *driverWithLogs) NewPage(c PageConfig) Page {
	WhenDebug(func() {
		config, _ := json.MarshalIndent(c, "", "  ")
		Logf("creating page: %s", config)
	})

	p := d.Driver.NewPage(c)
	if p.Err() != nil {
		Logf("creating page failed: %s", p.Err())
	}

	return &pageWithLogs{Page: p}
}

func (d *driverWithLogs) NewContextMenu(c MenuConfig) Menu {
	WhenDebug(func() {
		config, _ := json.MarshalIndent(c, "", "  ")
		Logf("creating context menu: %s", config)
	})

	m := d.Driver.NewContextMenu(c)
	if m.Err() != nil {
		Logf("creating context menu failed: %s", m.Err())
	}

	return &menuWithLogs{Menu: m}
}

func (d *driverWithLogs) NewFilePanel(c FilePanelConfig) Elem {
	WhenDebug(func() {
		config, _ := json.MarshalIndent(c, "", "  ")
		Logf("creating file panel: %s", config)
	})

	p := d.Driver.NewFilePanel(c)
	if p.Err() != nil {
		Logf("creating file panel failed: %s", p.Err())
	}

	return p
}

func (d *driverWithLogs) NewSaveFilePanel(c SaveFilePanelConfig) Elem {
	WhenDebug(func() {
		config, _ := json.MarshalIndent(c, "", "  ")
		Logf("creating save file panel: %s", config)
	})

	p := d.Driver.NewSaveFilePanel(c)
	if p.Err() != nil {
		Logf("creating save file panel failed: %s", p.Err())
	}

	return p
}

func (d *driverWithLogs) NewShare(v interface{}) Elem {
	WhenDebug(func() {
		Logf("creating share: %v", v)
	})

	s := d.Driver.NewShare(v)
	if s.Err() != nil {
		Logf("creating share failed: %s", s.Err())
	}

	return s
}

func (d *driverWithLogs) NewNotification(c NotificationConfig) Elem {
	WhenDebug(func() {
		config, _ := json.MarshalIndent(c, "", "  ")
		Logf("creating notification: %s", config)
	})

	n := d.Driver.NewNotification(c)
	if n.Err() != nil {
		Logf("creating notification failed: %s", n.Err())
	}

	return n
}

func (d *driverWithLogs) MenuBar() Menu {
	WhenDebug(func() {
		Logf("getting menubar")
	})

	m := d.Driver.MenuBar()
	if m.Err() != nil {
		Logf("getting menubar failed: %s", m.Err())
	}

	return &menuWithLogs{Menu: m}
}

func (d *driverWithLogs) NewStatusMenu(c StatusMenuConfig) StatusMenu {
	WhenDebug(func() {
		config, _ := json.MarshalIndent(c, "", "  ")
		Logf("creating status menu: %s", config)
	})

	m := d.Driver.NewStatusMenu(c)
	if m.Err() != nil {
		Logf("getting status menu failed: %s", m.Err())
	}

	return &statusMenuWithLogs{StatusMenu: m}
}

func (d *driverWithLogs) DockTile() DockTile {
	WhenDebug(func() {
		Logf("getting dock tile")
	})

	dt := d.Driver.DockTile()
	if dt.Err() != nil {
		Logf("getting dock tile failed: %s", dt.Err())
	}

	return &dockWithLogs{DockTile: dt}
}

func (d *driverWithLogs) Stop() {
	WhenDebug(func() {
		Logf("stopping driver")
	})

	d.Driver.Stop()
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
		Logf("window %s is loading %s",
			w.ID(),
			parsedURL,
		)
	})

	w.Window.Load(url, v...)
	if w.Err() != nil {
		Logf("window %s failed to load %s: %s",
			w.ID(),
			parsedURL,
			w.Err(),
		)
	}
}

func (w *windowWithLogs) Render(c Compo) {
	WhenDebug(func() {
		Logf("window %s is rendering %T",
			w.ID(),
			c,
		)
	})

	w.Window.Render(c)
	if w.Err() != nil {
		Logf("window %s failed to render %T: %s",
			w.ID(),
			c,
			w.Err(),
		)
	}
}

func (w *windowWithLogs) Reload() {
	WhenDebug(func() {
		Logf("window %s is reloading", w.ID())
	})

	w.Window.Reload()
	if w.Err() != nil {
		Logf("window %s failed to reload: %s",
			w.ID(),
			w.Err(),
		)
	}
}

func (w *windowWithLogs) Previous() {
	WhenDebug(func() {
		Logf("window %s is loading previous", w.ID())
	})

	w.Window.Previous()
	if w.Err() != nil {
		Logf("window %s failed to load previous: %s",
			w.ID(),
			w.Err(),
		)
	}
}

func (w *windowWithLogs) Next() {
	WhenDebug(func() {
		Logf("window %s is loading next", w.ID())
	})

	w.Window.Next()
	if w.Err() != nil {
		Logf("window %s failed to load next: %s",
			w.ID(),
			w.Err(),
		)
	}
}

func (w *windowWithLogs) Move(x, y float64) {
	WhenDebug(func() {
		Logf("window %s is moving to x:%.2f y:%.2f",
			w.ID(),
			x,
			y,
		)
	})

	w.Window.Move(x, y)
}

func (w *windowWithLogs) Center() {
	WhenDebug(func() {
		Logf("window %s is moving to center", w.ID())
	})

	w.Window.Center()
}

func (w *windowWithLogs) Resize(width, height float64) {
	WhenDebug(func() {
		Logf("window %s is resizing to width:%.2f height:%.2f",
			w.ID(),
			width,
			height,
		)
	})

	w.Window.Resize(width, height)
}

func (w *windowWithLogs) Focus() {
	WhenDebug(func() {
		Logf("window %s is getting focus", w.ID())
	})

	w.Window.Focus()
}

func (w *windowWithLogs) FullScreen() {
	WhenDebug(func() {
		Logf("window %s is entering full screen", w.ID())
	})

	w.Window.FullScreen()
}

func (w *windowWithLogs) ExitFullScreen() {
	WhenDebug(func() {
		Logf("window %s is exiting full screen", w.ID())
	})

	w.Window.ExitFullScreen()
}

func (w *windowWithLogs) Minimize() {
	WhenDebug(func() {
		Logf("window %s is minimizing", w.ID())
	})

	w.Window.Minimize()
}

func (w *windowWithLogs) Deminimize() {
	WhenDebug(func() {
		Logf("window %s is deminimizing", w.ID())
	})

	w.Window.Deminimize()
}

func (w *windowWithLogs) Close() {
	WhenDebug(func() {
		Logf("window %s is closing", w.ID())
	})

	w.Window.Close()
	if w.Err() != nil {
		Logf("window %s failed to close: %s",
			w.ID(),
			w.Err(),
		)
	}
}

// Page logs.
type pageWithLogs struct {
	Page
}

func (p *pageWithLogs) WhenPage(f func(Page)) {
	f(p)
}

func (p *pageWithLogs) WhenNavigator(f func(Navigator)) {
	f(p)
}

func (p *pageWithLogs) Load(url string, v ...interface{}) {
	parsedURL := fmt.Sprintf(url, v...)

	WhenDebug(func() {
		Logf("page %s is loading %s",
			p.ID(),
			parsedURL,
		)
	})

	p.Page.Load(url, v...)
	if p.Err() != nil {
		Logf("page %s failed to load %s: %s",
			p.ID(),
			parsedURL,
			p.Err(),
		)
	}
}

func (p *pageWithLogs) Render(c Compo) {
	WhenDebug(func() {
		Logf("page %s is rendering %T",
			p.ID(),
			c,
		)
	})

	p.Page.Render(c)
	if p.Err() != nil {
		Logf("page %s failed to render %T: %s",
			p.ID(),
			c,
			p.Err(),
		)
	}
}

func (p *pageWithLogs) Reload() {
	WhenDebug(func() {
		Logf("page %s is reloading", p.ID())
	})

	p.Page.Reload()
	if p.Err() != nil {
		Logf("page %s failed to reload: %s",
			p.ID(),
			p.Err(),
		)
	}
}

func (p *pageWithLogs) Previous() {
	WhenDebug(func() {
		Logf("page %s is loading previous", p.ID())
	})

	p.Page.Previous()
}

func (p *pageWithLogs) Next() {
	WhenDebug(func() {
		Logf("page %s is loading next", p.ID())
	})

	p.Page.Next()
}

// Menu logs.
type menuWithLogs struct {
	Menu
}

func (m *menuWithLogs) Load(url string, v ...interface{}) {
	parsedURL := fmt.Sprintf(url, v...)

	WhenDebug(func() {
		Logf("%s %s is loading %s",
			m.Type(),
			m.ID(),
			parsedURL,
		)
	})

	m.Menu.Load(url, v...)
	if m.Err() != nil {
		Logf("%s %s failed to load %s: %s",
			m.Type(),
			m.ID(),
			parsedURL,
			m.Err(),
		)
	}
}

func (m *menuWithLogs) Render(c Compo) {
	WhenDebug(func() {
		Logf("%s %s is rendering %T",
			m.Type(),
			m.ID(),
			c,
		)
	})

	m.Menu.Render(c)
	if m.Err() != nil {
		Logf("%s %s failed to render %T: %s",
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

func (d *dockWithLogs) WhenDockTile(f func(DockTile)) {
	f(d)
}

func (d *dockWithLogs) Load(url string, v ...interface{}) {
	parsedURL := fmt.Sprintf(url, v...)

	WhenDebug(func() {
		Logf("dock tile is loading %s", parsedURL)
	})

	d.DockTile.Load(url, v...)
	if d.Err() != nil {
		Logf("dock tile failed to load %s: %s",
			parsedURL,
			d.Err(),
		)
	}
}

func (d *dockWithLogs) Render(c Compo) {
	WhenDebug(func() {
		Logf("dock tile is rendering %T", c)
	})

	d.DockTile.Render(c)
	if d.Err() != nil {
		Logf("dock tile failed to render %T: %s",
			c,
			d.Err(),
		)
	}
}

func (d *dockWithLogs) SetIcon(name string) {
	WhenDebug(func() {
		Logf("dock tile is setting its icon to %s", name)
	})

	d.DockTile.SetIcon(name)
	if d.Err() != nil {
		Logf("dock tile failed to set its icon: %s", d.Err())
	}
}

func (d *dockWithLogs) SetBadge(v interface{}) {
	WhenDebug(func() {
		Logf("dock tile is setting its badge to %v", v)
	})

	d.DockTile.SetBadge(v)
	if d.Err() != nil {
		Logf("dock tile failed to set its badge: %s", d.Err())
	}
}

// Status menu logs.
type statusMenuWithLogs struct {
	StatusMenu
}

func (s *statusMenuWithLogs) WhenStatusMenu(f func(StatusMenu)) {
	f(s)
}

func (s *statusMenuWithLogs) Load(url string, v ...interface{}) {
	parsedURL := fmt.Sprintf(url, v...)

	WhenDebug(func() {
		Logf("status menu %s is loading %s",
			s.ID(),
			parsedURL,
		)
	})

	s.StatusMenu.Load(url, v...)
	if s.Err() != nil {
		Logf("status menu %T failed to load %s: %s",
			s.ID(),
			parsedURL,
			s.Err(),
		)
	}
}

func (s *statusMenuWithLogs) Render(c Compo) {
	WhenDebug(func() {
		Logf("status menu %s is rendering %T",
			s.ID(),
			c,
		)
	})

	s.StatusMenu.Render(c)
	if s.Err() != nil {
		Logf("status menu %s failed to render %T: %s",
			s.ID(),
			c,
			s.Err(),
		)
	}
}

func (s *statusMenuWithLogs) SetIcon(name string) {
	WhenDebug(func() {
		Logf("status menu %s is setting icon to %s",
			s.ID(),
			name,
		)
	})

	s.StatusMenu.SetIcon(name)
	if s.Err() != nil {
		Logf("status menu %s failed to set icon: %s",
			s.ID(),
			s.Err(),
		)
	}
}

func (s *statusMenuWithLogs) SetText(text string) {
	WhenDebug(func() {
		Logf("status menu %s is setting text to %s",
			s.ID(),
			text,
		)
	})

	s.StatusMenu.SetText(text)
	if s.Err() != nil {
		Logf("status menu %s failed to set text: %s",
			s.ID(),
			s.Err(),
		)
	}
}

func (s *statusMenuWithLogs) Close() {
	WhenDebug(func() {
		Logf("status menu %s is closing", s.ID())
	})

	s.StatusMenu.Close()
	if s.Err() != nil {
		Logf("status menu %s failed to close: %s",
			s.ID(),
			s.Err(),
		)
	}
}
