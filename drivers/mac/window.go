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
	"github.com/murlokswarm/app/internal/core"
)

// Window implements the app.Window interface.
type Window struct {
	core.ElementWithComponent

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
	markup = app.ConcurrentMarkup(markup)

	history := app.NewHistory()
	history = app.ConcurrentHistory(history)

	win := &Window{
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

	if err := driver.macRPC.Call("windows.New", nil, in); err != nil {
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

// Load satisfies the app.Window interface.
func (w *Window) Load(rawurl string, v ...interface{}) error {
	rawurl = fmt.Sprintf(rawurl, v...)
	u, err := url.Parse(rawurl)
	if err != nil {
		return err
	}

	compoName := app.ComponentNameFromURL(u)
	isRegisteredCompo := driver.factory.Registered(compoName)
	currentURL, err := w.history.Current()

	if isRegisteredCompo && (err != nil || currentURL != u.String()) {
		w.history.NewEntry(u.String())
	}
	return w.load(u)
}

func (w *Window) load(u *url.URL) error {
	compoName := app.ComponentNameFromURL(u)
	if len(compoName) == 0 {
		// Redirect web page to default web browser.
		return exec.
			Command("open", u.String()).
			Run()
	}

	compo, err := driver.factory.New(app.ComponentNameFromURL(u))
	if err != nil {
		return err
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

	return driver.macRPC.Call("windows.Load", nil, struct {
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
	})
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

	return driver.macRPC.Call("windows.Render", nil, struct {
		ID     string
		Render string
	}{
		ID:     w.ID().String(),
		Render: string(render),
	})
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

	return driver.macRPC.Call("windows.RenderAttributes", nil, struct {
		ID     string
		Render string
	}{
		ID:     w.ID().String(),
		Render: string(render),
	})
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
	var out struct {
		X float64
		Y float64
	}

	if err := driver.macRPC.Call("windows.Position", &out, struct {
		ID string
	}{
		ID: w.ID().String(),
	}); err != nil {
		panic(err)
	}
	return out.X, out.Y
}

// Move satisfies the app.Window interface.
func (w *Window) Move(x, y float64) error {
	return driver.macRPC.Call("windows.Move", nil, struct {
		ID string
		X  float64
		Y  float64
	}{
		ID: w.ID().String(),
		X:  x,
		Y:  y,
	})
}

func onWindowMove(w *Window, in map[string]interface{}) interface{} {
	if w.onMove != nil {
		w.onMove(in["X"].(float64), in["Y"].(float64))
	}
	return nil
}

// Center satisfies the app.Window interface.
func (w *Window) Center() error {
	return driver.macRPC.Call("windows.Center", nil, struct {
		ID string
	}{
		ID: w.ID().String(),
	})
}

// Size satisfies the app.Window interface.
func (w *Window) Size() (width, height float64) {
	var out struct {
		Width  float64
		Heigth float64
	}

	if err := driver.macRPC.Call("windows.Size", &out, struct {
		ID string
	}{
		ID: w.ID().String(),
	}); err != nil {
		panic(err)
	}
	return out.Width, out.Heigth
}

// Resize satisfies the app.Window interface.
func (w *Window) Resize(width, height float64) error {
	return driver.macRPC.Call("windows.Resize", nil, struct {
		ID     string
		Width  float64
		Height float64
	}{
		ID:     w.ID().String(),
		Width:  width,
		Height: height,
	})
}

func onWindowResize(w *Window, in map[string]interface{}) interface{} {
	if w.onResize != nil {
		w.onResize(in["Width"].(float64), in["Height"].(float64))
	}
	return nil
}

// Focus satisfies the app.Window interface.
func (w *Window) Focus() error {
	return driver.macRPC.Call("windows.Focus", nil, struct {
		ID string
	}{
		ID: w.ID().String(),
	})
}

func onWindowFocus(w *Window, in map[string]interface{}) interface{} {
	w.lastFocus = time.Now()

	if w.onFocus != nil {
		w.onFocus()
	}
	return nil
}

func onWindowBlur(w *Window, in map[string]interface{}) interface{} {
	if w.onBlur != nil {
		w.onBlur()
	}
	return nil
}

// ToggleFullScreen satisfies the app.Window interface.
func (w *Window) ToggleFullScreen() error {
	return driver.macRPC.Call("windows.ToggleFullScreen", nil, struct {
		ID string
	}{
		ID: w.ID().String(),
	})
}

func onWindowFullScreen(w *Window, in map[string]interface{}) interface{} {
	if w.onFullScreen != nil {
		w.onFullScreen()
	}
	return nil
}

func onWindowExitFullScreen(w *Window, in map[string]interface{}) interface{} {
	if w.onExitFullScreen != nil {
		w.onExitFullScreen()
	}
	return nil
}

// ToggleMinimize satisfies the app.Window interface.
func (w *Window) ToggleMinimize() error {
	return driver.macRPC.Call("windows.ToggleMinimize", nil, struct {
		ID string
	}{
		ID: w.ID().String(),
	})
}

func onWindowMinimize(w *Window, in map[string]interface{}) interface{} {
	if w.onMinimize != nil {
		w.onMinimize()
	}
	return nil
}

func onWindowDeminimize(w *Window, in map[string]interface{}) interface{} {
	if w.onDeminimize != nil {
		w.onDeminimize()
	}
	return nil
}

// Close satisfies the app.Window interface.
func (w *Window) Close() error {
	return driver.macRPC.Call("windows.Close", nil, struct {
		ID string
	}{
		ID: w.ID().String(),
	})
}

func onWindowClose(w *Window, in map[string]interface{}) interface{} {
	shouldClose := true
	if w.onClose != nil {
		shouldClose = w.onClose()
	}

	if shouldClose {
		if w.component != nil {
			w.markup.Dismount(w.component)
		}
		driver.elements.Remove(w)
	}

	return struct {
		ShouldClose bool
	}{
		ShouldClose: shouldClose,
	}
}

// WhenWindow calls the given handler.
// It satisfies the app.ElementWithComponent interface.
func (w *Window) WhenWindow(f func(app.Window)) {
	f(w)
}

func onWindowCallback(w *Window, in map[string]interface{}) interface{} {
	mappingString := in["Mapping"].(string)

	var mapping app.Mapping
	if err := json.Unmarshal([]byte(mappingString), &mapping); err != nil {
		app.Log("window callback failed: %s", err)
		return nil
	}

	if mapping.Override == "Files" {
		data, _ := json.Marshal(driver.droppedFiles)
		driver.droppedFiles = nil

		mapping.JSONValue = strings.Replace(
			mapping.JSONValue,
			`"FileOverride":"xxx"`,
			fmt.Sprintf(`"Files":%s`, data),
			1,
		)
	}

	function, err := w.markup.Map(mapping)
	if err != nil {
		app.Log("window callback failed: %s", err)
		return nil
	}

	if function != nil {
		function()
		return nil
	}

	var compo app.Component
	if compo, err = w.markup.Component(mapping.CompoID); err != nil {
		app.Log("window callback failed: %s", err)
		return nil
	}

	if err = w.Render(compo); err != nil {
		app.Log("window callback failed: %s", err)
	}
	return nil
}

func onWindowNavigate(w *Window, in map[string]interface{}) interface{} {
	win, _ := app.WindowByComponent(w.Component())
	win.Load(in["URL"].(string))
	return nil
}

func handleWindow(h func(w *Window, in map[string]interface{}) interface{}) bridge.GoRPCHandler {
	return func(in map[string]interface{}) interface{} {
		id, _ := uuid.Parse(in["ID"].(string))

		elem, err := driver.elements.Element(id)
		if err != nil {
			return nil
		}

		win := elem.(*Window)
		return h(win, in)
	}
}
