package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/maxence-charriere/go-app/v9/pkg/ui"
)

type flowExample struct {
	app.Compo
}

func (e *flowExample) Render() app.UI {
	return ui.Flow().
		Class("fill").
		StretchItems().
		Spacing(3).
		Content(
			app.Div().
				Style("width", "100%").
				Style("height", "100%").
				Style("min-height", "240px").
				Style("background-color", "deepskyblue"),
			app.Div().
				Style("width", "100%").
				Style("height", "100%").
				Style("min-height", "240px").
				Style("background-color", "deeppink"),
			app.Div().
				Style("width", "100%").
				Style("height", "100%").
				Style("min-height", "240px").
				Style("background-color", "orange"),
			app.Div().
				Style("width", "100%").
				Style("height", "100%").
				Style("min-height", "240px").
				Style("background-color", "deepskyblue"),
			app.Div().
				Style("width", "100%").
				Style("height", "100%").
				Style("min-height", "240px").
				Style("background-color", "deeppink"),
			app.Div().
				Style("width", "100%").
				Style("height", "100%").
				Style("min-height", "240px").
				Style("background-color", "orange"),
			app.Div().
				Style("width", "100%").
				Style("height", "100%").
				Style("min-height", "240px").
				Style("background-color", "deepskyblue"),
		)
}
