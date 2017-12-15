// +build darwin,amd64

package mac

import (
	"bytes"
	"fmt"
	"html/template"
	"math"
	"net/url"
	"time"

	"github.com/murlokswarm/app/appjs"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/bridge"
	"github.com/murlokswarm/app/html"
	"github.com/pkg/errors"
)

// Window implements the app.Window interface.
type Window struct {
	driver    *Driver
	id        uuid.UUID
	markup    app.Markup
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

func newWindow(driver *Driver, config app.WindowConfig) (w *Window, err error) {
	w = &Window{
		id:        uuid.New(),
		driver:    driver,
		markup:    app.NewConcurrentMarkup(html.NewMarkup(driver.factory)),
		lastFocus: time.Now(),

		onMove:           config.OnMove,
		onResize:         config.OnResize,
		onFocus:          config.OnFocus,
		onBlur:           config.OnBlur,
		onFullScreen:     config.OnFullScreen,
		onExitFullScreen: config.OnExitFullScreen,
		onMinimize:       config.OnMinimize,
		onDeminimize:     config.OnDeminimize,
		onClose:          config.OnClose,
	}

	var page string
	if page, err = w.makePage(config); err != nil {
		return
	}

	app.DefaultLogger.Log(page)

	payload := struct {
		Title           string              `json:"title"`
		X               float64             `json:"x"`
		Y               float64             `json:"y"`
		Width           float64             `json:"width"`
		MinWidth        float64             `json:"min-width"`
		MaxWidth        float64             `json:"max-width"`
		Height          float64             `json:"height"`
		MinHeight       float64             `json:"min-height"`
		MaxHeight       float64             `json:"max-height"`
		BackgroundColor string              `json:"background-color"`
		NoResizable     bool                `json:"no-resizable"`
		NoClosable      bool                `json:"no-closable"`
		NoMinimizable   bool                `json:"no-minimizable"`
		TitlebarHidden  bool                `json:"titlebar-hidden"`
		Page            string              `json:"page"`
		BaseURL         string              `json:"base-url"`
		Mac             app.MacWindowConfig `json:"mac"`
	}{
		Title:           config.Title,
		X:               config.X,
		Y:               config.Y,
		Width:           config.Width,
		BackgroundColor: config.BackgroundColor,
		NoResizable:     config.NoResizable,
		NoClosable:      config.NoClosable,
		NoMinimizable:   config.NoMinimizable,
		TitlebarHidden:  config.TitlebarHidden,
		Page:            page,
		BaseURL:         driver.Resources(),
		Mac:             config.Mac,
	}

	payload.MinWidth, payload.MaxWidth = normalizeWidowSize(config.MinWidth, config.MaxWidth)
	payload.MinHeight, payload.MaxHeight = normalizeWidowSize(config.MinHeight, config.MaxHeight)

	if _, err = driver.macos.Request(
		fmt.Sprintf("/window/new?id=%s", w.id),
		bridge.NewPayload(payload),
	); err != nil {
		return
	}

	err = driver.elements.Add(w)
	return
}

func normalizeWidowSize(min, max float64) (float64, float64) {
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

func (w *Window) makePage(config app.WindowConfig) (page string, err error) {
	var u *url.URL
	var renderedCompo string

	if u, err = url.Parse(config.DefaultURL); err != nil {
		return
	}

	if compoName := app.ComponentNameFromURL(u); len(compoName) != 0 {
		var compo app.Component
		var root app.Tag
		var buffer bytes.Buffer

		if compo, err = w.driver.factory.NewComponent(compoName); err != nil {
			return
		}

		if root, err = w.markup.Mount(compo); err != nil {
			return
		}

		enc := html.NewEncoder(&buffer, w.markup)
		if err = enc.Encode(root); err != nil {
			return
		}
		renderedCompo = buffer.String()
	}

	page = html.NewPage(html.PageConfig{
		Title:            config.Title,
		DefaultComponent: template.HTML(renderedCompo),
		AppJS:            appjs.AppJS("window.webkit.messageHandlers.golangRequest.postMessage"),
	})
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
func (w *Window) Contains(compo app.Component) bool {
	panic("not implemented")
}

// Render satisfies the app.ElementWithComponent interface.
func (w *Window) Render(compo app.Component) error {
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

	var pos point
	res.Unmarshal(&pos)
	return pos.X, pos.Y
}

// Move satisfies the app.Window interface.
func (w *Window) Move(x, y float64) {
	rawurl := fmt.Sprintf("/window/move?id=%s", w.id)
	payload := bridge.NewPayload(point{
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

	var pos point
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

	var size size
	res.Unmarshal(&size)
	return size.Width, size.Height
}

// Resize satisfies the app.Window interface.
func (w *Window) Resize(width, height float64) {
	rawurl := fmt.Sprintf("/window/resize?id=%s", w.id)
	payload := bridge.NewPayload(size{
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

	var size size
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
	w.lastFocus = time.Now()

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
