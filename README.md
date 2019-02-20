<p align="center">
    <a href="https://app-demo-232021.appspot.com"><img alt="ui demo" src="https://github.com/maxence-charriere/app/wiki/assets/ui.png"></a>
</p>

# app

<p align="center">
	<a href="https://circleci.com/gh/maxence-charriere/app"><img src="https://circleci.com/gh/maxence-charriere/app.svg?style=svg" alt="Circle CI Go build"></a>
    <a href="https://goreportcard.com/report/github.com/maxence-charriere/app"><img src="https://goreportcard.com/badge/github.com/maxence-charriere/app" alt="Go Report Card"></a>
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

### App

The app is the Go code compiled in web assembly and executed in the browser.

```go
// root/cmd/demo/main.go

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

### Server

The server serves the web assembly Go program and the other resources.

```go
// root/cmd/demo-server/main.go

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

# Open http://localhost:3000 to see the result.
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

See the [full example code](https://github.com/maxence-charriere/app/tree/master/demo) and the [online demo](https://app-demo-232021.appspot.com).

## Support

Requires [Go 1.11](https://golang.org/doc/go1.11).

|Platform|Chrome|Edge|Firefox|Safari|
|:-|:-:|:-:|:-:|:-:|
|Desktop|✔|✖|✔|✔|
|Mobile|✖|✖|✖|✖|

Issues:

- Go wasm currently triggers out of memory errors. This will be fix with Go 1.12.
- Edge support is worked on.