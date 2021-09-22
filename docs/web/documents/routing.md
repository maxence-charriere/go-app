## Intro

Routing is about **associating a component with an URL path**.

## Define a route

Defining a route is done by **associating a URL path with a given component type**.

When a page is requested, its URL path is compared with the defined routes. Then **a new instance of the component type associated with the route is created and displayed**.

Routes are defined by using a simple pattern or by a regular expression.

### Simple route

Simple routes are when a component type matches an exact URL path. They are defined with the [Route()](/reference#Route) function:

```go
func main() {
	app.Route("/", &hello{})  // hello component type is associated with default path "/".
	app.Route("/foo", &foo{}) // foo component type is associated with "/foo".
	app.RunWhenOnBrowser()    // Launches the app when in a web browser.
}
```

### Route with regular expression

Routes with regular expressions are when a component type matches an URL path with a given pattern. They are defined with the [RouteWithRegexp()](/reference#RouteWithRegexp)function:

```go
func main() {
	app.RouteWithRegexp("^/bar.*", &bar) // bar component is associated with all paths that start with /bar.
	app.RunWhenOnBrowser()               // Launches the app when in a web browser.
}
```

Regular expressions follow [Go standard syntax](https://github.com/google/re2/wiki/Syntax).

## How it works?

Progressive web apps created with the **go-app** package are working as a [single page application](https://en.wikipedia.org/wiki/Single-page_application). At first navigation, the app is loaded in the browser. Once loaded, each time a page is requested, the navigation event is intercepted and **go-app**'s routing mechanism reads the URL path, then loads a new instance of the associated [component](/components).

![routing.png](/web/images/routing.svg)

## Detect navigation

Some scenarios may require additional actions to be done when a page is navigated on. Components can detect when a page is navigated on by implementing the [Navigator](/reference#Navigator) interface:

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
