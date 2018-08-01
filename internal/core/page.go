package core

import (
	"net/url"

	"github.com/murlokswarm/app"
)

// Page is a base struct to embed in app.Page implementations.
type Page struct {
	Elem
}

// WhenPage satisfies the app.Page interface.
func (p *Page) WhenPage(f func(app.Page)) {
	f(p)
}

// WhenNavigator satisfies the app.Page interface.
func (p *Page) WhenNavigator(f func(app.Navigator)) {
	f(p)
}

// Load satisfies the app.Page interface.
func (p *Page) Load(url string, v ...interface{}) {
	p.SetErr(app.ErrNotSupported)
}

// Compo satisfies the app.Page interface.
func (p *Page) Compo() app.Compo {
	return nil
}

// Contains satisfies the app.Page interface.
func (p *Page) Contains(c app.Compo) bool {
	return false
}

// Render satisfies the app.Page interface.
func (p *Page) Render(c app.Compo) {
	p.SetErr(app.ErrNotSupported)
}

// Reload satisfies the app.Page interface.
func (p *Page) Reload() {
	p.SetErr(app.ErrNotSupported)
}

// CanPrevious satisfies the app.Page interface.
func (p *Page) CanPrevious() bool {
	return false
}

// Previous satisfies the app.Page interface.
func (p *Page) Previous() {
	p.SetErr(app.ErrNotSupported)
}

// CanNext satisfies the app.Page interface.
func (p *Page) CanNext() bool {
	return false
}

// Next satisfies the app.Page interface.
func (p *Page) Next() {
	p.SetErr(app.ErrNotSupported)
}

// URL satisfies the app.Page interface.
func (p *Page) URL() *url.URL {
	return nil
}

// Referer satisfies the app.Page interface.
func (p *Page) Referer() *url.URL {
	return nil
}

// Close satisfies the app.Page interface.
func (p *Page) Close() {
	p.SetErr(app.ErrNotSupported)
}
