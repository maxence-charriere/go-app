// +build js

// Package web is the driver to be used for web applications.
package web

import (
	"context"
	"os"
	"path"

	"github.com/gopherjs/gopherjs/js"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/core"
	"github.com/pkg/errors"
)

var (
	driver *Driver
)

func init() {
	app.Loggers = []app.Logger{
		app.NewLogger(os.Stdout, os.Stderr, true, false),
	}
}

// Run satisfies the app.Driver interface.
func (d *Driver) Run(f app.Factory) error {
	d.factory = f
	d.elems = core.NewElemDB()
	d.uichan = make(chan func(), 255)
	driver = d

	var ctx context.Context
	ctx, d.cancel = context.WithCancel(context.Background())

	go d.runLoop(ctx)

	var err error
	page, err := newPage(app.PageConfig{})
	if err != nil {
		return err
	}
	d.page = page

	return nil
}

func (d *Driver) runLoop(ctx context.Context) {
	for {
		select {
		case f := <-d.uichan:
			f()

		case <-ctx.Done():
			return
		}
	}
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

func (d *Driver) NewPage(c app.PageConfig) error {
	js.Global.Get("location").Set("href", c.DefaultURL)
	return nil
}

// Render satisfies the app.Driver interface.
func (d *Driver) Render(c app.Component) error {
	e := d.elems.GetByCompo(c)
	if e.IsNotSet() {
		return errors.New("element not set")
	}

	return e.Render(c)
}

// ElemByCompo satisfies the app.Driver interface.
func (d *Driver) ElemByCompo(c app.Component) app.Elem {
	return d.elems.GetByCompo(c)
}

// CallOnUIGoroutine satisfies the app.Driver interface.
func (d *Driver) CallOnUIGoroutine(f func()) {
	d.uichan <- f
}
