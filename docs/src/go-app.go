package main

import (
	_ "embed"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// intro is a component that displays a simple "intro World!". A component is a
// customizable, independent, and reusable UI element. It is created by
// embedding app.Compo into a struct.
type intro struct {
	app.Compo
}

// The Render method is where the component appearance is defined. Here, a
// markdown file is displayed as content.
//
//go:embed documents/entry2.md
var entry2Content string

func (h *intro) Render() app.UI {
	return newPage().
		Title("Building a Blog with Go-App").
		Icon(rocketSVG).
		Index(
			newIndexLink().Title("Intro").Href("/intro"),
			app.Div().Class("separator"),
		).
		Content(
			newMarkdownDoc().MD(entry2Content), // Use embedded content directly
		)
}
