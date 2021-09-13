package main

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type declarativeSyntaxPage struct {
	app.Compo
}

func newDeclarativeSyntaxPage() *declarativeSyntaxPage {
	return &declarativeSyntaxPage{}
}

func (p *declarativeSyntaxPage) Render() app.UI {
	return newPage().
		Title("Declarative Syntax").
		Icon(keyboardSVG)
}
