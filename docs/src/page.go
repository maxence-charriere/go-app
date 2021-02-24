package main

import (
	"github.com/maxence-charriere/go-app/v8/pkg/app"
)

const (
	pageItemWidth = 276
)

type page2 struct {
	app.Compo

	Iindex   []app.UI
	Icontent []app.UI

	isAppUpdateAvailable bool
}

func newPage2() *page2 {
	return &page2{}
}

func (p *page2) Index(v ...app.UI) *page2 {
	p.Iindex = app.FilterUIElems(v...)
	return p
}

func (p *page2) Content(v ...app.UI) *page2 {
	p.Icontent = app.FilterUIElems(v...)
	return p
}

func (p *page2) OnNav(ctx app.Context) {
	p.setAvailableUpdate(ctx)
}

func (p *page2) OnAppUpdate(ctx app.Context) {
	p.setAvailableUpdate(ctx)
}

func (p *page2) setAvailableUpdate(ctx app.Context) {
	p.isAppUpdateAvailable = ctx.AppUpdateAvailable
	p.Update()
}

func (p *page2) Render() app.UI {
	return app.Shell().
		Class("background").
		MenuWidth(pageItemWidth).
		Menu(newNav()).
		OverlayMenu(newNav().Class("overlay-menu")).
		Submenu(
			app.Nav().
				Class("header-out").
				Class("content").
				Body(
					app.Div().
						Class("hspace-out").
						Body(
							app.Header().
								Class("h1").
								Class("vspace-top").
								Text("Index"),
							app.Div().
								Class("vspace-top").
								Body(
									app.Range(p.Iindex).Slice(func(i int) app.UI {
										return p.Iindex[i]
									}),
								),
							newIndex().
								Class("vspace-top").
								Class("vspace-bottom").
								Links(
									"Report issue",
									"Support go-app",
								),
						),
				),
		).
		Content(
			app.Stack().
				Class("header").
				Center().
				Content(
					app.If(p.isAppUpdateAvailable,
						newLink().
							Class("hspace-out").
							Class("right").
							Class("link-update").
							Label("Update").
							Icon(newSVGIcon().RawSVG(downloadSVG)).
							OnClick(p.onUpdateClick),
					),
				),
			app.Main().
				Class("content").
				Body(
					app.Range(p.Icontent).Slice(func(i int) app.UI {
						return p.Icontent[i]
					}),
					app.Aside().
						Class("vspace-section").
						Class("vspace-bottom").
						Body(newSupportUs()),
				),
		)
}

func (p *page2) onUpdateClick() {
	p.Defer(func(ctx app.Context) {
		ctx.Reload()
	})
}
