# Components

Components are customizable, independent, and reusable UI elements.
They allow your UI to be split into independent and reusable pieces.

## Create

Creating a component is done by embedding [app.Compo](/reference#Compo) into a struct:

```go
type hello struct {
    app.Compo
}
```

## Customize

Once the component declared, the next thing to do is to customize its appearance.

This is done by implementing the [Render](/reference#Composer) method.

```go
func (h *hello) Render() app.UI {
	return app.H1().Text("Hello World!")
}
```

In the code above, the component is described as a simple [H1](/reference#H1) HTML tag that displays a `Hello World!` text.

`Render()` returns an [UI element](/reference#UI) that can be either an HTML element or another component. Refer to the [Declarative Syntax](/syntax) topic to know more about how to describe a component.

## Update

The `Render()` method defines the component appearance.

In the hello world example above, the rendering uses the `name` component field to define the title and the input default value.

It also set up an event handler that is called when the input change:

```go
func (h *hello) OnInputChange(ctx app.Context, e app.Event) {
    h.name = ctx.JSSrc().Get("value").String()
    h.Update()
}
```

At each change, the input value is assigned to the component `name` field.

Changing the value of a variable used in the `Render()` method does not update the UI.

The way to tell the browser that the component appearance has to be updated is by calling the [Compo.Update()](https://pkg.go.dev/github.com/maxence-charriere/go-app/v6/pkg/app#Compo.Update) method:

```go
h.Update()
```

Under the hood, the update method creates a new state that is compared to the current one. Then it performs a diff and updates only the HTML nodes where differences are found.

## Composing Components

Components can refer to other components in their `Render()` method.

Here is an example that shows `"Foo, Bar!"` by using a component that embeds another one:

```go
// foo component
type foo struct {
    app.Compo
}

func (f *foo) Render() app.UI {
    return app.P().Body(
        app.Text("Foo, "),
        &bar{},            // <-- bar component
    )
}

// bar component
type bar struct {
    app.Compo
}

func (b *bar) Render() app.UI {
    return app.Text("Bar!")
}

```

## Lifecycle

Components often use other resources to represent a UI. Those resources might require to be initialized or released.

The **go-app** package provides interfaces that allow calling functions at different times during the component lifecycle.

![lifecycle](https://storage.googleapis.com/murlok-github/lifecycle.png)

Implementing them in a component is a good place for initializing or free resources.

### OnMount

A component is mounted when it is inserted into the webpage DOM.

When the [Mounter](https://pkg.go.dev/github.com/maxence-charriere/go-app/v6/pkg/app#Mounter) interface is implemented, the `OnMount()` method is called right after the component is mounted.

```go
type foo struct {
    app.Compo
}

func (f *foo) OnMount(ctx app.Context) {
    fmt.Println("component mounted")
}
```

### OnNav

A component is navigated when a page where it is the body root is loaded, reloaded or navigated from an anchor link or an HREF change.

When the [Navigator](https://pkg.go.dev/github.com/maxence-charriere/go-app/v6/pkg/app#Navigator) interface is implemented, the `OnNav()` method is called each time the component is navigated.

```go
type foo struct {
    app.Compo
}

func (f *foo) OnNav(ctx app.Context, u *url.URL) {
    fmt.Println("component navigated:", u)
}
```

### OnDismount

A component is dismounted when it is removed from the webpage DOM.

When the [Dismounter](https://pkg.go.dev/github.com/maxence-charriere/go-app/v6/pkg/app#Dismounter) interface is implemented, the `OnDismount()` method is called right after the component is dismounted.

```go
type foo struct {
    app.Compo
}

func (f *foo) OnDismount() {
    fmt.Println("component dismounted")
}
```
