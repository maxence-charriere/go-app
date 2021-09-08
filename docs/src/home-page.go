package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/analytics"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/maxence-charriere/go-app/v9/pkg/ui"
)

type homePage struct {
	app.Compo
}

func newHomePage() *homePage {
	return &homePage{}
}

func (p *homePage) OnPreRender(ctx app.Context) {
	p.initPage(ctx)
}

func (p *homePage) OnNav(ctx app.Context) {
	p.initPage(ctx)
}

func (p *homePage) initPage(ctx app.Context) {
	ctx.Page().SetTitle(defaultTitle)
	ctx.Page().SetDescription(defaultDescription)
	analytics.Page("home", nil)
}

func (p *homePage) Render() app.UI {
	return newPage().
		Title("go-app").
		Icon("https://storage.googleapis.com/murlok-github/icon-192.png").
		Index(
			app.A().
				Class("index-link").
				Class(fragmentFocus("what-is-go-app")).
				Href("#what-is-go-app").
				Text("What is go-app?"),
			app.A().
				Class("index-link").
				Class(fragmentFocus("updates")).
				Href("#updates").
				Text("Updates"),
			app.A().
				Class("index-link").
				Class(fragmentFocus("declarative-syntax")).
				Href("#declarative-syntax").
				Text("Declarative Syntax"),
			app.A().
				Class("index-link").
				Class(fragmentFocus("standard-http-server")).
				Href("#standard-http-server").
				Text("Standard HTTP Server"),
			app.A().
				Class("index-link").
				Class(fragmentFocus("other-features")).
				Href("#other-features").
				Text("Other features"),
			app.A().
				Class("index-link").
				Class(fragmentFocus("built-with-goapp")).
				Href("#built-with-goapp").
				Text("Built With go-app"),

			app.Div().Class("separator"),

			app.A().
				Class("index-link").
				Class(fragmentFocus("next")).
				Href("#next").
				Text("Next"),
		).
		Content(
			ui.Flow().
				StretchItems().
				Spacing(84).
				Content(
					newRemoteMarkdownDoc().
						Class("fill").
						Src("/web/documents/what-is-go-app.md"),
					newRemoteMarkdownDoc().
						Class("fill").
						Class("updates").
						Src("/web/documents/updates.md"),
				),

			app.Div().Class("separator"),

			newRemoteMarkdownDoc().Src("/web/documents/home.md"),

			app.Div().Class("separator"),

			newBuiltWithGoapp().ID("built-with-goapp"),

			app.Div().Class("separator"),

			newRemoteMarkdownDoc().Src("/web/documents/home-next.md"),
		)
}
