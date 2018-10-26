package dom

import (
	"encoding/json"
	"html/template"
	"time"
)

var converters = map[string]interface{}{
	"raw":   rawHTML,
	"compo": compoHTMLTag,
	"time":  timeFormat,
	"json":  jsonFormat,
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
