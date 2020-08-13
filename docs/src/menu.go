package main

import "github.com/maxence-charriere/go-app/v7/pkg/app"

// Menu creates the app main menu.
func Menu() app.UI {
	return &menu{}
}

type menu struct {
	app.Compo
}

func (m *menu) Render() app.UI {
	return app.Nav().
		Class("menu").
		Body(
			app.Div().
				Class("title").
				Text("go-app"),
		)
}
