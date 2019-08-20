package maestro

// Compo is the interface that describes a component.
type Compo interface {
	// Render must return a HTML5 string. It supports standard Go html/template
	// API. The pipeline is based on the component struct.
	// See https://golang.org/pkg/text/template and
	// https://golang.org/pkg/html/template for template usage.
	Render() string
}
