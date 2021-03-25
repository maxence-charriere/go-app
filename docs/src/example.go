package main

import (
	"fmt"
	"path"
	"strings"

	"github.com/maxence-charriere/go-app/v8/pkg/app"
)

type exampleData struct {
	component   app.UI
	title       string
	description string
	path        string
	code        string
}

type example struct {
	app.Compo

	examples []exampleData
	current  exampleData
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
	defer e.Update()

	e.examples = []exampleData{
		{
			component:   newHello(),
			path:        "/examples/hello",
			title:       "Hello World Example",
			description: "A Hello World example built with go-app.",
			code:        egHelloCode,
		},
		{
			component:   newFoodList(),
			path:        "/examples/list",
			title:       "List Example",
			description: "An example that shows how to deal with a list using go-app.",
			code:        egListCode,
		},
	}

	eg := ctx.Page.URL().Path
	for _, ed := range e.examples {
		if eg == ed.path {
			e.current = ed
			ctx.Page.SetTitle(ed.title)
			ctx.Page.SetDescription(ed.description)
			return
		}
	}

	e.current = exampleData{}
	ctx.Page.SetTitle("Examples built with go-app")
	ctx.Page.SetDescription("Examples built with go-app.")
}

func (e *example) Render() app.UI {
	return newPage().
		Index(
			newIndexLink().
				Label("Hello").
				Href("/examples/hello").
				Focus("/examples/hello" == e.current.path),
			newIndexLink().
				Label("List").
				Href("/examples/list").
				Focus("/examples/list" == e.current.path),
		).
		Content(
			app.Div().
				Class("hspace-out-stretch").
				Body(
					app.If(e.current.path == "",
						app.H1().
							Class("h1").
							Class("header-separator").
							Text("Examples built with go-app"),
						app.P().
							Class("center-content").
							Body(
								app.Range(e.examples).Slice(func(i int) app.UI {
									eg := e.examples[i]
									return app.A().
										Class("button").
										Href(eg.path).
										Text(strings.TrimSuffix(eg.title, " Example"))
								}),
							),
					).
						Else(
							e.current.component,
							app.Div().
								Class("vspace-section").
								Body(
									app.H2().
										Class("h2").
										Class("header-separator").
										Text("Code"),
									newMarkdownContent().Markdown(e.current.code),
								),
						),
				),
		).
		IssueTitle(fmt.Sprintf("%s example", path.Base(e.current.path)))
}
