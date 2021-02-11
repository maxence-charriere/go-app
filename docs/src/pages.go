package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/maxence-charriere/go-app/v8/pkg/app"
)

const (
	buyMeACoffeeURL   = "https://www.buymeacoffee.com/maxence"
	openCollectiveURL = "https://opencollective.com/go-app"
	githubURL         = "https://github.com/maxence-charriere/go-app"
	githubSponsorURL  = "https://github.com/sponsors/maxence-charriere"
	twitterURL        = "https://twitter.com/jonhymaxoo"
)

type pageInfo struct {
	MarkdownPath    string
	TableOfContents []string
}

func pages() map[string]pageInfo {
	start := pageInfo{
		MarkdownPath: "/web/documents/start.md",
		TableOfContents: []string{
			"Getting started",
			"Prerequisite",
			"Install",
			"User interface",
			"Server",
			"Build and run",
			"Tips",
			"Next",
		},
	}

	return map[string]pageInfo{
		"/":      start,
		"/start": start,
		"/architecture": {
			MarkdownPath: "/web/documents/architecture.md",
			TableOfContents: []string{
				"Architecture",
				"Web browser",
				"Server",
				"App",
				"Static resources",
				"Next",
			},
		},
		"/components": {
			MarkdownPath: "/web/documents/components.md",
			TableOfContents: []string{
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
			MarkdownPath: "/web/documents/concurrency.md",
			TableOfContents: []string{
				"Concurrency",
				"UI goroutine",
				"Standard goroutines",
				"    When to use?",
				"Dispatch()",
				"Next",
			},
		},
		"/syntax": {
			MarkdownPath: "/web/documents/syntax.md",
			TableOfContents: []string{
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
			MarkdownPath: "/web/documents/js.md",
			TableOfContents: []string{
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
			MarkdownPath: "/web/documents/routing.md",
			TableOfContents: []string{
				"Routing",
				"Define a route",
				"    Simple route",
				"    Route with regular expression",
				"Detect navigation",
				"Next",
			},
		},
		"/static-resources": {
			MarkdownPath: "/web/documents/static-resources.md",
			TableOfContents: []string{
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
			MarkdownPath: "/web/documents/built-with.md",
			TableOfContents: []string{
				"Built with go-app",
				"Lofimusic.app",
				"Murlok.io",
				"Astextract",
				"Next",
			},
		},
		"/install": {
			MarkdownPath: "/web/documents/install.md",
			TableOfContents: []string{
				"Install",
				"Desktop",
				"IOS",
				"Android",
				"Next",
			},
		},
		"/lifecycle": {
			MarkdownPath: "/web/documents/lifecycle.md",
			TableOfContents: []string{
				"Lifecycle",
				"    First loading",
				"    Recurrent loadings",
				"    Loading after app update",
				"Listen for updates",
			},
		},
	}
}

type page struct {
	app.Compo

	markdownPath    string
	tableOfContents []string
}

func newPage() *page {
	return &page{}
}

func (p *page) OnPreRender(ctx app.Context) {
	p.init(ctx)
}

func (p *page) OnNav(ctx app.Context) {
	p.init(ctx)
}

func (p *page) init(ctx app.Context) {
	path := ctx.Page.URL().Path
	info := pages()[path]

	p.markdownPath = info.MarkdownPath
	if app.IsServer {
		u := *ctx.Page.URL()
		u.Path = info.MarkdownPath
		p.markdownPath = u.String()
	}

	p.tableOfContents = info.TableOfContents

	title := strings.TrimPrefix(path, "/")
	title = strings.ReplaceAll(title, "-", " ")
	title = strings.Title(title)
	title = fmt.Sprintf("go-app • %s Documentation", title)
	ctx.Page.SetTitle(title)

	p.Update()
}

func (p *page) Render() app.UI {
	return app.Shell().
		Class("app-background").
		Menu(&menu{}).
		Submenu(
			newTableOfContents().
				Links(p.tableOfContents...),
		).
		OverlayMenu(&overlayMenu{}).
		Content(
			newDocument(p.markdownPath).
				Description(filepath.Base(p.markdownPath)),
		)
}
