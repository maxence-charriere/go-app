package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/analytics"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type jsPage struct {
	app.Compo
}

func newJSPage() *jsPage {
	return &jsPage{}
}

func (p *jsPage) OnPreRender(ctx app.Context) {
	p.initPage(ctx)
}

func (p *jsPage) OnNav(ctx app.Context) {
	p.initPage(ctx)
}

func (p *jsPage) initPage(ctx app.Context) {
	ctx.Page().SetTitle("JavaScript Interoperability")
	ctx.Page().SetDescription("Documentation about how to call JavaScript from Go or Go from JavaScript.")
	analytics.Page("js", nil)
}

func (p *jsPage) Render() app.UI {
	return newPage().
		Title("JavaScript Interoperability").
		Icon(jsSVG).
		Index(
			app.Div().Class("separator"),

			newIndexLink().Title("Next"),
		).
		Content(
			newRemoteMarkdownDoc().Src("/web/documents/js.md"),
		)
}
