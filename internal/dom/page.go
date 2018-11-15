//go:generate go run page_gen.go
//go:generate go fmt

package dom

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/murlokswarm/app"
)

// Page generate an HTML page from the given configuration.
func Page(c app.HTMLConfig, bridge, loadedCompo string) string {
	var w bytes.Buffer

	c.CSS = cleanWindowsPath(c.CSS)
	c.Javascripts = cleanWindowsPath(c.Javascripts)

	tmpl := template.Must(template.New(c.Title).Parse(htmlTmpl))
	tmpl.Execute(&w, struct {
		Title       string
		Metas       []app.Meta
		CSS         []string
		LoadedCompo string
		JS          string
		Javascripts []string
	}{
		Title:       c.Title,
		Metas:       c.Metas,
		CSS:         c.CSS,
		LoadedCompo: loadedCompo,
		JS:          js(bridge),
		Javascripts: c.Javascripts,
	})

	return w.String()
}

func js(bridge string) string {
	return fmt.Sprintf(`
	var golangRequest = function (payload) {
		%s(payload);
	}
	%s`, bridge, jsTmpl)
}

func cleanWindowsPath(paths []string) []string {
	c := make([]string, len(paths))

	for i, p := range paths {
		c[i] = strings.Replace(p, `\`, "/", -1)
	}

	return c
}
