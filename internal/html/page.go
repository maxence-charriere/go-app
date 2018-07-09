package html

//go:generate go run gen.go
//go:generate go fmt

import (
	"bytes"
	"text/template"

	"github.com/murlokswarm/app"
)

// Page generate an HTML page from the given configuration.
func Page(c app.HTMLConfig) string {
	var w bytes.Buffer

	tmpl := template.Must(template.New(c.Title).Parse(htmlTemplate))
	tmpl.Execute(&w, struct {
		Title       string
		Metas       []app.Meta
		CSS         []string
		JS          string
		Javascripts []string
	}{
		Title:       c.Title,
		Metas:       c.Metas,
		CSS:         c.CSS,
		JS:          jsTemplate,
		Javascripts: c.Javascripts,
	})

	return w.String()
}
