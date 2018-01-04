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
		<li><button onclick="OnNext">Next</button></li>
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
}

// PageConfig return allow to set page information like title or meta when the
// component is mounted as the root component.
func (c *WebviewComponent) PageConfig() html.PageConfig {
	return html.PageConfig{
		Title: fmt.Sprintf("Test component %v", c.Page),
	}
}

// OnNext is the function to be called when the Next button is clicked.
func (c *WebviewComponent) OnNext() {
	win, err := app.Context(c)
	if err != nil {
		app.DefaultLogger.Error(err)
		return
	}

	page := c.Page
	page++

	win.Load("/webviewcomponent?page=%v", page)
}

// OnLink is the function to be called when the External link button is clicked.
func (c *WebviewComponent) OnLink() {
	app.DefaultLogger.Log("Onlink Clicked")

	win, err := app.Context(c)
	if err != nil {
		app.DefaultLogger.Error(err)
	}
	win.Load("http://www.github.com")
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
