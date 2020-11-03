# Declarative syntax

Customizing [components](/components) requires to describe how their UI looks like.

The main way to do it is to use the HTML elements defined in the package API.

By using them with a chaining mechanism and the Go syntax, writing a UI is done in a declarative fashion, without using another language.

Here is an example that describes a title and its text:

```go
func (c *myCompo) Render() app.UI {
	return app.Div().Body(
		app.H1().
			Class("title").
			Text("Build a GUI with Go"),
		app.P().
			Class("text").
			Text("Just because Go and this package are really awesome!"),
	)
}
```

## HTML elements

The package provides an interface for each standard HTML element.

Here is a simplified version of the interface for a [\<div>](/reference#HTMLDiv):

```go
type HTMLDiv interface {
    // Attributes:
    Body(nodes ...Node) HTMLDiv
    Class(v string) HTMLDiv
    ID(v string) HTMLDiv
    Style(k, v string) HTMLDiv

    // Event handlers:
    OnClick(h EventHandler) HTMLDiv
    OnKeyPress(h EventHandler) HTMLDiv
    OnMouseOver(h EventHandler) HTMLDiv
}
```

### Create

Creating an HTML element is done by calling a function named after its name. The example bellow create a [\<div>](/reference#Div):

```go
func (c *myCompo) Render() app.UI {
	return app.Div()
}
```

### Standard elements

Standard elements are elements that can contain other elements. To do so they provide the `Body()` method which takes other elements as parameters:

```go
func (c *myCompo) Render() app.UI {
	return app.Div().Body(        // Div Container
		app.H1().Text("Title"),   // First child
		app.P(), Text("Content"), // Second child
	)
}
```

### Self closing elements

Self-closing elements are elements that cannot contain other elements.

```go
func (c *myCompo) Render() app.UI {
	return app.Img().Src("/myImage.png")
}
```

### Attributes

HTML element interfaces provide methods to set element attributes:

```go
func (c *myCompo) Render() app.UI {
	return app.Div().Class("my-class")
}
```

Multiple attributes are set by using the chaining mechanism:

```go
func (c *myCompo) Render() app.UI {
	return app.Div().
		ID("id-name").
		Class("class-name")
}
```

### Style

Style is an attribute that sets the element style with CSS.

```go
func (c *myCompo) Render() app.UI {
	return app.Div().Style("width", "400px")
}
```

Multiple styles can be set by calling the `Style()` method successively:

```go
func (c *myCompo) Render() app.UI {
	return app.Div().
		Style("width", "400px").
		Style("height", "200px").
		Style("background-color", "deepskyblue")
}
```

### Event handlers

[Event handlers](/reference#EventHandler) are functions that are called when an HTML event occurs. They must have the following signature:

```go
func(ctx app.Context, e app.Event)
```

Like attributes, element interfaces provide methods to bind event handlers to given functions:

```go
func (c *myCompo) Render() app.UI {
	return app.Div().OnClick(c.onClick)
}

func (c *myCompo) onClick(ctx app.Context, e app.Event) {
	fmt.Println("onClick is called")
}

```

[Context](/reference#Context) is a context that can be used with any function accepting a [context.Context](https://golang.org/pkg/context/#Context). It is cancelled when the source of the event is dismounted. Source element value can be retrieved by the JSSrc method:

```go
func (c *myCompo) Render() app.UI {
	return app.Div().OnChange(c.onChange)
}

func (c *myCompo) onChange(ctx app.Context, e app.Event) {
	v := ctx.JSSrc().Get("value")
}

```

`ctx.JSSrc()` and [Event](/reference#Event) are [JavaScript objects wrapped in Go interfaces](/js).

## Text

Text represents simple HTML text. They are created by calling the [Text()](/reference#Text) function:

```go
func (c *myCompo) Render() app.UI {
	return app.Div().Body( // Container
		app.Text("Hello"), // First text
		app.Text("World"), // Second text
	)
}
```

## Raw elements

There is some case where some raw HTML code might be required. It can be done by using the [Raw()](/reference#Raw) function with HTML code as argument:

```go
func (c *myCompo) Render() app.UI {
	return app.Raw(`
	<svg width="100" height="100">
		<circle cx="50" cy="50" r="40" stroke="green" stroke-width="4" fill="yellow" />
	</svg>
	`)
}
```

Be aware that using this function is **unsafe** since there is no check on HTML construct or format.

## Nested components

[Components](/components) are structs that let you split the UI into independent and reusable pieces. They can be used within other components to achieve more complex UIs.

```go
// foo component:
type foo struct {
	app.Compo
}

func (f *foo) Render() app.UI {
	return app.P().Body(
		app.Text("Foo, "), // Simple HTML text
		&bar{},            // Nested bar component
	)
}

// bar component:
type bar struct {
	app.Compo
}

func (b *bar) Render() app.UI {
	return app.Text("Bar!")
}
```

In the example above, a `bar` component is used in the `Render()` method of `foo`, so that it will be displayed as content of `foo`.

## Condition

A [Condition](/reference#Condition) is a construct that selects the UI elements that satisfy a condition. They are created by calling the [If()](/reference#If) function.

### If

Here is an If example that shows a title only when the `showTitle` value is `true`:

```go
type myCompo struct {
	app.Compo

	showTitle bool
}

func (c *myCompo) Render() app.UI {
	return app.Div().Body(
		app.If(c.showTitle,
			app.H1().Text("hello"),
		),
	)
}
```

### ElseIf

Here is an ElseIf example that shows a title in different colors depending on an `int` value:

```go
type myCompo struct {
	app.Compo

	color int
}

func (c *myCompo) Render() app.UI {
	return app.Div().Body(
		app.If(c.color > 7,
			app.H1().
				Style("color", "green").
				Text("Good!"),
		).ElseIf(c.color < 4,
			app.H1().
				Style("color", "red").
				Text("Bad!"),
		).Else(
			app.H1().
				Style("color", "orange").
				Text("So so!"),
		),
	)
}
```

### Else

Here is an Else example that shows a simple text when the `showTitle` value is `false`:

```go
type myCompo struct {
	app.Compo

	showTitle bool
}

func (c *myCompo) Render() app.UI {
	return app.Div().Body(
		app.If(c.showTitle,
			app.H1().Text("hello"),
		).Else(
			app.Text("world"), // Shown when showTitle == false
		),
	)
}
```

## Range

Range represents a [range loop](/reference#RangeLoop) that shows UI elements generated from a [slice](#slice) or [map](#map). They are created by calling the [Range()](/reference#Range) function.

### Slice

Here is a slice example that shows an unordered list from a `[]string`:

```go
func (c *myCompo) Render() app.UI {
	data := []string{
		"hello",
		"go-app",
		"is",
		"sexy",
	}

	return app.Ul().Body(
		app.Range(data).Slice(func(i int) app.UI {
			return app.Li().Text(data[i])
		}),
	)
}
```

### Map

Here is a map example that shows an unordered list from a `map[string]int`:

```go
func (c *myCompo) Render() app.UI {
	data := map[string]int{
		"Go":         10,
		"JavaScript": 4,
		"Python":     6,
		"C":          8,
	}

	return app.Ul().Body(
		app.Range(data).Map(func(k string) app.UI {
			s := fmt.Sprintf("%s: %v/10", k, data[k])

			return app.Li().Text(s)
		}),
	)
}
```
