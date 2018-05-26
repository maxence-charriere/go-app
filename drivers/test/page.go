package test

import (
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/html"
	"github.com/pkg/errors"
)

type page struct {
	id        uuid.UUID
	factory   app.Factory
	markup    app.Markup
	history   app.History
	lastFocus time.Time
	component app.Component

	onClose func()
}

func newPage(d *Driver, c app.PageConfig) (app.Page, error) {
	var markup app.Markup = html.NewMarkup(d.factory)
	markup = app.ConcurrentMarkup(markup)

	history := app.NewHistory()
	history = app.ConcurrentHistory(history)

	page := &page{
		id:        uuid.New(),
		factory:   d.factory,
		markup:    markup,
		history:   history,
		lastFocus: time.Now(),
	}

	d.elements.Add(page)
	page.onClose = func() {
		d.elements.Remove(page)
	}

	var err error
	if len(c.DefaultURL) != 0 {
		err = page.Load(c.DefaultURL)
	}
	return page, err
}

// ID satisfies the app.Page interface.
func (p *page) ID() uuid.UUID {
	return p.id
}

// Component satisfies the app.Page interface.
func (p *page) Component() app.Component {
	return p.component
}

// Contains satisfies the app.Page interface.
func (p *page) Contains(c app.Component) bool {
	return p.markup.Contains(c)
}

// Load satisfies the app.Page interface.
func (p *page) Load(rawurl string, v ...interface{}) error {
	rawurl = fmt.Sprintf(rawurl, v...)
	u, err := url.Parse(rawurl)
	if err != nil {
		return err
	}

	var currentURL string
	if currentURL, err = p.history.Current(); err != nil || currentURL != u.String() {
		p.history.NewEntry(u.String())
	}
	return p.load(u)
}

func (p *page) load(u *url.URL) error {
	if p.component != nil {
		p.markup.Dismount(p.component)
	}

	compo, err := p.factory.New(app.ComponentNameFromURL(u))
	if err != nil {
		return err
	}

	if _, err = p.markup.Mount(compo); err != nil {
		return errors.Wrapf(err, "loading %s in test page %p failed", u, p)
	}

	p.component = compo
	return nil
}

// Render satisfies the app.Page interface.
func (p *page) Render(compo app.Component) error {
	_, err := p.markup.Update(compo)
	return err
}

// Reload satisfies the app.Page interface.
func (p *page) Reload() error {
	rawurl, err := p.history.Current()
	if err != nil {
		return err
	}

	u, err := url.Parse(rawurl)
	if err != nil {
		return err
	}
	return p.load(u)
}

// LastFocus satisfies the app.Page interface.
func (p *page) LastFocus() time.Time {
	return p.lastFocus
}

// CanPrevious satisfies the app.Page interface.
func (p *page) CanPrevious() bool {
	return p.history.CanPrevious()
}

// Previous satisfies the app.Page interface.
func (p *page) Previous() error {
	rawurl, err := p.history.Previous()
	if err != nil {
		return err
	}

	u, err := url.Parse(rawurl)
	if err != nil {
		return err
	}
	return p.load(u)
}

// CanNext satisfies the app.Page interface.
func (p *page) CanNext() bool {
	return p.history.CanNext()
}

// Next satisfies the app.Page interface.
func (p *page) Next() error {
	rawurl, err := p.history.Next()
	if err != nil {
		return err
	}

	u, err := url.Parse(rawurl)
	if err != nil {
		return err
	}
	return p.load(u)
}

func (p *page) URL() *url.URL {
	rawurl, _ := p.history.Current()
	u, _ := url.Parse(rawurl)
	return u
}

func (p *page) Referer() *url.URL {
	rawurl, err := p.history.Previous()
	if err != nil {
		return nil
	}
	u, _ := url.Parse(rawurl)

	p.history.Next()
	return u
}

func (p *page) Close() error {
	p.onClose()
	return nil
}
