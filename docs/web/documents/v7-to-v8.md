# V7 to V8 migration guide

Go-app V8 solves the [SEO](/seo) critical problem by providing server-side prerendering. Unfortunately, parts of the package had to be reworked, which resulted in breaking changes.

This document is here to help to migrate V7 to V8, by enumerating things that have changed.

## Import

Replace `v7` imports by `v8`.

```go
import (
	"github.com/maxence-charriere/go-app/v8/pkg/app"
)
```

## Build directives

V8 does not require build instructions anymore. Build instructions such as `// +build wasm` or `// +build !wasm` should be removed and their `main()` functions merged.

Here is how it is done for the hello example:

```go
// +build wasm

package main

import "github.com/maxence-charriere/go-app/v7/pkg/app"

func main() {
	app.Route("/", &hello{})
	app.Run()
}
```

and

```go
// +build !wasm

package main

import (
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v7/pkg/app"
)

func main() {
	http.Handle("/", &app.Handler{
		Name:        "Hello",
		Description: "An Hello World! example",
	})

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
```

should be merged into:

```go
package main

import (
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v8/pkg/app"
)

func main() {
	app.Route("/", &hello{})
	app.RunWhenOnBrowser()

	http.Handle("/", &app.Handler{
		Name:        "Hello",
		Description: "An Hello World! example",
	})

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
```

**See the [Getting Started](/start) section for the version with comments**.

## Routing

Since components are now also rendered on the server-side, calls to **[`Route()`](/reference#Route) and [`RouteWithRegexp()`](/reference#RouteWithRegexp) must also be in the server code**.

It should not be a problem if the `main()` functions were merged like in the [Build directives](#build-directives) section, but some may have it still separated in different packages or binaries. In that case, **don't forget to route components on the server-side too!**

Furthermore, `Route()` and `RouteWithRegexp()` calls now only register the type of a given [component](/components). When the routing is done, a fresh instance of the associated component is created. Initialization should be now with [component lifecycle](/components#lifecycle) interfaces.

## Package functions

| V7                                               | V8                                                                                                                                                              | Description                                                                                                                                                                                                                                                                                                                  |
| ------------------------------------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `func Run()`                                     | [`func RunWhenOnBrowser()`](/reference#RunWhenOnBrowser)                                                                                                        | Starting the app on the client-side can now be called in the same code as the server-side. Build instructions and server/client code separation are not required anymore. See the [Getting started](/start#code) article.                                                                                                    |
| `func Route(path string, node UI)`               | [`func Route(path string, c Composer)`](/reference#Route)                                                                                                       | A route can only be associated with a [component](/components) now. It now associates the type of the given component rather than an instance. When navigated on, a new component is created and initialized with its zero value. Initializing component shall now be done with the [Mounter](/reference#Mounter) interface. |
| `func RouteWithRegexp(pattern string, node UI)`  | [`func RouteWithRegexp(pattern string, c Composer)`](/reference#RouteWithRegexp)                                                                                | Same as [Route()](/reference#Route).                                                                                                                                                                                                                                                                                         |
| `func Navigate(rawurl string)`                   | - [`func (ctx Context) Navigate(rawURL string)`](/reference#Context.Navigate)<br>- [`func (ctx Context) NavigateTo(u *url.URL)`](/reference#Context.NavigateTo) | Navigating to another page is now a [Context](/reference#Context) method.                                                                                                                                                                                                                                                    |
| `func Reload()`                                  | [`func (ctx Context) Reload()`](/reference#Context.Reload)                                                                                                      | Reloading the current page is now a [Context](/reference#Context) method.                                                                                                                                                                                                                                                    |
| `func Dispatch(fn func())`                       | - [`func (ctx Context) Dispatch(fn func())`](/reference#Context.Dispatch)()<br> - [`func (c *Compo) Defer(fn func(Context))`](/reference#Compo.Defer)           | Executing a function on the UI goroutine is now done from a [Context](/reference#Context) or a [Component](http://localhost:7777/reference#Composer).                                                                                                                                                                        |
| `func StaticResource(path string) string`        | [`func (ctx Context) ResolveStaticResource(path string) string`](/reference#Context.ResolveStaticResource)                                                      | Resolving a [static resource](/static-resources) is now a [Context](/reference#Context) method.                                                                                                                                                                                                                              |
| `func NewContextMenu(menuItems ...MenuItemNode)` | **Removed**                                                                                                                                                     | Context menus have been removed. This may come back eventually under another form.                                                                                                                                                                                                                                           |
| `func Log(format string, v ...interface{})`      | [`func Logf(format string, v ...interface{})`](/reference#Logf)                                                                                                 | `Log()` has been renamed `Logf()` to match the [fmt](https://golang.org/pkg/fmt/) package. The new [`Log()`](/reference#Log) function now have a similar behavior as `fmt.Println()`.                                                                                                                                        |

## Component interfaces

| Interface                         | V7                         | V8                  |
| --------------------------------- | -------------------------- | ------------------- |
| [Navigator](/reference#Navigator) | `OnNav(Context, *url.URL)` | `OnNav(Context)`    |
| [Resizer](/reference#Resizer)     | `OnAppResize(Context)`     | `OnResize(Context)` |

## Resource provider

The [ResourceProvider](/reference#ResourceProvider) interface has been changed:

```go
type ResourceProvider interface {
	// Package returns the path where the package resources are located.
	Package() string

	// Static returns the path where the static resources directory (/web) is
	// located.
	Static() string

	// AppWASM returns the app.wasm file path.
	AppWASM() string
}
```

## Concurrency

Goroutines launched from components should now be created with [`Context.Async()`](/reference#Context.Async). See [concurrency](/concurrency#async) topic.

## Next

- [Getting started with v8](/start)
- [API reference](/reference)
