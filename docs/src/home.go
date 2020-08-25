package main

import (
	"github.com/maxence-charriere/go-app/v7/pkg/app"
)

type home struct {
	app.Compo
}

func newHome() app.UI {
	return &home{}
}

func (h *home) Render() app.UI {
	return app.Shell().
		Class("app-background").
		Menu(Menu()).
		Submenu(bloc("transparent", "")).
		OverlayMenu(bloc("deepskyblue", "")).
		Content(bloc("transparent", ""))
}

func bloc(color, text string) app.UI {
	return app.Article().
		Style("width", "auto").
		Style("height", "100%").
		Style("background-color", color).
		Text(text)
}

type cbloc struct {
	app.Compo

	Color string
	Text  string
}

func (c *cbloc) Render() app.UI {
	return app.Article().
		Style("width", "auto").
		Style("height", "100%").
		Style("background-color", c.Color).
		Text(c.Text)
}
