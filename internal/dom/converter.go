package dom

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/url"
	"strings"
	"time"
)

var converters = map[string]interface{}{
	"bind":  bind,
	"compo": urlToHTMLTag,
	"json":  jsonFormat,
	"raw":   rawHTML,
	"time":  timeFormat,
}

func bind(v ...interface{}) template.JS {
	targets := make([]string, len(v))

	for i, t := range v {
		targets[i] = fmt.Sprint(t)
	}

	return template.JS(strings.Join(targets, "."))
}

func urlToHTMLTag(s string) template.HTML {
	u, _ := url.Parse(s)

	b := make([]byte, 0, len(s)+len(s)/2)
	b = append(b, '<')

	tag := strings.TrimPrefix(u.Path, "/")
	b = append(b, tag...)
	b = append(b, ' ')

	for k, v := range u.Query() {
		k = strings.ToLower(k)
		b = append(b, k...)

		if len(v) != 0 {
			b = append(b, `="`...)
			b = append(b, v[0]...)
			b = append(b, `" `...)
		}
	}

	if b[len(b)-1] == ' ' {
		b = b[:len(b)-1]
	}

	b = append(b, '>')
	return template.HTML(b)
}

func jsonFormat(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	return string(b), err
}

func rawHTML(s string) template.HTML {
	return template.HTML(s)
}

func timeFormat(t time.Time, layout string) string {
	return t.Format(layout)
}
