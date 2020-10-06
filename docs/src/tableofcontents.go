package main

import "github.com/maxence-charriere/go-app/v7/pkg/app"

type contentLink struct {
	Name string
	URL  string
}

type tableOfContents struct {
	app.Compo

	Ilinks []contentLink
}

func newTableOfContents() *tableOfContents {
	return &tableOfContents{}
}

func (t *tableOfContents) Links(v ...contentLink) *tableOfContents {
	t.Ilinks = v
	return t
}

func (t *tableOfContents) Render() app.UI {
	return app.Aside().
		Class("pane").
		Class("index").
		Body(
			app.H1().Text("Index"),
			app.Section().Body(
				app.Range(t.Ilinks).Slice(func(i int) app.UI {
					link := t.Ilinks[i]

					return app.A().
						Href(link.URL).
						Text(link.Name)
				}),
			),
		)
}
