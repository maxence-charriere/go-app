package test

import (
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
	config    app.WindowConfig
	id        uuid.UUID
	factory   app.Factory
	markup    app.Markup
	lastFocus time.Time

	onLoad  func(compo app.Component)
	onClose func()
}

// NewWindow creates a new widnow.
func NewWindow(driver *Driver, config app.WindowConfig) *Window {
	window := &Window{
		driver:    driver,
		config:    config,
		id:        uuid.New(),
		factory:   driver.factory,
		markup:    app.NewConcurrentMarkup(html.NewMarkup(driver.factory)),
		lastFocus: time.Now(),
	}

	driver.elements.Add(window)
	window.onClose = func() {
		driver.elements.Remove(window)
	}

	if driver.OnWindowLoad != nil {
		window.onLoad = func(compo app.Component) {
			driver.OnWindowLoad(window, compo)
		}
	}

	if len(config.DefaultURL) != 0 {
		if err := window.Load(config.DefaultURL); err != nil {
			driver.Test.Log(err)
		}
	}
	return window
}

// ID satisfies the app.Element interface.
func (w *Window) ID() uuid.UUID {
	return w.id
}

// Contains satisfies the app.ElementWithComponent interface.
func (w *Window) Contains(compo app.Component) bool {
	return w.markup.Contains(compo)
}

// Load satisfies the app.ElementWithComponent interface.
func (w *Window) Load(rawurl string) error {
	u, err := url.Parse(rawurl)
	if err != nil {
		return err
	}

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

// LastFocus satisfies the app.ElementWithComponent interface.
func (w *Window) LastFocus() time.Time {
	return w.lastFocus
}

// CanPrevious satisfies the app.ElementWithNavigation interface.
func (w *Window) CanPrevious() bool {
	return false
}

// Previous satisfies the app.ElementWithNavigation interface.
func (w *Window) Previous() error {
	return nil
}

// CanNext satisfies the app.ElementWithNavigation interface.
func (w *Window) CanNext() bool {
	return false
}

// Next satisfies the app.ElementWithNavigation interface.
func (w *Window) Next() error {
	return nil
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
