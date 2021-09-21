package main

import (
	"fmt"
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/analytics"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type statesPage struct {
	app.Compo
}

func newStatesPage() *statesPage {
	return &statesPage{}
}

func (p *statesPage) OnPreRender(ctx app.Context) {
	p.initPage(ctx)
}

func (p *statesPage) OnNav(ctx app.Context) {
	p.initPage(ctx)
}

func (p *statesPage) initPage(ctx app.Context) {
	ctx.Page().SetTitle("State Management")
	ctx.Page().SetDescription("Documentation about how to set and observe states.")
	analytics.Page("states", nil)
}

func (p *statesPage) Render() app.UI {
	return newPage().
		Title("State Management").
		Icon(stateSVG).
		Index(
			newIndexLink().Title("What is a state?"),
			newIndexLink().Title("Set"),
			newIndexLink().Title("    Options"),
			newIndexLink().Title("Observe"),
			newIndexLink().Title("    Conditional Observation"),
			newIndexLink().Title("    Additional Instructions"),
			newIndexLink().Title("Get"),

			app.Div().Class("separator"),

			newIndexLink().Title("Next"),
		).
		Content(
			newRemoteMarkdownDoc().Src("/web/documents/states.md"),
		)
}

func handleGreet(ctx app.Context, a app.Action) {
	var name string
	ctx.GetState("greet-name", &name)

	// ...
}

type hello struct {
	app.Compo
	name string
}

func (h *hello) OnMount(ctx app.Context) {
	ctx.ObserveState("greet-name").
		OnChange(func() {
			fmt.Println("greet-name was changed at", time.Now())
		}).
		Value(&h.name)
}
