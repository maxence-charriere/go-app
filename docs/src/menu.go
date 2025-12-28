package main

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/maxence-charriere/go-app/v10/pkg/ui"
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
							Href("https://go-app.dev/").
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

				app.Div().Class("separator"),

				ui.Link().
					Class(linkClass).
					Icon(imgFolderSVG).
					Label("Intro").
					Href("/intro").
					Class(isFocus("/intro")),
				app.Div().Class("separator"),
				ui.Link().
					Class(linkClass).
					Icon(keyboardSVG).
					Label("Words").
					Href("https://my-go-app-bvk.pages.dev/"),
				ui.Link().
					Class(linkClass).
					Icon(keyboardSVG).
					Label("Destroy Math").
					Href("https://destroy-math.pages.dev/"),
				ui.Link().
					Class(linkClass).
					Icon(keyboardSVG).
					Label("Pieboy").
					Href("https://pieboy.melamday.workers.dev/"),
				ui.Link().
					Class(linkClass).
					Icon(keyboardSVG).
					Label("Hackrooms").
					Href("https://hackrooms.pages.dev/"),
				app.Div().Class("separator"),

				app.If(m.appInstallable, func() app.UI {
					return ui.Link().
						Class(linkClass).
						Icon(downloadSVG).
						Label("Install").
						OnClick(m.installApp)
				}),
				ui.Link().
					Class(linkClass).
					Icon(userLockSVG).
					Label("Resume").
					Href("http://mamday-resume.s3-website.us-east-2.amazonaws.com/"),

				app.Div().Class("separator"),
			),
		)
}

func (m *menu) installApp(ctx app.Context, e app.Event) {
	ctx.NewAction(installApp)
}
