package app

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"syscall/js"

	"github.com/maxence-charriere/go-app/v7/pkg/errors"
)

var (
	body        *htmlBody
	content     UI
	contextMenu = &contextMenuLayout{}
	rootPrefix  string
	window      = &browserWindow{value: value{Value: js.Global()}}
)

func run() {
	defer func() {
		err := recover()
		displayLoadError(err)
		panic(err)
	}()

	staticResourcesURL = Getenv("GOAPP_STATIC_RESOURCES_URL")
	rootPrefix = Getenv("GOAPP_ROOT_PREFIX")

	initBody()
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
		panic(errors.New("navigating to page failed").
			Tag("url", url).
			Wrap(err),
		)
	}

	for {
		select {
		case f := <-uiChan:
			f()
		}
	}
}

func initBody() {
	ctx, cancel := context.WithCancel(context.Background())

	body = &htmlBody{
		elem: elem{
			ctx:       ctx,
			ctxCancel: cancel,
			jsvalue:   Window().Get("document").Get("body"),
			tag:       "body",
		},
	}

	body.setSelf(body)
}

func initContent() {
	ctx, cancel := context.WithCancel(context.Background())

	content := &htmlDiv{
		elem: elem{
			ctx:       ctx,
			ctxCancel: cancel,
			jsvalue:   body.JSValue().Get("firstElementChild"),
			tag:       "div",
		},
	}

	content.setSelf(content)
	content.setParent(body)
	body.body = append(body.body, content)
}

func initContextMenu() {
	ctx, cancel := context.WithCancel(context.Background())

	tmp := &htmlDiv{
		elem: elem{
			attrs:     map[string]string{"id": "app-context-menu"},
			ctx:       ctx,
			ctxCancel: cancel,
			jsvalue: Window().
				Get("document").
				Call("getElementById", "app-context-menu"),
			tag: "div",
		},
	}

	tmp.setSelf(tmp)
	tmp.setParent(body)
	body.body = append(body.body, tmp)

	body.replaceChildAt(1, contextMenu)
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
	dispatch(func() {
		navigate(Window().URL(), false)
	})
	return nil
}

func navigate(u *url.URL, updateHistory bool) error {
	contextMenu.hide()

	if isExternalNavigation(u) {
		Window().Get("location").Set("href", u.String())
		return nil
	}

	path := u.Path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	root, ok := routes.ui(strings.TrimPrefix(u.Path, rootPrefix))
	if !ok {
		root = NotFound
	}

	if content != root {
		if err := body.replaceChildAt(0, root); err != nil {
			return errors.New("replacing content failed").Wrap(err)
		}
		content = root
	}

	if updateHistory {
		Window().Get("history").Call("pushState", nil, "", u.String())
	}

	if isFragmentNavigation(u) {
		dispatch(func() {
			Window().ScrollToID(u.Fragment)
		})
	}

	root.onNav(u)
	return nil
}

func isExternalNavigation(u *url.URL) bool {
	return u.Host != "" && u.Host != Window().URL().Host
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
