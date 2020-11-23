package main

import (
	"fmt"
	"path"

	"github.com/maxence-charriere/go-app/v7/pkg/app"
)

func issue(src string) app.UI {
	return app.Div().Body(
		app.H2().
			ID("report-issue").
			Text("Report issue"),
		app.P().Body(
			app.Text("Found something incorrect, a typo or have suggestions to improve this article? "),
			app.A().
				Href(fmt.Sprintf(
					"%s/issues/new?title=Documentation issues in %s",
					githubURL,
					path.Base(src),
				)).
				Target("_black").
				Text("Let us know :)"),
		),
	)
}
