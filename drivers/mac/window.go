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
	"github.com/murlokswarm/app/internal/file"
	html "github.com/murlokswarm/app/internal/html-v2"
	"github.com/pkg/errors"
)

// Window implements the app.Window interface.
type Window struct {
	core.Window

	dom          html.DOM
	history      *core.History
	id           string
	compo        app.Compo
	isFullscreen bool
	isMinimized  bool

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

func newWindow(c app.WindowConfig) *Window {
	id := uuid.New().String()

	w := &Window{
		history: core.NewHistory(),
		id:      id,

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
		ID:                 w.id,
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

	if w.compo != nil {
		// Close dom.
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
		htmlConf.CSS = file.CSS(driver.Resources("css"))
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
		Page:    html.Page(htmlConf),
		LoadURL: u,
		BaseURL: driver.Resources(),
	}); err != nil {
		return
	}

	w.dom = html.NewDOM(driver.factory, w.id)

	w.Render(c)
	if w.Err() != nil {
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
	var err error
	defer func() {
		w.SetErr(err)
	}()

	var changes []html.Change
	if changes, err = w.dom.Render(c); err != nil {
		return
	}

	var b []byte
	if b, err = json.Marshal(changes); err != nil {
		return
	}

	err = driver.macRPC.Call("windows.Render", nil, struct {
		ID     string
		Render string
	}{
		ID:     w.id,
		Render: string(b),
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
	var out struct {
		X float64
		Y float64
	}

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
	var out struct {
		Width  float64
		Heigth float64
	}

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

// Close satisfies the app.Window interface.
func (w *Window) Close() {
	err := driver.macRPC.Call("windows.Close", nil, struct {
		ID string
	}{
		ID: w.id,
	})

	w.SetErr(err)
}

// func (w *Window) render(sync app.TagSync) error {
// 	var buffer bytes.Buffer

// 	enc := html.NewEncoder(&buffer, w.markup, true)
// 	if err := enc.Encode(sync.Tag); err != nil {
// 		return err
// 	}

// 	render, err := json.Marshal(struct {
// 		ID    string `json:"id"`
// 		Compo string `json:"component"`
// 	}{
// 		ID:    sync.Tag.ID,
// 		Compo: buffer.String(),
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	return driver.macRPC.Call("windows.Render", nil, struct {
// 		ID     string
// 		Render string
// 	}{
// 		ID:     w.id,
// 		Render: string(render),
// 	})
// }

// func (w *Window) renderAttributes(sync app.TagSync) error {
// 	attrs := make(app.AttributeMap, len(sync.Tag.Attributes))
// 	for name, val := range sync.Tag.Attributes {
// 		attrs[name] = html.AttrValueFormatter{
// 			Name:       name,
// 			Value:      val,
// 			FormatHref: true,
// 			CompoID:    sync.Tag.CompoID,
// 			Factory:    driver.factory,
// 		}.Format()
// 	}

// 	render, err := json.Marshal(struct {
// 		ID         string           `json:"id"`
// 		Attributes app.AttributeMap `json:"attributes"`
// 	}{
// 		ID:         sync.Tag.ID,
// 		Attributes: attrs,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	return driver.macRPC.Call("windows.RenderAttributes", nil, struct {
// 		ID     string
// 		Render string
// 	}{
// 		ID:     w.id,
// 		Render: string(render),
// 	})
// }

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

	// function, err := w.markup.Map(mapping)
	// if err != nil {
	// 	app.Log("window callback failed: %s", err)
	// 	return nil
	// }

	// if function != nil {
	// 	function()
	// 	return nil
	// }

	// var c app.Compo
	// if c, err = w.markup.Compo(mapping.CompoID); err != nil {
	// 	app.Log("window callback failed: %s", err)
	// 	return nil
	// }

	// w.Render(c)
	// if w.Err() != nil {
	// 	app.Log("window callback failed: %s", w.Err())
	// }

	return nil
}

func onWindowNavigate(w *Window, in map[string]interface{}) interface{} {
	e := app.ElemByCompo(w.Compo())

	e.WhenWindow(func(w app.Window) {
		w.Load(in["URL"].(string))
	})

	return nil
}

func onWindowMove(w *Window, in map[string]interface{}) interface{} {
	if w.onMove != nil {
		w.onMove(
			in["X"].(float64),
			in["Y"].(float64),
		)
	}

	return nil
}

func onWindowResize(w *Window, in map[string]interface{}) interface{} {
	if w.onResize != nil {
		w.onResize(
			in["Width"].(float64),
			in["Height"].(float64),
		)
	}

	return nil
}

func onWindowFocus(w *Window, in map[string]interface{}) interface{} {
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

func onWindowFullScreen(w *Window, in map[string]interface{}) interface{} {
	if w.onFullScreen != nil {
		w.onFullScreen()
	}

	w.isFullscreen = true
	return nil
}

func onWindowExitFullScreen(w *Window, in map[string]interface{}) interface{} {
	if w.onExitFullScreen != nil {
		w.onExitFullScreen()
	}

	w.isFullscreen = false
	return nil
}

func onWindowMinimize(w *Window, in map[string]interface{}) interface{} {
	if w.onMinimize != nil {
		w.onMinimize()
	}

	w.isMinimized = true
	return nil
}

func onWindowDeminimize(w *Window, in map[string]interface{}) interface{} {
	if w.onDeminimize != nil {
		w.onDeminimize()
	}

	w.isMinimized = false
	return nil
}

func onWindowClose(w *Window, in map[string]interface{}) interface{} {
	shouldClose := true
	if w.onClose != nil {
		shouldClose = w.onClose()
	}

	if shouldClose {
		// dom.Close()
		driver.elems.Delete(w)
	}

	return struct {
		ShouldClose bool
	}{
		ShouldClose: shouldClose,
	}
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
