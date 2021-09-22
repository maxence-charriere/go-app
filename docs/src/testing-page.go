package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/analytics"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type testingPage struct {
	app.Compo
}

func newTestingPage() *testingPage {
	return &testingPage{}
}

func (p *testingPage) OnPreRender(ctx app.Context) {
	p.initPage(ctx)
}

func (p *testingPage) OnNav(ctx app.Context) {
	p.initPage(ctx)
}

func (p *testingPage) initPage(ctx app.Context) {
	ctx.Page().SetTitle("Testing Components")
	ctx.Page().SetDescription("Documentation about how to unit test components created with go-app.")
	analytics.Page("testing", nil)
}

func (p *testingPage) Render() app.UI {
	return newPage().
		Title("Testing").
		Icon(testSVG).
		Index(
			newIndexLink().Title("Intro"),
			newIndexLink().Title("Component server prerendering"),
			newIndexLink().Title("Component client lifecycle"),
			newIndexLink().Title("Asynchronous operations"),
			newIndexLink().Title("UI elements"),

			app.Div().Class("separator"),

			newIndexLink().Title("Next"),
		).
		Content(
			newRemoteMarkdownDoc().Src("/web/documents/testing.md"),
		)
}
