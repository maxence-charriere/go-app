package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/analytics"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type routingPage struct {
	app.Compo
}

func newRoutingPage() *routingPage {
	return &routingPage{}
}

func (p *routingPage) OnPreRender(ctx app.Context) {
	p.initPage(ctx)
}

func (p *routingPage) OnNav(ctx app.Context) {
	p.initPage(ctx)
}

func (p *routingPage) initPage(ctx app.Context) {
	ctx.Page().SetTitle("Routing URL paths to Components")
	ctx.Page().SetDescription("Documentation about how to associate URL paths to go-app components.")
	analytics.Page("routing", nil)
}

func (p *routingPage) Render() app.UI {
	return newPage().
		Title("Routing").
		Icon(routeSVG).
		Index(
			newIndexLink().Title("Intro"),
			newIndexLink().Title("Define a route"),
			newIndexLink().Title("    Simple route"),
			newIndexLink().Title("    Route with regular expression"),
			newIndexLink().Title("How it works?"),
			newIndexLink().Title("Detect navigation"),

			app.Div().Class("separator"),

			newIndexLink().Title("Next"),
		).
		Content(
			newRemoteMarkdownDoc().Src("/web/documents/routing.md"),
		)
}
