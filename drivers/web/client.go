// +build js

// Package web is the driver to be used for web applications.
package web

import (
	"net/url"
	"path"

	"github.com/gopherjs/gopherjs/js"
	"github.com/murlokswarm/app"
)

// Name satisfies the app.Driver interface.
func (d *Driver) Name() string {
	return "Web"
}

// Run satisfies the app.Driver interface.
func (d *Driver) Run(f app.Factory) error {
	elements := app.NewElemDB()
	elements = app.NewConcurrentElemDB(elements)
	d.elements = elements

	d.uichan = make(chan func(), 4096)
	defer close(d.uichan)

	rawurl := js.Global.Get("location").Get("href").String()

	u, err := url.Parse(rawurl)
	if err != nil {
		return err
	}

	if len(u.Path) == 0 || u.Path == "/" {
		u.Path = d.DefaultURL
	}

	return nil
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

func (d *Driver) NewPage(c app.PageConfig) (app.Page, error) {
	panic("not implemented")
}

// Render satisfies the app.Driver interface.
func (d *Driver) Render(c app.Component) error {
	panic("not implemented")
}

// ElementByComponent satisfies the app.Driver interface.
func (d *Driver) ElementByComponent(c app.Component) (app.ElementWithComponent, error) {
	panic("not implemented")
}

// CallOnUIGoroutine satisfies the app.Driver interface.
func (d *Driver) CallOnUIGoroutine(f func()) {
	panic("not implemented")
}
