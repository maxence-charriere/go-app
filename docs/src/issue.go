package main

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v8/pkg/app"
)

type issue struct {
	app.Compo

	Iclass string
	Ititle string
}

func newIssue() *issue {
	return &issue{}
}

func (i *issue) Class(v string) *issue {
	if v == "" {
		return i
	}
	if i.Iclass != "" {
		i.Iclass += " "
	}
	i.Iclass += v
	return i
}

func (i *issue) Title(v string) *issue {
	i.Ititle = v
	return i
}

func (i *issue) Render() app.UI {
	return app.Div().
		Class(i.Iclass).
		Body(
			app.Div().
				ID("report-issue").
				Class("h2").
				Class("header-separator").
				Text("Report issue"),
			app.P().Body(
				app.Text("Found something incorrect, a typo or have suggestions to improve this article? "),
				app.A().
					Href(fmt.Sprintf(
						"%s/issues/new?title=Documentation issues in %s",
						githubURL,
						i.Ititle,
					)).
					Text("Let me know :)"),
			),
		)
}
