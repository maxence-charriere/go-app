package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/analytics"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type architecturePage struct {
	app.Compo
}

func newArchitecturePage() *architecturePage {
	return &architecturePage{}
}

func (p *architecturePage) OnPreRender(ctx app.Context) {
	p.initPage(ctx)
}

func (p *architecturePage) OnNav(ctx app.Context) {
	p.initPage(ctx)
}

func (p *architecturePage) initPage(ctx app.Context) {
	ctx.Page().SetTitle("Understanding go-app Architecture")
	ctx.Page().SetDescription("Documentation about how go-app parts are working together to form a Progressive Web App (PWA).")
	analytics.Page("architecture", nil)
}

func (p *architecturePage) Render() app.UI {
	return newPage().
		Title("Architecture").
		Icon(fileTreeSVG).
		Index(
			app.A().
				Class("index-link").
				Class(fragmentFocus("overview")).
				Href("#overview").
				Text("Overview"),
			app.A().
				Class("index-link").
				Class(fragmentFocus("web-browser")).
				Href("#web-browser").
				Text("Web Browser"),
			app.A().
				Class("index-link").
				Class(fragmentFocus("server")).
				Href("#server").
				Text("Server"),
			app.A().
				Class("index-link").
				Class(fragmentFocus("html-pages")).
				Href("#html-pages").
				Text("HTML Pages"),
			app.A().
				Class("index-link").
				Class(fragmentFocus("package-resources")).
				Href("#package-resources").
				Text("Package Resources"),
			app.A().
				Class("index-link").
				Class(fragmentFocus("app-wasm")).
				Href("#app-wasm").
				Text("app.wasm"),
			app.A().
				Class("index-link").
				Class(fragmentFocus("static-resources")).
				Href("#static-resources").
				Text("Static Resources"),

			app.Div().Class("separator"),

			app.A().
				Class("index-link").
				Class(fragmentFocus("next")).
				Href("#next").
				Text("Next"),
		).
		Content(
			newRemoteMarkdownDoc().Src("/web/documents/architecture.md"),
		)
}
