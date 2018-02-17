package test

import (
	"context"
	"path/filepath"

	"github.com/murlokswarm/app"
	"github.com/pkg/errors"
)

var (
	// ErrSimulated is an error that is set to simulated a return error
	// behavior.
	ErrSimulated = errors.New("simulated error")
)

// Driver is an app.Driver implementation for testing.
type Driver struct {
	SimulateErr bool
	Ctx         context.Context

	OnRun func()

	factory  app.Factory
	elements app.ElementDB
	dock     app.DockTile
	menubar  app.Menu
	uichan   chan func()
}

// Name satisfies the app.Driver interface.
func (d *Driver) Name() string {
	return "Test"
}

// Run satisfies the app.Driver interface.
func (d *Driver) Run(f app.Factory) error {
	if d.SimulateErr {
		return ErrSimulated
	}

	d.factory = f

	elements := app.NewElemDB()
	elements = app.NewConcurrentElemDB(elements)
	d.elements = elements

	d.dock = newDockTile(d)

	menubar, err := newMenu(d, "menu bar", app.MenuConfig{})
	if err != nil {
		return err
	}
	d.menubar = menubar

	d.uichan = make(chan func(), 256)

	if d.OnRun != nil {
		d.uichan <- d.OnRun
	}

	if d.Ctx == nil {
		return errors.New("driver.Ctx is nil")
	}

	for {
		select {
		case f := <-d.uichan:
			f()

		case <-d.Ctx.Done():
			return nil
		}
	}
}

// AppName satisfies the app.Driver interface.
func (d *Driver) AppName() string {
	return "Driver unit tests"
}

// Resources satisfies the app.Driver interface.
func (d *Driver) Resources(path ...string) string {
	resources := filepath.Join(path...)
	return filepath.Join("resources", resources)
}

// Storage satisfies the app.Driver interface.
func (d *Driver) Storage(path ...string) string {
	storage := filepath.Join(path...)
	return filepath.Join("storage", storage)
}

// NewWindow satisfies the app.Driver interface.
func (d *Driver) NewWindow(c app.WindowConfig) (app.Window, error) {
	if d.SimulateErr {
		return nil, ErrSimulated
	}
	return newWindow(d, c)
}

// NewContextMenu satisfies the app.Driver interface.
func (d *Driver) NewContextMenu(c app.MenuConfig) (app.Menu, error) {
	if d.SimulateErr {
		return nil, ErrSimulated
	}
	return newMenu(d, "context menu", c)
}

// Render satisfies the app.Driver interface.
func (d *Driver) Render(compo app.Component) error {
	if d.SimulateErr {
		return ErrSimulated
	}

	elem, err := d.elements.ElementByComponent(compo)
	if err != nil {
		return err
	}
	return elem.Render(compo)
}

// ElementByComponent satisfies the app.Driver interface.
func (d *Driver) ElementByComponent(c app.Component) (app.ElementWithComponent, error) {
	return d.elements.ElementByComponent(c)
}

// NewFilePanel satisfies the app.Driver interface.
func (d *Driver) NewFilePanel(c app.FilePanelConfig) error {
	return app.NewErrNotSupported("file panels")
}

// NewShare satisfies the app.Driver interface.
func (d *Driver) NewShare(v interface{}) error {
	return app.NewErrNotSupported("share")
}

// NewNotification satisfies the app.Driver interface.
func (d *Driver) NewNotification(c app.NotificationConfig) error {
	return app.NewErrNotSupported("notifications")
}

// MenuBar satisfies the app.Driver interface.
func (d *Driver) MenuBar() app.Menu {
	return d.menubar
}

// Dock satisfies the app.Driver interface.
func (d *Driver) Dock() app.DockTile {
	return d.dock
}

// CallOnUIGoroutine satisfies the app.Driver interface.
func (d *Driver) CallOnUIGoroutine(f func()) {
	d.uichan <- f
}
