package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/analytics"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type componentsPage struct {
	app.Compo
}

func newComponentsPage() *componentsPage {
	return &componentsPage{}
}

func (p *componentsPage) OnPreRender(ctx app.Context) {
	p.initPage(ctx)
}

func (p *componentsPage) OnNav(ctx app.Context) {
	p.initPage(ctx)
}

func (p *componentsPage) initPage(ctx app.Context) {
	ctx.Page().SetTitle("Building Components: Customizable, Independent, and Reusable UI Elements")
	ctx.Page().SetDescription("Documentation about building customizable, independent, and reusable UI elements.")
	analytics.Page("components", nil)
}

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
			newIndexLink().Title("    Reference"),
			newIndexLink().Title("Updates"),
			newIndexLink().Title("    Manually Trigger an Update"),

			app.Div().Class("separator"),

			newIndexLink().Title("Next"),
		).
		Content(
			newRemoteMarkdownDoc().Src("/web/documents/components.md"),
		)
}
