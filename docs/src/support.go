package main

import "github.com/maxence-charriere/go-app/v7/pkg/app"

func support() app.UI {
	return app.Div().Body(
		app.H2().
			ID("support-go-app").
			Text("Support go-app"),
		app.P().Body(
			app.Text("Hello there, I'am Maxence, the creator of "),
			app.B().Text("go-app"),
			app.Text(". I hope you like what you are seeing here and will consider or already use this package to build your next project."),
		),
		app.P().Body(
			app.Text("A lot of hard work is put to develop this package. Not mandatory but heartly appreciated, "),
			app.A().
				Href(buyMeACoffeeURL).
				Target("_blank").
				Text("buying me a coffee "),
			app.Text("or be part of our "),
			app.A().
				Href(openCollectiveURL).
				Target("_blank").
				Text("open collective "),
			app.Text("is always a boost and bring great motivation to keep the good work."),
			app.Flow().
				StrechtOnSingleRow().
				ItemsBaseWidth(192).
				Content(
					newSupportPartner().
						Name("Buy me a coffee").
						URL(buyMeACoffeeURL).
						Icon(`
						<svg style="width:48px;height:48px" viewBox="0 0 24 24">
							<path fill="currentColor" d="M2,21H20V19H2M20,8H18V5H20M20,3H4V13A4,4 0 0,0 8,17H14A4,4 0 0,0 18,13V10H20A2,2 0 0,0 22,8V5C22,3.89 21.1,3 20,3Z" />
						</svg>
						`),
					newSupportPartner().
						Name("Open Collective").
						URL(openCollectiveURL).
						Icon(`
						<svg style="width:48px;height:48px" viewBox="0 0 24 24">
							<path fill="currentColor" d="M15.41,22C15.35,22 15.28,22 15.22,22C15.1,21.95 15,21.85 14.96,21.73L12.74,15.93C12.65,15.69 12.77,15.42 13,15.32C13.71,15.06 14.28,14.5 14.58,13.83C15.22,12.4 14.58,10.73 13.15,10.09C11.72,9.45 10.05,10.09 9.41,11.5C9.11,12.21 9.09,13 9.36,13.69C9.66,14.43 10.25,15 11,15.28C11.24,15.37 11.37,15.64 11.28,15.89L9,21.69C8.96,21.81 8.87,21.91 8.75,21.96C8.63,22 8.5,22 8.39,21.96C3.24,19.97 0.67,14.18 2.66,9.03C4.65,3.88 10.44,1.31 15.59,3.3C18.06,4.26 20.05,6.15 21.13,8.57C22.22,11 22.29,13.75 21.33,16.22C20.32,18.88 18.23,21 15.58,22C15.5,22 15.47,22 15.41,22M12,3.59C7.03,3.46 2.9,7.39 2.77,12.36C2.68,16.08 4.88,19.47 8.32,20.9L10.21,16C8.38,15 7.69,12.72 8.68,10.89C9.67,9.06 11.96,8.38 13.79,9.36C15.62,10.35 16.31,12.64 15.32,14.47C14.97,15.12 14.44,15.65 13.79,16L15.68,20.93C17.86,19.95 19.57,18.16 20.44,15.93C22.28,11.31 20.04,6.08 15.42,4.23C14.33,3.8 13.17,3.58 12,3.59Z" />
						</svg>
						`),
				),
		),
	)
}

type supportPartner struct {
	app.Compo

	Iname string
	Iurl  string
	Iicon string
}

func (p *supportPartner) Name(v string) *supportPartner {
	p.Iname = v
	return p
}

func (p *supportPartner) URL(v string) *supportPartner {
	p.Iurl = v
	return p
}

func (p *supportPartner) Icon(v string) *supportPartner {
	p.Iicon = v
	return p
}

func newSupportPartner() *supportPartner {
	return &supportPartner{}
}

func (p *supportPartner) Render() app.UI {
	return app.A().
		Class("support-partner").
		Href(p.Iurl).
		Target("_blank").
		Title("Contribute").
		Body(
			app.Stack().
				Center().
				Vertical().
				Content(
					app.Raw(p.Iicon),
					app.Div().
						Class("support-partner-name").
						Text(p.Iname),
				),
		)
}
