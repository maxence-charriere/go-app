package main

import "github.com/maxence-charriere/go-app/v7/pkg/app"

type start struct {
	app.Compo
}

func newStart() *start {
	return &start{}
}

func (s *start) Render() app.UI {
	return app.Shell().
		Class("app-background").
		Menu(Menu()).
		Submenu(
			newTableOfContents().
				Links(
					contentLink{
						Name: "Prerequisite",
						URL:  "#Prerequisite",
					},
					contentLink{
						Name: "Install",
						URL:  "#Install",
					},
					contentLink{
						Name: "User interface",
						URL:  "#UserInterface",
					},
					contentLink{
						Name: "Server",
						URL:  "#Server",
					},
					contentLink{
						Name: "Build and run",
						URL:  "#BuildAndRun",
					},
					contentLink{
						Name: "Tips",
						URL:  "#Tips",
					},
				),
		).
		OverlayMenu(Menu()).
		Content(
			newDocument("/web/start.md").
				Description("start.md"),
		)
}
