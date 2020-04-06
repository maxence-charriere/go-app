package app

import (
	"fmt"
	"net/url"
	"strings"
	"syscall/js"

	"github.com/maxence-charriere/go-app/v6/pkg/log"
)

var (
	window         = &browserWindow{value: value{Value: js.Global()}}
	body           = Body()
	content     UI = Div()
	contextMenu    = &contextMenuLayout{}
)

func init() {
	log.DefaultColor = ""
	log.InfoColor = ""
	log.ErrorColor = ""
	log.WarnColor = ""
	log.DebugColor = ""
	log.CurrentLevel = log.DebugLevel

	LocalStorage = newJSStorage("localStorage")
	SessionStorage = newJSStorage("sessionStorage")
}

func run() {
	defer func() {
		err := recover()
		displayLoadError(err)
		panic(err)
	}()

	initRemoteRootDir()
	initContent()
	initContextMenu()

	onnav := FuncOf(onNavigate)
	defer onnav.Release()
	Window().Set("onclick", onnav)

	onpopstate := FuncOf(onPopState)
	defer onpopstate.Release()
	Window().Set("onpopstate", onpopstate)

	url := Window().URL()

	if err := navigate(url, false); err != nil {
		log.Error("loading page failed").
			T("error", err).
			T("url", url).
			Panic()
	}

	for {
		select {
		case f := <-uiChan:
			f()
		}
	}
}

func displayLoadError(err interface{}) {
	loadingLabel := Window().
		Get("document").
		Call("getElementById", "app-wasm-loader-label")
	if !loadingLabel.Truthy() {
		return
	}
	loadingLabel.Set("innerText", fmt.Sprint(err))
}

func initRemoteRootDir() {
	remoteRootDir = Getenv("GOAPP_REMOTE_ROOT_DIR")
}

func initContent() {
	body.(*htmlBody).value = Window().Get("document").Get("body")
	content.(*htmlDiv).value = body.JSValue().Get("firstElementChild")
	content.setParent(body)
	body.appendChild(content)
}

func initContextMenu() {
	rawContextMenu := Div().ID("app-context-menu")
	rawContextMenu.(*htmlDiv).value = Window().
		Get("document").
		Call("getElementById", "app-context-menu")
	rawContextMenu.setParent(body)
	body.appendChild(rawContextMenu)

	if err := update(rawContextMenu, contextMenu); err != nil {
		log.Error("initializing context menu failed").
			T("error", err).
			Panic()
	}
}

func onNavigate(this Value, args []Value) interface{} {
	url := ""
	event := Event{Value: args[0]}

	elem := event.Get("target")
	if !elem.Truthy() {
		elem = event.Get("srcElement")
	}

findAnchor:
	for {
		switch elem.Get("tagName").String() {
		case "A":
			url = elem.Get("href").String()
			break findAnchor

		case "BODY":
			return nil

		default:
			elem = elem.Get("parentElement")
			if !elem.Truthy() {
				return nil
			}
		}
	}

	event.PreventDefault()
	Navigate(url)
	return nil

}

func onPopState(this Value, args []Value) interface{} {
	if u := Window().URL(); u.Fragment == "" {
		dispatcher(func() {
			navigate(u, false)
		})
	}
	return nil
}

func navigate(u *url.URL, updateHistory bool) error {
	contextMenu.hide(nil, Event{Value: Null()})

	if !isPWANavigation(u) {
		Window().Get("location").Set("href", u.String())
		return nil
	}

	path := u.Path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	root, ok := routes.ui(path)
	if !ok {
		root = NotFound
	}

	defer func() {
		if nav, ok := root.(Navigator); ok {
			nav.OnNav(u)
		}

		if updateHistory {
			Window().Get("history").Call("pushState", nil, "", u.String())
		}
	}()

	if content == root {
		return nil
	}
	if err := replace(content, root); err != nil {
		return err
	}
	content = root

	return nil
}

func isPWANavigation(u *url.URL) bool {
	externalNav := u.Host != "" && u.Host != Window().URL().Host
	fragmentNav := u.Fragment != ""
	return !externalNav && !fragmentNav
}

func reload() {
	Window().Get("location").Call("reload")
}

func newContextMenu(menuItems ...MenuItemNode) {
	contextMenu.show(menuItems...)
}

func getenv(k string) string {
	goappEnv := Window().Get("goappEnv")
	if !goappEnv.Truthy() {
		log.Error("goappEnv not found")
		return ""
	}
	return goappEnv.Get(k).String()
}
