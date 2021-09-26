## What is a component?

A component is a customizable, independent, and reusable UI element that allows your UI to be split into independent and reusable pieces.

## Create

Creating a component is done by embedding [Compo](/reference#Compo) into a struct:

```go
type hello struct {
    app.Compo
}
```

## Customize Look

Customizing a component look is done by implementing the [Render()](/reference#Composer) method.

```go
func (h *hello) Render() app.UI {
	return app.H1().Text("Hello World!")
}
```

The code above displays an [H1](/reference#H1) HTML element that shows the `Hello World!` text.

`Render()` returns a [UI element](/reference#UI) that can be either an HTML element or another component. Refer to the [Declarative Syntax](/declarative-syntax) topic to know more about how to customize a component look.

## Fields

Fields are struct fields that store data that can be used to customize a component when rendering. The example below shows a component that displays a name stored in a component field:

```go
type hello struct {
	app.Compo
	Name string // Exported field.
}

func (h *hello) Render() app.UI {
	return app.Div().Text("Hello, " + h.Name) // The Name field is display after "Hello, "
}
```

### Exported vs Unexported

In addition to the [Go distinction between exported and unexported fields](https://stackoverflow.com/questions/40256161/exported-and-unexported-fields-in-go-language), go-app uses that distinction to define whether a component needs to be updated.

When a UI element update is triggered (done internally), a UI element tree is rendered and compared to the currently displayed one. When 2 child components of the same type are compared to check differences, the comparison is based on the value of exported fields.

Here is a pseudo-Go code that illustrates how it works internally:

```go
type hello struct {
	app.Compo

	ExportedName   string
	unexportedName string
}

func updateFromExportedField() {
	current := &hello{
		ExportedName:   "Max",
		unexportedName: "Eric",
	}

	new := &hello{
		ExportedName:   "Maxence",
		unexportedName: "Erin",
	}

	update(app.Div().Body(current), app.Div().Body(new))

    // Current component exported field is updated:
	fmt.Println("current exported name:" + current.ExportedName)     // Updated:     "Maxence"
	fmt.Println("current unexported name:" + current.unexportedName) // Not Updated: "Eric"
}

func updateFromUnexportedField() {
	current := &hello{
		ExportedName:   "Max",
		unexportedName: "Eric",
	}

	new := &hello{
		ExportedName:   "Max",
		unexportedName: "Erin",
	}

	update(app.Div().Body(current), app.Div().Body(new))

	// Current component is not updated (no different exported field value):
	fmt.Println("current exported name:" + current.ExportedName)     // Not Updated: "Max"
	fmt.Println("current unexported name:" + current.unexportedName) // Not Updated: "Eric"
}
```

**Child components are updated only when there is diff with their exported fields**, and **only exported field are updated**.

### How chose between exported and unexported?

| Component field type | Triggers update | Value change | Usecase             |
| -------------------- | --------------- | ------------ | ------------------- |
| **Exported field**   | Yes             | Yes          | Component attribute |
| **Unexported field** | No              | No           | Component state     |

## Lifecycle Events

During its life, a component goes through several steps where actions can be performed to initialize or release data and resources.

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

### Nav

A component is navigated when a page is loaded, reloaded, or navigated from an anchor link or an HREF change. It can occur multiple times during a component life.

Custom actions can be performed by implementing the [Navigator](/reference#Navigator) interface:

```go
type foo struct {
    app.Compo
}

func (f *foo) OnNav(ctx app.Context) {
    fmt.Println("component navigated:", u)
}
```

### Dismount

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

### Reference

Here is a list of all the component lifecycle events available:

| Interface                               | Description                                               | Frequency                        |
| --------------------------------------- | --------------------------------------------------------- | -------------------------------- |
| [PreRenderer](/reference#PreRenderer)   | Listen to component prerendering.                         | Once on server-side              |
| [Mounter](/reference#Mounter)           | Listen to component mounting.                             | Once on client-side              |
| [Dismounter](/reference#Dismounter)     | Listen to component dismounting.                          | Once                             |
| [Navigator](/reference#Navigator)       | Listen to page navigation.                                | Once                             |
| [Updater](/reference#Updater)           | Listen to component update triggered by a parent element. | Can occur multiple times         |
| [AppUpdater](/reference#AppUpdater)     | Listen to available app update.                           | Can occur once                   |
| [AppInstaller](/reference#AppInstaller) | Listen to whether an app is installable.                  | Can occur once                   |
| [Resizer](/reference#Resizer)           | Listen to the app and parent components resizes.          | Each time a component is resized |

## Updates

Components are meant to be responsive to different events, modifying their appearance when they occur.

When this is happening, go-app internally starts an update mechanism that checks modifications in the currently displayed UI element tree and, performs the necessary modifications to achieve the desired state.

**This update mechanism is automatically trigerred when the following scenario occurs:**

- [Component lifecycle events](#lifecycle-events-reference)
- [HTML event handlers](/declarative-syntax#event-handlers)
- [Context.Dispatch](/reference#Context.Dispatch)
- [Context.Handle](/reference#Context.Handler)
- [Context.ObserveState](/reference#Context.ObserveState)

### Manually Trigger an Update

In the event where it is not automatically triggered with your use case, the component update mechanism can be manually launched by using [Compo.Update](/reference#Compo.Update).

```go
type myCompo struct {
	app.Compo

	Number int
}

func (c *myCompo) Render() app.UI {
	return app.Div().Text(c.Number)
}

func (c *myCompo) customTrigger() {
	c.Number = rand.Intn(42)
	c.Update() // Manual updated trigger
}
```

## Next

- [Customize components with the declarative syntax](/declarative-syntax)
- [Reference](/reference)
