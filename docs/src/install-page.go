package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/analytics"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type installPage struct {
	app.Compo
}

func newInstallPage() *installPage {
	return &installPage{}
}

func (p *installPage) OnPreRender(ctx app.Context) {
	p.initPage(ctx)
}

func (p *installPage) OnNav(ctx app.Context) {
	p.initPage(ctx)
}

func (p *installPage) initPage(ctx app.Context) {
	ctx.Page().SetTitle("Handling App Install")
	ctx.Page().SetDescription("Documentation about how to install an app created with go-app.")
	analytics.Page("install", nil)
}

func (p *installPage) Render() app.UI {
	return newPage().
		Title("Install").
		Icon(downloadSVG).
		Index(
			newIndexLink().Title("Intro"),
			newIndexLink().Title("Desktop"),
			newIndexLink().Title("IOS"),
			newIndexLink().Title("Android"),
			newIndexLink().Title("Programmatically"),
			newIndexLink().Title("    Detect Install Support"),
			newIndexLink().Title("    Display Install Popup"),

			app.Div().Class("separator"),

			newIndexLink().Title("Next"),
		).
		Content(
			newRemoteMarkdownDoc().Src("/web/documents/install.md"),
		)
}
