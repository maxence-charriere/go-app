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
						URL:  "#prerequisite",
					},
					contentLink{
						Name: "Install",
						URL:  "#install",
					},
					contentLink{
						Name: "User interface",
						URL:  "#user-interface",
					},
					contentLink{
						Name: "Server",
						URL:  "#server",
					},
					contentLink{
						Name: "Build and run",
						URL:  "#build-and-run",
					},
					contentLink{
						Name: "Tips",
						URL:  "#tips",
					},
				),
		).
		OverlayMenu(Menu()).
		Content(
			newDocument("/web/start.md").
				Description("start.md"),
		)
}
