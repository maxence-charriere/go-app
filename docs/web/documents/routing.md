# Routing

Progressive web apps created with the **go-app** package are working as a [single page application](https://en.wikipedia.org/wiki/Single-page_application).

At first navigation, the app is loaded in the browser. Then, each time a page is requested, the navigation event is intercepted, then **go-app** routing system read the URL path and display the corresponding component.

![routing.png](/web/images/routing.png)

## Define a route

Defining a route is done by **associating a URL path with a given component**. Routes can be defined by using a simple pattern or by a regular expression.

### Simple route

Simple routes are when a component matches an exact URL path. They are defined with the [Route()](/reference#Route) function:

```go
func main() {
	app.Route("/", &hello{})  // hello component is associated with default path "/".
	app.Route("/foo", &foo{}) // foo component is associated with "/foo".
	app.Run()                 // Launches the app in the web browser.
}
```

### Route with regular expression

Routes with regular expressions are when a component matches an URL path with a given pattern. They are defined with the [RouteWithRegexp()](/reference#RouteWithRegexp)function:

```go
func main() {
	app.RouteWithRegexp("^/bar.*", &bar) // bar component is associated with all paths that start with /bar.
	app.Run()                            // Launches the app in the web browser.
}
```

Regular expressions follow [Go standard syntax](https://github.com/google/re2/wiki/Syntax).

## Detect navigation

Some scenarios may require actions to be done when a page is navigated on. Components can detect when a page is navigated on by implementing the [Navigator](/reference#Navigator) interface.

See [component lifecyle](/components#onnav).

## Next

- [Understand go-app architecture](/architecture)
- [How to create a component](/components)
- [API reference](/reference)
