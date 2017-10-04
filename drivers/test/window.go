package test

import (
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/markup"
	"github.com/pkg/errors"
)

// A Window implementation for tests.
type Window struct {
	driver       *Driver
	config       app.WindowConfig
	id           uuid.UUID
	compoBuilder markup.CompoBuilder
	env          markup.Env
	lastFocus    time.Time

	onLoad  func(c markup.Component)
	onClose func()
}

// NewWindow creates a new widnow.
func NewWindow(d *Driver, c app.WindowConfig) *Window {
	window := &Window{
		driver:       d,
		config:       c,
		id:           uuid.New(),
		compoBuilder: d.compoBuilder,
		env:          markup.NewEnv(d.compoBuilder),
		lastFocus:    time.Now(),
	}

	d.elements.Add(window)
	window.onClose = func() {
		d.elements.Remove(window)
	}

	if d.OnWindowLoad != nil {
		window.onLoad = func(c markup.Component) {
			d.OnWindowLoad(window, c)
		}
	}

	if len(c.DefaultURL) != 0 {
		if err := window.Load(c.DefaultURL); err != nil {
			d.Test.Log(err)
		}
	}
	return window
}

// ID satisfies the app.Element interface.
func (w *Window) ID() uuid.UUID {
	return w.id
}

// Contains satisfies the app.ElementWithComponent interface.
func (w *Window) Contains(c markup.Component) bool {
	return w.env.Contains(c)
}

// Load satisfies the app.ElementWithComponent interface.
func (w *Window) Load(rawurl string) error {
	u, err := url.Parse(rawurl)
	if err != nil {
		return err
	}

	componame, ok := markup.ComponentNameFromURL(u)
	if !ok {
		return nil
	}

	compo, err := w.compoBuilder.New(componame)
	if err != nil {
		return err
	}

	if _, err = w.env.Mount(compo); err != nil {
		return errors.Wrapf(err, "loading %s in test window %p failed", u, w)
	}

	if w.onLoad != nil {
		w.onLoad(compo)
	}
	return nil
}

// Render satisfies the app.ElementWithComponent interface.
func (w *Window) Render(c markup.Component) error {
	_, err := w.env.Update(c)
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
