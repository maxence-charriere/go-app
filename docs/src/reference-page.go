package main

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type referencePage struct {
	app.Compo
}

func newReferencePage() *referencePage {
	return &referencePage{}
}

func (p *referencePage) OnNav(ctx app.Context) {
}

func (p *referencePage) Render() app.UI {
	return newPage().
		Title("Reference").
		Icon(golangSVG).
		Index(
			app.A().
				Class("index-link").
				Class(fragmentFocus("pkg-overview")).
				Href("#pkg-overview").
				Text("Overview"),
			newReferenceContent().
				Class("reference-index").
				Index(true),
		).
		Content(
			newReferenceContent().Class("reference"),
		)
}
