<h1 align="center">
    <a href="https://go-app.dev">
        <img alt="go-app"  width="150" height="150" src="https://storage.googleapis.com/murlok-github/icon-192.png">
    </a>
</h1>

<p align="center">
	<a href="https://circleci.com/gh/maxence-charriere/go-app"><img src="https://circleci.com/gh/maxence-charriere/go-app.svg?style=svg" alt="Circle CI Go build"></a>
    <a href="https://goreportcard.com/report/github.com/maxence-charriere/go-app"><img src="https://goreportcard.com/badge/github.com/maxence-charriere/go-app" alt="Go Report Card"></a>
	<a href="https://GitHub.com/maxence-charriere/go-app/releases/"><img src="https://img.shields.io/github/release/maxence-charriere/go-app.svg" alt="GitHub release"></a>
	<a href="https://pkg.go.dev/github.com/maxence-charriere/go-app/v7/pkg/app"><img src="https://img.shields.io/badge/dev-reference-007d9c?logo=go&logoColor=white&style=flat" alt="pkg.go.dev docs"></a>
    <a href="https://twitter.com/jonhymaxoo"><img alt="Twitter URL" src="https://img.shields.io/badge/twitter-@jonhymaxoo-35A9F8?logo=twitter&style=flat"></a>
    <a href="https://opencollective.com/go-app" alt="Financial Contributors on Open Collective"><img src="https://opencollective.com/go-app/all/badge.svg?label=open+collective&color=4FB9F6" /></a>
</p>

**go-app** is a package to build [progressive web apps (PWA)](https://developers.google.com/web/progressive-web-apps/) with [Go programming language](https://golang.org) and [WebAssembly](https://webassembly.org).

It uses a [declarative syntax](#declarative-syntax) that allows creating and dealing with HTML elements only by using Go, and without writing any HTML markup.

The package also provides an [http.handler](#http-handler) ready to serve all the required resources to run Go-based progressive web apps.

## Documentation

[![go-app documentation](docs/web/images/go-app.png)](https://go-app.dev)

- [Getting started](https://go-app.dev/start)
- [Architecture](https://go-app.dev/architecture)
- [Reference](https://go-app.dev/reference) | [Go.dev](https://pkg.go.dev/github.com/maxence-charriere/go-app/v7/pkg/app)

- [Components](https://go-app.dev/components)
- [Declarative syntax](https://go-app.dev/syntax)
- [Routing](https://go-app.dev/routing)
- [Concurrency](https://go-app.dev/concurrency)
- [Javascript and DOM](https://go-app.dev/js)
- [Static resources](https://go-app.dev/static-resources)

- [Built with go-app](https://go-app.dev/built-with)
- [Intall](https://go-app.dev/install)

## Sneak peek

### Install

**go-app** requirements:

- [Go 1.14](https://golang.org/doc/go1.14) or newer
- [Go module](https://github.com/golang/go/wiki/Modules)

```sh
go mod init
go get -u github.com/maxence-charriere/go-app/v7/pkg/app
```

### Declarative syntax

**go-app** uses a declarative syntax so you can write component-based UI elements just by using the Go programming language.

```go
package main

import "github.com/maxence-charriere/go-app/v7/pkg/app"

type hello struct {
    app.Compo
    name string
}

func (h *hello) Render() app.UI {
    return app.Div().Body(
        app.Main().Body(
            app.H1().Body(
                app.Text("Hello, "),
                app.If(h.name != "",
                    app.Text(h.name),
                ).Else(
                    app.Text("World"),
                ),
            ),
            app.Input().
                Value(h.name).
                Placeholder("What is your name?").
                AutoFocus(true).
                OnChange(h.OnInputChange),
        ),
    )
}

func (h *hello) OnInputChange(ctx app.Context, e app.Event) {
    h.name = ctx.JSSrc.Get("value").String()
    h.Update()
}

func main() {
    app.Route("/", &hello{})
    app.Route("/hello", &hello{})
    app.Run()
}

```

The app is built with the Go build tool by specifying WebAssembly as architecture and javascript as operating system:

```sh
GOARCH=wasm GOOS=js go build -o app.wasm
```

Note that the build output is named `app.wasm` because the HTTP handler expects the wasm app to be named that way in order to serve its content.

### HTTP handler

Once the wasm app is built, the next step is to serve it.

This package provides an [http.Handler implementation](https://pkg.go.dev/github.com/maxence-charriere/go-app/pkg/app#Handler) ready to serve a PWA and all the required resources to make it work in a web browser.

The handler can be used to create either a web server or a cloud function (AWS Lambda, Google Cloud function or Azure function).

```go
package main

import (
    "net/http"

    "github.com/maxence-charriere/go-app/v7/pkg/app"
)

func main() {
    h := &app.Handler{
        Title:  "Hello Demo",
        Author: "Maxence Charriere",
    }

    if err := http.ListenAndServe(":7777", h); err != nil {
        panic(err)
    }
}
```

The server is built as a standard Go program:

```sh
go build
```

Once the server and the wasm app built, `app.wasm` must be moved in the `web` directory, located by the side of the server binary. The web directory is where to put static resources, such as the wasm app, styles, scripts, or images.

The directory should look like as the following:

```sh
.                   # Directory root
├── hello           # Server binary
├── main.go         # Server source code
└── web             # Directory for static resources
    └── app.wasm    # Wasm binary
```

Then launch the `hello` server and open the web browser to see the result.

Demo code and live demo are available on the following links:

- [Demo repository](https://github.com/maxence-charriere/go-app-demo/tree/v6)
- [Live Hello world](https://go-app-demo-42.appspot.com)

## Contributors

### Code Contributors

This project exists thanks to all the people who contribute. [[Contribute](CONTRIBUTING.md)].

<a href="https://github.com/maxence-charriere/go-app/graphs/contributors"><img src="https://opencollective.com/go-app/contributors.svg?width=890&button=false" /></a>

### Financial Contributors

Become a financial contributor and help us sustain [go-app](https://github.com/maxence-charriere/go-app) development. [[Contribute](https://opencollective.com/go-app/contribute)]

#### Individuals

<a href="https://opencollective.com/go-app"><img src="https://opencollective.com/go-app/individuals.svg?width=890"></a>

#### Organizations

Support this project with your organization. Your logo will show up here with a link to your website. [[Contribute](https://opencollective.com/go-app/contribute)]

<a href="https://opencollective.com/go-app/organization/0/website"><img src="https://opencollective.com/go-app/organization/0/avatar.svg"></a>
<a href="https://opencollective.com/go-app/organization/1/website"><img src="https://opencollective.com/go-app/organization/1/avatar.svg"></a>
<a href="https://opencollective.com/go-app/organization/2/website"><img src="https://opencollective.com/go-app/organization/2/avatar.svg"></a>
<a href="https://opencollective.com/go-app/organization/3/website"><img src="https://opencollective.com/go-app/organization/3/avatar.svg"></a>
<a href="https://opencollective.com/go-app/organization/4/website"><img src="https://opencollective.com/go-app/organization/4/avatar.svg"></a>
<a href="https://opencollective.com/go-app/organization/5/website"><img src="https://opencollective.com/go-app/organization/5/avatar.svg"></a>
<a href="https://opencollective.com/go-app/organization/6/website"><img src="https://opencollective.com/go-app/organization/6/avatar.svg"></a>
<a href="https://opencollective.com/go-app/organization/7/website"><img src="https://opencollective.com/go-app/organization/7/avatar.svg"></a>
<a href="https://opencollective.com/go-app/organization/8/website"><img src="https://opencollective.com/go-app/organization/8/avatar.svg"></a>
<a href="https://opencollective.com/go-app/organization/9/website"><img src="https://opencollective.com/go-app/organization/9/avatar.svg"></a>
