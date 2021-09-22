package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/analytics"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type staticResourcesPage struct {
	app.Compo
}

func newStaticResourcePage() *staticResourcesPage {
	return &staticResourcesPage{}
}

func (p *staticResourcesPage) OnPreRender(ctx app.Context) {
	p.initPage(ctx)
}

func (p *staticResourcesPage) OnNav(ctx app.Context) {
	p.initPage(ctx)
}

func (p *staticResourcesPage) initPage(ctx app.Context) {
	ctx.Page().SetTitle("Images and Static Resources")
	ctx.Page().SetDescription("Documentation about how to deal with images and other static resources.")
	analytics.Page("static-resources", nil)
}

func (p *staticResourcesPage) Render() app.UI {
	return newPage().
		Title("Images and Static Resources").
		Icon(imgFolderSVG).
		Index(
			app.Div().Class("separator"),

			newIndexLink().Title("Next"),
		).
		Content(
			newRemoteMarkdownDoc().Src("/web/documents/static-resources.md"),
		)
}
