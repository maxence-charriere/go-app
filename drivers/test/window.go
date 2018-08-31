package test

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/core"
	"github.com/murlokswarm/app/internal/dom"
	"github.com/pkg/errors"
)

// Window is a test window that implements the app.Window interface.
type Window struct {
	core.Window

	driver  *Driver
	dom     *dom.DOM
	history *core.History
	id      string
	compo   app.Compo
	x       float64
	y       float64
	width   float64
	height  float64

	onClose func() bool
}

func newWindow(d *Driver, c app.WindowConfig) *Window {
	w := &Window{
		driver:  d,
		dom:     dom.NewDOM(d.factory, dom.JsToGoHandler, dom.HrefCompoFmt),
		history: core.NewHistory(),
		id:      uuid.New().String(),

		onClose: c.OnClose,
	}

	d.elems.Put(w)

	if len(c.URL) != 0 {
		w.Load(c.URL)
	}

	return w
}

// ID satisfies the app.Window interface.
func (w *Window) ID() string {
	return w.id
}

// Load satisfies the app.Window interface.
func (w *Window) Load(urlFmt string, v ...interface{}) {
	var err error
	defer func() {
		w.SetErr(err)
	}()

	u := fmt.Sprintf(urlFmt, v...)
	n := core.CompoNameFromURLString(u)

	var c app.Compo
	if c, err = w.driver.factory.NewCompo(n); err != nil {
		return
	}

	if w.compo != nil {
		w.dom.Clean()
	}

	w.compo = c

	if u != w.history.Current() {
		w.history.NewEntry(u)
	}

	_, err = w.dom.New(c)
}

// Compo satisfies the app.Window interface.
func (w *Window) Compo() app.Compo {
	return w.compo
}

// Contains satisfies the app.Window interface.
func (w *Window) Contains(c app.Compo) bool {
	return w.dom.Contains(c)
}

// Render satisfies the app.Window interface.
func (w *Window) Render(c app.Compo) {
	_, err := w.dom.Update(c)
	w.SetErr(err)
}

// Reload satisfies the app.Window interface.
func (w *Window) Reload() {
	u := w.history.Current()

	if len(u) == 0 {
		w.SetErr(errors.New("no component loaded"))
		return
	}

	w.Load(u)
}

// CanPrevious satisfies the app.Window interface.
func (w *Window) CanPrevious() bool {
	return w.history.CanPrevious()
}

// Previous satisfies the app.Window interface.
func (w *Window) Previous() {
	u := w.history.Previous()

	if len(u) == 0 {
		w.SetErr(errors.New("no previous component"))
		return
	}

	w.Load(u)
}

// CanNext satisfies the app.Window interface.
func (w *Window) CanNext() bool {
	return w.history.CanNext()
}

// Next satisfies the app.Window interface.
func (w *Window) Next() {
	u := w.history.Next()

	if len(u) == 0 {
		w.SetErr(errors.New("no next component"))
		return
	}

	w.Load(u)
}

// Position satisfies the app.Window interface.
func (w *Window) Position() (x, y float64) {
	w.SetErr(nil)
	return w.x, w.y
}

// Move satisfies the app.Window interface.
func (w *Window) Move(x, y float64) {
	w.x = x
	w.y = y
	w.SetErr(nil)
}

// Center satisfies the app.Window interface.
func (w *Window) Center() {
	w.SetErr(nil)
}

// Size satisfies the app.Window interface.
func (w *Window) Size() (width, height float64) {
	w.SetErr(nil)
	return w.width, w.height
}

// Resize satisfies the app.Window interface.
func (w *Window) Resize(width, height float64) {
	w.height = height
	w.width = width
	w.SetErr(nil)
}

// Focus satisfies the app.Window interface.
func (w *Window) Focus() {
	w.SetErr(nil)
}

// FullScreen satisfies the app.Window interface.
func (w *Window) FullScreen() {
	w.SetErr(nil)
}

// ExitFullScreen satisfies the app.Window interface.
func (w *Window) ExitFullScreen() {
	w.SetErr(nil)
}

// Minimize satisfies the app.Window interface.
func (w *Window) Minimize() {
	w.SetErr(nil)
}

// Deminimize satisfies the app.Window interface.
func (w *Window) Deminimize() {
	w.SetErr(nil)
}

// Close satisfies the app.Window interface.
func (w *Window) Close() {
	if w.onClose != nil && !w.onClose() {
		return
	}

	w.driver.elems.Delete(w)
	w.SetErr(nil)
	w.driver.setElemErr(w)
}
