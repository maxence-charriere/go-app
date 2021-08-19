package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/maxence-charriere/go-app/v9/pkg/ui"
)

const (
	headerHeight = 72
)

type page struct {
	app.Compo

	Iclass string

	updateAvailable bool
}

func newPage() *page {
	return &page{}
}

func (p *page) OnNav(ctx app.Context) {
	p.updateAvailable = ctx.AppUpdateAvailable()
}

func (p *page) OnAppUpdate(ctx app.Context) {
	p.updateAvailable = ctx.AppUpdateAvailable()
}

func (p *page) Render() app.UI {
	return ui.Shell().
		Class("fill").
		Class("background").
		HamburgerMenu(newMenu().Class("fill")).
		Menu(newMenu().Class("fill")).
		Index(
			ui.Scroll().
				Class("fill").
				HeaderHeight(headerHeight),
		).
		Content(
			ui.Scroll().
				Class("fill").
				Header(
					app.Nav().
						Class("fill").
						Body(
							ui.Stack().
								Class("fill").
								Right().
								Middle().
								Content(
									app.If(p.updateAvailable,
										app.Div().
											Class("link-update").
											Body(
												ui.Link().
													Class("link").
													Class("heading").
													Class("fit").
													Class("unselectable").
													Icon(downloadSVG).
													Label("Update").
													OnClick(p.updateApp),
											),
									),
								),
						),
				).
				HeaderHeight(headerHeight),
		).
		Ads(
			ui.Flyer().
				Class("fill").
				HeaderHeight(headerHeight).
				PremiumHeight(200).
				Premium(
					newGithubSponsor().Class("fill"),
				),
		)
}

func (p *page) updateApp(ctx app.Context, e app.Event) {
	ctx.NewAction(updateApp)
}
