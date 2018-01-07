package main

import (
	"fmt"
	"math/rand"
	"net/url"
	"strconv"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/html"
)

func init() {
	app.Import(&WebviewComponent{})
}

// WebviewComponent is a component to test html in webview based elements.
// It implements the app.Component interface.
type WebviewComponent struct {
	Title       string
	Page        int
	SquareColor string
	Number      int
	CanPrevious bool
	CanNext     bool
}

// Render statisfies the app.Component interface.
func (c *WebviewComponent) Render() string {
	return `
<div class="root">
	<h1>Test Window</h1>
	<p>
		Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod 
		tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam,
		quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo
		consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse
		cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat
		non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
	</p>
	
	<ul>
		<li><a href="webviewcomponent?page=42">To page 42</a></li>
		<li><a href="unknown?page=42">Unknown compopent</a></li>
		<li><a href="http://theverge.com">external hyperlink</a></li>
		<li><button onclick="OnNextPage">Next Page</button></li>
		<li><button onclick="OnLink">External link</button></li>
		<li><button onclick="NotMapped">Not mapped</button></li>
		<li>
			<button onclick="OnChangeSquareColor">Render Attributes: change square color</button>
			<div class="square {{.SquareColor}}"></div>
		</li>
		<li>
			<button onclick="OnChangeNumber">Render: change number</button>
			<div>{{.Number}}</div>
		</li>
		<li>
		<button {{if not .CanPrevious}}disabled{{end}} onclick="OnPrevious">Previous</button>
		<button onclick="OnReload">Reload</button>
		<button {{if not .CanNext}}disabled{{end}} onclick="OnNext">Next</button>
		</li>
	</ul>
	
	<p>Page: {{.Page}}</p>
</div>
	`
}

// OnNavigate is the function that is called when a component is navigated.
func (c *WebviewComponent) OnNavigate(u *url.URL) {
	if pagevals := u.Query()["page"]; len(pagevals) != 0 {
		c.Page, _ = strconv.Atoi(pagevals[0])
	}

	if c.Page == 0 {
		c.Page = 1
	}

	if win, err := app.WindowFromComponent(c); err == nil {
		c.CanPrevious = win.CanPrevious()
		c.CanNext = win.CanNext()
	}

	app.Render(c)
}

// PageConfig return allow to set page information like title or meta when the
// component is mounted as the root component.
func (c *WebviewComponent) PageConfig() html.PageConfig {
	return html.PageConfig{
		Title: fmt.Sprintf("Test component %v", c.Page),
	}
}

// OnNextPage is the function to be called when the Next page button is clicked.
func (c *WebviewComponent) OnNextPage() {
	page := c.Page
	page++

	if win, err := app.WindowFromComponent(c); err == nil {
		win.Load("/webviewcomponent?page=%v", page)
	}
}

// OnLink is the function to be called when the External link button is clicked.
func (c *WebviewComponent) OnLink() {
	if win, err := app.WindowFromComponent(c); err == nil {
		win.Load("http://www.github.com")
	}
}

// OnChangeSquareColor is the function to be called when the change color button
// is clicked.
func (c *WebviewComponent) OnChangeSquareColor() {
	switch c.SquareColor {
	case "blue":
		c.SquareColor = "pink"
	case "pink":
		c.SquareColor = ""
	default:
		c.SquareColor = "blue"
	}
	app.Render(c)
}

// OnChangeNumber is the function to be called when the change number button is
// clicked.
func (c *WebviewComponent) OnChangeNumber() {
	c.Number = rand.Int()
	app.Render(c)
}

// OnPrevious is the function that is called when the previous button is
// clicked.
func (c *WebviewComponent) OnPrevious() {
	if win, err := app.WindowFromComponent(c); err == nil {
		win.Previous()
	}
}

// OnReload is the function that is called when the reload button is clicked.
func (c *WebviewComponent) OnReload() {
	if win, err := app.WindowFromComponent(c); err == nil {
		win.Reload()
	}
}

// OnNext is the function that is called when the next button is clicked.
func (c *WebviewComponent) OnNext() {
	if win, err := app.WindowFromComponent(c); err == nil {
		win.Next()
	}
}
