package main

import "github.com/maxence-charriere/go-app/pkg/app"

type hello struct {
	app.Compo
	name string
}

func (h *hello) Render() app.UI {
	return app.Div().Body(
		app.Main().Body(
			app.H1().Body(
				app.Text("Hello, "),
				app.If(h.name != "",
					app.Text(h.name),
				).Else(
					app.Text("World"),
				),
			),
			app.Input().
				Value(h.name).
				Placeholder("What is your name?").
				AutoFocus(true).
				OnChange(h.OnInputChange),
		),
	)
}

func (h *hello) OnInputChange(src app.Value, e app.Event) {
	h.name = src.Get("value").String()
	h.Update()
}

func main() {
	app.Route("/", &hello{})
	app.Run()
}
