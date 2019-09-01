package app

import (
	"context"
	"net/url"
	"reflect"
	"syscall/js"

	"github.com/maxence-charriere/app/internal/maestro"
	"github.com/maxence-charriere/app/pkg/log"
)

var (
	// DefaultPath is the path to the component to be  loaded when no path is
	// specified.
	DefaultPath string

	// NotFoundPath is the path to the component to be  loaded when an non
	// imported component is requested.
	NotFoundPath = "/app.notfound"

	ui         = make(chan func(), 256)
	components = make(maestro.CompoBuilder)
	msgs       = &messenger{}
	dom        *maestro.Dom
	cursorX    int
	cursorY    int
)

func init() {
	dom = &maestro.Dom{
		CompoBuilder:        components,
		CallOnUI:            UI,
		TrackCursorPosition: trackCursorPosition,
		ContextMenu:         &ContextMenu{},
	}

	log.DefaultColor = ""
	log.InfoColor = ""
	log.ErrorColor = ""
	log.WarnColor = ""
	log.DebugColor = ""
}

// Import imports the given components into the app.
// Components must be imported in order the be used by the app package.
// This allows components to be created dynamically when they are found into
// markup.
func Import(c ...Compo) {
	for _, compo := range c {
		if err := components.Import(compo); err != nil {
			panic(err)
		}
	}
}

// Run runs the app with the loaded URL.
func Run() {
	go func() {
		url, err := getURL()
		if err != nil {
			return
		}

		compo, err := components.New(compoNameFromURL(url))
		if err != nil {
			log.Error("creating component failed").
				T("reason", err).
				T("url", url.String())
			return
		}

		if err := dom.NewBody(compo); err != nil {
			log.Error("creating page failed").
				T("reason", err).
				T("url", url.String()).
				T("component", reflect.TypeOf(compo))
			return
		}

		if nav, ok := compo.(Navigable); ok {
			UI(func() {
				nav.OnNavigate(url)
			})
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		for {
			select {
			case <-ctx.Done():
				return

			case f := <-ui:
				f()
			}
		}
	}()

	select {}
	log.Info("wasm exit")
	return
}

// Render renders the given component. It should be called whenever a component
// is modified. Render is always excecuted on the UI goroutine.
//
// It panics if called before Run.
func Render(c Compo) {
	if err := dom.Render(c); err != nil {
		log.Error("rendering component failed").
			T("reason", err).
			T("component", reflect.TypeOf(c))
	}
}

// UI calls a function on the UI goroutine.
func UI(f func()) {
	ui <- f
}

// Navigate navigates to the given URL.
func Navigate(url string) {
	js.Global().Get("location").Set("href", url)
}

// Reload reloads the current page.
func Reload(s, e js.Value) {
	js.Global().Get("location").Call("reload")

}

// Bind creates a binding between a message and the given component.
func Bind(msg string, c Compo) *Binding {
	b, close := msgs.bind(msg, c)

	if err := dom.SetBindingClose(c, close); err != nil {
		log.Error("creating a binding failed").
			T("reason", err).
			T("component", reflect.TypeOf(c)).
			T("message", msg).
			Panic()
	}

	return b
}

// Emit emits a message that triggers the associated bindings.
func Emit(msg string, args ...interface{}) {
	go msgs.emit(msg, args...)
}

// NewContextMenu displays a context menu filled with the given menu items.
func NewContextMenu(items ...MenuItem) {
	Emit("__app.NewContextMenu", items)
}

// MenuItem represents a menu item.
type MenuItem struct {
	Disabled  bool
	Keys      string
	Icon      string
	Label     string
	OnClick   func(s, e js.Value)
	Separator bool
}

func getURL() (*url.URL, error) {
	rawurl := js.Global().
		Get("location").
		Get("href").
		String()

	url, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	if url.Path == "" || url.Path == "/" {
		url.Path = DefaultPath
	}
	if !components.IsImported(compoNameFromURL(url)) {
		url.Path = NotFoundPath
	}
	return url, nil
}

func trackCursorPosition(e js.Value) {
	x := e.Get("clientX")
	if !x.Truthy() {
		return
	}
	cursorX = x.Int()

	y := e.Get("clientY")
	if !y.Truthy() {
		return
	}
	cursorY = y.Int()
}
