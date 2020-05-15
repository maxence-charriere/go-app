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
	currentURL  url.URL
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

	for {
		switch elem.Get("tagName").String() {
		case "A":
			event.PreventDefault()
			url = elem.Get("href").String()
			Navigate(url)
			return nil

		case "BODY":
			return nil

		default:
			elem = elem.Get("parentElement")
			if !elem.Truthy() {
				return nil
			}
		}
	}

	return nil

}

func onPopState(this Value, args []Value) interface{} {
	dispatcher(func() {
		navigate(Window().URL(), false)
	})
	return nil
}

func navigate(u *url.URL, updateHistory bool) error {
	contextMenu.hide(nil, Event{Value: Null()})

	if isExternalNavigation(u) {
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

	if content != root {
		if err := replace(content, root); err != nil {
			return err
		}
		content = root
	}

	currentURL = *u
	triggerOnNav(root, u)

	if updateHistory {
		Window().Get("history").Call("pushState", nil, "", u.String())
	}

	if isFragmentNavigation(u) {
		dispatcher(func() {
			Window().ScrollToID(u.Fragment)
		})
	}

	return nil
}

func isExternalNavigation(u *url.URL) bool {
	return u.Host != Window().URL().Host
}

func isFragmentNavigation(u *url.URL) bool {
	return u.Fragment != ""
}

func reload() {
	Window().Get("location").Call("reload")
}

func newContextMenu(menuItems ...MenuItemNode) {
	contextMenu.show(menuItems...)
}

func getenv(k string) string {
	env := Window().Call("goappGetenv", k)
	if !env.Truthy() {
		return ""
	}
	return env.String()
}

func keepBodyClean() func() {
	close := Window().Call("goappKeepBodyClean")

	return func() {
		close.Invoke()
	}
}
