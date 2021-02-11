//go:generate go run gen/html.go
//go:generate go run gen/scripts.go
//go:generate go fmt

// Package app is a package to build progressive web apps (PWA) with Go
// programming language and WebAssembly.
// It uses a declarative syntax that allows creating and dealing with HTML
// elements only by using Go, and without writing any HTML markup.
// The package also provides an http.handler ready to serve all the required
// resources to run Go-based progressive web apps.
package app

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/maxence-charriere/go-app/v8/pkg/errors"
)

const (
	// IsClient reports whether the code is running as a client in the
	// WebAssembly binary (app.wasm).
	IsClient = runtime.GOARCH == "wasm" && runtime.GOOS == "js"

	// IsServer reports whether the code is running on a server for
	// pre-rendering purposes.
	IsServer = runtime.GOARCH != "wasm" || runtime.GOOS != "js"

	orientationChangeDelay = time.Millisecond * 500
)

var (
	staticResourcesURL string
	rootPrefix         string
	appUpdateAvailable bool
)

// Getenv retrieves the value of the environment variable named by the key. It
// returns the value, which will be empty if the variable is not present.
func Getenv(k string) string {
	if IsServer {
		return os.Getenv(k)
	}

	env := Window().Call("goappGetenv", k)
	if !env.Truthy() {
		return ""
	}
	return env.String()
}

// KeepBodyClean prevents third-party Javascript libraries to add nodes to the
// body element.
func KeepBodyClean() (close func()) {
	if IsServer {
		return func() {}
	}

	release := Window().Call("goappKeepBodyClean")
	return func() {
		release.Invoke()
	}
}

// StaticResource makes a static resource path point to the right
// location whether the root directory is remote or not.
//
// Static resources are resources located in the web directory.
//
// This call is used internally to resolve paths within Cite, Data, Href, Src,
// and SrcSet. Paths already resolved are skipped.
func StaticResource(path string) string {
	if !strings.HasPrefix(path, "/web/") &&
		!strings.HasPrefix(path, "web/") {
		return path
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return staticResourcesURL + path
}

// Window returns the JavaScript "window" object.
func Window() BrowserWindow {
	return window
}

// RunWhenOnBrowser starts the app, displaying the component associated with the
// current URL path.
//
// This call is skipped when the program is not run on a web browser. This
// allows writing client and server-side code without separation or
// pre-compilation flags.
//
// Eg:
//  func main() {
// 		// Define app routes.
// 		app.Route("/", myComponent{})
// 		app.Route("/other-page", myOtherComponent{})
//
// 		// Run the application when on a web browser (only executed on client side).
// 		app.RunWhenOnBrowser()
//
// 		// Launch the server that serves the app (only executed on server side):
// 		http.Handle("/", &app.Handler{Name: "My app"})
// 		http.ListenAndServe(":8080", nil)
//  }
func RunWhenOnBrowser() {
	if IsServer {
		return
	}

	defer func() {
		err := recover()
		displayLoadError(err)
		panic(err)
	}()

	staticResourcesURL = Getenv("GOAPP_STATIC_RESOURCES_URL")
	rootPrefix = Getenv("GOAPP_ROOT_PREFIX")

	disp := newUIDispatcher(IsServer)
	defer disp.Close()
	disp.body = newClientBody(disp)
	window.setBody(disp.body)

	onAchorClick := FuncOf(onAchorClick(disp))
	defer onAchorClick.Release()
	Window().Set("onclick", onAchorClick)

	onPopState := FuncOf(onPopState(disp))
	defer onPopState.Release()
	Window().Set("onpopstate", onPopState)

	onAppUpdate := FuncOf(onAppUpdate(disp))
	defer onAppUpdate.Release()
	Window().Set("goappOnUpdate", onAppUpdate)

	closeAppResize := Window().AddEventListener("resize", onAppResize)
	defer closeAppResize()

	closeAppOrientationChange := Window().AddEventListener("orientationchange", onAppOrientationChange)
	defer closeAppOrientationChange()

	performNavigate(disp, Window().URL(), false)
	disp.start(context.Background())
}

func displayLoadError(err interface{}) {
	loadingLabel := Window().
		Get("document").
		Call("getElementById", "app-wasm-loader-label")
	if !loadingLabel.Truthy() {
		return
	}
	loadingLabel.setInnerText(fmt.Sprint(err))
}

func newClientBody(d Dispatcher) *htmlBody {
	ctx, cancel := context.WithCancel(context.Background())
	body := &htmlBody{
		elem: elem{
			ctx:       ctx,
			ctxCancel: cancel,
			jsvalue:   Window().Get("document").Get("body"),
			tag:       "body",
			disp:      d,
		},
	}
	body.setSelf(body)

	ctx, cancel = context.WithCancel(context.Background())
	content := &htmlDiv{
		elem: elem{
			ctx:       ctx,
			ctxCancel: cancel,
			jsvalue:   body.JSValue().firstElementChild(),
			tag:       "div",
			disp:      d,
		},
	}
	content.setSelf(content)
	content.setParent(body)

	body.body = append(body.body, content)
	return body
}

func onAchorClick(d *uiDispatcher) func(Value, []Value) interface{} {
	return func(this Value, args []Value) interface{} {
		event := Event{Value: args[0]}
		elem := event.Get("target")

		for {
			switch elem.Get("tagName").String() {
			case "A":
				if target := elem.Get("target"); target.Truthy() && target.String() == "_blank" {
					return nil
				}

				u := elem.Get("href").String()
				if u, _ := url.Parse(u); isExternalNavigation(u) {
					elem.Set("target", "_blank")
					return nil
				}

				if meta := event.Get("metaKey"); meta.Truthy() && meta.Bool() {
					return nil
				}

				if ctrl := event.Get("ctrlKey"); ctrl.Truthy() && ctrl.Bool() {
					return nil
				}

				event.PreventDefault()
				navigate(d, u)
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
	}
}

func onPopState(d Dispatcher) func(this Value, args []Value) interface{} {
	return func(this Value, args []Value) interface{} {
		d.Dispatch(func() {
			performNavigate(d, Window().URL(), false)
		})
		return nil
	}
}

func navigate(d Dispatcher, rawURL string) {
	u, err := url.Parse(rawURL)
	if err != nil {
		Log("%s", errors.New("navigating to URL failed").
			Tag("url", rawURL).
			Wrap(err))
		return
	}
	navigateTo(d, u)
}

func navigateTo(d Dispatcher, u *url.URL) {
	if u.String() == Window().URL().String() {
		fmt.Println("same url")
		return
	}
	performNavigate(d, u, true)
}

func performNavigate(d Dispatcher, u *url.URL, updateHistory bool) {
	if IsServer {
		return
	}

	if isExternalNavigation(u) {
		Window().Get("location").Set("href", u.String())
		return
	}

	compo, ok := routes.createComponent(strings.TrimPrefix(u.Path, rootPrefix))
	if !ok {
		compo = &notFound{}
	}

	disp := d.(*uiDispatcher)
	disp.Mount(compo)

	if updateHistory {
		Window().Get("history").Call("pushState", nil, "", u.String())
	}

	disp.Nav(u)
	if isFragmentNavigation(u) {
		disp.Dispatch(func() {
			Window().ScrollToID(u.Fragment)
		})
	}
}

func isExternalNavigation(u *url.URL) bool {
	return u.Host != "" && u.Host != Window().URL().Host
}

func isFragmentNavigation(u *url.URL) bool {
	return u.Fragment != ""
}

func onAppUpdate(d *uiDispatcher) func(this Value, args []Value) interface{} {
	return func(this Value, args []Value) interface{} {
		d.Dispatch(func() {
			appUpdateAvailable = true
		})
		d.AppUpdate()
		d.Dispatch(func() {
			fmt.Println("app has been updated, reload to see changes")
		})
		return nil
	}
}

func onAppResize(ctx Context, e Event) {
	if d, ok := ctx.dispatcher.(*uiDispatcher); ok {
		d.AppResize()
	}
}

func onAppOrientationChange(ctx Context, e Event) {
	if d, ok := ctx.dispatcher.(*uiDispatcher); ok {
		go func() {
			time.Sleep(orientationChangeDelay)
			d.AppResize()
		}()
	}
}
