<!-- wiki:ignore -->

## Intro

Routing is about **associating a component with an URL path**.

## Define a route

Defining a route is done by **associating a URL path with a function that returns a component**. The requested URL path is matched and the function is used to create an instance of the component to display. Paths can be defined with a simple pattern or a regular expression.

### Simple route

Simple routes are when the requested URL path matches a given one. They are defined with the [Route()](/reference#Route) function:

```go
func main() {
	app.Route("/", func() app.Composer { return &component{} })    // component is created for the root path
	app.Route("/foo", func() app.Composer { return &component{} }) // component is created when the path is /foo
	app.RunWhenOnBrowser()                                         // Launches the app when in a web browser
}
```

### Route with regular expression

Routes with regular expressions are used when the requested URL path matches a given pattern. They are defined using the [RouteWithRegexp()](/reference#RouteWithRegexp) function:

```go
func main() {
	app.RouteWithRegexp("^/bar/(code|tender).*", func() app.Composer { return &component{} }) // component is created when the path is /bar/code or /bar/tender
	app.RunWhenOnBrowser()                                                                    // Launches the app when in a web browser.
}
```

Regular expressions follow the [Go standard syntax](https://github.com/google/re2/wiki/Syntax).

## How it works?

Progressive web apps created with the **go-app** package function as a [single-page application](https://en.wikipedia.org/wiki/Single-page_application). On the first navigation, the app is loaded in the browser. Once loaded, each time a page is requested, the navigation event is intercepted, and **go-app**'s routing mechanism reads the URL path, then loads the [component](/components) returned by the associated function.

![routing.png](/web/images/routing.svg)

## Detect navigation

Some scenarios may require additional actions when a page is navigated. Components can detect when a page is navigated by implementing the [Navigator](/reference#Navigator) interface:

```go
type foo struct {
    app.Compo
}

func (f *foo) OnNav(ctx app.Context) {
    fmt.Println("component navigated:", u)
}
```

See [component lifecycle](/components#nav).

## Next

- [Images and Static Resources](/static-resources)
- [Reference](/reference)
