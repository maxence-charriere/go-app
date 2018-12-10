package main

// Hello is a hello world component.
type Hello struct {
	Name string
}

// Render returns a string that describes the component markup.
func (h *Hello) Render() string {
	return `
<div class="Layout">
	<navpane>
	<div class="Hello-Content">
		<h1>
			Hello
			{{if .Name}}
				{{.Name}}
			{{else}}
				world
			{{end}}!
		</h1>
		<input class="Hello-Input" value="{{.Name}}" placeholder="Say something..." onchange="Name" autofocus>
	</div>
</div>
	`
}
