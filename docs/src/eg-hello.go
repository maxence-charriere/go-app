package main

import (
	"github.com/maxence-charriere/go-app/v8/pkg/app"
)

type hello struct {
	app.Compo

	name string
}

func newHello() *hello {
	return &hello{}
}

// func (h *hello) OnPreRender(ctx app.Context) {
// 	ctx.Page.SetTitle("A Hello World written with go-app")
// 	ctx.Page.SetAuthor("Maxence")
// 	h.Update()
// }

// func (h *hello) Render() app.UI {
// 	return app.H1().Text("Hello " + h.name)
// }

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
