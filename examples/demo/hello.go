package main

// Hello is a hello world component.
type Hello struct {
	Name string
}

// Render returns a html string that describes the component.
func (h *Hello) Render() string {
	return `
<div class="Layout Hello">
	<navpane current="hello">
	<div class="Hello-Content">
		<h1>
			Hello
			{{if .Name}}
				{{.Name}}
			{{else}}
				world
			{{end}}!
		</h1>
		<input class="Hello-Input" value="{{.Name}}" placeholder="Say something..." onchange="{{to "Name"}}" autofocus>
	</div>
</div>
	`
}
