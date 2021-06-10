package main

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type shellExample struct {
	app.Compo
}

func (e *shellExample) Render() app.UI {
	return app.Shell2().
		Class("fill").
		HamburgerMenu(
			app.Div().
				Style("width", "100%").
				Style("height", "100%").
				Style("background-color", "green"),
		).
		Menu(
			app.Div().
				Style("width", "100%").
				Style("height", "100%").
				Style("background-color", "deepskyblue"),
		).
		Index(
			app.Div().
				Style("width", "100%").
				Style("height", "100%").
				Style("background-color", "deeppink"),
		).
		Content(
			app.Div().
				Style("width", "100%").
				Style("height", "100%").
				Style("background-color", "orange"),
		)
}
