package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/analytics"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type statesPage struct {
	app.Compo
}

func newStatesPage() *statesPage {
	return &statesPage{}
}

func (p *statesPage) OnPreRender(ctx app.Context) {
	p.initPage(ctx)
}

func (p *statesPage) OnNav(ctx app.Context) {
	p.initPage(ctx)
}

func (p *statesPage) initPage(ctx app.Context) {
	ctx.Page().SetTitle("State Management")
	ctx.Page().SetDescription("Documentation about how to set and observe states.")
	analytics.Page("states", nil)
}

func (p *statesPage) Render() app.UI {
	return newPage().
		Title("State Management").
		Icon(stateSVG).
		Index(
			newIndexLink().Title("What is a state?"),
			newIndexLink().Title("Set"),
			newIndexLink().Title("    Options"),
			newIndexLink().Title("Observe"),
			newIndexLink().Title("    Conditional Observation"),
			newIndexLink().Title("    Additional Instructions"),
			newIndexLink().Title("Get"),

			app.Div().Class("separator"),

			newIndexLink().Title("Next"),
		).
		Content(
			newRemoteMarkdownDoc().Src("/web/documents/states.md"),
		)
}
