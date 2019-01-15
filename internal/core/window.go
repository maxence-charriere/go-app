package core

import (
	"encoding/json"
	"fmt"
	"math"
	"net/url"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/dom"
	"github.com/murlokswarm/app/internal/file"
	"github.com/pkg/errors"
)

// Window is a modular implementation of the app.Window interface that can be
// configured address the different drivers needs.
type Window struct {
	Elem

	ConvertHTMLPaths func([]string) []string
	DefaultWidth     float64
	DefaultHeight    float64
	DOM              dom.Engine
	Driver           *Driver
	History          History

	compo        app.Compo
	x            float64
	y            float64
	width        float64
	height       float64
	isFocus      bool
	isFullScreen bool
	isMinimized  bool
}

// Create creates and display the window.
func (w *Window) Create(c app.WindowConfig) {
	w.id = uuid.New().String()
	w.DOM.Sync = w.render

	if w.ConvertHTMLPaths == nil {
		w.ConvertHTMLPaths = func(paths []string) []string {
			return paths
		}
	}

	if c.Width == 0 {
		c.Width = w.DefaultWidth
	}

	if c.Height == 0 {
		c.Height = w.DefaultHeight
	}

	if c.MaxWidth == 0 {
		c.MaxWidth = math.MaxFloat64
	}

	if c.MaxHeight == 0 {
		c.MaxHeight = math.MaxFloat64
	}

	out := struct {
		X      float64
		Y      float64
		Width  float64
		Height float64
	}{}

	if w.err = w.Driver.Platform.Call("windows.New", &out, struct {
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
	}); w.err != nil {
		return
	}

	w.x = out.X
	w.y = out.Y
	w.width = out.Width
	w.height = out.Height

	w.Driver.Elems.Put(w)

	if len(c.URL) != 0 {
		w.Load(c.URL)
	}
}

// Contains satisfies the app.Window interface.
func (w *Window) Contains(c app.Compo) bool {
	return w.DOM.Contains(c)
}

// WhenView satisfies the app.Window interface.
func (w *Window) WhenView(f func(app.View)) {
	f(w)
}

// WhenWindow satisfies the app.Window interface.
func (w *Window) WhenWindow(f func(app.Window)) {
	f(w)
}

// Load satisfies the app.Window interface.
func (w *Window) Load(rawurl string, v ...interface{}) {
	rawurl = fmt.Sprintf(rawurl, v...)
	compoName := CompoNameFromURLString(rawurl)

	if !w.Driver.Factory.IsCompoRegistered(compoName) {
		w.Driver.OpenDefaultBrowser(rawurl)
		return
	}

	if w.compo, w.err = w.Driver.Factory.NewCompo(compoName); w.err != nil {
		return
	}

	if rawurl != w.History.Current() {
		w.History.NewEntry(rawurl)
	}

	htmlConf := app.HTMLConfig{}
	if conf, ok := w.compo.(app.Configurator); ok {
		htmlConf = conf.Config()
	}

	if len(htmlConf.CSS) == 0 {
		htmlConf.CSS = file.Filenames(w.Driver.Resources("css"), ".css")
	}

	if len(htmlConf.Javascripts) == 0 {
		htmlConf.Javascripts = file.Filenames(w.Driver.Resources("js"), ".js")
	}

	page := dom.Page{
		Title:         htmlConf.Title,
		Metas:         htmlConf.Metas,
		CSS:           w.ConvertHTMLPaths(htmlConf.CSS),
		Javascripts:   w.ConvertHTMLPaths(htmlConf.Javascripts),
		GoRequest:     w.Driver.JSToPlatform,
		RootCompoName: compoName,
	}

	if w.err = w.Driver.Platform.Call("windows.Load", nil, struct {
		ID      string
		Title   string
		Page    string
		LoadURL string
		BaseURL string
	}{
		ID:      w.id,
		Title:   htmlConf.Title,
		Page:    page.String(),
		LoadURL: rawurl, // TODO: MacOS try to get rid.
		BaseURL: w.Driver.Resources(),
	}); w.err != nil {
		return
	}

	if w.err = w.DOM.New(w.compo); w.err != nil {
		return
	}

	if nav, ok := w.compo.(app.Navigable); ok {
		u, _ := url.Parse(rawurl)

		w.Driver.UI(func() {
			nav.OnNavigate(u)
		})
	}
}

// Reload satisfies the app.Window interface.
func (w *Window) Reload() {
	url := w.History.Current()
	if len(url) == 0 {
		w.err = errors.New("no component to reload")
		return
	}

	w.Load(url)
}

// CanPrevious satisfies the app.Window interface.
func (w *Window) CanPrevious() bool {
	return w.History.CanPrevious()
}

// Previous satisfies the app.Window interface.
func (w *Window) Previous() {
	url := w.History.Previous()
	if len(url) == 0 {
		w.err = errors.New("no previous component to load")
		return
	}

	w.Load(url)
}

// CanNext satisfies the app.Window interface.
func (w *Window) CanNext() bool {
	return w.History.CanNext()
}

// Next satisfies the app.Window interface.
func (w *Window) Next() {
	url := w.History.Next()
	if len(url) == 0 {
		w.err = errors.New("no next component to load")
		return
	}

	w.Load(url)
}

// Compo satisfies the app.Window interface.
func (w *Window) Compo() app.Compo {
	return w.compo
}

// Render satisfies the app.Window interface.
func (w *Window) Render(c app.Compo) {
	w.err = w.DOM.Render(c)
}

func (w *Window) render(changes interface{}) error {
	b, err := json.Marshal(changes)
	if err != nil {
		return errors.Wrap(err, "encoding changes failed")
	}

	return w.Driver.Platform.Call("windows.Render", nil, struct {
		ID      string
		Changes string
	}{
		ID:      w.id,
		Changes: string(b),
	})
}

// Position satisfies the app.Window interface.
func (w *Window) Position() (x, y float64) {
	return w.x, w.y
}

// Move satisfies the app.Window interface.
func (w *Window) Move(x, y float64) {
	w.err = w.Driver.Platform.Call("windows.Move", nil, struct {
		ID string
		X  float64
		Y  float64
	}{
		ID: w.id,
		X:  x,
		Y:  y,
	})
}

// Center satisfies the app.Window interface.
func (w *Window) Center() {
	w.err = w.Driver.Platform.Call("windows.Center", nil, struct {
		ID string
	}{
		ID: w.id,
	})
}

// SetPosition set the position state to the given values.
func (w *Window) SetPosition(x, y float64) {
	w.x = x
	w.y = y
	w.Driver.Events.Emit(app.WindowMoved, w)
}

// Size satisfies the app.Window interface.
func (w *Window) Size() (width, height float64) {
	return w.width, w.height
}

// Resize satisfies the app.Window interface.
func (w *Window) Resize(width, height float64) {
	w.err = w.Driver.Platform.Call("windows.Resize", nil, struct {
		ID     string
		Width  float64
		Height float64
	}{
		ID:     w.id,
		Width:  width,
		Height: height,
	})
}

// SetSize set the size state to the given values.
func (w *Window) SetSize(width, height float64) {
	w.width = width
	w.height = height
	w.Driver.Events.Emit(app.WindowResized, w)
}

// Focus satisfies the app.Window interface.
func (w *Window) Focus() {
	w.err = w.Driver.Platform.Call("windows.Focus", nil, struct {
		ID string
	}{
		ID: w.id,
	})
}

// IsFocus satisfies the app.Window interface.
func (w *Window) IsFocus() bool {
	return w.isFocus
}

// SetIsFocus set the focus state to the given value.
func (w *Window) SetIsFocus(v bool) {
	w.isFocus = v

	if w.isFocus {
		w.Driver.Events.Emit(app.WindowFocused, w)
	} else {
		w.Driver.Events.Emit(app.WindowBlurred, w)
	}
}

// FullScreen satisfies the app.Window interface.
func (w *Window) FullScreen() {
	w.err = nil

	if w.isFullScreen {
		return
	}

	w.err = w.Driver.Platform.Call("windows.SetFullScreen", nil, struct {
		ID     string
		Enable bool
	}{
		ID:     w.id,
		Enable: true,
	})
}

// ExitFullScreen satisfies the app.Window interface.
func (w *Window) ExitFullScreen() {
	w.err = nil

	if !w.isFullScreen {
		return
	}

	w.err = w.Driver.Platform.Call("windows.SetFullScreen", nil, struct {
		ID     string
		Enable bool
	}{
		ID:     w.id,
		Enable: false,
	})
}

// IsFullScreen satisfies the app.Window interface.
func (w *Window) IsFullScreen() bool {
	return w.isFullScreen
}

// SetIsFullScreen set the full screen state to the given value.
func (w *Window) SetIsFullScreen(v bool) {
	w.isFullScreen = v

	if w.isFullScreen {
		w.Driver.Events.Emit(app.WindowEnteredFullScreen, w)
	} else {
		w.Driver.Events.Emit(app.WindowExitedFullScreen, w)
	}
}

// Minimize satisfies the app.Window interface.
func (w *Window) Minimize() {
	w.err = nil

	if w.isMinimized {
		return
	}

	w.err = w.Driver.Platform.Call("windows.SetMinimize", nil, struct {
		ID     string
		Enable bool
	}{
		ID:     w.id,
		Enable: true,
	})
}

// Deminimize satisfies the app.Window interface.
func (w *Window) Deminimize() {
	w.err = nil

	if !w.isMinimized {
		return
	}

	w.err = w.Driver.Platform.Call("windows.SetMinimize", nil, struct {
		ID     string
		Enable bool
	}{
		ID:     w.id,
		Enable: false,
	})
}

// IsMinimized satisfies the app.Window interface.
func (w *Window) IsMinimized() bool {
	return w.isMinimized
}

// SetIsMinimized set the minimized state to the given value.
func (w *Window) SetIsMinimized(v bool) {
	w.isMinimized = v

	if w.isMinimized {
		w.Driver.Events.Emit(app.WindowMinimized, w)
	} else {
		w.Driver.Events.Emit(app.WindowDeminimized, w)
	}
}

// Close satisfies the app.Window interface.
func (w *Window) Close() {
	w.err = w.Driver.Platform.Call("windows.Close", nil, struct {
		ID string
	}{
		ID: w.id,
	})
}

// Release release the resources allocated to display the window.
func (w *Window) Release() {
	w.Driver.Events.Emit(app.WindowClosed, w)
	w.DOM.Close()
	w.Driver.Elems.Delete(w)
}
