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

In addition to being visible outside a Go package, **exported fields are triggers that says to go-app if a component needs to have its display updated**.

Here is a modified version of the hello component that modifies its `Name` field depending on user input:

```go
type hello struct {
	app.Compo
	Name string // Exported field.
}

func (h *hello) Render() app.UI {
	return app.Div().Body(
		app.H1().Text("Hello, "+h.Name), // Exported field used as tilte value.
		app.Input().
			Type("text").
			Value(h.Name). // Exported field used a input value.
			OnChange(h.onChange),
	)
}

func (h *hello) onChange(ctx app.Context, e app.Event) {
	h.Name = e.Get("value").String() // Name field is modified.
}
```

The `Name` field is modified with the changed input value. After the `onChange` function returns, go-app renders a new UI element tree that is compared with the currently displayed one. When a different name is found, the currently displayed tree is updated with the new value:

```sh
0 -> input value: ""        => initial state
1 -> input value: "Max"     => component is updated     ("Max" != "")
2 -> input value: "Maxence" => component is updated     ("Maxence" != "Max")
3 -> input value: "Maxence" => component is not updated ("Maxence" == "Maxence")
```

**Exported values are usually used as an equivalent of HTML attributes but for components**.

### Unexported fields

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
			Value(h.name).             // The name field used as current input value
			OnChange(h.OnInputChange), // The event handler that will store the username
	)
}

func (h *hello) OnInputChange(ctx app.Context, e app.Event) {
	h.name = ctx.JSSrc.Get("value").String() // Name field is modified
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
| Different exported fields on a same component type         | Current component fields are updated      |
| Different non-exported fields on a same component type     | No modifications                          |
| Extra node in the new the tree                             | Node added to the current tree            |
| Missing node in the new tree                               | Extra node is the current tree is removed |

## Fields

Component fields are used to store component data and state. When a component is [updated](#update), field behavior depends on whether it is exported or not.

Here is a test component with an exported and non exported field:

```go
type myCompo struct {
	app.Compo

	ExportedField string
	internalField string
}
```

### Exported fields

Exported fields are like HTML attributes but for components. When [updated](#update), a component is compared to a newly rendered version. Each component fields are compared and updated if different from their newer version counterpart.

Here is a pseudo-Go code that illustrates how it works internally:

```go
a := myCompo{
	ExportedField: "A",
}

b := myCompo{
	ExportedField: "B",
}

// update is the internal function that updates a UI element with a newer
// rendered version.
update(a, b)

fmt.Println("a.ExportedField:", a.ExportedField) // => "B" (updated)
```

### Internal fields

Internal (or non exported) fields are like a state. They are not modified when a component is [updated](#update).

Here is a pseudo-Go code that illustrates how it works internally:

```go
a := myCompo{
	ExportedField: "A",
	internalField: "a",
}

b := myCompo{
	ExportedField: "B",
	internalField: "b",
}

// update is the internal function that updates a UI element with a newer
// rendered version.
update(a, b)

fmt.Println("a.ExportedField:", a.ExportedField) // => "B" (updated)
fmt.Println("a.internalField:", a.internalField) // => "a" (not updated)
```

## Lifecycle

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
