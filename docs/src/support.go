package main

import "github.com/maxence-charriere/go-app/v8/pkg/app"

const (
	supportUsIconSize = 60
)

type supportUs struct {
	app.Compo
}

func newSupportUs() *supportUs {
	return &supportUs{}
}

func (s *supportUs) Render() app.UI {
	return app.Div().
		ID("support-go-app").
		Body(
			app.Div().
				Class("hspace-out").
				Body(
					app.Header().
						Class("h2").
						Class("header-separator").
						Text("Support go-app"),
					app.P().Body(
						app.Text("Hello there, I'am Maxence, the creator of "),
						app.B().Text("go-app"),
						app.Text("."),
					),
					app.P().Body(
						app.Text("Let's go straight to the point: "),
						app.Strong().Text("I want to make this package become something big!"),
					),
					app.P().Body(
						app.A().
							Href(buyMeACoffeeURL).
							Target("_blank").
							Text("Buying me a coffee"),
						app.Text(", being part of the "),
						app.A().
							Href(openCollectiveURL).
							Target("_blank").
							Text("Open Collective"),
						app.Text(", sponsoring me on "),
						app.A().
							Href(githubSponsorURL).
							Target("_blank").
							Text("GitHub"),
						app.Text(", or giving me some cryptocurrencies, all would help me reach that goal, sustain the development, and boost motivation during long coding sessions."),
					),
				),
			app.Flow().
				Class("space-flow").
				StrechtOnSingleRow().
				Content(
					app.Div().
						Class("space-flow-item").
						Body(
							newSupportUsItem().
								Class("fill").
								Label("Buy me a coffee").
								Href(buyMeACoffeeURL).
								Icon(
									newSVGIcon().
										RawSVG(coffeeSVG).
										Size(supportUsIconSize),
								),
						),
					app.Div().
						Class("space-flow-item").
						Body(
							newSupportUsItem().
								Class("fill").
								Label("Donate cryptos").
								Href(coinbaseBusinessURL).
								Icon(
									newSVGIcon().
										RawSVG(bitcoinSVG).
										Size(supportUsIconSize),
								),
						),
					app.Div().
						Class("space-flow-item").
						Body(
							newSupportUsItem().
								Class("fill").
								Label("Github sponsor").
								Href(githubSponsorURL).
								Icon(
									newSVGIcon().
										RawSVG(githubSVG).
										Size(supportUsIconSize),
								),
						),
					app.Div().
						Class("space-flow-item").
						Body(
							newSupportUsItem().
								Class("fill").
								Label("Open collective").
								Href(openCollectiveURL).
								Icon(
									newSVGIcon().
										RawSVG(opensourceSVG).
										Size(supportUsIconSize),
								),
						),
				),
		)
}

type supportUsItem struct {
	app.Compo

	Iclass string
	Ilabel string
	Ihref  string
	Iicon  app.UI
}

func newSupportUsItem() *supportUsItem {
	return &supportUsItem{}
}

func (i *supportUsItem) Class(v string) *supportUsItem {
	if v == "" {
		return i
	}
	if i.Iclass != "" {
		i.Iclass += " "
	}
	i.Iclass += v
	return i
}

func (i *supportUsItem) Label(v string) *supportUsItem {
	i.Ilabel = v
	return i
}

func (i *supportUsItem) Href(v string) *supportUsItem {
	i.Ihref = v
	return i
}

func (i *supportUsItem) Icon(v app.UI) *supportUsItem {
	i.Iicon = v
	return i
}

func (i *supportUsItem) Render() app.UI {
	return app.A().
		Class("support-us-item").
		Class("hspace-in-stretch").
		Class("vspace-in").
		Class("vignette").
		Class("magnify").
		Class(i.Iclass).
		Href(i.Ihref).
		Body(
			app.Stack().
				Class("fill").
				Center().
				Content(
					app.Div().
						Class("fit").
						Class("center").
						Body(
							app.Div().
								Class("fit").
								Class("center").
								Body(i.Iicon),
							app.Div().
								Class("support-us-item-label").
								Class("heading").
								Class("default").
								Class("center").
								Text(i.Ilabel),
						),
				),
		)
}
