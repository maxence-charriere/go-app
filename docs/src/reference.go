package main

import "github.com/maxence-charriere/go-app/v7/pkg/app"

func Reference() app.UI {
	return &reference{}
}

type reference struct {
	app.Compo
}

func (r *reference) Render() app.UI {
	return app.Shell().
		Class("app-background").
		Menu(bloc("transparent", "")).
		Submenu(GodocMenu()).
		OverlayMenu(bloc("deeppink", "")).
		Content(Godoc())
}
