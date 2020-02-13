package main

import (
	"github.com/maxence-charriere/app/pkg/app"
)

type hello struct {
	app.Compo

	Name string
}

func (h *hello) Render() app.UI {
	return app.Div().
		Body(
			app.Div().
				Class("menu-button").
				OnClick(h.OnMenuClick).
				Body(
					app.Text("☰"),
				),
			app.Main().
				Class("hello").
				Body(
					app.H1().
						Class("hello-title").
						Body(
							app.Text("Hello, "),
							app.If(h.Name != "",
								app.Text(h.Name),
							).Else(
								app.Text("World"),
							),
						),
					app.Input().
						Class("hello-input").
						Value(h.Name).
						Placeholder("What is your name?").
						AutoFocus(true).
						OnChange(h.OnInputChange),
				),
		)
}

func (h *hello) OnMenuClick(src app.Value, e app.Event) {
	app.NewContextMenu(
		app.MenuItem().
			Label("Reload").
			Keys("cmdorctrl+r").
			OnClick(func(src app.Value, e app.Event) {
				app.Reload()
			}),
		app.MenuItem().Separator(),
		app.MenuItem().
			Label("City demo").
			OnClick(func(src app.Value, e app.Event) {
				app.Navigate("/city")
			}),
		app.MenuItem().Separator(),
		app.MenuItem().
			Icon(icon).
			Label("Go to repository").
			OnClick(func(src app.Value, e app.Event) {
				app.Navigate("https://github.com/maxence-charriere/app")
			}),
		app.MenuItem().
			Icon(icon).
			Label("Sources").
			OnClick(func(src app.Value, e app.Event) {
				app.Navigate("https://github.com/maxence-charriere/app/blob/master/demo/cmd/demo-wasm/hello.go")
			}),
	)
}

func (h *hello) OnInputChange(src app.Value, e app.Event) {
	h.Name = src.Get("value").String()
	h.Update()
}
