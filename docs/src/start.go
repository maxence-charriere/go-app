package main

import "github.com/maxence-charriere/go-app/v7/pkg/app"

type startMenu struct {
	app.Compo
}

func newStartMenu() *startMenu {
	return &startMenu{}
}

func (m *startMenu) Render() app.UI {
	return app.Aside().
		Class("layout").
		Body(
			app.Div().Class("header"),
			app.Div().
				Class("content").
				Body(
					app.Section().Body(
						app.H1().Text("Table of contents"),
					),
				),
		)
}

type start struct {
	app.Compo
}

func newStart() *start {
	return &start{}
}

func (s *start) Render() app.UI {
	return app.Shell().
		Class("app-background").
		Menu(Menu()).
		Submenu(
			newStartMenu(),
		).
		OverlayMenu(Menu()).
		Content(
			newDocument("/web/start.md").
				Description("start.md"),
		)
}
