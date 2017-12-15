package main

import "github.com/murlokswarm/app"

func init() {
	app.Import(&WebviewComponent{})
}

// WebviewComponent is a component to test html in webview based elements.
// It implements the app.Component interface.
type WebviewComponent struct {
	Title string
}

// Render statisfies the app.Component interface.
func (c *WebviewComponent) Render() string {
	return `
<div>
	<h1>Test Window</h1>
	<p>
		Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod 
		tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam,
		quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo
		consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse
		cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat
		non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
	</p>
</div>
	`
}
