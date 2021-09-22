package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/analytics"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type lifecyclePage struct {
	app.Compo
}

func newLifecyclePage() *lifecyclePage {
	return &lifecyclePage{}
}

func (p *lifecyclePage) OnPreRender(ctx app.Context) {
	p.initPage(ctx)
}

func (p *lifecyclePage) OnNav(ctx app.Context) {
	p.initPage(ctx)
}

func (p *lifecyclePage) initPage(ctx app.Context) {
	ctx.Page().SetTitle("App Lifecycle and Updates")
	ctx.Page().SetDescription("Documentation that describes how a web browser installs and updates a go-app Progressive Web App.")
	analytics.Page("lifecycle", nil)
}

func (p *lifecyclePage) Render() app.UI {
	return newPage().
		Title("Lifecycle and Updates").
		Icon(arrowSVG).
		Index(
			newIndexLink().Title("Lifecycle Overview"),
			newIndexLink().Title("    First loading"),
			newIndexLink().Title("    Recurrent loadings"),
			newIndexLink().Title("    Loading after an app update"),
			newIndexLink().Title("Listen for App Updates"),

			app.Div().Class("separator"),

			newIndexLink().Title("Next"),
		).
		Content(
			newRemoteMarkdownDoc().Src("/web/documents/lifecycle.md"),
		)
}
