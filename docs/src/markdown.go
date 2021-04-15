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
			description: defaultDescription,
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
			title:       "Getting started building a Progressive Web App with Go and WebAssembly",
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
		"/v7-to-v8": {
			path:        "/web/documents/v7-to-v8.md",
			title:       "Migrating a go-app Progressive Web App from V7 to V8",
			description: "Guide about how to migrate a Progressive Web App built with go-app V7 to V8.",
			index: []string{
				"V7 to V8 migration guide",
				"Build directives",
				"Routing",
				"Package functions",
				"Component interfaces",
				"Resource provider",
				"Concurrency",
				"Next",
			},
		},
		"/architecture": {
			path:        "/web/documents/architecture.md",
			title:       "Understanding go-app architecture",
			description: "Documentation about how go-app parts are working together to deliver Progressive Web Apps out of the box?",
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
			title:       "Building components: customizable, independent, and reusable UI elements",
			description: "Documentation about building customizable, independent, and reusable UI elements.",
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
			title:       "Building responsive Progressive Web Apps",
			description: "Documentation about go-app tools that help to build reactive and concurrency safe Progressive Web Apps.",
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
			title:       "A Go syntax for building beautiful UIs",
			description: "Documentation about the Go (Golang) syntax to customize go-app components look, and craft beautiful UIs only with the Go Programming Language.",
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
			title:       "Interoperability between Go and JavaScript",
			description: "Documentation about how to interact with the webpage DOM and JavaScript libraries from Go (Golang).",
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
			title:       "Understanding app lifecycle in the web browser",
			description: "Documentation that describes how a web browser installs and updates a go-app Progressive Web App.",
			index: []string{
				"Lifecycle",
				"    First loading",
				"    Recurrent loadings",
				"    Loading after app update",
				"Listen for app updates",
				"Next",
			},
		},
		"/routing": {
			path:        "/web/documents/routing.md",
			title:       "Routing pages to go-app components",
			description: "Documentation about how to associate URL paths to go-app components.",
			index: []string{
				"Routing",
				"Define a route",
				"    Simple route",
				"    Route with regular expression",
				"Detect navigation",
				"Next",
			},
		},
		"/seo": {
			path:        "/web/documents/seo.md",
			title:       "Building an SEO friendly Progressive Web App",
			description: "Documentation about how to make a Progressive Web App indexable by search engines with go-app package.",
			index: []string{
				"SEO",
				"Prerendering",
				"    Customizing prerendering",
				"    Customizing page metadata",
				"    Caching",
				"Next",
			},
		},
		"/static-resources": {
			path:        "/web/documents/static-resources.md",
			title:       "Dealing with static resources",
			description: "Documentation that describes what are static resources, how to interact, and where to host them.",
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
		"/testing": {
			path:        "/web/documents/testing.md",
			title:       "Testing components",
			description: "Documentation about how to unit test components created with go-app.",
			index: []string{
				"Testing",
				"Component server prerendering",
				"Component client lifecycle",
				"Asynchronous operations",
				"UI elements",
				"Next",
			},
		},
		"/built-with": {
			path:        "/web/documents/built-with.md",
			title:       "Progressive Web Apps created with go-app",
			description: "An index that lists Progressive Web Apps built with the go-app package.",
			index: []string{
				"Built with go-app",
				"Lofimusic.app",
				"Murlok.io",
				"Astextract",
				"Liwasc",
				"Next",
			},
		},
		"/install": {
			path:        "/web/documents/install.md",
			title:       "Installing Progressive Web Apps on user devices",
			description: "Documentation about how to install a Progressive Web App on a user device, from Chromium-based web browsers to mobile devices.",
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
	markdown  string
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
			d.markdown = string(md)
			d.err = err
			d.isLoading = false
			d.Update()

			fragment := ctx.Page.URL().Fragment
			if fragment == "" {
				fragment = "top"
			}
			ctx.ScrollTo(fragment)
		})
	})
}

func (d *markdownDoc) Render() app.UI {
	return newPage().
		Index(
			newIndex().Links(d.index...),
		).
		Content(
			newMarkdownContent().
				Class("hspace-out-stretch").
				Markdown(d.markdown),
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

type markdownContent struct {
	app.Compo

	Iclass string
	Imd    string
	md     string
}

func newMarkdownContent() *markdownContent {
	return &markdownContent{}
}

func (m *markdownContent) Class(v string) *markdownContent {
	if v == "" {
		return m
	}
	if m.Iclass != "" {
		m.Iclass += " "
	}
	m.Iclass += v
	return m
}

func (m *markdownContent) Markdown(v string) *markdownContent {
	m.Imd = v
	return m
}

func (m *markdownContent) Render() app.UI {
	if m.Imd != m.md {
		m.Defer(m.highlightCode)
	}

	return app.Div().
		Class("markdown").
		Class(m.Iclass).
		Body(
			app.Raw(fmt.Sprintf("<div>%s</div>", parseMarkdown([]byte(m.Imd)))),
		)
}

func (d *markdownContent) highlightCode(ctx app.Context) {
	app.Window().Get("Prism").Call("highlightAll")
}

func parseMarkdown(md []byte) []byte {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)
	return markdown.ToHTML(md, parser, nil)
}
