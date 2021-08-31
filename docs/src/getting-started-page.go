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
			app.A().
				Class("index-link").
				Class(fragmentFocus("intro")).
				Href("#intro").
				Text("Intro"),
			app.A().
				Class("index-link").
				Class(fragmentFocus("prerequisite")).
				Href("#prerequisite").
				Text("Prerequisite"),
			app.A().
				Class("index-link").
				Class(fragmentFocus("install")).
				Href("#install").
				Text("Install"),
			app.A().
				Class("index-link").
				Class(fragmentFocus("code")).
				Href("#code").
				Text("Code"),
			app.A().
				Class("index-link").
				Class(fragmentFocus("hello-component")).
				Href("#hello-component").
				Text("    Hello component"),
			app.A().
				Class("index-link").
				Class(fragmentFocus("main")).
				Href("#main").
				Text("    Main"),
			app.A().
				Class("index-link").
				Class(fragmentFocus("build-and-run")).
				Href("#build-and-run").
				Text("Build and Run"),
			app.A().
				Class("index-link").
				Class(fragmentFocus("build-the-client")).
				Href("#build-the-client").
				Text("    Build the Client"),
			app.A().
				Class("index-link").
				Class(fragmentFocus("build-the-server")).
				Href("#build-the-server").
				Text("    Build the Server"),
			app.A().
				Class("index-link").
				Class(fragmentFocus("running-the-app")).
				Href("#run-the-app").
				Text("    Run the App"),
			app.A().
				Class("index-link").
				Class(fragmentFocus("using-a-makefile")).
				Href("#use-a-makefile").
				Text("    Use a Makefile"),

			app.Div().Class("separator"),

			app.A().
				Class("index-link").
				Class(fragmentFocus("next")).
				Href("#next").
				Text("Next"),
		).
		Content(
			newRemoteMarkdownDoc().Src("/web/documents/getting-started.md"),
		)
}
