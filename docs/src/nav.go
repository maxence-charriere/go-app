package main

import "github.com/maxence-charriere/go-app/v8/pkg/app"

type nav struct {
	app.Compo

	Iclass      string
	currentPath string
}

func newNav() *nav {
	return &nav{}
}

func (n *nav) Class(v string) *nav {
	if v == "" {
		return n
	}
	if n.Iclass != "" {
		n.Iclass += " "
	}
	n.Iclass += v
	return n
}

func (n *nav) OnNav(ctx app.Context) {
	n.currentPath = ctx.Page.URL().Path
	n.Update()
}

func (n *nav) Render() app.UI {
	return app.Div().
		Class(n.Iclass).
		Class("fill").
		Class("unselectable").
		Body(
			app.Stack().
				Class("header").
				Class("hspace-out").
				Center().
				Content(
					app.Header().Body(
						app.A().
							Class("app-title").
							Text("go-app").
							Href("/"),
					),
				),
			app.Nav().
				Class("content").
				Body(
					app.Div().
						Class("hspace-out").
						Body(
							app.Div().
								Class("vspace-top").
								Body(
									newLink().
										Label("Home").
										Icon(newSVGIcon().RawSVG(homeSVG)).
										Href("/").
										Focus(n.currentPath == "/"),
									newLink().
										Label("Getting started").
										Icon(newSVGIcon().RawSVG(rocketSVG)).
										Href("/start").
										Focus(n.currentPath == "/start"),
									newLink().
										Label("V7 to V8 guide").
										Icon(newSVGIcon().RawSVG(upgradeSVG)).
										Href("/v7-to-v8").
										Focus(n.currentPath == "/v7-to-v8"),
								),
							app.Div().
								Class("vspace-top").
								Body(
									newLink().
										Label("API reference").
										Icon(newSVGIcon().RawSVG(golangSVG)).
										Href("/reference").
										Focus(n.currentPath == "/reference"),
									newLink().
										Label("Architecture").
										Icon(newSVGIcon().RawSVG(fileTreeSVG)).
										Href("/architecture").
										Focus(n.currentPath == "/architecture"),
									newLink().
										Label("Components").
										Icon(newSVGIcon().RawSVG(gridSVG)).
										Href("/components").
										Focus(n.currentPath == "/components"),
									newLink().
										Label("Concurrency").
										Icon(newSVGIcon().RawSVG(concurrecySVG)).
										Href("/concurrency").
										Focus(n.currentPath == "/concurrency"),
									newLink().
										Label("Declarative syntax").
										Icon(newSVGIcon().RawSVG(keyboardSVG)).
										Href("/syntax").
										Focus(n.currentPath == "/syntax"),
									newLink().
										Label("JS/Dom").
										Icon(newSVGIcon().RawSVG(jsSVG)).
										Href("/js").
										Focus(n.currentPath == "/js"),
									newLink().
										Label("Lifecycle").
										Icon(newSVGIcon().RawSVG(arrowSVG)).
										Href("/lifecycle").
										Focus(n.currentPath == "/lifecycle"),
									newLink().
										Label("Routing").
										Icon(newSVGIcon().RawSVG(routeSVG)).
										Href("/routing").
										Focus(n.currentPath == "/routing"),
									newLink().
										Label("SEO").
										Icon(newSVGIcon().RawSVG(seoSVG)).
										Href("/seo").
										Focus(n.currentPath == "/seo"),
									newLink().
										Label("Static resources").
										Icon(newSVGIcon().RawSVG(fileSVG)).
										Href("/static-resources").
										Focus(n.currentPath == "/static-resources"),
									newLink().
										Label("Testing").
										Icon(newSVGIcon().RawSVG(testSVG)).
										Href("/testing").
										Focus(n.currentPath == "/testing"),
								),
							app.Div().
								Class("vspace-top").
								Body(
									newLink().
										Label("Built with go-app").
										Icon(newSVGIcon().RawSVG(hammerSVG)).
										Href("/built-with").
										Focus(n.currentPath == "/built-with"),
									newLink().
										Label("Examples").
										Icon(newSVGIcon().RawSVG(schoolSVG)).
										Href("/examples").
										Focus(n.currentPath == "/examples"),
									newLink().
										Label("Install").
										Icon(newSVGIcon().RawSVG(downloadSVG)).
										Href("/install").
										Focus(n.currentPath == "/install"),
								),
							app.Div().
								Class("vspace-top").
								Class("vspace-bottom").
								Body(
									newLink().
										Label("GitHub").
										Icon(newSVGIcon().RawSVG(githubSVG)).
										Href(githubURL),
									newLink().
										Label("Twitter").
										Icon(newSVGIcon().RawSVG(twitterSVG)).
										Href(twitterURL),
								),
							app.Div().
								Class("vspace-top").
								Class("vspace-bottom").
								Body(
									newLink().
										Label("Buy me a coffee").
										Icon(newSVGIcon().RawSVG(coffeeSVG)).
										Href(buyMeACoffeeURL),
									newLink().
										Label("Donate cryptos").
										Icon(newSVGIcon().RawSVG(bitcoinSVG)).
										Href(coinbaseBusinessURL),
									newLink().
										Label("GitHub Sponsor").
										Icon(newSVGIcon().RawSVG(githubSVG)).
										Href(githubSponsorURL),
									newLink().
										Label("Open Collective").
										Icon(newSVGIcon().RawSVG(opensourceSVG)).
										Href(openCollectiveURL),
								),
						),
				),
		)
}
