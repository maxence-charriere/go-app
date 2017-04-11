package app

import (
	"bytes"
	"text/template"

	"github.com/satori/go.uuid"
)

const (
	htmlContextTmpl = `
<!DOCTYPE html>
<html lang="{{if .Lang}}{{.Lang}}{{else}}en{{end}}">
<head>
    <meta charset="UTF-8">

    <style media="all" type="text/css">
        html {
            height: 100%;
            width: 100%;
            margin: 0;
        }
        
        body {
            height: 100%;
            width: 100%;
            margin: 0;
            font-family: "Helvetica Neue", "Segoe UI";
        }
    </style>

    {{range .CSS}}
    <link type="text/css" rel="stylesheet" href="css/{{.}}" />{{end}}

<title>{{.Title}}</title>
</head>

<body oncontextmenu="event.preventDefault()">
    <div data-murlok-root="{{.ID}}"></div>

    <script>{{.MurlokJS}}</script>

    {{range .JS}}
    <script src="js/{{.}}"></script>{{end}}
</body>
</html>
    `
)

// HTMLContext contains the data required to generate the minimum HTML to
// setup a webview based context.
// Should be used only in drivers implementations.
type HTMLContext struct {
	ID       uuid.UUID
	Title    string
	Lang     string
	MurlokJS string
	JS       []string
	CSS      []string
}

// HTML generate the HTML based on the data of c.
func (c HTMLContext) HTML() string {
	var b bytes.Buffer
	t := template.Must(template.New("").Parse(htmlContextTmpl))
	t.Execute(&b, c)
	return b.String()
}
