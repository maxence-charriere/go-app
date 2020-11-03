package main

import (
	"strings"

	"github.com/maxence-charriere/go-app/v7/pkg/app"
)

type tableOfContents struct {
	app.Compo

	Ilinks []string
}

func newTableOfContents() *tableOfContents {
	return &tableOfContents{}
}

func (t *tableOfContents) Links(v ...string) *tableOfContents {
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
						Href(githubIndex(link)).
						Text(link)
				}),
			),
		)
}

func githubIndex(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	return "#" + strings.ReplaceAll(s, " ", "-")
}
