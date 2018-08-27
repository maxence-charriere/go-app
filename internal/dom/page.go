package dom

//go:generate go run gen.go
//go:generate go fmt

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/murlokswarm/app"
)

// Page generate an HTML page from the given configuration.
func Page(c app.HTMLConfig, bridge, loadedCompo string) string {
	var w bytes.Buffer

	c = sanitizeHTMLConfigPaths(c)

	tmpl := template.Must(template.New(c.Title).Parse(htmlTemplate))
	tmpl.Execute(&w, struct {
		Title       string
		Metas       []app.Meta
		CSS         []string
		LoadedCompo template.JS
		JS          template.JS
		Javascripts []string
	}{
		Title:       c.Title,
		Metas:       c.Metas,
		CSS:         c.CSS,
		LoadedCompo: template.JS(loadedCompo),
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

func sanitizeHTMLConfigPaths(c app.HTMLConfig) app.HTMLConfig {
	for i, css := range c.CSS {
		c.CSS[i] = strings.Replace(css, `\`, "/", -1)
	}

	for i, js := range c.Javascripts {
		c.Javascripts[i] = strings.Replace(js, `\`, "/", -1)
	}

	return c
}
