package mac

import (
	"fmt"
	"math"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/bridge"
	"github.com/murlokswarm/app/geom"
	"github.com/murlokswarm/app/markup"
	"github.com/pkg/errors"
)

// Window implements the app.Window interface.
type Window struct {
	driver    *Driver
	id        uuid.UUID
	env       markup.Env
	lastFocus time.Time

	onMove           func(x, y float64)
	onResize         func(width, height float64)
	onFocus          func()
	onBlur           func()
	onFullScreen     func()
	onExitFullScreen func()
	onMinimize       func()
	onDeminimize     func()
	onClose          func() bool
}

func newWindow(d *Driver, c app.WindowConfig) (w *Window, err error) {
	w = &Window{
		driver:    d,
		id:        uuid.New(),
		env:       markup.NewEnv(d.components),
		lastFocus: time.Now(),

		onMove:           c.OnMove,
		onResize:         c.OnResize,
		onFocus:          c.OnFocus,
		onBlur:           c.OnBlur,
		onFullScreen:     c.OnFullScreen,
		onExitFullScreen: c.OnExitFullScreen,
		onMinimize:       c.OnMinimize,
		onDeminimize:     c.OnDeminimize,
		onClose:          c.OnClose,
	}

	normalizeSize := func(min, max float64) (float64, float64) {
		min = math.Max(0, min)
		min = math.Min(min, 10000)

		if max == 0 {
			max = 10000
		}
		max = math.Max(0, max)
		max = math.Min(max, 10000)

		min = math.Min(min, max)
		return min, max
	}

	c.MinWidth, c.MaxWidth = normalizeSize(c.MinWidth, c.MaxWidth)
	c.MinHeight, c.MaxHeight = normalizeSize(c.MinHeight, c.MaxHeight)

	rawurl := fmt.Sprintf("/window/new?id=%s", w.id)
	if _, err = d.macos.Request(rawurl, bridge.NewPayload(c)); err != nil {
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
		panic(errors.Wrapf(err, "retrieving positon of window %v", w.ID()))
	}

	var pos geom.Point
	res.Unmarshal(&pos)
	return pos.X, pos.Y
}

// Move satisfies the app.Window interface.
func (w *Window) Move(x, y float64) {
	rawurl := fmt.Sprintf("/window/move?id=%s", w.id)
	payload := bridge.NewPayload(geom.Point{
		X: x,
		Y: y,
	})

	_, err := w.driver.macos.Request(rawurl, payload)
	if err != nil {
		panic(errors.Wrapf(err, "moving window %v failed", w.ID()))
	}
}

func onWindowMove(w *Window, u *url.URL, p bridge.Payload) (res bridge.Payload) {
	if w.onMove == nil {
		return
	}

	var pos geom.Point
	p.Unmarshal(&pos)
	w.onMove(pos.X, pos.Y)
	return
}

// Center satisfies the app.Window interface.
func (w *Window) Center() {
	rawurl := fmt.Sprintf("/window/center?id=%s", w.id)

	_, err := w.driver.macos.Request(rawurl, nil)
	if err != nil {
		panic(errors.Wrapf(err, "centering window %v failed", w.ID()))
	}
}

// Size satisfies the app.Window interface.
func (w *Window) Size() (width, height float64) {
	rawurl := fmt.Sprintf("/window/size?id=%s", w.id)

	res, err := w.driver.macos.RequestWithAsyncResponse(rawurl, nil)
	if err != nil {
		panic(errors.Wrapf(err, "retrieving size of window %v failed", w.ID()))
	}

	var size geom.Size
	res.Unmarshal(&size)
	return size.Width, size.Height
}

// Resize satisfies the app.Window interface.
func (w *Window) Resize(width, height float64) {
	rawurl := fmt.Sprintf("/window/resize?id=%s", w.id)
	payload := bridge.NewPayload(geom.Size{
		Width:  width,
		Height: height,
	})

	_, err := w.driver.macos.Request(rawurl, payload)
	if err != nil {
		panic(errors.Wrapf(err, "resizing window %v failed", w.ID()))
	}
}

func onWindowResize(w *Window, u *url.URL, p bridge.Payload) (res bridge.Payload) {
	if w.onResize == nil {
		return
	}

	var size geom.Size
	p.Unmarshal(&size)
	w.onResize(size.Width, size.Height)
	return
}

// Focus satisfies the app.Window interface.
func (w *Window) Focus() {
	rawurl := fmt.Sprintf("/window/focus?id=%s", w.id)

	_, err := w.driver.macos.Request(rawurl, nil)
	if err != nil {
		panic(errors.Wrapf(err, "focusing window %v failed", w.ID()))
	}
}

func onWindowFocus(w *Window, u *url.URL, p bridge.Payload) (res bridge.Payload) {
	if w.onFocus == nil {
		return
	}

	w.onFocus()
	return
}

func onWindowBlur(w *Window, u *url.URL, p bridge.Payload) (res bridge.Payload) {
	if w.onBlur == nil {
		return
	}

	w.onBlur()
	return
}

// ToggleFullScreen satisfies the app.Window interface.
func (w *Window) ToggleFullScreen() {
	rawurl := fmt.Sprintf("/window/togglefullscreen?id=%s", w.id)

	_, err := w.driver.macos.Request(rawurl, nil)
	if err != nil {
		panic(errors.Wrapf(err, "toggling full screen on window %v failed", w.ID()))
	}
}

func onWindowFullScreen(w *Window, u *url.URL, p bridge.Payload) (res bridge.Payload) {
	if w.onFullScreen == nil {
		return
	}

	w.onFullScreen()
	return
}

func onWindowExitFullScreen(w *Window, u *url.URL, p bridge.Payload) (res bridge.Payload) {
	if w.onExitFullScreen == nil {
		return
	}

	w.onExitFullScreen()
	return
}

// ToggleMinimize satisfies the app.Window interface.
func (w *Window) ToggleMinimize() {
	rawurl := fmt.Sprintf("/window/toggleminimize?id=%s", w.id)

	_, err := w.driver.macos.Request(rawurl, nil)
	if err != nil {
		panic(errors.Wrapf(err, "toggling minimize on window %v failed", w.ID()))
	}
}

func onWindowMinimize(w *Window, u *url.URL, p bridge.Payload) (res bridge.Payload) {
	if w.onMinimize == nil {
		return
	}

	w.onMinimize()
	return
}

func onWindowDeminimize(w *Window, u *url.URL, p bridge.Payload) (res bridge.Payload) {
	if w.onDeminimize == nil {
		return
	}

	w.onDeminimize()
	return
}

// Close satisfies the app.Window interface.
func (w *Window) Close() {
	rawurl := fmt.Sprintf("/window/close?id=%s", w.id)

	_, err := w.driver.macos.Request(rawurl, nil)
	if err != nil {
		panic(errors.Wrapf(err, "closing window %v failed", w.ID()))
	}
}

func onWindowClose(w *Window, u *url.URL, p bridge.Payload) (res bridge.Payload) {
	shouldClose := true
	if w.onClose != nil {
		shouldClose = w.onClose()
	}
	res = bridge.NewPayload(shouldClose)

	if shouldClose {
		w.driver.elements.Remove(w)
	}
	return
}
