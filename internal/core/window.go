package core

import (
	"github.com/murlokswarm/app"
)

// Window is a base struct to embed in app.Window implementations.
type Window struct {
	Elem
}

// WhenWindow satisfies the app.Window interface.
func (w *Window) WhenWindow(f func(app.Window)) {
	f(w)
}

// WhenNavigator satisfies the app.Window interface.
func (w *Window) WhenNavigator(f func(app.Navigator)) {
	f(w)
}

// Load satisfies the app.Window interface.
func (w *Window) Load(url string, v ...interface{}) {
	w.SetErr(app.ErrNotSupported)
}

// Compo satisfies the app.Window interface.
func (w *Window) Compo() app.Compo {
	return nil
}

// Contains satisfies the app.Window interface.
func (w *Window) Contains(c app.Compo) bool {
	return false
}

// Render satisfies the app.Window interface.
func (w *Window) Render(c app.Compo) {
	w.SetErr(app.ErrNotSupported)
}

// Reload satisfies the app.Window interface.
func (w *Window) Reload() {
	w.SetErr(app.ErrNotSupported)
}

// CanPrevious satisfies the app.Window interface.
func (w *Window) CanPrevious() bool {
	return false
}

// Previous satisfies the app.Window interface.
func (w *Window) Previous() {
	w.SetErr(app.ErrNotSupported)
}

// CanNext satisfies the app.Window interface.
func (w *Window) CanNext() bool {
	return false
}

// Next satisfies the app.Window interface.
func (w *Window) Next() {
	w.SetErr(app.ErrNotSupported)
}

// Position satisfies the app.Window interface.
func (w *Window) Position() (x, y float64) {
	w.SetErr(app.ErrNotSupported)
	return 0, 0
}

// Move satisfies the app.Window interface.
func (w *Window) Move(x, y float64) {
	w.SetErr(app.ErrNotSupported)
}

// Center satisfies the app.Window interface.
func (w *Window) Center() {
	w.SetErr(app.ErrNotSupported)
}

// Size satisfies the app.Window interface.
func (w *Window) Size() (width, height float64) {
	w.SetErr(app.ErrNotSupported)
	return 0, 0
}

// Resize satisfies the app.Window interface.
func (w *Window) Resize(width, height float64) {
	w.SetErr(app.ErrNotSupported)
}

// Focus satisfies the app.Window interface.
func (w *Window) Focus() {
	w.SetErr(app.ErrNotSupported)
}

// FullScreen satisfies the app.Window interface.
func (w *Window) FullScreen() {
	w.SetErr(app.ErrNotSupported)
}

// ExitFullScreen satisfies the app.Window interface.
func (w *Window) ExitFullScreen() {
	w.SetErr(app.ErrNotSupported)
}

// Minimize satisfies the app.Window interface.
func (w *Window) Minimize() {
	w.SetErr(app.ErrNotSupported)
}

// Deminimize satisfies the app.Window interface.
func (w *Window) Deminimize() {
	w.SetErr(app.ErrNotSupported)
}

// Close satisfies the app.Window interface.
func (w *Window) Close() {
	w.SetErr(app.ErrNotSupported)
}
