//go:generate go run page_gen.go
//go:generate go fmt

package dom

import (
	"bytes"
	"text/template"

	"github.com/murlokswarm/app"
)

// Page represents a html page.
type Page struct {
	// The title.
	Title string

	// The metadata.
	Metas []app.Meta

	// The css file paths to include.
	CSS []string

	// The javascript file paths to include.
	Javascripts []string

	// The name of the javascript function to pass data to Go.
	GoRequest string

	// The name of the root component.
	RootCompoName string
}

func (p Page) String() string {
	var b bytes.Buffer

	tmpl := template.Must(template.New(p.Title).Parse(htmlTmpl))
	tmpl.Execute(&b, struct {
		Title         string
		Metas         []app.Meta
		CSS           []string
		Javascripts   []string
		PageJS        string
		GoRequest     string
		RootCompoName string
	}{
		Title:         p.Title,
		Metas:         p.Metas,
		CSS:           p.CSS,
		Javascripts:   p.Javascripts,
		PageJS:        jsTmpl,
		GoRequest:     p.GoRequest,
		RootCompoName: p.RootCompoName,
	})

	return b.String()
}
