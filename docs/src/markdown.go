package main

import (
	"fmt"
	"path/filepath"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/maxence-charriere/go-app/v8/pkg/app"
	"github.com/maxence-charriere/go-app/v8/pkg/logs"
)

type markdownPage struct {
	path        string
	title       string
	description string
	index       []string
}

func mardownPages() map[string]markdownPage {
	return map[string]markdownPage{
		"/": {
			path:        "/web/documents/home.md",
			title:       defaultTitle,
			description: "Introduction abouWorld with Go (Golang) and WebAssembly (WASM).",
			index: []string{
				"go-app",
				"Declarative syntax",
				"Standard HTTP",
				"Other features",
				"Next",
			},
		},
		"/start": {
			path:        "/web/documents/start.md",
			title:       "Getting started building a Progressive Web App with Golang and WebAssembly",
			description: "Introduction about how to create a Progressive Web App displaying a simple Hello World with the Go programming language (Golang) and WebAssembly (Wasm).",
			index: []string{
				"Getting started",
				"Prerequisite",
				"Install",
				"Code",
				"Build and run",
				"    Building the client",
				"    Building the server",
				"    Launching the app",
				"    Tips",
				"Next",
			},
		},
		"/architecture": {
			path:        "/web/documents/architecture.md",
			title:       defaultTitle,
			description: defaultDescription,
			index: []string{
				"Architecture",
				"Web browser",
				"Server",
				"HTML pages",
				"Package resources",
				"app.wasm",
				"Static resources",
				"Next",
			},
		},
		"/components": {
			path:        "/web/documents/components.md",
			title:       defaultTitle,
			description: defaultDescription,
			index: []string{
				"Components",
				"Create",
				"Customize",
				"Update",
				"Fields",
				"    Exported fields",
				"    Internal fields",
				"Lifecycle",
				"    Prerender",
				"    Mount",
				"    Nav",
				"    Dismount",
				"Extensions",
				"Next",
			},
		},
		"/concurrency": {
			path:        "/web/documents/concurrency.md",
			title:       defaultTitle,
			description: defaultDescription,
			index: []string{
				"Concurrency",
				"UI goroutine",
				"Async",
				"Dispatch",
				"Defer",
				"Next",
			},
		},
		"/syntax": {
			path:        "/web/documents/syntax.md",
			title:       defaultTitle,
			description: defaultDescription,
			index: []string{
				"Declarative syntax",
				"HTML elements",
				"    Create",
				"    Standard elements",
				"    Self closing elements",
				"    Style",
				"    Attributes",
				"    Event handlers",
				"Text",
				"Raw elements",
				"Nested components",
				"Condition",
				"    If",
				"    ElseIf",
				"    Else",
				"Range",
				"    Slice",
				"    Map",
				"Form helpers",
				"    ValueTo",
				"Next",
			},
		},
		"/js": {
			path:        "/web/documents/js.md",
			title:       defaultTitle,
			description: defaultDescription,
			index: []string{
				"JavaScript and DOM",
				"Include JS files",
				"    Page's scope",
				"    Inlined in Components",
				"Using window global object",
				"    Get element by ID",
				"    Create JS object",
				"Cancel an event",
				"Get input value",
				"Next",
			},
		},
		"/lifecycle": {
			path:        "/web/documents/lifecycle.md",
			title:       defaultTitle,
			description: defaultDescription,
			index: []string{
				"Lifecycle",
				"    First loading",
				"    Recurrent loadings",
				"    Loading after app update",
				"Listen for app updates",
			},
		},
		"/routing": {
			path:        "/web/documents/routing.md",
			title:       defaultTitle,
			description: defaultDescription,
			index: []string{
				"Routing",
				"Define a route",
				"    Simple route",
				"    Route with regular expression",
				"Detect navigation",
				"Next",
			},
		},
		"/static-resources": {
			path:        "/web/documents/static-resources.md",
			title:       defaultTitle,
			description: defaultDescription,
			index: []string{
				"Static resources",
				"Access static resources",
				"    In Handler",
				"    In components",
				"    In CSS files",
				"Setup local web directory",
				"Setup remote web directory",
				"Fully static app",
				"Next",
			},
		},
		"/built-with": {
			path:        "/web/documents/built-with.md",
			title:       defaultTitle,
			description: defaultDescription,
			index: []string{
				"Built with go-app",
				"Lofimusic.app",
				"Murlok.io",
				"Astextract",
				"Next",
			},
		},
		"/install": {
			path:        "/web/documents/install.md",
			title:       defaultTitle,
			description: defaultDescription,
			index: []string{
				"Install",
				"Desktop",
				"IOS",
				"Android",
				"Next",
			},
		},
	}
}

type markdownDoc struct {
	app.Compo

	path      string
	title     string
	index     []string
	markdow   string
	isLoading bool
	err       error
}

func newMarkdownDoc() *markdownDoc {
	return &markdownDoc{}
}

func (d *markdownDoc) OnPreRender(ctx app.Context) {
	d.init(ctx)
}

func (d *markdownDoc) OnNav(ctx app.Context) {
	d.init(ctx)
}

func (d *markdownDoc) init(ctx app.Context) {
	path := ctx.Page.URL().Path
	if d.path == path {
		return
	}

	page, ok := mardownPages()[ctx.Page.URL().Path]
	if !ok {
		app.Log(logs.New("markdown page not found").
			Tag("url", ctx.Page.URL().String()),
		)
		return
	}

	d.path = path
	d.title = page.path
	d.index = page.index
	ctx.Page.SetTitle(page.title)
	ctx.Page.SetDescription(page.description)

	d.Update()
	d.load(ctx, page.path)
}

func (d *markdownDoc) load(ctx app.Context, path string) {
	d.isLoading = true
	d.err = nil
	d.Update()

	ctx.Async(func() {
		md, err := get(ctx, path)

		d.Defer(func(ctx app.Context) {
			d.markdow = fmt.Sprintf("<div>%s</div>", parseMarkdown(md))
			d.err = err
			d.isLoading = false
			d.Update()
			d.highlightCode(ctx)

			fragment := ctx.Page.URL().Fragment
			if fragment == "" {
				fragment = "top"
			}
			ctx.ScrollTo(fragment)
		})
	})
}

func (d *markdownDoc) highlightCode(ctx app.Context) {
	ctx.Dispatch(func() {
		app.Window().Get("Prism").Call("highlightAll")
	})
}

func (d *markdownDoc) Render() app.UI {
	return newPage().
		Index(
			newIndex().Links(d.index...),
		).
		Content(
			app.Div().
				Class("markdown").
				Class("hspace-out-stretch").
				Body(
					app.Raw(d.markdow),
				),
			newLoader().
				Class("page-loader").
				Class("fill").
				Title("Loading").
				Description(filepath.Base(d.title)).
				Loading(d.isLoading).
				Err(d.err),
		).
		IssueTitle(filepath.Base(d.title))
}

func parseMarkdown(md []byte) []byte {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)
	return markdown.ToHTML(md, parser, nil)
}
