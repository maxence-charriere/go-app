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

// A Window implementation for tests.
type Window struct {
	driver    *Driver
	id        uuid.UUID
	factory   app.Factory
	markup    app.Markup
	history   app.History
	lastFocus time.Time

	onLoad  func(compo app.Component)
	onClose func()
}

func newWindow(d *Driver, c app.WindowConfig) (*Window, error) {
	win := &Window{
		driver:    d,
		id:        uuid.New(),
		factory:   d.factory,
		markup:    html.NewMarkup(d.factory),
		history:   app.NewHistory(),
		lastFocus: time.Now(),
	}
	d.elements.Add(win)

	win.onClose = func() {
		d.elements.Remove(win)
	}

	var err error
	if len(c.DefaultURL) != 0 {
		err = win.Load(c.DefaultURL)
	}
	return win, err
}

// ID satisfies the app.Element interface.
func (w *Window) ID() uuid.UUID {
	return w.id
}

// Contains satisfies the app.ElementWithComponent interface.
func (w *Window) Contains(c app.Component) bool {
	return w.markup.Contains(c)
}

// Load satisfies the app.ElementWithComponent interface.
func (w *Window) Load(rawurl string, v ...interface{}) error {
	rawurl = fmt.Sprintf(rawurl, v...)

	u, err := url.Parse(rawurl)
	if err != nil {
		return err
	}

	w.history.NewEntry(u.String())

	compo, err := w.factory.NewComponent(app.ComponentNameFromURL(u))
	if err != nil {
		return err
	}

	if _, err = w.markup.Mount(compo); err != nil {
		return errors.Wrapf(err, "loading %s in test window %p failed", u, w)
	}

	if w.onLoad != nil {
		w.onLoad(compo)
	}
	return nil
}

// Render satisfies the app.ElementWithComponent interface.
func (w *Window) Render(compo app.Component) error {
	_, err := w.markup.Update(compo)
	return err
}

// Reload satisfies the app.ElementWithNavigation interface.
func (w *Window) Reload() error {
	rawurl, err := w.history.Current()
	if err != nil {
		return err
	}
	return w.Load(rawurl)
}

// LastFocus satisfies the app.ElementWithComponent interface.
func (w *Window) LastFocus() time.Time {
	return w.lastFocus
}

// CanPrevious satisfies the app.ElementWithNavigation interface.
func (w *Window) CanPrevious() bool {
	return w.history.CanPrevious()
}

// Previous satisfies the app.ElementWithNavigation interface.
func (w *Window) Previous() error {
	rawurl, err := w.history.Previous()
	if err != nil {
		return err
	}
	return w.Load(rawurl)
}

// CanNext satisfies the app.ElementWithNavigation interface.
func (w *Window) CanNext() bool {
	return w.history.CanNext()
}

// Next satisfies the app.ElementWithNavigation interface.
func (w *Window) Next() error {
	rawurl, err := w.history.Next()
	if err != nil {
		return err
	}
	return w.Load(rawurl)
}

// Position satisfies the app.Window interface.
func (w *Window) Position() (x, y float64) {
	return
}

// Move satisfies the app.Window interface.
func (w *Window) Move(x, y float64) {
}

// Center satisfies the app.Window interface.
func (w *Window) Center() {
}

// Size satisfies the app.Window interface.
func (w *Window) Size() (width, height float64) {
	return
}

// Resize satisfies the app.Window interface.
func (w *Window) Resize(width, height float64) {
}

// Focus satisfies the app.Window interface.
func (w *Window) Focus() {
	w.lastFocus = time.Now()
}

// ToggleFullScreen satisfies the app.Window interface.
func (w *Window) ToggleFullScreen() {
}

// ToggleMinimize satisfies the app.Window interface.
func (w *Window) ToggleMinimize() {
}

// Close satisfies the app.Window interface.
func (w *Window) Close() {
	w.onClose()
}
