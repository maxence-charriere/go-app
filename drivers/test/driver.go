package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/murlokswarm/app"
	"github.com/pkg/errors"
)

// Driver is an app.Driver implementation for testing.
type Driver struct {
	Test        *testing.T
	factory     app.Factory
	elements    app.ElementDB
	menubar     app.Menu
	dock        app.DockTile
	RunSouldErr bool
	UIchan      chan func()

	OnWindowLoad func(win app.Window, compo app.Component)
}

// Run satisfies the app.Driver interface.
func (d *Driver) Run(factory app.Factory) error {
	d.factory = factory
	d.elements = app.NewConcurrentElemDB(app.NewElementDB())
	d.menubar = newMenu(d, app.MenuConfig{})
	d.dock = newDockTile(d)
	d.UIchan = make(chan func(), 256)

	if d.RunSouldErr {
		return errors.New("simulating run error")
	}
	return nil
}

// Render satisfies the app.Driver interface.
func (d *Driver) Render(compo app.Component) error {
	elem, err := d.elements.ElementByComponent(compo)
	if err != nil {
		return errors.Wrap(err, "rendering component")
	}
	return elem.Render(compo)
}

// Context satisfies the app.Driver interface.
func (d *Driver) Context(compo app.Component) (e app.ElementWithComponent, err error) {
	if e, err = d.elements.ElementByComponent(compo); err != nil {
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
