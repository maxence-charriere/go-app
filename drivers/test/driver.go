package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/db"
	"github.com/murlokswarm/app/log"
	"github.com/murlokswarm/app/markup"
	"github.com/pkg/errors"
)

// Driver is an app.Driver implementation for testing.
type Driver struct {
	Test         *testing.T
	compoBuilder markup.CompoBuilder
	elements     app.ElementDB
	menubar      app.Menu
	dock         app.DockTile
	logger       app.Logger
	RunSouldErr  bool
	UIchan       chan func()

	OnWindowLoad func(w app.Window, c markup.Component)
}

// Run satisfies the app.Driver interface.
func (d *Driver) Run(b markup.CompoBuilder) error {
	d.compoBuilder = b
	d.elements = db.NewElementDB(256)
	d.menubar = newMenu(d, app.MenuConfig{})
	d.dock = newDockTile(d)
	d.logger = &log.Logger{}
	d.UIchan = make(chan func(), 256)

	if d.RunSouldErr {
		return errors.New("simulating run error")
	}
	return nil
}

// Render satisfies the app.Driver interface.
func (d *Driver) Render(c markup.Component) error {
	elem, err := d.elements.ElementByComponent(c)
	if err != nil {
		return errors.Wrap(err, "rendering component")
	}
	return elem.Render(c)
}

// Context satisfies the app.Driver interface.
func (d *Driver) Context(c markup.Component) (e app.ElementWithComponent, err error) {
	if e, err = d.elements.ElementByComponent(c); err != nil {
		err = errors.Wrap(err, "can't get context")
	}
	return
}

// NewContextMenu satisfies the app.Driver interface.
func (d *Driver) NewContextMenu(c app.MenuConfig) app.Menu {
	return newMenu(d, c)
}

// Resources satisfies the app.Driver interface.
func (d *Driver) Resources() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.Join(wd, "resources")
}

// Logs satisfies the app.Driver interface.
func (d *Driver) Logs() app.Logger {
	return d.logger
}

// CallOnUIGoroutine satisfies the app.Driver interface.
func (d *Driver) CallOnUIGoroutine(f func()) {
	d.UIchan <- f
}

// Storage satisfies the app.DriverWithStorage interface.
func (d *Driver) Storage() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.Join(wd, "storage")
}

// NewWindow satisfies the app.DriverWithWindows interface.
func (d *Driver) NewWindow(c app.WindowConfig) app.Window {
	return NewWindow(d, c)
}

// MenuBar satisfies the app.DriverWithMenuBar interface.
func (d *Driver) MenuBar() app.Menu {
	return d.menubar
}

// Dock satisfies the app.DriverWithDock interface.
func (d *Driver) Dock() app.DockTile {
	return d.dock
}

// Share satisfies the app.DriverWithShare interface.
func (d *Driver) Share(v interface{}) {
}

// NewFilePanel satisfies the app.DriverWithFilePanels interface.
func (d *Driver) NewFilePanel(c app.FilePanelConfig) app.Element {
	return NewElement(d)
}

// NewPopupNotification satisfies the app.DriverWithPopupNotifications
// interface.
func (d *Driver) NewPopupNotification(c app.PopupNotificationConfig) app.Element {
	return NewElement(d)
}
