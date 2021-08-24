package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/maxence-charriere/go-app/v9/pkg/ui"
)

type githubSponsor struct {
	app.Compo

	Iclass string
}

func newGithubSponsor() *githubSponsor {
	return &githubSponsor{}
}

func (s *githubSponsor) Class(v string) *githubSponsor {
	s.Iclass = app.AppendClass(s.Iclass, v)
	return s
}

func (s *githubSponsor) Render() app.UI {
	return ui.Stack().
		Class(s.Iclass).
		Center().
		Middle().
		Content(
			app.Aside().
				Class("magnify").
				Class("text-center").
				Body(
					app.A().
						Class("default").
						Href(githubSponsorURL).
						Body(
							ui.Icon().
								Class("center").
								Class("icon-top").
								Size(72).
								Src(githubSVG),
							app.Header().
								Class("h3").
								Class("default").
								Text("Support on GitHub"),
							app.P().
								Class("subtext").
								Text("Help with go-app development by sponsoring it on GitHub."),
						),
				),
		)
}
