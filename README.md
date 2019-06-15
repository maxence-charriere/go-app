<p align="center">
    <a href="https://app-demo-232021.appspot.com"><img alt="ui demo" src="https://github.com/maxence-charriere/app/wiki/assets/ui.png"></a>
</p>

# app

<p align="center">
	<a href="https://circleci.com/gh/maxence-charriere/app"><img src="https://circleci.com/gh/maxence-charriere/app.svg?style=svg" alt="Circle CI Go build"></a>
    <a href="https://goreportcard.com/report/github.com/maxence-charriere/app"><img src="https://goreportcard.com/badge/github.com/maxence-charriere/app" alt="Go Report Card"></a>
    <a href="https://godoc.org/github.com/maxence-charriere/app"><img src="https://godoc.org/github.com/maxence-charriere/app?status.svg" alt="GoDoc"></a>
    <a href="https://www.patreon.com/maxencecharriere"><img alt="Custom badge" src="https://img.shields.io/endpoint.svg?url=https%3A%2F%2Fshieldsio-patreon.herokuapp.com%2Fmaxencecharriere" alt="patreon"></a>
</p>

A [WebAssembly](https://webassembly.org) framework to build GUI with
[Go](https://golang.org), [HTML](https://en.wikipedia.org/wiki/HTML5) and
[CSS](https://en.wikipedia.org/wiki/Cascading_Style_Sheets).

## Install

```sh
# Package sources:
go get -u -v github.com/maxence-charriere/app

# Goapp cli tool:
go get -u -v github.com/maxence-charriere/app/cmd/goapp
```

## How it works

### Project layout

```bash
root
├── cmd
│   ├── demo-server
│   │   └── main.go
│   └── demo-wasm
│       └── main.go
└── web
    ├── wasm_exec.js
    ├── style sheets...
    ├── images...
    └── etc...
```

This layout follows the project layout defined in [golang-standards/project-layout](https://github.com/golang-standards/project-layout):

- The `cmd` directory contains the project main applications.
- The `demo` directory contains the app that is compiled in **wasm** and that will run in the browser.
- The `demo-server` directory contains the server that serves the **wasm** app and its resources.
- The `web` directory contrains the app resources like style sheets (css), images and other static resources.

Project layout can be initialized by running this command in the repository root.

```bash
goapp init -v
```

### App

The app is the Go code compiled in web assembly and executed in the browser.

```go
// root/cmd/demo-wasm/main.go

package main

import (
    "log"

    "github.com/maxence-charriere/app"
)

type Hello struct {
    Name string
}

func (h *Hello) Render() string {
    return `
<div class="Hello">
    <h1>
        Hello
        {{if .Name}}
            {{.Name}}
        {{else}}
            world
        {{end}}!
    </h1>
    <input value="{{.Name}}" placeholder="What is your name?" onchange="{{bind "Name"}}" autofocus>
</div>
    `
}

func main() {
    app.Import(&Hello{})

    app.DefaultPath = "/hello"

    if err := app.Run(); err != nil {
        log.Print(err)
    }
}
```

### Server

The server serves the web assembly Go program and the other resources.

```go
// root/cmd/demo-server/main.go

package main

import (
    "log"
    "net/http"
    "github.com/maxence-charriere/app"
)

func main() {
    http.Handle("/", &app.Handler{})

    if err := http.ListenAndServe(":3000", nil); err != nil {
        log.Fatal(err)
    }
}

```

### Build

The whole project is built with the
[goapp](https://github.com/maxence-charriere/app/tree/master/cmd/goapp/main.go)
CLI tool.
Goapp builds the server, the wasm app, imports the required javascript
support file and puts the pieces together to provide a ready to use project.

```bash
# Get the goapp CLI tool:
go get -u github.com/maxence-charriere/app/cmd/goapp

# Builds a server ready to serve the wasm app and its resources:
goapp build -v

# Launches the server and app in the default browser:
goapp run -v -b default
```

Once built, the directory tree should look like:

```bash
root
├── cmd
│   ├── demo-server
│   │   └── main.go
│   └── demo-wasm
│       └── main.go
├── demo-server (server)
└── web
    ├── goapp.wasm (app)
    ├── wasm_exec.js
    ├── style sheets...
    ├── images...
    └── etc...
```

See the [full example code](https://github.com/maxence-charriere/app/tree/master/demo) and the [online demo](https://app-demo-232021.appspot.com).

## Support

Requires [Go 1.12](https://golang.org/doc/go1.12).

|Platform|Chrome|Edge|Firefox|Safari|
|:-|:-:|:-:|:-:|:-:|
|Desktop|✔|✖|✔|✔|
|Mobile|✖|✖|✖|✖|

Issues:

- Go wasm currently triggers out of memory errors on mobile platforms (https://github.com/golang/go/issues/27462).
- Edge does not support `TextEncoder` which is used by the javascript support file.
