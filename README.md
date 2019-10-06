<p align="center">
    <a href="https://demo.murlok.io"><img alt="ui demo" src="https://github.com/maxence-charriere/app/wiki/assets/ui.png"></a>
</p>

# app

<p align="center">
	<a href="https://circleci.com/gh/maxence-charriere/app"><img src="https://circleci.com/gh/maxence-charriere/app.svg?style=svg" alt="Circle CI Go build"></a>
    <a href="https://goreportcard.com/report/github.com/maxence-charriere/app"><img src="https://goreportcard.com/badge/github.com/maxence-charriere/app" alt="Go Report Card"></a>
    <a href="https://godoc.org/github.com/maxence-charriere/app/pkg/app"><img src="https://godoc.org/github.com/maxence-charriere/app/pkg/app?status.svg" alt="GoDoc"></a>
    <a href="https://www.patreon.com/maxencecharriere"><img alt="Custom badge" src="https://img.shields.io/endpoint.svg?url=https%3A%2F%2Fshieldsio-patreon.herokuapp.com%2Fmaxencecharriere" alt="patreon"></a>
</p>

A [WebAssembly](https://webassembly.org) framework to build GUI with
[Go](https://golang.org), [HTML](https://en.wikipedia.org/wiki/HTML5) and
[CSS](https://en.wikipedia.org/wiki/Cascading_Style_Sheets).

It features:

- [PWA support](https://developers.google.com/web/progressive-web-apps/)
- [Build tool](https://github.com/maxence-charriere/app/tree/master/cmd/goapp/main.go) that removes the hassle of packaging wasm apps
- [React](https://reactjs.org) flavored API


## Install

Requires [Go 1.13](https://golang.org/doc/go1.13)

```sh
# Package sources + goapp CLI:
go get -u -v github.com/maxence-charriere/app/cmd/goapp

# Package sources only:
go get -u -v github.com/maxence-charriere/app/pkg/app

```

## Getting started

```sh
cd $GOPATH/src          # go to your gopath sources (optional)
mkdir demo && cd demo   # create and go to your go package
goapp init -v           # init project layout
goapp run -v -b chrome  # run the app and launch the main page on chrome
```

## How it works

### Project layout

```bash
demo
└── cmd
    ├── demo-server
    │   ├── main.go
    │   └── web
    │       ├── style sheets...
    │       ├── images...
    │       └── etc...
    └── demo-wasm
        └── main.go
```

- The `cmd` directory contains the project main applications.
- The `demo-wasm` directory contains the app that is compiled in **wasm** and that will run in the browser.
- The `demo-server` directory contains the server that serves the **wasm** app and its resources.
- The `web` directory contrains the app resources like style sheets (css), images and other static resources.

Project layout can be initialized by running this command in the repository root.

```bash
goapp init -v
```

### App

The app is the Go code compiled in web assembly and executed in the browser.

```go
// cmd/demo-wasm/main.go

package main

import (
    "log"

    "github.com/maxence-charriere/app/pkg/app"
    "github.com/maxence-charriere/app/pkg/log"
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
    <input value="{{.Name}}" placeholder="What is your name?" onchange="Name" autofocus>
</div>
    `
}

func main() {
    app.Import(&Hello{})

    app.DefaultPath = "/hello"
    app.Run()
}
```

### Server

The server serves the web assembly Go program and the other resources.

```go
// cmd/demo-server/main.go

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

The whole project can be built with the
[goapp](https://github.com/maxence-charriere/app/tree/master/cmd/goapp/main.go)
CLI tool.
**goapp** builds the server, the wasm app, imports the required javascript
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
demo
└── cmd
    ├── demo-server
    │   ├── demo-server (server)
    │   ├── main.go
    │   └── web
    │       ├── goapp.wasm (app)
    │       ├── wasm_exec.js
    │       ├── images...
    │       └── etc...
    └── demo-wasm
        └── main.go
```

See a [full example](https://github.com/maxence-charriere/app/tree/master/demo) and its online demo:

- [Hello](https://demo.murlok.io)
- [City](https://demo.murlok.io/city)

## Support
|Platform|Chrome|Edge|Firefox|Safari|
|:-|:-:|:-:|:-:|:-:|
|Desktop|✔|✔*|✔|✔|
|Mobile|✔|✔|✔|✔|

Issues:

- Non Chromiun based Edge does not support `TextEncoder` which is used by the javascript support file provided by Go.
