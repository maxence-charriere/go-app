package main

import (
	"fmt"

	"github.com/murlokswarm/app"
)

// NavPane is a component that describes the navigation pane.
type NavPane struct {
	Current  string
	Examples []string
}

// OnMount is the func called when the component is mounted.
func (n *NavPane) OnMount() {
	n.Examples = []string{
		"hello",
	}

	app.Render(n)
	fmt.Println("bwa")
}

// Render returns the html that describes the nav pane content.
func (n *NavPane) Render() string {
	return `
<div class="NavPane">
	<h1 class="NavPane-Title">Demo</h1>
	<ul>
		{{range .Examples}}
		<li class="NavPane-Elem {{if eq . .Current}}NavPane-Selected{{end}}"><{{.}}/li>
		{{end}}
	</ul>
</div>
	`
}
