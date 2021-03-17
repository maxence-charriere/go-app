package main

import (
	"fmt"
	"path"

	"github.com/maxence-charriere/go-app/v8/pkg/app"
)

type example struct {
	app.Compo

	example     app.UI
	examplePath string
}

func newExample() *example {
	return &example{}
}

func (e *example) OnPreRender(ctx app.Context) {
	e.init(ctx)
}

func (e *example) OnNav(ctx app.Context) {
	e.init(ctx)
}

func (e *example) init(ctx app.Context) {
	eg := ctx.Page.URL().Path
	if eg == e.examplePath {
		return
	}

	switch eg {
	case "/examples/hello":
		e.example = newHello()

	case "/examples/list":
		e.example = newFoodList()
	}
	e.examplePath = eg

	e.Update()
}

func (e *example) Render() app.UI {
	return newPage().
		Index(
			newIndexLink().
				Label("Hello").
				Href("/examples/hello").
				Focus("/examples/hello" == e.examplePath),
			newIndexLink().
				Label("List").
				Href("/examples/list").
				Focus("/examples/list" == e.examplePath),
		).
		Content(
			app.Div().
				Class("hspace-out-stretch").
				Body(e.example),
		).
		IssueTitle(fmt.Sprintf("%s example", path.Base(e.examplePath)))
}
