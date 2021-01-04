package main

import (
	"path/filepath"

	"github.com/maxence-charriere/go-app/v7/pkg/app"
)

const (
	buyMeACoffeeURL   = "https://www.buymeacoffee.com/maxence"
	openCollectiveURL = "https://opencollective.com/go-app"
	githubURL         = "https://github.com/maxence-charriere/go-app"
	githubSponsorURL  = "https://github.com/sponsors/maxence-charriere"
	twitterURL        = "https://twitter.com/jonhymaxoo"
)

func pages() map[string]func() app.UI {
	return map[string]func() app.UI{
		"":                 newStart,
		"start":            newStart,
		"architecture":     newArchitecture,
		"reference":        newReference,
		"components":       newCompo,
		"concurrency":      newConcurrency,
		"syntax":           newSyntax,
		"js":               newJS,
		"routing":          newRouting,
		"static-resources": newStaticResources,
		"built-with":       newBuiltWith,
		"install":          newInstall,
		"lifecycle":        newLifecycle,
	}
}

func newStart() app.UI {
	return newPage().
		Path("/web/documents/start.md").
		TableOfContents(
			"Getting started",
			"Prerequisite",
			"Install",
			"User interface",
			"Server",
			"Build and run",
			"Tips",
			"Next",
		)
}

func newArchitecture() app.UI {
	return newPage().
		Path("/web/documents/architecture.md").
		TableOfContents(
			"Architecture",
			"Web browser",
			"Server",
			"App",
			"Static resources",
			"Next",
		)
}

func newCompo() app.UI {
	return newPage().
		Path("/web/documents/components.md").
		TableOfContents(
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
		)
}

func newConcurrency() app.UI {
	return newPage().
		Path("/web/documents/concurrency.md").
		TableOfContents(
			"Concurrency",
			"UI goroutine",
			"Standard goroutines",
			"    When to use?",
			"Dispatch()",
			"Next",
		)
}

func newSyntax() app.UI {
	return newPage().
		Path("/web/documents/syntax.md").
		TableOfContents(
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
		)
}

func newJS() app.UI {
	return newPage().
		Path("/web/documents/js.md").
		TableOfContents(
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
		)
}

func newRouting() app.UI {
	return newPage().
		Path("/web/documents/routing.md").
		TableOfContents(
			"Routing",
			"Define a route",
			"    Simple route",
			"    Route with regular expression",
			"Detect navigation",
			"Next",
		)
}

func newStaticResources() app.UI {
	return newPage().
		Path("/web/documents/static-resources.md").
		TableOfContents(
			"Static resources",
			"Access static resources",
			"    In Handler",
			"    In components",
			"    In CSS files",
			"Setup local web directory",
			"Setup remote web directory",
			"Fully static app",
			"Next",
		)
}

func newBuiltWith() app.UI {
	return newPage().
		Path("/web/documents/built-with.md").
		TableOfContents(
			"Built with go-app",
			"Lofimusic.app",
			"Murlok.io",
			"Astextract",
			"Next",
		)
}

func newInstall() app.UI {
	return newPage().
		Path("/web/documents/install.md").
		TableOfContents(
			"Install",
			"Desktop",
			"IOS",
			"Android",
			"Next",
		)
}

func newLifecycle() app.UI {
	return newPage().
		Path("/web/documents/lifecycle.md").
		TableOfContents(
			"Lifecycle",
			"    First loading",
			"    Recurrent loadings",
			"    Loading after app update",
			"Listen for updates",
		)
}

type page struct {
	app.Compo

	path  string
	links []string
}

func newPage() *page {
	return &page{}
}

func (p *page) Path(v string) *page {
	p.path = v
	return p
}

func (p *page) TableOfContents(v ...string) *page {
	p.links = v
	return p
}

func (p *page) OnAppUpdate(ctx app.Context) {
}

func (p *page) Render() app.UI {
	return app.Shell().
		Class("app-background").
		Menu(&menu{}).
		Submenu(
			newTableOfContents().
				Links(p.links...),
		).
		OverlayMenu(&overlayMenu{}).
		Content(
			newDocument(p.path).
				Description(filepath.Base(p.path)),
		)
}
