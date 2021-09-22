package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/analytics"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type gettingStartedPage struct {
	app.Compo
}

func newGettingStartedPage() *gettingStartedPage {
	return &gettingStartedPage{}
}

func (p *gettingStartedPage) OnPreRender(ctx app.Context) {
	p.initPage(ctx)
}

func (p *gettingStartedPage) OnNav(ctx app.Context) {
	p.initPage(ctx)
}

func (p *gettingStartedPage) initPage(ctx app.Context) {
	ctx.Page().SetTitle("Start building a PWA with Go and WASM")
	ctx.Page().SetDescription("Documentation that shows how to start building a Progressive Web App (PWA) with Go (Golang) and WebAssembly (WASM).")
	analytics.Page("getting-started", nil)
}

func (p *gettingStartedPage) Render() app.UI {
	return newPage().
		Title("Getting Started").
		Icon(rocketSVG).
		Index(
			newIndexLink().Title("Intro"),
			newIndexLink().Title("Prerequisite"),
			newIndexLink().Title("Install"),
			newIndexLink().Title("Code"),
			newIndexLink().Title("    Hello component"),
			newIndexLink().Title("    Main"),
			newIndexLink().Title("Build and Run"),
			newIndexLink().Title("    Build the Client"),
			newIndexLink().Title("    Build the Server"),
			newIndexLink().Title("    Run the App"),
			newIndexLink().Title("    Use a Makefile"),

			app.Div().Class("separator"),

			newIndexLink().Title("Next"),
		).
		Content(
			newRemoteMarkdownDoc().Src("/web/documents/getting-started.md"),
		)
}
