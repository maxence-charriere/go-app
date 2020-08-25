package main

import "github.com/maxence-charriere/go-app/v7/pkg/app"

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
		// Submenu().
		OverlayMenu(Menu()).
		Content(&complexScenario{})
}
