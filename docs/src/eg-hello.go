package main

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v8/pkg/app"
)

type hello struct {
	app.Compo

	name string
}

func newHello() *hello {
	return &hello{}
}

func (h *hello) Render() app.UI {
	return app.Div().Body(
		app.H1().Body(
			app.Text("Hello, "),
			app.If(h.name != "",
				app.Text(h.name),
			).Else(
				app.Text("World!"),
			),
		),
		app.P().Body(
			app.Input().
				Type("text").
				Value(h.name).
				Placeholder("What is your name?").
				AutoFocus(true).
				OnChange(h.ValueTo(&h.name)),
		),
	)
}

// This will be replaced later by Go 1.16 embed once it will be more widely s
// upported.
var egHelloCode = fmt.Sprintf("```go\n%s\n```", `
package main

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v8/pkg/app"
)

type hello struct {
	app.Compo

	name string
}

func newHello() *hello {
	return &hello{}
}

func (h *hello) Render() app.UI {
	return app.Div().Body(
		app.H1().Body(
			app.Text("Hello, "),
			app.If(h.name != "",
				app.Text(h.name),
			).Else(
				app.Text("World!"),
			),
		),
		app.P().Body(
			app.Input().
				Type("text").
				Value(h.name).
				Placeholder("What is your name?").
				AutoFocus(true).
				OnChange(h.ValueTo(&h.name)),
		),
	)
}
`)
