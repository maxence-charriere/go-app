package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/analytics"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type concurrencyPage struct {
	app.Compo
}

func newConcurrencyPage() *concurrencyPage {
	return &concurrencyPage{}
}

func (p *concurrencyPage) OnPreRender(ctx app.Context) {
	p.initPage(ctx)
}

func (p *concurrencyPage) OnNav(ctx app.Context) {
	p.initPage(ctx)
}

func (p *concurrencyPage) initPage(ctx app.Context) {
	ctx.Page().SetTitle("Building a Concurrency-Safe PWA")
	ctx.Page().SetDescription("Documentation about how to build a concurrency-safe and reactive progressive web app (PWA).")
	analytics.Page("concurrency", nil)
}

func (p *concurrencyPage) Render() app.UI {
	return newPage().
		Title("Concurrency").
		Icon(concurrencySVG).
		Index(
			newIndexLink().Title("Intro"),
			newIndexLink().Title("UI Goroutine"),
			newIndexLink().Title("Async"),
			newIndexLink().Title("Dispatch"),
			newIndexLink().Title("Defer"),

			app.Div().Class("separator"),

			newIndexLink().Title("Next"),
		).
		Content(
			newRemoteMarkdownDoc().Src("/web/documents/concurrency.md"),
		)
}
