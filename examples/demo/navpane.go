package main

import "github.com/murlokswarm/app"

// NavPane is a component that describes the navigation pane.
type NavPane struct {
	Current  string
	Examples []string
}

// OnMount is the func called when the component is mounted.
func (n *NavPane) OnMount() {
	n.Examples = []string{
		"hello",
		"driver",
		"window",
	}

	app.Render(n)
}

// Render returns the html that describes the nav pane content.
func (n *NavPane) Render() string {
	return `
<div class="NavPane">
	<h1>Demo</h1>
	<ul>
		{{$current := .Current}}

		{{range .Examples}}
		<li class="{{if eq . $current}}Selected{{end}}">
			<a href="{{.}}">{{.}}</a>
		</li>
		{{end}}
	</ul>
</div>
	`
}
