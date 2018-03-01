// +build js

// Package web is the driver to be used for web applications.
package web

import (
	"context"
	"path"

	"github.com/murlokswarm/app"
)

var (
	driver *Driver
)

// Name satisfies the app.Driver interface.
func (d *Driver) Name() string {
	return "Web"
}

// Run satisfies the app.Driver interface.
func (d *Driver) Run(f app.Factory) error {
	d.factory = f
	elements := app.NewElemDB()
	elements = app.NewConcurrentElemDB(elements)
	d.elements = elements
	driver = d

	var err error
	page, err := newPage(app.PageConfig{})
	if err != nil {
		return err
	}
	d.page = page

	d.uichan = make(chan func(), 255)

	var ctx context.Context
	ctx, d.cancel = context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case f := <-d.uichan:
				f()

			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

// AppName satisfies the app.Driver interface.
func (d *Driver) AppName() string {
	return "goapp"
}

// Resources satisfies the app.Driver interface.
func (d *Driver) Resources(p ...string) string {
	resources := path.Join(p...)
	resources = path.Join("resources", resources)
	return resources
}

// Storage satisfies the app.Driver interface.
func (d *Driver) Storage(p ...string) string {
	return ""
}

// Render satisfies the app.Driver interface.
func (d *Driver) Render(c app.Component) error {
	elem, err := d.ElementByComponent(c)
	if err != nil {
		return err
	}
	return elem.Render(c)
}

// ElementByComponent satisfies the app.Driver interface.
func (d *Driver) ElementByComponent(c app.Component) (app.ElementWithComponent, error) {
	return d.ElementByComponent(c)
}

// CallOnUIGoroutine satisfies the app.Driver interface.
func (d *Driver) CallOnUIGoroutine(f func()) {
	d.uichan <- f
}
