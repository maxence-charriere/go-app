package main

import (
	"path/filepath"

	"github.com/maxence-charriere/go-app/v7/pkg/app"
)

func newStart() app.UI {
	return newPage().
		Path("/web/documents/start.md").
		TableOfContents(
			contentLink{
				Name: "Getting started",
				URL:  "#getting-started",
			},
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
			contentLink{
				Name: "Next",
				URL:  "#next",
			},
		)
}

func newArchitecture() app.UI {
	return newPage().
		Path("/web/documents/architecture.md").
		TableOfContents(
			contentLink{
				Name: "Architecture",
				URL:  "#architecture",
			},
			contentLink{
				Name: "Web browser",
				URL:  "#web-browser",
			},
			contentLink{
				Name: "Server",
				URL:  "#server",
			},
			contentLink{
				Name: "App",
				URL:  "#app",
			},
			contentLink{
				Name: "Static resources",
				URL:  "#static-resources",
			},
			contentLink{
				Name: "Next",
				URL:  "#next",
			},
		)
}

type page struct {
	app.Compo

	path  string
	links []contentLink
}

func newPage() *page {
	return &page{}
}

func (p *page) Path(v string) *page {
	p.path = v
	return p
}

func (p *page) TableOfContents(v ...contentLink) *page {
	p.links = v
	return p
}

func (p *page) Render() app.UI {
	return app.Shell().
		Class("app-background").
		Menu(Menu()).
		Submenu(
			newTableOfContents().
				Links(p.links...),
		).
		OverlayMenu(Menu()).
		Content(
			newDocument(p.path).
				Description(filepath.Base(p.path)),
		)
}
