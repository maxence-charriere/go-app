// +build darwin,amd64

package mac

import (
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"os/exec"
	"strings"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/bridge"
	"github.com/murlokswarm/app/internal/core"
	"github.com/murlokswarm/app/internal/dom"
	"github.com/murlokswarm/app/internal/file"
	"github.com/pkg/errors"
)

// Window implements the app.Window interface.
type Window struct {
	core.Window

	id           string
	dom          dom.Engine
	history      core.History
	compo        app.Compo
	isFocus      bool
	isFullscreen bool
	isMinimized  bool
}

func newWindow(c app.WindowConfig) *Window {
	id := uuid.New().String()

	w := &Window{
		id: id,
		dom: dom.Engine{
			Factory:   driver.factory,
			Resources: driver.Resources,
			AttrTransforms: []dom.Transform{
				dom.JsToGoHandler,
				dom.HrefCompoFmt,
			},
			UI: driver.UI,
		},
	}

	w.dom.Sync = w.render

	in := struct {
		ID                string
		Title             string
		X                 float64
		Y                 float64
		Width             float64
		MinWidth          float64
		MaxWidth          float64
		Height            float64
		MinHeight         float64
		MaxHeight         float64
		BackgroundColor   string
		FrostedBackground bool
		FixedSize         bool
		CloseHidden       bool
		MinimizeHidden    bool
		TitlebarHidden    bool
	}{
		ID:                w.id,
		Title:             c.Title,
		X:                 c.X,
		Y:                 c.Y,
		Width:             c.Width,
		MinWidth:          c.MinWidth,
		MaxWidth:          c.MaxWidth,
		Height:            c.Height,
		MinHeight:         c.MinHeight,
		MaxHeight:         c.MaxHeight,
		BackgroundColor:   c.BackgroundColor,
		FrostedBackground: c.FrostedBackground,
		FixedSize:         c.FixedSize,
		CloseHidden:       c.CloseHidden,
		MinimizeHidden:    c.MinimizeHidden,
	}

	if in.Width == 0 {
		in.Width = 1280
	}

	if in.Height == 0 {
		in.Height = 720
	}

	in.MinWidth, in.MaxWidth = normalizeWidowSize(in.MinWidth, in.MaxWidth)
	in.MinHeight, in.MaxHeight = normalizeWidowSize(in.MinHeight, in.MaxHeight)

	if err := driver.macRPC.Call("windows.New", nil, in); err != nil {
		w.SetErr(err)
		return w
	}

	driver.elems.Put(w)

	if len(c.URL) != 0 {
		w.Load(c.URL)
	}

	return w
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

	// Redirect web page to default web browser.
	if !driver.factory.IsCompoRegistered(n) {
		err = exec.Command("open", u).Run()
		return
	}

	var c app.Compo
	if c, err = driver.factory.NewCompo(n); err != nil {
		return
	}

	w.compo = c

	if u != w.history.Current() {
		w.history.NewEntry(u)
	}

	htmlConf := app.HTMLConfig{}
	if configurator, ok := c.(app.Configurator); ok {
		htmlConf = configurator.Config()
	}

	if len(htmlConf.CSS) == 0 {
		htmlConf.CSS = file.Filenames(driver.Resources("css"), ".css")
	}

	if len(htmlConf.Javascripts) == 0 {
		htmlConf.Javascripts = file.Filenames(driver.Resources("js"), ".js")
	}

	page := dom.Page{
		Title:         htmlConf.Title,
		Metas:         htmlConf.Metas,
		CSS:           htmlConf.CSS,
		Javascripts:   htmlConf.Javascripts,
		GoRequest:     "window.webkit.messageHandlers.golangRequest.postMessage",
		RootCompoName: n,
	}

	if err = driver.macRPC.Call("windows.Load", nil, struct {
		ID      string
		Title   string
		Page    string
		LoadURL string
		BaseURL string
	}{
		ID:      w.id,
		Title:   htmlConf.Title,
		Page:    page.String(),
		LoadURL: u,
		BaseURL: driver.Resources(),
	}); err != nil {
		return
	}

	err = w.dom.New(c)
	if err != nil {
		return
	}

	if nav, ok := c.(app.Navigable); ok {
		navURL, _ := url.Parse(u)
		nav.OnNavigate(navURL)
	}
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
	w.SetErr(w.dom.Render(c))
}

func (w *Window) render(changes interface{}) error {
	b, err := json.Marshal(changes)
	if err != nil {
		return errors.Wrap(err, "encode changes failed")
	}

	return driver.macRPC.Call("windows.Render", nil, struct {
		ID      string
		Changes string
	}{
		ID:      w.id,
		Changes: string(b),
	})
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
	out := struct {
		X float64
		Y float64
	}{}

	err := driver.macRPC.Call("windows.Position", &out, struct {
		ID string
	}{
		ID: w.id,
	})

	w.SetErr(err)
	return out.X, out.Y
}

// Move satisfies the app.Window interface.
func (w *Window) Move(x, y float64) {
	err := driver.macRPC.Call("windows.Move", nil, struct {
		ID string
		X  float64
		Y  float64
	}{
		ID: w.id,
		X:  x,
		Y:  y,
	})

	w.SetErr(err)
}

// Center satisfies the app.Window interface.
func (w *Window) Center() {
	err := driver.macRPC.Call("windows.Center", nil, struct {
		ID string
	}{
		ID: w.id,
	})

	w.SetErr(err)
}

// Size satisfies the app.Window interface.
func (w *Window) Size() (width, height float64) {
	out := struct {
		Width  float64
		Heigth float64
	}{}

	err := driver.macRPC.Call("windows.Size", &out, struct {
		ID string
	}{
		ID: w.id,
	})

	w.SetErr(err)
	return out.Width, out.Heigth
}

// Resize satisfies the app.Window interface.
func (w *Window) Resize(width, height float64) {
	err := driver.macRPC.Call("windows.Resize", nil, struct {
		ID     string
		Width  float64
		Height float64
	}{
		ID:     w.id,
		Width:  width,
		Height: height,
	})

	w.SetErr(err)
}

// Focus satisfies the app.Window interface.
func (w *Window) Focus() {
	err := driver.macRPC.Call("windows.Focus", nil, struct {
		ID string
	}{
		ID: w.id,
	})

	w.SetErr(err)
}

// IsFocus satisfies the app.Window interface.
func (w *Window) IsFocus() bool {
	return w.isFocus
}

// FullScreen satisfies the app.Window interface.
func (w *Window) FullScreen() {
	if w.isFullscreen {
		w.SetErr(nil)
		return
	}

	err := driver.macRPC.Call("windows.ToggleFullScreen", nil, struct {
		ID string
	}{
		ID: w.id,
	})

	w.SetErr(err)
}

// ExitFullScreen satisfies the app.Window interface.
func (w *Window) ExitFullScreen() {
	if !w.isFullscreen {
		w.SetErr(nil)
		return
	}

	err := driver.macRPC.Call("windows.ToggleFullScreen", nil, struct {
		ID string
	}{
		ID: w.id,
	})

	w.SetErr(err)
}

// IsFullScreen satisfies the app.Window interface.
func (w *Window) IsFullScreen() bool {
	return w.isFullscreen
}

// Minimize satisfies the app.Window interface.
func (w *Window) Minimize() {
	if w.isMinimized {
		w.SetErr(nil)
		return
	}

	err := driver.macRPC.Call("windows.ToggleMinimize", nil, struct {
		ID string
	}{
		ID: w.id,
	})

	w.SetErr(err)
}

// Deminimize satisfies the app.Window interface.
func (w *Window) Deminimize() {
	if !w.isMinimized {
		w.SetErr(nil)
		return
	}

	err := driver.macRPC.Call("windows.ToggleMinimize", nil, struct {
		ID string
	}{
		ID: w.id,
	})

	w.SetErr(err)
}

// IsMinimized satisfies the app.Window interface.
func (w *Window) IsMinimized() bool {
	return w.isMinimized
}

// Close satisfies the app.Window interface.
func (w *Window) Close() {
	err := driver.macRPC.Call("windows.Close", nil, struct {
		ID string
	}{
		ID: w.id,
	})

	w.SetErr(err)
}

// WhenWindow satisfies the app.Window interface.
func (w *Window) WhenWindow(f func(app.Window)) {
	f(w)
}

// WhenNavigator satisfies the app.Window interface.
func (w *Window) WhenNavigator(f func(app.Navigator)) {
	f(w)
}

func onWindowCallback(w *Window, in map[string]interface{}) interface{} {
	mappingStr := in["Mapping"].(string)

	var m dom.Mapping
	if err := json.Unmarshal([]byte(mappingStr), &m); err != nil {
		app.Logf("window callback failed: %s", err)
		return nil
	}

	if m.Override == "Files" {
		data, _ := json.Marshal(driver.droppedFiles)
		driver.droppedFiles = nil

		m.JSONValue = strings.Replace(
			m.JSONValue,
			`"FileOverride":"xxx"`,
			fmt.Sprintf(`"Files":%s`, data),
			1,
		)
	}

	c, err := w.dom.CompoByID(m.CompoID)
	if err != nil {
		app.Logf("window callback failed: %s", err)
		return nil
	}

	var f func()
	if f, err = m.Map(c); err != nil {
		app.Logf("window callback failed: %s", err)
		return nil
	}

	if f != nil {
		f()
		return nil
	}

	app.Render(c)
	return nil
}

func onWindowNavigate(w *Window, in map[string]interface{}) interface{} {
	e := app.ElemByCompo(w.Compo())

	e.WhenWindow(func(w app.Window) {
		w.Load(in["URL"].(string))
	})

	return nil
}

func onWindowAlert(w *Window, in map[string]interface{}) interface{} {
	app.Logf("%s", in["Alert"])
	return nil
}

func onWindowMove(w *Window, in map[string]interface{}) interface{} {
	driver.events.Emit(app.WindowMoved, w)
	return nil
}

func onWindowResize(w *Window, in map[string]interface{}) interface{} {
	driver.events.Emit(app.WindowResized, w)
	return nil
}

func onWindowFocus(w *Window, in map[string]interface{}) interface{} {
	w.isFocus = true
	driver.events.Emit(app.WindowFocused, w)
	return nil
}

func onWindowBlur(w *Window, in map[string]interface{}) interface{} {
	w.isFocus = false
	driver.events.Emit(app.WindowBlurred, w)
	return nil
}

func onWindowFullScreen(w *Window, in map[string]interface{}) interface{} {
	w.isFullscreen = true
	driver.events.Emit(app.WindowEnteredFullScreen, w)
	return nil
}

func onWindowExitFullScreen(w *Window, in map[string]interface{}) interface{} {
	w.isFullscreen = false
	driver.events.Emit(app.WindowExitedFullScreen, w)
	return nil
}

func onWindowMinimize(w *Window, in map[string]interface{}) interface{} {
	w.isMinimized = true
	driver.events.Emit(app.WindowMinimized, w)
	return nil
}

func onWindowDeminimize(w *Window, in map[string]interface{}) interface{} {
	w.isMinimized = false
	driver.events.Emit(app.WindowDeminimized, w)
	return nil
}

func onWindowClose(w *Window, in map[string]interface{}) interface{} {
	driver.events.Emit(app.WindowClosed, w)
	w.dom.Close()
	driver.elems.Delete(w)
	return nil
}

func handleWindow(h func(w *Window, in map[string]interface{}) interface{}) bridge.GoRPCHandler {
	return func(in map[string]interface{}) interface{} {
		id, _ := in["ID"].(string)

		e := driver.elems.GetByID(id)
		if e.Err() == app.ErrElemNotSet {
			return nil
		}

		return h(e.(*Window), in)
	}
}
