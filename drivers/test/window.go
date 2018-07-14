package test

import (
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/html"
	"github.com/murlokswarm/app/internal/core"
	"github.com/pkg/errors"
)

// A Window implementation for tests.
type Window struct {
	core.Elem

	id          uuid.UUID
	factory     app.Factory
	markup      app.Markup
	history     *core.History
	lastFocus   time.Time
	component   app.Compo
	x           float64
	y           float64
	width       float64
	height      float64
	simulateErr bool

	onClose func() error
}

func newWindow(d *Driver, c app.WindowConfig) (app.Window, error) {
	var markup app.Markup = html.NewMarkup(d.factory)
	markup = app.ConcurrentMarkup(markup)

	win := &Window{
		id:          uuid.New(),
		factory:     d.factory,
		markup:      markup,
		history:     core.NewHistory(),
		lastFocus:   time.Now(),
		simulateErr: d.SimulateElemErr,
	}

	d.elems.Put(win)
	win.onClose = func() error {
		d.elems.Delete(win)
		return nil
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

// Compo satisfies the app.Window interface.
func (w *Window) Compo() app.Compo {
	return w.component
}

// Contains satisfies the app.Window interface.
func (w *Window) Contains(c app.Compo) bool {
	return w.markup.Contains(c)
}

// Load satisfies the app.Window interface.
func (w *Window) Load(rawurl string, v ...interface{}) error {
	if w.simulateErr {
		return ErrSimulated
	}

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
	if w.simulateErr {
		return ErrSimulated
	}

	if w.component != nil {
		w.markup.Dismount(w.component)
	}

	compo, err := w.factory.New(app.CompoNameFromURL(u))
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
func (w *Window) Render(compo app.Compo) error {
	if w.simulateErr {
		return ErrSimulated
	}

	_, err := w.markup.Update(compo)
	return err
}

// Reload satisfies the app.Window interface.
func (w *Window) Reload() error {
	if w.simulateErr {
		return ErrSimulated
	}

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
	if w.simulateErr {
		return ErrSimulated
	}

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
func (w *Window) Move(x, y float64) error {
	if w.simulateErr {
		return ErrSimulated
	}

	w.x = x
	w.y = y
	return nil
}

// Center satisfies the app.Window interface.
func (w *Window) Center() error {
	if w.simulateErr {
		return ErrSimulated
	}

	return w.Move(500, 500)
}

// Size satisfies the app.Window interface.
func (w *Window) Size() (width, height float64) {
	return w.width, w.height
}

// Resize satisfies the app.Window interface.
func (w *Window) Resize(width, height float64) error {
	if w.simulateErr {
		return ErrSimulated
	}

	w.width = width
	w.height = height
	return nil
}

// Focus satisfies the app.Window interface.
func (w *Window) Focus() error {
	if w.simulateErr {
		return ErrSimulated
	}

	w.lastFocus = time.Now()
	return nil
}

// ToggleFullScreen satisfies the app.Window interface.
func (w *Window) ToggleFullScreen() error {
	if w.simulateErr {
		return ErrSimulated
	}

	return nil
}

// ToggleMinimize satisfies the app.Window interface.
func (w *Window) ToggleMinimize() error {
	if w.simulateErr {
		return ErrSimulated
	}
	return nil
}

// Close satisfies the app.Window interface.
func (w *Window) Close() error {
	if w.simulateErr {
		return ErrSimulated
	}
	return w.onClose()
}
