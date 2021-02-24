package main

import "github.com/maxence-charriere/go-app/v8/pkg/app"

type markdownDoc struct {
	app.Compo
}

func newMarkdownDoc() *markdownDoc {
	return &markdownDoc{}
}

func (d *markdownDoc) Render() app.UI {
	return newPage2()
}
