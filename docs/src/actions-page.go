package main

import (
	"github.com/maxence-charriere/go-app/v10/pkg/analytics"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type actionPage struct {
	app.Compo
}

func newActionPage() *actionPage {
	return &actionPage{}
}

func (p *actionPage) OnNav(ctx app.Context) {
	p.initPage(ctx)
}

func (p *actionPage) initPage(ctx app.Context) {
	ctx.Page().SetTitle("Creating and Listening to Actions")
	ctx.Page().SetDescription("Documentation about how to create and listen to actions.")
	analytics.Page("actions", nil)
}

func (p *actionPage) Render() app.UI {
	return newPage().
		Title("Actions").
		Icon(actionSVG).
		Index(
			newIndexLink().Title("What is an Action?"),
			newIndexLink().Title("Create"),
			newIndexLink().Title("Handling"),
			newIndexLink().Title("    Global Level"),
			newIndexLink().Title("    Component Level"),

			app.Div().Class("separator"),

			newIndexLink().Title("Next"),
		).
		Content(
			newRemoteMarkdownDoc().Src("/web/documents/actions.md"),
		)
}
