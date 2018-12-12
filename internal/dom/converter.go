package dom

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strings"
	"time"
)

var converters = map[string]interface{}{
	"raw":   rawHTML,
	"compo": compoHTMLTag,
	"time":  timeFormat,
	"json":  jsonFormat,
	"to":    target,
}

func rawHTML(s string) template.HTML {
	return template.HTML(s)
}

func compoHTMLTag(s string) template.HTML {
	return template.HTML("<" + s + ">")
}

func timeFormat(t time.Time, layout string) string {
	return t.Format(layout)
}

func jsonFormat(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	return string(b), err
}

func target(v ...interface{}) template.JS {
	targets := make([]string, len(v))

	for i, t := range v {
		targets[i] = fmt.Sprint(t)
	}

	return template.JS(strings.Join(targets, "."))
}
