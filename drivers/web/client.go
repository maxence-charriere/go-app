// +build js

// Package web is the driver to be used for web applications.
package web

import (
	"context"
	"os"

	"github.com/gopherjs/gopherjs/js"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/core"
	"github.com/murlokswarm/app/internal/logs"
)

var (
	driver *Driver
)

func init() {
	logger := logs.ToWriter(os.Stderr)
	logger = logs.WithPrompt(logger)
	app.Logger = logger
}

// Run satisfies the app.Driver interface.
func (d *Driver) Run(f *app.Factory) error {
	d.factory = f
	d.elems = core.NewElemDB()
	d.uichan = make(chan func(), 256)
	driver = d

	go func() {
		defer close(d.uichan)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		d.stop = cancel

		for {
			select {
			case <-ctx.Done():
				return

			case fn := <-d.uichan:
				fn()
			}
		}
	}()

	p := newPage(app.PageConfig{})
	return p.Err()
}

// AppName satisfies the app.Driver interface.
func (d *Driver) AppName() string {
	return "go webapp"
}

// Storage satisfies the app.Driver interface.
func (d *Driver) Storage(p ...string) string {
	return ""
}

func (d *Driver) NewPage(c app.PageConfig) app.Page {
	js.Global.Get("location").Set("href", c.URL)
	return d.Driver.NewPage(c)
}

// Render satisfies the app.Driver interface.
func (d *Driver) Render(c app.Compo) {
	e := d.ElemByCompo(c)
	if e.Err() == nil {
		e.(app.ElemWithCompo).Render(c)
	}
}

// ElemByCompo satisfies the app.Driver interface.
func (d *Driver) ElemByCompo(c app.Compo) app.Elem {
	return d.elems.GetByCompo(c)
}

// CallOnUIGoroutine satisfies the app.Driver interface.
func (d *Driver) CallOnUIGoroutine(f func()) {
	d.uichan <- f
}
