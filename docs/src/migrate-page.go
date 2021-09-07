package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/analytics"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type migratePage struct {
	app.Compo
}

func newMigratePage() *migratePage {
	return &migratePage{}
}

func (p *migratePage) OnPreRender(ctx app.Context) {
	p.initPage(ctx)
}

func (p *migratePage) OnNav(ctx app.Context) {
	p.initPage(ctx)
}

func (p *migratePage) initPage(ctx app.Context) {
	ctx.Page().SetTitle("Migrate Codebase From go-app v8 To v9")
	ctx.Page().SetDescription("Documentation about what changed between go-app v8 and v9.")
	analytics.Page("migrate", nil)
}

func (p *migratePage) Render() app.UI {
	return newPage().
		Title("Migrate From v8 to v9").
		Icon(swapSVG).
		Index(
			app.A().
				Class("index-link").
				Class(fragmentFocus("intro")).
				Href("#intro").
				Text("Intro"),
			app.A().
				Class("index-link").
				Class(fragmentFocus("changes")).
				Href("#changes").
				Text("Changes"),
			app.A().
				Class("index-link").
				Class(fragmentFocus("general")).
				Href("#general").
				Text("    General"),
			app.A().
				Class("index-link").
				Class(fragmentFocus("components")).
				Href("#components").
				Text("    Components"),
			app.A().
				Class("index-link").
				Class(fragmentFocus("context")).
				Href("#context").
				Text("    Context"),
			app.A().
				Class("index-link").
				Class(fragmentFocus("api-design-decisions")).
				Href("#api-design-decisions").
				Text("API Design Decisions"),

			app.Div().Class("separator"),
		).
		Content(
			newRemoteMarkdownDoc().Src("/web/documents/migrate.md"),
		)
}
