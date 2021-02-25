package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/maxence-charriere/go-app/v8/pkg/app"
	"github.com/maxence-charriere/go-app/v8/pkg/errors"
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
		"/start": {
			path:        "/web/documents/start.md",
			title:       defaultTitle,
			description: defaultDescription,
			index: []string{
				"Getting started",
				"Prerequisite",
				"Install",
				"User interface",
				"Server",
				"Build and run",
				"Tips",
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
				"App",
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
				"    Update mechanism",
				"Lifecycle",
				"    OnMount",
				"    OnNav",
				"    OnDismount",
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
				"Standard goroutines",
				"    When to use?",
				"Dispatch()",
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
				"Next",
			},
		},
		"/js": {
			path:        "/web/documents/js.md",
			title:       defaultTitle,
			description: defaultDescription,
			index: []string{
				"Javascript and DOM",
				"Include JS files",
				"    Handler",
				"    Inline",
				"Window",
				"    Get element by ID",
				"    Create JS object",
				"Cancel an event",
				"Get input value",
				"Next",
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
		"/lifecycle": {
			path:        "/web/documents/lifecycle.md",
			title:       defaultTitle,
			description: defaultDescription,
			index: []string{
				"Lifecycle",
				"    First loading",
				"    Recurrent loadings",
				"    Loading after app update",
				"Listen for updates",
			},
		},
	}
}

type markdownDoc struct {
	app.Compo

	path      string
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
	if d.path == ctx.Page.URL().Path {
		return
	}

	page, ok := mardownPages()[ctx.Page.URL().Path]
	if !ok {
		app.Log("%s", logs.New("markdown page not found").
			Tag("url", ctx.Page.URL().String()),
		)
		return
	}

	d.path = page.path
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
		md, err := d.get(ctx, path)

		d.Defer(func(ctx app.Context) {
			d.markdow = md
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

func (d *markdownDoc) get(ctx app.Context, path string) (string, error) {
	url := *ctx.Page.URL()
	url.Path = path
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return "", errors.New("creating request failed").
			Tag("url", url).
			Wrap(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.New("getting document failed").
			Tag("url", url).
			Wrap(err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return "", errors.New(res.Status)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.New("reading document failed").Wrap(err)
	}
	return fmt.Sprintf("<div>%s</div>", parseMarkdown(b)), nil
}

func (d *markdownDoc) highlightCode(ctx app.Context) {
	ctx.Dispatch(func() {
		app.Window().Get("Prism").Call("highlightAll")
	})
}

func (d *markdownDoc) Render() app.UI {
	return newPage2().
		Index(
			newIndex().Links(d.index...),
		).
		Content(
			app.Article().
				Class("markdown").
				Class("hspace-out-stretch").
				Body(
					app.Raw(d.markdow),
				),
			newLoader().
				Class("page-loader").
				Class("fill").
				Title("Loading").
				Description(filepath.Base(d.path)).
				Loading(d.isLoading).
				Err(d.err),
		)
}

func parseMarkdown(md []byte) []byte {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)
	return markdown.ToHTML(md, parser, nil)
}
