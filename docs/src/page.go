package main

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/maxence-charriere/go-app/v10/pkg/ui"
)

const (
	headerHeight  = 72
	adsenseClient = "ca-pub-1013306768105236"
	adsenseSlot   = "9307554044"
)

type page struct {
	app.Compo

	Iclass   string
	Iindex   []app.UI
	Iicon    string
	Ititle   string
	Icontent []app.UI

	updateAvailable bool
}

func newPage() *page {
	return &page{}
}

func (p *page) Index(v ...app.UI) *page {
	p.Iindex = app.FilterUIElems(v...)
	return p
}

func (p *page) Icon(v string) *page {
	p.Iicon = v
	return p
}

func (p *page) Title(v string) *page {
	p.Ititle = v
	return p
}

func (p *page) Content(v ...app.UI) *page {
	p.Icontent = app.FilterUIElems(v...)
	return p
}

func (p *page) OnNav(ctx app.Context) {
	p.updateAvailable = ctx.AppUpdateAvailable()
	ctx.Defer(scrollTo)
}

func (p *page) OnAppUpdate(ctx app.Context) {
	p.updateAvailable = ctx.AppUpdateAvailable()
}

func (p *page) Render() app.UI {
	return ui.Shell().
		Class("fill").
		Class("background").
		HamburgerButton(app.Div().
			Class("hamburger-menu-icon").
			Body(app.Raw(`
			<svg class="hamburger-menu-icon" viewBox="0 0 24 24">
				<path fill="currentColor" d="M3,6H21V8H3V6M3,11H21V13H3V11M3,16H21V18H3V16Z" />
			</svg>`))).
		HamburgerMenu(
			newMenu().
				Class("fill").
				Class("menu-hamburger-background"),
		).
		Menu(
			newMenu().Class("fill"),
		).
		Index(
			app.If(len(p.Iindex) != 0, func() app.UI {
				return ui.Scroll().
					Class("fill").
					HeaderHeight(headerHeight).
					Content(
						app.Nav().
							Class("index").
							Body(
								app.Div().Class("separator"),
								app.Header().
									Class("h2").
									Text("Index"),
								app.Div().Class("separator"),
								app.Range(p.Iindex).Slice(func(i int) app.UI {
									return p.Iindex[i]
								}),
								newIndexLink().Title("Report an Issue"),
								app.Div().Class("separator"),
							),
					)
			}),
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
									app.If(p.updateAvailable, func() app.UI {
										return app.Div().
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
											)
									}),
								),
						),
				).
				HeaderHeight(headerHeight).
				Content(
					app.Main().Body(
						app.Article().Body(
							app.Header().
								ID("page-top").
								Class("page-title").
								Class("center").
								Body(
									ui.Stack().
										Center().
										Middle().
										Content(
											ui.Icon().
												Class("icon-left").
												Class("unselectable").
												Size(90).
												Src(p.Iicon),
											app.H1().Text(p.Ititle),
										),
								),
							app.Div().Class("separator"),
							app.Range(p.Icontent).Slice(func(i int) app.UI {
								return p.Icontent[i]
							}),

							app.Div().Class("separator"),
							app.Aside().Body(
								app.Header().
									ID("report-an-issue").
									Class("h2").
									Text("Report an issue"),
								app.P().Body(
									app.Text("Found something incorrect, a typo or have suggestions to improve this page? "),
									app.A().
										Href(fmt.Sprintf(
											"%s/issues/new?title=Documentation issue in %s page",
											githubURL,
											p.Ititle,
										)).
										Text("🚀 Submit a GitHub issue!"),
								),
							),
							app.Div().Class("separator"),

							// Testing space
							// app.H2().Text("Test"),
							// app.Form().
							// 	Method("post").
							// 	Action("http://localhost:9600/api/test").
							// 	EncType("multipart/form-data").
							// 	Body(
							// 		app.Input().Placeholder("What is your first name?").AutoFocus(true),
							// 		app.Input().Placeholder("What is your last name?"),
							// 		app.Input().Type("submit").Value("Submit"),
							// 	),
						),
					),
				),
		).
		Ads(
			ui.Flyer().
				Class("fill").
				HeaderHeight(headerHeight).
				Banner(
					app.Aside().
						Class("fill").
						Body(
							ui.AdsenseDisplay().
								Class("fill").
								Class("no-scroll").
								Client(adsenseClient).
								Slot(adsenseSlot),
						),
				).
				PremiumHeight(200).
				Premium(
					newGithubSponsor().Class("fill"),
				),
		)
}

func (p *page) updateApp(ctx app.Context, e app.Event) {
	ctx.NewAction(updateApp)
}

func scrollTo(ctx app.Context) {
	id := ctx.Page().URL().Fragment
	if id == "" {
		id = "page-top"
	}
	ctx.ScrollTo(id)
}

func fragmentFocus(fragment string) string {
	if fragment == app.Window().URL().Fragment {
		return "focus"
	}
	return ""
}
