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
	id        uuid.UUID
	factory   app.Factory
	markup    app.Markup
	history   app.History
	lastFocus time.Time
	component app.Component
	x         float64
	y         float64
	width     float64
	height    float64

	onClose func()
}

func newWindow(d *Driver, c app.WindowConfig) (app.Window, error) {
	var markup app.Markup = html.NewMarkup(d.factory)
	markup = app.ConcurrentMarkup(markup)

	history := app.NewHistory()
	history = app.ConcurrentHistory(history)

	rawWin := &Window{
		id:        uuid.New(),
		factory:   d.factory,
		markup:    markup,
		history:   history,
		lastFocus: time.Now(),
	}

	win := app.WindowWithLogs(rawWin)

	d.elements.Add(win)
	rawWin.onClose = func() {
		d.elements.Remove(win)
	}

	var err error
	if len(c.DefaultURL) != 0 {
		err = win.Load(c.DefaultURL)
	}
	return win, err
}

// ID satisfies the app.Window interface.
func (w *Window) ID() uuid.UUID {
	return w.id
}

// Base satisfies the app.Window interface.
func (w *Window) Base() app.Window {
	return w
}

// Component satisfies the app.Window interface.
func (w *Window) Component() app.Component {
	return w.component
}

// Contains satisfies the app.Window interface.
func (w *Window) Contains(c app.Component) bool {
	return w.markup.Contains(c)
}

// Load satisfies the app.Window interface.
func (w *Window) Load(rawurl string, v ...interface{}) error {
	rawurl = fmt.Sprintf(rawurl, v...)
	u, err := url.Parse(rawurl)
	if err != nil {
		return err
	}

	var currentURL string
	if currentURL, err = w.history.Current(); err != nil || currentURL != u.String() {
		w.history.NewEntry(u.String())
	}
	return w.load(u)
}

func (w *Window) load(u *url.URL) error {
	if w.component != nil {
		w.markup.Dismount(w.component)
	}

	compo, err := w.factory.New(app.ComponentNameFromURL(u))
	if err != nil {
		return err
	}

	if _, err = w.markup.Mount(compo); err != nil {
		return errors.Wrapf(err, "loading %s in test window %p failed", u, w)
	}

	w.component = compo
	return nil
}

// Render satisfies the app.Window interface.
func (w *Window) Render(compo app.Component) error {
	_, err := w.markup.Update(compo)
	return err
}

// Reload satisfies the app.Window interface.
func (w *Window) Reload() error {
	rawurl, err := w.history.Current()
	if err != nil {
		return err
	}

	u, err := url.Parse(rawurl)
	if err != nil {
		return err
	}
	return w.load(u)
}

// LastFocus satisfies the app.Window interface.
func (w *Window) LastFocus() time.Time {
	return w.lastFocus
}

// CanPrevious satisfies the app.Window interface.
func (w *Window) CanPrevious() bool {
	return w.history.CanPrevious()
}

// Previous satisfies the app.Window interface.
func (w *Window) Previous() error {
	rawurl, err := w.history.Previous()
	if err != nil {
		return err
	}

	u, err := url.Parse(rawurl)
	if err != nil {
		return err
	}
	return w.load(u)
}

// CanNext satisfies the app.Window interface.
func (w *Window) CanNext() bool {
	return w.history.CanNext()
}

// Next satisfies the app.Window interface.
func (w *Window) Next() error {
	rawurl, err := w.history.Next()
	if err != nil {
		return err
	}

	u, err := url.Parse(rawurl)
	if err != nil {
		return err
	}
	return w.load(u)
}

// Position satisfies the app.Window interface.
func (w *Window) Position() (x, y float64) {
	return w.x, w.y
}

// Move satisfies the app.Window interface.
func (w *Window) Move(x, y float64) {
	w.x = x
	w.y = y
}

// Center satisfies the app.Window interface.
func (w *Window) Center() {
	w.Move(500, 500)
}

// Size satisfies the app.Window interface.
func (w *Window) Size() (width, height float64) {
	return w.width, w.height
}

// Resize satisfies the app.Window interface.
func (w *Window) Resize(width, height float64) {
	w.width = width
	w.height = height
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
