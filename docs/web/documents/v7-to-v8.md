# V7 to V8 migration guide

Go-app V8 solves the [SEO](/seo) critical problem by providing server-side prerendering. Unfortunately, parts of the package had to be reworked, which resulted in breaking changes.

This document is here to help to migrate V7 to V8, by enumerating things that have changed.

## Import

Replace V7 imports by V8.

```go
import (
	"github.com/maxence-charriere/go-app/v8/pkg/app"
)
```

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
