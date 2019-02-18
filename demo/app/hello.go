package main

import (
	"github.com/murlokswarm/app"
)

// Hello is a hello world component.
type Hello struct {
	Name string
}

// Render returns a string that describes the component markup.
func (h *Hello) Render() string {
	return `
<div class="Hello">
	<button class="Menu"  onclick="OnMenuClick" oncontextmenu="OnMenuClick">☰</button>
	<app.contextmenu>

	<h1>
		Hello
		{{if .Name}}
			{{.Name}}
		{{else}}
			world
		{{end}}!
	</h1>
	<input value="{{.Name}}" placeholder="What is your name?" onchange="Name" autofocus>
</div>
	`
}

// OnMenuClick creates a context menu when the menu button is clicked.
func (h *Hello) OnMenuClick() {
	app.NewContextMenu(
		app.MenuItem{
			Label:   "Reload",
			Keys:    "cmdorctrl+r",
			OnClick: app.Reload},
		app.MenuItem{Separator: true},
		app.MenuItem{
			Label: "Go to repository",
			OnClick: func() {
				app.Navigate("https://github.com/maxence-charriere/app")
			}},
		app.MenuItem{
			Label: "Source code",
			OnClick: func() {
				app.Navigate("https://github.com/maxence-charriere/app/tree/master/demo/app")
			}},
	)
}
