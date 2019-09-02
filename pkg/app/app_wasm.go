package app

import (
	"context"
	"net/url"
	"reflect"
	"syscall/js"

	"github.com/maxence-charriere/app/pkg/log"
)

var (
	page    *dom
	cursorX int
	cursorY int
)

func init() {
	page = &dom{
		compoBuilder:        components,
		callOnUI:            UI,
		trackCursorPosition: trackCursorPosition,
		contextMenu:         &contextMenu{},
	}

	log.DefaultColor = ""
	log.InfoColor = ""
	log.ErrorColor = ""
	log.WarnColor = ""
	log.DebugColor = ""
	log.CurrentLevel = log.DebugLevel
}

func run() {
	go func() {
		defer page.clean()

		overrideAnchorClick := js.FuncOf(overrideAnchorClick)
		defer overrideAnchorClick.Release()
		js.Global().Set("onclick", overrideAnchorClick)

		onpopstate := js.FuncOf(onPopState)
		defer onpopstate.Release()
		js.Global().Set("onpopstate", onpopstate)

		url := getURL()

		if err := renderPage(url); err != nil {
			log.Error("rendering page failed").
				T("reason", err).
				T("url", url).
				Panic()
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

func render(c Compo) {
	if err := page.render(c); err != nil {
		log.Error("rendering component failed").
			T("reason", err).
			T("component", reflect.TypeOf(c))
	}
}

func reload() {
	js.Global().Get("location").Call("reload")
}

func bind(msg string, c Compo) *Binding {
	b, close := msgs.bind(msg, c)

	if err := page.setBindingClose(c, close); err != nil {
		log.Error("creating a binding failed").
			T("reason", err).
			T("component", reflect.TypeOf(c)).
			T("message", msg).
			Panic()
	}

	return b
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

func getURL() *url.URL {
	rawurl := js.Global().
		Get("location").
		Get("href").
		String()

	url, err := url.Parse(rawurl)
	if err != nil {
		log.Error("getting current url failed").
			T("reason", err).
			Panic()
	}

	return url
}

func renderPage(url *url.URL) error {
	if url.Path == "" || url.Path == "/" {
		url.Path = DefaultPath
	}
	if !components.isImported(compoNameFromURL(url)) {
		url.Path = NotFoundPath
	}

	compoName := compoNameFromURL(url)
	compo, err := components.new(compoName)
	if err != nil {
		return err
	}
	mapCompoFieldFromURLQuery(compo, url.Query())

	return page.newBody(compo)
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

// Navigate navigates to the given URL.
func Navigate(rawurl string) {
	UI(func() {
		navigate(rawurl, true)
	})
}

func navigate(rawurl string, updateHistory bool) {
	currentURL := getURL()

	u, err := url.Parse(rawurl)
	if err != nil {
		log.Error("navigating failed").
			T("reason", err).
			T("url", rawurl)
		return
	}

	if u.Host == "" && u.Scheme == "" {
		u.Scheme = currentURL.Scheme
		u.Host = currentURL.Host
	}

	fragmentNav := u.Host == currentURL.Host &&
		u.Path == currentURL.Path &&
		u.Fragment != currentURL.Fragment &&
		u.Fragment != ""

	otherHostNav := u.Host != currentURL.Host

	if otherHostNav || fragmentNav {
		js.Global().Get("location").Set("href", u.String())
		return
	}

	page.clean()
	if err := renderPage(u); err != nil {
		log.Error("rendering page failed").
			T("reason", err).
			T("url", u).
			Panic()
		return
	}

	if updateHistory {
		js.Global().Get("history").Call("pushState", nil, "", rawurl)
	}
}

func overrideAnchorClick(this js.Value, args []js.Value) interface{} {
	event := args[0]

	elem := event.Get("target")
	if !elem.Truthy() {
		elem = event.Get("srcElement")
	}

	if elem.Get("tagName").String() != "A" {
		return nil
	}

	event.Call("preventDefault")
	Navigate(elem.Get("href").String())
	return nil
}

func onPopState(this js.Value, args []js.Value) interface{} {
	UI(func() {
		rawurl := js.Global().Get("location").Get("href").String()
		navigate(rawurl, false)
	})
	return nil
}
