package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/analytics"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type seoPage struct {
	app.Compo
}

func newSEOPage() *seoPage {
	return &seoPage{}
}

func (p *seoPage) OnPreRender(ctx app.Context) {
	p.initPage(ctx)
}

func (p *seoPage) OnNav(ctx app.Context) {
	p.initPage(ctx)
}

func (p *seoPage) initPage(ctx app.Context) {
	ctx.Page().SetTitle("Building SEO-friendly PWA")
	ctx.Page().SetDescription("Documentation about how to make a Progressive Web App indexable by search engines with go-app package.")
	analytics.Page("seo", nil)
}

func (p *seoPage) Render() app.UI {
	return newPage().
		Title("SEO").
		Icon(seoSVG).
		Index(
			newIndexLink().Title("Intro"),
			newIndexLink().Title("Prerendering"),
			newIndexLink().Title("    Customizing prerendering"),
			newIndexLink().Title("    Customizing page metadata"),
			newIndexLink().Title("    Caching"),

			app.Div().Class("separator"),

			newIndexLink().Title("Next"),
		).
		Content(
			newRemoteMarkdownDoc().Src("/web/documents/seo.md"),
		)
}
