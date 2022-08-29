//go:generate go run gen/html.go
//go:generate go run gen/scripts.go
//go:generate go fmt

// Package app is a package to build progressive web apps (PWA) with Go
// programming language and WebAssembly.
// It uses a declarative syntax that allows creating and dealing with HTML
// elements only by using Go, and without writing any HTML markup.
// The package also provides an http.Handler ready to serve all the required
// resources to run Go-based progressive web apps.
package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

const (
	// IsClient reports whether the code is running as a client in the
	// WebAssembly binary (app.wasm).
	IsClient = runtime.GOARCH == "wasm" && runtime.GOOS == "js"

	// IsServer reports whether the code is running on a server for
	// pre-rendering purposes.
	IsServer = runtime.GOARCH != "wasm" || runtime.GOOS != "js"

	orientationChangeDelay = time.Millisecond * 500
	engineUpdateRate       = 120
	resizeInterval         = time.Millisecond * 250
)

var (
	rootPrefix         string
	isInternalURL      func(string) bool
	appUpdateAvailable bool
	lastURLVisited     *url.URL
	resizeTimer        *time.Timer
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
//
//	 func main() {
//			// Define app routes.
//			app.Route("/", myComponent{})
//			app.Route("/other-page", myOtherComponent{})
//
//			// Run the application when on a web browser (only executed on client side).
//			app.RunWhenOnBrowser()
//
//			// Launch the server that serves the app (only executed on server side):
//			http.Handle("/", &app.Handler{Name: "My app"})
//			http.ListenAndServe(":8080", nil)
//	 }
func RunWhenOnBrowser() {
	if IsServer {
		return
	}

	defer func() {
		err := recover()
		displayLoadError(err)
		panic(err)
	}()

	rootPrefix = Getenv("GOAPP_ROOT_PREFIX")
	isInternalURL = internalURLChecker()
	staticResourcesResolver := newClientStaticResourceResolver(Getenv("GOAPP_STATIC_RESOURCES_URL"))

	disp := engine{
		FrameRate:              engineUpdateRate,
		LocalStorage:           newJSStorage("localStorage"),
		SessionStorage:         newJSStorage("sessionStorage"),
		StaticResourceResolver: staticResourcesResolver,
		ActionHandlers:         actionHandlers,
	}
	disp.Page = browserPage{dispatcher: &disp}
	disp.Body = newClientBody(&disp)
	disp.init()
	defer disp.Close()

	window.setBody(disp.Body)

	onAchorClick := FuncOf(onAchorClick(&disp))
	defer onAchorClick.Release()
	Window().Set("onclick", onAchorClick)

	onPopState := FuncOf(onPopState(&disp))
	defer onPopState.Release()
	Window().Set("onpopstate", onPopState)

	goappNav := FuncOf(goappNav(&disp))
	defer goappNav.Release()
	Window().Set("goappNav", goappNav)

	onAppUpdate := FuncOf(onAppUpdate(&disp))
	defer onAppUpdate.Release()
	Window().Set("goappOnUpdate", onAppUpdate)

	onAppInstallChange := FuncOf(onAppInstallChange(&disp))
	defer onAppInstallChange.Release()
	Window().Set("goappOnAppInstallChange", onAppInstallChange)

	closeAppResize := Window().AddEventListener("resize", onResize)
	defer closeAppResize()

	closeAppOrientationChange := Window().AddEventListener("orientationchange", onAppOrientationChange)
	defer closeAppOrientationChange()

	performNavigate(&disp, Window().URL(), false)
	disp.start(context.Background())
}

func displayLoadError(err any) {
	loadingLabel := Window().
		Get("document").
		Call("getElementById", "app-wasm-loader-label")
	if !loadingLabel.Truthy() {
		return
	}
	loadingLabel.setInnerText(fmt.Sprint(err))
}

func newClientStaticResourceResolver(staticResourceURL string) func(string) string {
	return func(path string) string {
		if isRemoteLocation(path) || !isStaticResourcePath(path) {
			return path
		}

		var b strings.Builder
		b.WriteString(staticResourceURL)
		b.WriteByte('/')
		b.WriteString(strings.TrimPrefix(path, "/"))
		return b.String()
	}
}

func internalURLChecker() func(string) bool {
	var urls []string
	json.Unmarshal([]byte(Getenv("GOAPP_INTERNAL_URLS")), &urls)

	return func(url string) bool {
		for _, u := range urls {
			if strings.HasPrefix(url, u) {
				return true
			}
		}
		return false
	}
}

func newClientBody(d Dispatcher) *htmlBody {
	ctx, cancel := context.WithCancel(context.Background())
	body := &htmlBody{
		htmlElement: htmlElement{
			tag:           "body",
			context:       ctx,
			contextCancel: cancel,
			dispatcher:    d,
			jsElement:     Window().Get("document").Get("body"),
		},
	}
	body.setSelf(body)

	ctx, cancel = context.WithCancel(context.Background())
	content := &htmlDiv{
		htmlElement: htmlElement{
			tag:           "div",
			context:       ctx,
			contextCancel: cancel,
			dispatcher:    d,
			jsElement:     body.JSValue().firstElementChild(),
		},
	}
	content.setSelf(content)
	content.setParent(body)

	body.children = append(body.children, content)
	return body
}

func onAchorClick(d Dispatcher) func(Value, []Value) any {
	return func(this Value, args []Value) any {
		event := Event{Value: args[0]}
		elem := event.Get("target")

		for {
			switch elem.Get("tagName").String() {
			case "A":
				if meta := event.Get("metaKey"); meta.Truthy() && meta.Bool() {
					return nil
				}

				if ctrl := event.Get("ctrlKey"); ctrl.Truthy() && ctrl.Bool() {
					return nil
				}

				if download := elem.Call("getAttribute", "download"); !download.IsNull() {
					return nil
				}

				event.PreventDefault()
				if href := elem.Get("href"); href.Truthy() {
					navigate(d, elem.Get("href").String())
				}
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

func onPopState(d Dispatcher) func(this Value, args []Value) any {
	return func(this Value, args []Value) any {
		d.Dispatch(Dispatch{
			Mode: Update,
			Function: func(ctx Context) {
				navigateTo(d, Window().URL(), false)
			},
		})
		return nil
	}
}

func goappNav(d Dispatcher) func(this Value, args []Value) any {
	return func(this Value, args []Value) any {
		navigate(d, args[0].String())
		return nil
	}
}

func navigate(d Dispatcher, rawURL string) {
	u, err := url.Parse(rawURL)
	if err != nil {
		Log(errors.New("navigating to URL failed").
			Tag("url", rawURL).
			Wrap(err))
		return
	}
	navigateTo(d, u, true)
}

func navigateTo(d Dispatcher, u *url.URL, updateHistory bool) {
	if IsServer {
		return
	}

	if isExternalNavigation(u) {
		if rawurl := u.String(); isInternalURL(rawurl) || isMailTo(u) {
			Window().Get("location").Set("href", u.String())
		} else {
			Window().Call("open", rawurl)
		}
		return
	}

	luv := lastURLVisited

	if u.String() == luv.String() {
		return
	}

	if u.Path == luv.Path && u.Fragment != luv.Fragment {
		if updateHistory {
			Window().addHistory(u)
		} else {
			lastURLVisited = u
		}

		d, ok := d.(ClientDispatcher)
		if !ok {
			return
		}
		d.Nav(u)

		if isFragmentNavigation(u) {
			d.Dispatch(Dispatch{
				Mode: Defer,
				Function: func(ctx Context) {
					Window().ScrollToID(u.Fragment)
				},
			})
		}
		return
	}

	performNavigate(d, u, updateHistory)
}

func performNavigate(d Dispatcher, u *url.URL, updateHistory bool) {
	if IsServer {
		return
	}

	path := strings.TrimPrefix(u.Path, rootPrefix)
	if path == "" {
		path = "/"
	}
	compo, ok := routes.createComponent(path)
	if !ok {
		compo = &notFound{}
	}

	disp, ok := d.(ClientDispatcher)
	if !ok {
		return
	}
	disp.Mount(compo)

	if updateHistory {
		Window().addHistory(u)
	} else {
		lastURLVisited = u
	}

	disp.Nav(u)
	if isFragmentNavigation(u) {
		d.Dispatch(Dispatch{
			Mode: Defer,
			Function: func(ctx Context) {
				Window().ScrollToID(u.Fragment)
			},
		})
	}
}

func isExternalNavigation(u *url.URL) bool {
	switch {
	case u.Host != "" && u.Host != Window().URL().Host,
		isMailTo(u):
		return true

	default:
		return false
	}
}

func isMailTo(u *url.URL) bool {
	return u.Scheme == "mailto"
}

func isFragmentNavigation(u *url.URL) bool {
	return u.Fragment != ""
}

func onAppUpdate(d ClientDispatcher) func(this Value, args []Value) any {
	return func(this Value, args []Value) any {
		d.Dispatch(Dispatch{
			Mode: Update,
			Function: func(ctx Context) {
				appUpdateAvailable = true
				d.AppUpdate()
				ctx.Defer(func(Context) {
					Log("app has been updated, reload to see changes")
				})
			},
		})
		return nil
	}
}

func onAppInstallChange(d ClientDispatcher) func(this Value, args []Value) any {
	return func(this Value, args []Value) any {
		d.AppInstallChange()
		return nil
	}
}

func onResize(ctx Context, e Event) {
	if resizeTimer != nil {
		resizeTimer.Stop()
		resizeTimer.Reset(resizeInterval)
		return
	}

	resizeTimer = time.AfterFunc(resizeInterval, func() {
		if d, ok := ctx.Dispatcher().(ClientDispatcher); ok {
			d.AppResize()
		}
	})
}

func onAppOrientationChange(ctx Context, e Event) {
	if d, ok := ctx.Dispatcher().(ClientDispatcher); ok {
		go func() {
			time.Sleep(orientationChangeDelay)
			d.AppResize()
		}()
	}
}
