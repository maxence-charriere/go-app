package test

import (
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/core"
	"github.com/murlokswarm/app/internal/dom"
	"github.com/pkg/errors"
)

// Page is a test page that implements the app.Page interface.
type Page struct {
	core.Page

	driver  *Driver
	dom     *dom.DOM
	history *core.History
	id      string
	compo   app.Compo
}

func newPage(d *Driver, c app.PageConfig) *Page {
	p := &Page{
		driver:  d,
		dom:     dom.NewDOM(d.factory, dom.JsToGoHandler),
		history: core.NewHistory(),
		id:      uuid.New().String(),
	}

	d.elems.Put(p)

	if len(c.URL) != 0 {
		p.Load(c.URL)
	}

	return p
}

// ID satisfies the app.Page interface.
func (p *Page) ID() string {
	return p.id
}

// Load satisfies the app.Page interface.
func (p *Page) Load(urlFmt string, v ...interface{}) {
	var err error
	defer func() {
		p.SetErr(err)
	}()

	u := fmt.Sprintf(urlFmt, v...)
	n := core.CompoNameFromURLString(u)

	var c app.Compo
	if c, err = p.driver.factory.NewCompo(n); err != nil {
		return
	}

	if p.compo != nil {
		p.dom.Clean()
	}

	p.compo = c

	if u != p.history.Current() {
		p.history.NewEntry(u)
	}

	_, err = p.dom.New(c)
}

// Compo satisfies the app.Page interface.
func (p *Page) Compo() app.Compo {
	return p.compo
}

// Contains satisfies the app.Page interface.
func (p *Page) Contains(c app.Compo) bool {
	return p.dom.Contains(c)
}

// Render satisfies the app.Page interface.
func (p *Page) Render(c app.Compo) {
	_, err := p.dom.Update(c)
	p.SetErr(err)
}

// Reload satisfies the app.Page interface.
func (p *Page) Reload() {
	u := p.history.Current()

	if len(u) == 0 {
		p.SetErr(errors.New("no component loaded"))
		return
	}

	p.Load(u)
}

// CanPrevious satisfies the app.Page interface.
func (p *Page) CanPrevious() bool {
	return p.history.CanPrevious()
}

// Previous satisfies the app.Page interface.
func (p *Page) Previous() {
	u := p.history.Previous()

	if len(u) == 0 {
		p.SetErr(nil)
		return
	}

	p.Load(u)
}

// CanNext satisfies the app.Page interface.
func (p *Page) CanNext() bool {
	return p.history.CanNext()
}

// Next satisfies the app.Page interface.
func (p *Page) Next() {
	u := p.history.Next()

	if len(u) == 0 {
		p.SetErr(nil)
		return
	}

	p.Load(u)
}

// URL satisfies the app.Page interface.
func (p *Page) URL() *url.URL {
	u, err := url.Parse(p.history.Current())
	p.SetErr(err)
	return u
}

// Referer satisfies the app.Page interface.
func (p *Page) Referer() *url.URL {
	p.SetErr(nil)
	return nil
}

// Close satisfies the app.Page interface.
func (p *Page) Close() {
	p.driver.elems.Delete(p)
	p.SetErr(nil)
}
