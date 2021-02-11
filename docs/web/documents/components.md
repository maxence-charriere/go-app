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

The appearance of the component is defined in the `Render()` method which is called by default when the component is mounted.

There are some scenarios where the appearance of a component can dynamically change, like when there is user input. In that case, we need to trigger a component rendering in order to see the modifications on the screen.

This can be done by calling the component [Update()](/reference#Composer) method.

```go
type hello struct {
	app.Compo

	Name string // Name field
}

func (h *hello) Render() app.UI {
	return app.Div().Body(
		app.H1().Body(
			app.Text("Hello "),
			app.Text(h.Name), // Name field is used to display who is greeted
		),
		app.Input().
			Value(h.Name). // Name field is used to display the current input value
			OnChange(h.OnInputChange),
	)
}

func (h *hello) OnInputChange(ctx app.Context, e app.Event) {
	h.Name = ctx.JSSrc.Get("value").String() // Name field is modified
	h.Update()                               // Render() is triggered
}
```

In the example above, `Update()` is called when the `input onchange` event is triggered, once the component `Name` field is modified with the input value.

### Update mechanism

When the component update is triggered, the `Render()` method is called and a new tree of UI element is generated. This new tree is then compared with the current component tree and only nonmatching nodes are modified or replaced.

Here are how the modifications are performed:

| Diff                                                       | Modification                              |
| ---------------------------------------------------------- | ----------------------------------------- |
| Different types of nodes (Text, HTML element or Component) | Current node is replaced                  |
| Different texts                                            | Current node text value is updated        |
| Different HTML elements                                    | Current node is replaced                  |
| Different HTML element attributes                          | Current node attributes are updated       |
| Different HTML element event handlers                      | Current node event handlers are updated   |
| Different component types                                  | Current node is replaced                  |
| Different component exported fields                        | Current component fields are updated      |
| Different component non exported fields                    | No modifications                          |
| Extra node in new the tree                                 | Node added to the current tree            |
| Missing node in the new tree                               | Extra node is the current tree is removed |

## Lifecycle

During its life, a component goes through several steps where actions could be performed to initialize or release data and resources.

![lifecycle](/web/images/lifecycle.png)

It is possible to trigger instructions when those different steps happen by implementing the corresponding interfaces in the component.

### OnMount

A component is mounted when it is inserted into the webpage DOM.

When the [Mounter](/reference#Mounter) interface is implemented, the `OnMount()` method is called right after the component is mounted.

```go
type foo struct {
    app.Compo
}

func (f *foo) OnMount(ctx app.Context) {
    fmt.Println("component mounted")
}
```

### OnNav

A component is navigated when a page is loaded, reloaded, or navigated from an anchor link or an HREF change.

When the [Navigator](/reference#Navigator) interface is implemented, the `OnNav()` method is called each time the component is navigated.

```go
type foo struct {
    app.Compo
}

func (f *foo) OnNav(ctx app.Context) {
    fmt.Println("component navigated:", u)
}
```

### OnDismount

A component is dismounted when it is removed from the webpage DOM.

When the [Dismounter](/reference#Dismounter) interface is implemented, the `OnDismount()` method is called right after the component is dismounted.

```go
type foo struct {
    app.Compo
}

func (f *foo) OnDismount() {
    fmt.Println("component dismounted")
}
```

## Next

- [Customize components with the declarative syntax](/syntax)
- [Associate components with URL paths](/routing)
- [API reference](/reference)
