package main

import (
	"github.com/maxence-charriere/go-app/v8/pkg/app"
)

const (
	pageItemWidth = 276
)

type page struct {
	app.Compo

	Iindex      []app.UI
	Icontent    []app.UI
	IissueTitle string

	isAppUpdateAvailable bool
}

func newPage() *page {
	return &page{}
}

func (p *page) Index(v ...app.UI) *page {
	p.Iindex = app.FilterUIElems(v...)
	return p
}

func (p *page) Content(v ...app.UI) *page {
	p.Icontent = app.FilterUIElems(v...)
	return p
}

func (p *page) IssueTitle(v string) *page {
	p.IissueTitle = v
	return p
}

func (p *page) OnNav(ctx app.Context) {
	p.setAvailableUpdate(ctx)
}

func (p *page) OnAppUpdate(ctx app.Context) {
	p.setAvailableUpdate(ctx)
}

func (p *page) setAvailableUpdate(ctx app.Context) {
	p.isAppUpdateAvailable = ctx.AppUpdateAvailable
	p.Update()
}

func (p *page) Render() app.UI {
	return app.Shell().
		Class("background").
		MenuWidth(pageItemWidth).
		Menu(newNav()).
		OverlayMenu(
			newNav().
				Class("overlay-menu").
				Class("unselectable"),
		).
		Submenu(
			app.Nav().
				Class("header-out").
				Class("content").
				Class("unselectable").
				Body(
					app.Div().
						Class("hspace-out").
						Body(
							app.Header().
								Class("h2").
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
			app.Div().
				Class("content").
				Body(
					app.Main().
						ID("top").
						Body(
							app.Article().
								Body(
									app.Range(p.Icontent).Slice(func(i int) app.UI {
										return p.Icontent[i]
									}),
									app.Footer().
										Class("vspace-section").
										Class("hspace-out").
										Body(newIssue().Title(p.IissueTitle)),
								),
							app.Aside().
								Class("vspace-section").
								Class("vspace-bottom").
								Body(newSupportUs()),
						),
				),
		)
}

func (p *page) onUpdateClick() {
	p.Defer(func(ctx app.Context) {
		ctx.Reload()
	})
}
