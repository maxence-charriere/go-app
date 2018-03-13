// +build darwin,amd64

package mac

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"net/url"
	"os/exec"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/appjs"
	"github.com/murlokswarm/app/bridge"
	"github.com/murlokswarm/app/html"
)

// Window implements the app.Window interface.
type Window struct {
	id        uuid.UUID
	markup    app.Markup
	component app.Component
	history   app.History
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

func newWindow(c app.WindowConfig) (app.Window, error) {
	var markup app.Markup = html.NewMarkup(driver.factory)
	markup = app.NewConcurrentMarkup(markup)

	history := app.NewHistory()
	history = app.NewConcurrentHistory(history)

	rawWin := &Window{
		id:        uuid.New(),
		markup:    markup,
		history:   history,
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
	win := app.NewWindowWithLogs(rawWin)

	in := struct {
		ID                 string
		Title              string
		X                  float64
		Y                  float64
		Width              float64
		MinWidth           float64
		MaxWidth           float64
		Height             float64
		MinHeight          float64
		MaxHeight          float64
		BackgroundColor    string
		FixedSize          bool
		CloseHidden        bool
		MinimizeHidden     bool
		TitlebarHidden     bool
		BackgroundVibrancy app.Vibrancy
	}{
		ID:                 win.ID().String(),
		Title:              c.Title,
		X:                  c.X,
		Y:                  c.Y,
		Width:              c.Width,
		MinWidth:           c.MinWidth,
		MaxWidth:           c.MaxWidth,
		Height:             c.Height,
		MinHeight:          c.MinHeight,
		MaxHeight:          c.MaxHeight,
		BackgroundColor:    c.BackgroundColor,
		FixedSize:          c.FixedSize,
		CloseHidden:        c.CloseHidden,
		MinimizeHidden:     c.MinimizeHidden,
		TitlebarHidden:     c.TitlebarHidden,
		BackgroundVibrancy: c.Mac.BackgroundVibrancy,
	}

	in.MinWidth, in.MaxWidth = normalizeWidowSize(in.MinWidth, in.MaxWidth)
	in.MinHeight, in.MaxHeight = normalizeWidowSize(in.MinHeight, in.MaxHeight)

	if err := driver.macRPC.Call("windows.New", in, nil); err != nil {
		return nil, err
	}

	if err := driver.elements.Add(win); err != nil {
		return nil, err
	}

	if len(c.DefaultURL) != 0 {
		if err := win.Load(c.DefaultURL); err != nil {
			return nil, err
		}
	}

	return win, nil
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

// ID satisfies the app.Window interface.
func (w *Window) ID() uuid.UUID {
	return w.id
}

// Base satisfies the app.Window interface.
func (w *Window) Base() app.Window {
	return w
}

// Load satisfies the app.Window interface.
func (w *Window) Load(rawurl string, v ...interface{}) error {
	rawurl = fmt.Sprintf(rawurl, v...)
	u, err := url.Parse(rawurl)
	if err != nil {
		return err
	}

	compoName := app.ComponentNameFromURL(u)
	isRegiteredCompo := driver.factory.Registered(compoName)
	currentURL, err := w.history.Current()

	if isRegiteredCompo && (err != nil || currentURL != u.String()) {
		w.history.NewEntry(u.String())
	}
	return w.load(u)
}

func (w *Window) load(u *url.URL) error {
	compo, err := driver.factory.New(app.ComponentNameFromURL(u))
	if err != nil {
		// Redirect web page to default web browser.
		return exec.
			Command("open", u.String()).
			Run()
	}

	if w.component != nil {
		w.markup.Dismount(w.component)
	}

	if _, err = w.markup.Mount(compo); err != nil {
		return err
	}
	w.component = compo

	if navigable, ok := compo.(app.Navigable); ok {
		navigable.OnNavigate(u)
	}

	var root app.Tag
	if root, err = w.markup.Root(compo); err != nil {
		return err
	}

	var buffer bytes.Buffer
	enc := html.NewEncoder(&buffer, w.markup, true)
	if err = enc.Encode(root); err != nil {
		return err
	}

	pageConfig := html.PageConfig{
		DisableDefaultContextMenu: true,
	}
	if page, ok := compo.(html.Page); ok {
		pageConfig = page.PageConfig()
	}

	if len(pageConfig.CSS) == 0 {
		pageConfig.CSS = app.CSSResources()
	}

	pageConfig.DefaultComponent = template.HTML(buffer.String())
	pageConfig.AppJS = appjs.AppJS("window.webkit.messageHandlers.golangRequest.postMessage")

	in := struct {
		ID      string
		Title   string
		Page    string
		LoadURL string
		BaseURL string
	}{
		ID:      w.ID().String(),
		Title:   pageConfig.Title,
		Page:    html.NewPage(pageConfig),
		LoadURL: u.String(),
		BaseURL: driver.Resources(),
	}

	return driver.macRPC.Call("windows.Load", in, nil)
}

// Contains satisfies the app.Window interface.
func (w *Window) Contains(compo app.Component) bool {
	return w.markup.Contains(compo)
}

// Component satisfies the app.Window interface.
func (w *Window) Component() app.Component {
	return w.component
}

// Render satisfies the app.Window interface.
func (w *Window) Render(compo app.Component) error {
	syncs, err := w.markup.Update(compo)
	if err != nil {
		return err
	}

	for _, sync := range syncs {
		if sync.Replace {
			err = w.render(sync)
		} else {
			err = w.renderAttributes(sync)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (w *Window) render(sync app.TagSync) error {
	var buffer bytes.Buffer

	enc := html.NewEncoder(&buffer, w.markup, true)
	if err := enc.Encode(sync.Tag); err != nil {
		return err
	}

	render, err := json.Marshal(struct {
		ID        string `json:"id"`
		Component string `json:"component"`
	}{
		ID:        sync.Tag.ID.String(),
		Component: buffer.String(),
	})
	if err != nil {
		return err
	}

	in := struct {
		ID     string
		Render string
	}{
		ID:     w.ID().String(),
		Render: string(render),
	}

	return driver.macRPC.Call("windows.Render", in, nil)
}

func (w *Window) renderAttributes(sync app.TagSync) error {
	attrs := make(app.AttributeMap, len(sync.Tag.Attributes))
	for name, val := range sync.Tag.Attributes {
		attrs[name] = html.AttrValueFormatter{
			Name:       name,
			Value:      val,
			FormatHref: true,
			CompoID:    sync.Tag.CompoID,
			Factory:    driver.factory,
		}.Format()
	}

	render, err := json.Marshal(struct {
		ID         string           `json:"id"`
		Attributes app.AttributeMap `json:"attributes"`
	}{
		ID:         sync.Tag.ID.String(),
		Attributes: attrs,
	})
	if err != nil {
		return err
	}

	in := struct {
		ID     string
		Render string
	}{
		ID:     w.ID().String(),
		Render: string(render),
	}

	return driver.macRPC.Call("windows.RenderAttributes", in, nil)
}

// LastFocus satisfies the app.Window interface.
func (w *Window) LastFocus() time.Time {
	return w.lastFocus
}

// Reload satisfies the app.Window interface.
func (w *Window) Reload() error {
	var rawurl string
	var u *url.URL
	var err error

	if rawurl, err = w.history.Current(); err != nil {
		return err
	}

	if u, err = url.Parse(rawurl); err != nil {
		return err
	}

	return w.load(u)
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

	var u *url.URL
	if u, err = url.Parse(rawurl); err != nil {
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

	var u *url.URL
	if u, err = url.Parse(rawurl); err != nil {
		return err
	}
	return w.load(u)
}

// Position satisfies the app.Window interface.
func (w *Window) Position() (x, y float64) {
	in := struct {
		ID string
	}{
		ID: w.ID().String(),
	}

	var out struct {
		X float64
		Y float64
	}

	if err := driver.macRPC.Call("windows.Position", in, &out); err != nil {
		panic(err)
	}
	return out.X, out.Y
}

// Move satisfies the app.Window interface.
func (w *Window) Move(x, y float64) {
	in := struct {
		ID string
		X  float64
		Y  float64
	}{
		ID: w.ID().String(),
		X:  x,
		Y:  y,
	}

	if err := driver.macRPC.Call("windows.Move", in, nil); err != nil {
		panic(err)
	}
}

func onWindowMove(w *Window, u *url.URL, p bridge.Payload) (res bridge.Payload) {
	if w.onMove == nil {
		return nil
	}

	var pos point
	p.Unmarshal(&pos)
	w.onMove(pos.X, pos.Y)
	return nil
}

// Center satisfies the app.Window interface.
func (w *Window) Center() {
	in := struct {
		ID string
	}{
		ID: w.ID().String(),
	}

	if err := driver.macRPC.Call("windows.Center", in, nil); err != nil {
		panic(err)
	}
}

// Size satisfies the app.Window interface.
func (w *Window) Size() (width, height float64) {
	in := struct {
		ID string
	}{
		ID: w.ID().String(),
	}

	var out struct {
		Width  float64
		Heigth float64
	}

	if err := driver.macRPC.Call("windows.Size", in, &out); err != nil {
		panic(err)
	}
	return out.Width, out.Heigth
}

// Resize satisfies the app.Window interface.
func (w *Window) Resize(width, height float64) {
	in := struct {
		ID     string
		Width  float64
		Height float64
	}{
		ID:     w.ID().String(),
		Width:  width,
		Height: height,
	}

	if err := driver.macRPC.Call("windows.Resize", in, nil); err != nil {
		panic(err)
	}
}

func onWindowResize(w *Window, u *url.URL, p bridge.Payload) (res bridge.Payload) {
	if w.onResize == nil {
		return nil
	}

	var size size
	p.Unmarshal(&size)
	w.onResize(size.Width, size.Height)
	return nil
}

// Focus satisfies the app.Window interface.
func (w *Window) Focus() {
	in := struct {
		ID string
	}{
		ID: w.ID().String(),
	}

	if err := driver.macRPC.Call("windows.Focus", in, nil); err != nil {
		panic(err)
	}
}

func onWindowFocus(w *Window, u *url.URL, p bridge.Payload) (res bridge.Payload) {
	w.lastFocus = time.Now()

	if w.onFocus == nil {
		return nil
	}

	w.onFocus()
	return nil
}

func onWindowBlur(w *Window, u *url.URL, p bridge.Payload) (res bridge.Payload) {
	if w.onBlur == nil {
		return nil
	}

	w.onBlur()
	return nil
}

// ToggleFullScreen satisfies the app.Window interface.
func (w *Window) ToggleFullScreen() {
	in := struct {
		ID string
	}{
		ID: w.ID().String(),
	}

	if err := driver.macRPC.Call("windows.ToggleFullScreen", in, nil); err != nil {
		panic(err)
	}
}

func onWindowFullScreen(w *Window, u *url.URL, p bridge.Payload) (res bridge.Payload) {
	if w.onFullScreen == nil {
		return nil
	}

	w.onFullScreen()
	return nil
}

func onWindowExitFullScreen(w *Window, u *url.URL, p bridge.Payload) (res bridge.Payload) {
	if w.onExitFullScreen == nil {
		return nil
	}

	w.onExitFullScreen()
	return nil
}

// ToggleMinimize satisfies the app.Window interface.
func (w *Window) ToggleMinimize() {
	in := struct {
		ID string
	}{
		ID: w.ID().String(),
	}

	if err := driver.macRPC.Call("windows.ToggleMinimize", in, nil); err != nil {
		panic(err)
	}
}

func onWindowMinimize(w *Window, u *url.URL, p bridge.Payload) (res bridge.Payload) {
	if w.onMinimize == nil {
		return nil
	}

	w.onMinimize()
	return nil
}

func onWindowDeminimize(w *Window, u *url.URL, p bridge.Payload) (res bridge.Payload) {
	if w.onDeminimize == nil {
		return nil
	}

	w.onDeminimize()
	return nil
}

// Close satisfies the app.Window interface.
func (w *Window) Close() {
	in := struct {
		ID string
	}{
		ID: w.ID().String(),
	}

	if err := driver.macRPC.Call("windows.Close", in, nil); err != nil {
		panic(err)
	}
}

func onWindowClose(w *Window, u *url.URL, p bridge.Payload) (res bridge.Payload) {
	shouldClose := true
	if w.onClose != nil {
		shouldClose = w.onClose()
	}
	res = bridge.NewPayload(shouldClose)

	if w.component != nil {
		w.markup.Dismount(w.component)
	}

	if shouldClose {
		driver.elements.Remove(w)
	}
	return res
}

func onWindowCallback(w *Window, u *url.URL, p bridge.Payload) (res bridge.Payload) {
	var mapping app.Mapping
	p.Unmarshal(&mapping)

	if mapping.Override == "Files" {
		data, _ := json.Marshal(driver.droppedFiles)
		driver.droppedFiles = nil

		mapping.JSONValue = strings.Replace(
			mapping.JSONValue,
			`"file-override":"xxx"`,
			fmt.Sprintf(`"Files":%s`, data),
			1,
		)
	}

	function, err := w.markup.Map(mapping)
	if err != nil {
		app.DefaultLogger.Error(err)
		return nil
	}

	if function != nil {
		function()
		return nil
	}

	var compo app.Component
	if compo, err = w.markup.Component(mapping.CompoID); err != nil {
		app.Error(err)
		return nil
	}

	if err = w.Render(compo); err != nil {
		app.Error(err)
	}

	return nil
}

func onWindowNavigate(w *Window, u *url.URL, p bridge.Payload) (res bridge.Payload) {
	var rawurl string
	p.Unmarshal(&rawurl)
	w.Load(rawurl)
	return nil
}
