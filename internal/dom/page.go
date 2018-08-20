package dom

//go:generate go run gen.go
//go:generate go fmt

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/murlokswarm/app"
)

// Page generate an HTML page from the given configuration.
func Page(c app.HTMLConfig, bridge string) string {
	var w bytes.Buffer

	tmpl := template.Must(template.New(c.Title).Parse(htmlTemplate))
	tmpl.Execute(&w, struct {
		Title       string
		Metas       []app.Meta
		CSS         []string
		JS          template.JS
		Javascripts []string
	}{
		Title:       c.Title,
		Metas:       c.Metas,
		CSS:         c.CSS,
		JS:          js(bridge),
		Javascripts: c.Javascripts,
	})

	return w.String()
}

func js(bridge string) template.JS {
	return template.JS(fmt.Sprintf(`
	var golangRequest = function (payload) {
		%s(payload);
	}
	%s`, bridge, jsTemplate))
}
