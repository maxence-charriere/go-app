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
	dom     *maestro.Dom
	msgs    *messenger
	cursorX int
	cursorY int
)

func init() {
	dom = &maestro.Dom{
		CompoBuilder:        components,
		CallOnUI:            UI,
		TrackCursorPosition: trackCursorPosition,
		ContextMenu:         &ContextMenu{},
	}
	msgs = &messenger{}

	log.DefaultColor = ""
	log.InfoColor = ""
	log.ErrorColor = ""
	log.WarnColor = ""
	log.DebugColor = ""
}

func navigate(url string) {
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

func render(c Compo) error {
	return dom.Render(c)
}

func run() error {
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
	return nil
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
