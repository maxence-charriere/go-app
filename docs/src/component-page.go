package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type componentsPage struct {
	app.Compo
}

func newComponentsPage() *componentsPage {
	return &componentsPage{}
}

func (p *componentsPage) OnNav(ctx app.Context) {}

func (p *componentsPage) Render() app.UI {
	return newPage().
		Title("Components").
		Icon(gridSVG).
		Index(
			newIndexLink().Title("What is a Component?"),
			newIndexLink().Title("Create"),
			newIndexLink().Title("Customize Look"),
			newIndexLink().Title("Fields"),
			newIndexLink().Title("    Exported vs Unexported"),
			newIndexLink().Title("    How chose between Exported and Unexported?"),
			newIndexLink().Title("Lifecycle Events"),
			newIndexLink().Title("    PreRender"),
			newIndexLink().Title("    Mount"),
			newIndexLink().Title("    Nav"),
			newIndexLink().Title("    Dismount"),
			newIndexLink().Title("    Lifecycle Events Reference"),

			newIndexLink().Title("Fields"),
			newIndexLink().Title("Fields"),

			app.Div().Class("separator"),

			newIndexLink().Title("Next"),
		).
		Content(
			newRemoteMarkdownDoc().Src("/web/documents/components.md"),
		)
}
