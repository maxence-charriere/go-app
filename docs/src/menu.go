package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/maxence-charriere/go-app/v9/pkg/ui"
)

type menu struct {
	app.Compo

	Iclass string

	appInstallable bool
}

func newMenu() *menu {
	return &menu{}
}

func (m *menu) Class(v string) *menu {
	m.Iclass = app.AppendClass(m.Iclass, v)
	return m
}

func (m *menu) OnNav(ctx app.Context) {
	m.appInstallable = ctx.IsAppInstallable()
}

func (m *menu) OnAppInstallChange(ctx app.Context) {
	m.appInstallable = ctx.IsAppInstallable()
}

func (m *menu) Render() app.UI {
	linkClass := "link heading fit unselectable"

	isFocus := func(path string) string {
		if app.Window().URL().Path == path {
			return "focus"
		}
		return ""
	}

	return ui.Scroll().
		Class("menu").
		Class(m.Iclass).
		HeaderHeight(headerHeight).
		Header(
			ui.Stack().
				Class("fill").
				Middle().
				Content(
					app.Header().Body(
						app.A().
							Class("heading").
							Class("app-title").
							Href("/").
							Text("Go-App"),
					),
				),
		).
		Content(
			app.Nav().Body(
				app.Div().Class("separator"),

				ui.Link().
					Class(linkClass).
					Icon(homeSVG).
					Label("Home").
					Href("/").
					Class(isFocus("/")),
				ui.Link().
					Class(linkClass).
					Icon(rocketSVG).
					Label("Getting Started").
					Href("/getting-started").
					Class(isFocus("/getting-started")),
				ui.Link().
					Class(linkClass).
					Icon(fileTreeSVG).
					Label("Architecture").
					Href("/architecture").
					Class(isFocus("/architecture")),
				ui.Link().
					Class(linkClass).
					Icon(golangSVG).
					Label("Reference").
					Href("/reference").
					Class(isFocus("/reference")),

				app.Div().Class("separator"),

				ui.Link().
					Class(linkClass).
					Icon(swapSVG).
					Label("Migrate From v8 to v9").
					Href("/migrate").
					Class(isFocus("/migrate")),
				app.If(m.appInstallable,
					ui.Link().
						Class(linkClass).
						Icon(downloadSVG).
						Label("Install").
						OnClick(m.installApp),
				),

				app.Div().Class("separator"),

				ui.Link().
					Class(linkClass).
					Icon(githubSVG).
					Label("GitHub").
					Href(githubURL),
				ui.Link().
					Class(linkClass).
					Icon(twitterSVG).
					Label("Twitter").
					Href(twitterURL),

				app.Div().Class("separator"),
			),
		)
}

func (m *menu) installApp(ctx app.Context, e app.Event) {
	ctx.NewAction(installApp)
}
