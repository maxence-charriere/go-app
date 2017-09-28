package mac

import (
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/bridge"
	"github.com/murlokswarm/app/geom"
	"github.com/murlokswarm/app/markup"
	"github.com/pkg/errors"
)

type Window struct {
	driver    *Driver
	id        uuid.UUID
	env       markup.Env
	lastFocus time.Time
}

func newWindow(d *Driver, c app.WindowConfig) (w *Window, err error) {
	w = &Window{
		driver:    d,
		id:        uuid.New(),
		env:       markup.NewEnv(d.components),
		lastFocus: time.Now(),
	}

	rawurl := fmt.Sprintf("/window/new?id=%s", w.id)
	if _, err = d.macos.RequestWithAsyncResponse(rawurl, bridge.NewPayload(c)); err != nil {
		return
	}

	err = d.elements.Add(w)
	return
}

// ID satisfies the app.Element interface.
func (w *Window) ID() uuid.UUID {
	return w.id
}

// Load satisfies the app.ElementWithComponent interface.
func (w *Window) Load(url string) error {
	panic("not implemented")
}

// Contains satisfies the app.ElementWithComponent interface.
func (w *Window) Contains(c markup.Component) bool {
	panic("not implemented")
}

// Render satisfies the app.ElementWithComponent interface.
func (w *Window) Render(c markup.Component) error {
	panic("not implemented")
}

// LastFocus satisfies the app.ElementWithComponent interface.
func (w *Window) LastFocus() time.Time {
	return w.lastFocus
}

// CanPrevious satisfies the app.ElementWithNavigation interface.
func (w *Window) CanPrevious() bool {
	panic("not implemented")
}

// Previous satisfies the app.ElementWithNavigation interface.
func (w *Window) Previous() error {
	panic("not implemented")
}

// CanNext satisfies the app.ElementWithNavigation interface.
func (w *Window) CanNext() bool {
	panic("not implemented")
}

// Next satisfies the app.ElementWithNavigation interface.
func (w *Window) Next() error {
	panic("not implemented")
}

// Position satisfies the app.Window interface.
func (w *Window) Position() (x, y float64) {
	rawurl := fmt.Sprintf("/window/position?id=%s", w.id)
	res, err := w.driver.macos.RequestWithAsyncResponse(rawurl, nil)
	if err != nil {
		panic(errors.Wrap(err, "window position unavailable"))
	}

	var pos geom.Point
	res.Unmarshal(&pos)
	return pos.X, pos.Y
}

// Move satisfies the app.Window interface.
func (w *Window) Move(x, y float64) {
	panic("not implemented")
}

func onWindowMove(w *Window, u *url.URL, p bridge.Payload) {
	panic("not implemented")
}

// Center satisfies the app.Window interface.
func (w *Window) Center() {
	panic("not implemented")
}

// Size satisfies the app.Window interface.
func (w *Window) Size() (width, height float64) {
	panic("not implemented")
}

// Resize satisfies the app.Window interface.
func (w *Window) Resize(width, height float64) {
	panic("not implemented")
}

func onWindowResize(w *Window, u *url.URL, p bridge.Payload) {
	panic("not implemented")
}

// Focus satisfies the app.Window interface.
func (w *Window) Focus() {
	panic("not implemented")
}

func onWindowFocus(w *Window, u *url.URL, p bridge.Payload) {
	panic("not implemented")
}

// Close satisfies the app.Window interface.
func (w *Window) Close() {
	panic("not implemented")
}

func onWindowClose(w *Window, u *url.URL, p bridge.Payload) {
	w.driver.elements.Remove(w)
}
