<p align="center">
    <img alt="ui demo" src="https://github.com/maxence-charriere/app/wiki/assets/ui.png">
</p>

# app

<p align="center">
	<a href="https://travis-ci.org/maxence-charriere/app"><img src="https://travis-ci.org/maxence-charriere/app.svg?branch=master" alt="Build Status"></a>
    <a href="https://goreportcard.com/report/github.com/maxence-charriere/app"><img src="https://goreportcard.com/badge/github.com/maxence-charriere/app" alt="Go Report Card"></a>
    <a href="https://github.com/avelino/awesome-go#readme"><img src="https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg" alt="awesome-go"></a>
    <a href="https://godoc.org/github.com/maxence-charriere/app"><img src="https://godoc.org/github.com/maxence-charriere/app?status.svg" alt="GoDoc"></a>
</p>

A [WebAssembly](https://webassembly.org) framework to build GUI with
[Go](https://golang.org), [HTML](https://en.wikipedia.org/wiki/HTML5) and
[CSS](https://en.wikipedia.org/wiki/Cascading_Style_Sheets).

## Install

```sh
go get -u github.com/maxence-charriere/app
```

## How it works

### Project layout

```bash
root
├── cmd
│   ├── demo
│   │   └── main.go
│   └── demo-server
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

### App - *root/cmd/demo/main.go*

The app is the Go code compiled in web assembly and executed in the browser.

```go
package main

import (
    "log"

    "github.com/maxence-charriere/app"
)

// Hello is a component that describes a hello world. It implements the
// app.Compo interface.
type Hello struct {
    Name string
}

// Render returns UI to display.
//
// The onchange="{{bind "Name"}}" binds the onchange value to the Hello.Name
// field.
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

// The app entry point.
func main() {
    // Imports the hello component declared above in order to make it loadable
    // in a page or usable in other components.
    //
    // Imported component can be use as URL or html tags by referencing them by
    // their lowercased names.
    // E.g:
    //  Hello   => hello
    //  foo.Bar => foo.bar
    app.Import(&Hello{})

    // Defines the component to load when an URL without path is loaded.
    app.DefaultPath = "/hello"

    // Runs the app in the browser.
    if err := app.Run(); err != nil {
        log.Print(err)
    }
}
```

### Server - *root/cmd/demo-server/main.go*

The server serves the web assembly Go program and the other resources.

```go

package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/maxence-charriere/app"
)

func main() {
    // Setup the http handler to serve the web assembly app.
    http.Handle("/", &app.Handler{
        // The path of the directory that contains the wasm app file and the
        // other resources like .css files.
        WebDir:  "web",

        // The name of the wasm file that contains the app.
        Wasm:    "demo.wasm",
    })

    // Launches the server.
    if err := http.ListenAndServe(":3000", nil); err != nil {
        log.Fatal(err)
    }
}
```

### Build

Assuming the working directory is the root directory:

```bash
# Build the app:
GOOS=js GOARCH=wasm go build -o web/demo.wasm ./cmd/demo

# Copy the javascript support file:
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ./web

# Build the server:
go build ./cmd/demo-server

# Launch the server:
./demo-server
```

Once built, the directory tree should look like:

```bash
root
├── cmd
│   ├── demo
│   │   └── main.go
│   └── demo-server
│       └── main.go
├── demo-server (server)
└── web
    ├── demo.wasm (app)
    ├── wasm_exec.js
    ├── style sheets...
    ├── images...
    └── etc...
```

## Support

Requires [Go 1.11](https://golang.org/doc/go1.11).

|Platform|Chrome|Edge|Firefox|Safari|
|:-|:-:|:-:|:-:|:-:|
|Desktop|✔|✖|✔|✔|
|Mobile|✖|✖|✖|✖|

Issues:
- Go wasm currently trigger out of memory errors. This will be fix with Go 1.12.
- Edge support is worked on.