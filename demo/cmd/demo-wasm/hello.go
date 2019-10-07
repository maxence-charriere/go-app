// +build wasm

package main

import (
	"syscall/js"

	"github.com/maxence-charriere/app/pkg/app"
)

// Hello is a component that describes a hello world. It implements the
// app.Compo interface.
type Hello struct {
	Name string
}

// Render returns what to display.
//
// The onchange="Name" binds the onchange value to the Hello.Name
// field.
func (h *Hello) Render() string {
	return `
<div class="Hello">
	<div class="Menu" onclick="OnMenuClick" oncontextmenu="OnMenuClick">☰</div>

	<main class="content">
		<h1>
			Hello
			{{if .Name}}
				{{.Name}}
			{{else}}
				world
			{{end}}!
		</h1>
		<input value="{{.Name}}" placeholder="What is your name?" onchange="Name" autofocus>
	</main>
</div>
	`
}

// OnMenuClick creates a context menu when the menu button is clicked.
func (h *Hello) OnMenuClick(s, e js.Value) {
	app.NewContextMenu(
		app.MenuItem{
			Label: "Reload",
			Keys:  "cmdorctrl+r",
			OnClick: func(s, e js.Value) {
				app.Reload()
			},
		},
		app.MenuItem{Separator: true},
		app.MenuItem{
			Label: "Go to repository",
			OnClick: func(s, e js.Value) {
				app.Navigate("https://github.com/maxence-charriere/app")
			}},
		app.MenuItem{
			Label: "Source code",
			OnClick: func(s, e js.Value) {
				app.Navigate("https://github.com/maxence-charriere/app/blob/master/demo/cmd/demo-wasm/hello.go")
			}},
		app.MenuItem{Separator: true},
		app.MenuItem{
			Label: "City example",
			OnClick: func(s, e js.Value) {
				app.Navigate("city")
			}},
	)
}
