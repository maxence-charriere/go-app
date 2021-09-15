package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/analytics"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type declarativeSyntaxPage struct {
	app.Compo
}

func newDeclarativeSyntaxPage() *declarativeSyntaxPage {
	return &declarativeSyntaxPage{}
}

func (p *declarativeSyntaxPage) OnPreRender(ctx app.Context) {
	p.initPage(ctx)
}

func (p *declarativeSyntaxPage) OnNav(ctx app.Context) {
	p.initPage(ctx)
}

func (p *declarativeSyntaxPage) initPage(ctx app.Context) {
	ctx.Page().SetTitle("Customize Components with go-app Declarative Syntax")
	ctx.Page().SetDescription("Documentation about how to customize components with go-app declarative syntax.")
	analytics.Page("declarative-syntax", nil)
}

func (p *declarativeSyntaxPage) Render() app.UI {
	return newPage().
		Title("Declarative Syntax").
		Icon(keyboardSVG).
		Index(
			newIndexLink().Title("Intro"),
			newIndexLink().Title("HTML Elements"),
			newIndexLink().Title("    Create"),
			newIndexLink().Title("    Standard Elements"),
			newIndexLink().Title("    Self Closing Elements"),
			newIndexLink().Title("    Attributes"),
			newIndexLink().Title("    Style"),
			newIndexLink().Title("    Event handlers"),
			newIndexLink().Title("Raw elements"),
			newIndexLink().Title("Nested Components"),
			newIndexLink().Title("Condition"),
			newIndexLink().Title("    If"),
			newIndexLink().Title("    ElseIf"),
			newIndexLink().Title("    Else"),
			newIndexLink().Title("Range"),
			newIndexLink().Title("    Slice"),
			newIndexLink().Title("    Map"),
			newIndexLink().Title("Form helpers"),
			newIndexLink().Title("    ValueTo"),

			app.Div().Class("separator"),

			newIndexLink().Title("Next"),
		).
		Content(
			newRemoteMarkdownDoc().Src("/web/documents/declarative-syntax.md"),
		)
}
