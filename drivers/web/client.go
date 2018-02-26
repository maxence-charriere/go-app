// +build js

// Package web is the driver to be used for web applications.
package web

import (
	"path"

	"github.com/murlokswarm/app"
)

// Name satisfies the app.Driver interface.
func (d *Driver) Name() string {
	return "Web"
}

// Run satisfies the app.Driver interface.
func (d *Driver) Run(f app.Factory) error {
	panic("not implemented")
}

// AppName satisfies the app.Driver interface.
func (d *Driver) AppName() string {
	return "goapp"
}

// Resources satisfies the app.Driver interface.
func (d *Driver) Resources(p ...string) string {
	p = append([]string{"resources"}, p...)
	return path.Join(p...)
}

// Storage satisfies the app.Driver interface.
func (d *Driver) Storage(p ...string) string {
	return ""
}

// NewWindow satisfies the app.Driver interface.
func (d *Driver) NewWindow(c app.WindowConfig) (app.Window, error) {
	return nil, app.NewErrNotSupported("window")
}

// NewContextMenu satisfies the app.Driver interface.
func (d *Driver) NewContextMenu(c app.MenuConfig) (app.Menu, error) {
	return nil, app.NewErrNotSupported("context menu")
}

// Render satisfies the app.Driver interface.
func (d *Driver) Render(c app.Component) error {
	panic("not implemented")
}

// ElementByComponent satisfies the app.Driver interface.
func (d *Driver) ElementByComponent(c app.Component) (app.ElementWithComponent, error) {
	panic("not implemented")
}

// NewFilePanel satisfies the app.Driver interface.
func (d *Driver) NewFilePanel(c app.FilePanelConfig) error {
	return app.NewErrNotSupported("file panel")
}

// NewSaveFilePanel satisfies the app.Driver interface.
func (d *Driver) NewSaveFilePanel(c app.SaveFilePanelConfig) error {
	return app.NewErrNotSupported("save file panel")
}

// NewShare satisfies the app.Driver interface.
func (d *Driver) NewShare(v interface{}) error {
	return app.NewErrNotSupported("share")
}

// NewNotification satisfies the app.Driver interface.
func (d *Driver) NewNotification(c app.NotificationConfig) error {
	return app.NewErrNotSupported("notification")
}

// MenuBar satisfies the app.Driver interface.
func (d *Driver) MenuBar() app.Menu {
	panic("not implemented")
}

// Dock satisfies the app.Driver interface.
func (d *Driver) Dock() app.DockTile {
	panic("not implemented")
}

// CallOnUIGoroutine satisfies the app.Driver interface.
func (d *Driver) CallOnUIGoroutine(f func()) {
	panic("not implemented")
}
