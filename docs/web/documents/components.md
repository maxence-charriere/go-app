# Components

Components are customizable, independent, and reusable UI elements.
They allow your UI to be split into independent and reusable pieces.

## Create

Creating a component is done by embedding [Compo](/reference#Compo) into a struct:

```go
type hello struct {
    app.Compo
}
```

## Customize

Customizing a component appearance is done by implementing the [Render()](/reference#Composer) method.

```go
func (h *hello) Render() app.UI {
	return app.H1().Text("Hello World!")
}
```

The code above displays an [H1](/reference#H1) HTML element that shows `Hello World!` text.

`Render()` returns an [UI element](/reference#UI) that can be either an HTML element or another component. Refer to the [Declarative Syntax](/syntax) topic to know more about how to describe a component.

## Update

In some scenarios, the component appearance can dynamically change.

Let's update the hello component to make it display the name of the user:

```go
type hello struct {
	app.Compo

	name string // Field where the username is stored
}

func (h *hello) Render() app.UI {
	return app.Div().Body(
		app.H1().Body(
			app.Text("Hello "),
			app.Text(h.name), // The name field used in the title
		),

		// The input HTML element that get the username.
		app.Input().
			Value(h.Name).             // The name field used as current input value
			OnChange(h.OnInputChange), // The event handler that will store the username
	)
}

func (h *hello) OnInputChange(ctx app.Context, e app.Event) {
	h.name = ctx.JSSrc.Get("value").String() // Name field is modified
	h.Update()                               // Update the component UI
}
```

The component now displays the username in its title and provides input for the user to type his/her name. When the user does so, an event handler is called and the name is stored in the component field named `name`.

The **[Update()](/reference#Composer) method call is what tells the component that its state changed and that its appearance must be updated**.

It internally triggers the `Render()` method and performs a diff with the current component state in order to define and process the changes. Here is how rendering diff behave:

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

![component lifecycle diagram](/web/images/compo-lifecycle.svg)

It is possible to trigger instructions when those different steps happen by implementing the corresponding interfaces in the component.

### Prerender

A component is prerendered when it is used on the server-side to generate HTML markup that is included in a requested HTML page, allowing search engines to index contents created with go-app.

Custom actions can be performed by implementing the [PreRenderer](/reference#PreRenderer) interface:

```go
type foo struct {
    app.Compo
}

func (f *foo) OnPreRender(ctx app.Context) {
    fmt.Println("component prerendered")
}
```

### Mount

A component is mounted when it is inserted into the webpage DOM.

Custom actions can be performed by implementing the [Mounter](/reference#Mounter) interface:

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

Custom actions can be performed by implementing the [Navigator](/reference#Navigator) interface:

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

Custom actions can be performed by implementing the [Dismounter](/reference#Dismounter) interface:

```go
type foo struct {
    app.Compo
}

func (f *foo) OnDismount() {
    fmt.Println("component dismounted")
}
```

## Extensions

Extensions are interfaces that when implemented, allow components to react to various events.

| Interface                             | Description                                      |
| ------------------------------------- | ------------------------------------------------ |
| [PreRenderer](/reference#PreRenderer) | Listen to component prerendering.                |
| [Mounter](/reference#Mounter)         | Listen to component mounting.                    |
| [Dismounter](/reference#Dismounter)   | Listen to component dismounting.                 |
| [Navigator](/reference#Navigator)     | Listen to page navigation.                       |
| [Updater](/reference#Updater)         | Listen to available app update.                  |
| [Resizer](/reference#Resizer)         | Listen to the app and parent components resizes. |

Here is an example where a component reacts to page navigation and updates the page title.

```go
type foo struct {
	app.Compo
}

func (f *foo) OnNav(ctx app.Context) {
	ctx.Page.SetTitle("Now the page is named Foo!")
}
```

## Next

- [Customize components with the declarative syntax](/syntax)
- [Associate components with URL paths](/routing)
- [API reference](/reference)
