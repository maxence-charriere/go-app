<h1 align="center">
    <a href="https://github.com/maxence-charriere/go-app">
        <img alt="go-app"  width="150" height="150" src="https://storage.googleapis.com/murlok-github/icon-192.png">
    </a>
</h1>

<p align="center">
	<a href="https://circleci.com/gh/maxence-charriere/go-app"><img src="https://circleci.com/gh/maxence-charriere/go-app.svg?style=svg" alt="Circle CI Go build"></a>
    <a href="https://goreportcard.com/report/github.com/maxence-charriere/go-app"><img src="https://goreportcard.com/badge/github.com/maxence-charriere/go-app" alt="Go Report Card"></a>
	<a href="https://GitHub.com/maxence-charriere/go-app/releases/"><img src="https://img.shields.io/github/release/maxence-charriere/go-app.svg" alt="GitHub release"></a>
	<a href="https://pkg.go.dev/github.com/maxence-charriere/go-app/pkg/app"><img src="https://img.shields.io/badge/dev-reference-007d9c?logo=go&logoColor=white&style=flat" alt="pkg.go.dev docs"></a>
	<a href="https://github.com/maxence-charriere/go-app/wiki"><img src="https://img.shields.io/badge/github-wiki-6E7AF8?logo=github&style=flat" alt="pkg.go.dev docs"></a>
    <a href="https://twitter.com/jonhymaxoo"><img alt="Twitter URL" src="https://img.shields.io/badge/twitter-@jonhymaxoo-35A9F8?logo=twitter&style=flat"></a>
</p>

**app** is a package to build [progressive web apps (PWA)](https://developers.google.com/web/progressive-web-apps/) with [Go programming language](https://golang.org) and [WebAssembly](https://webassembly.org).

It uses a [declarative syntax](#declarative-syntax) that allows creating and dealing with HTML elements only by using Go, and without writing any HTML markup.

The package also provides an [http.handler](#http-handler) ready to serve all the required resources to run Go-based progressive web apps.

## Install

**app** requires [Go 1.13](https://golang.org/doc/go1.13) or newer.

```sh
go get -u -v github.com/maxence-charriere/go-app/pkg/app
```

## How it works

<p align="center">
     <img alt="app diagram"  width="680" src="https://storage.googleapis.com/murlok-github/app.png">
</p>

- **Users:** The users of your app. They request pages and resources from their web browser.
- **[app.Handler](https://pkg.go.dev/github.com/maxence-charriere/go-app/pkg/app#Handler)**: An [http.Handler](https://golang.org/pkg/net/http/#Handler) used by your server or cloud function. It serves your app, its static resources, and all the required files to make it work on user browsers.
- **Application**: Your app built with this package. It is built as a WebAssembly (.wasm) binary and is served by the [app.Handler](https://pkg.go.dev/github.com/maxence-charriere/go-app/pkg/app#Handler).
- **Other static resources**: Styles, images, and scripts used by your app. They are also served by the [app.Handler](https://pkg.go.dev/github.com/maxence-charriere/go-app/pkg/app#Handler).

## Declarative syntax

**go-app** uses a declarative syntax so you can write component-based UI elements just by using the Go programming language.

```go
package main

import "github.com/maxence-charriere/go-app/pkg/app"

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

func (h *hello) OnInputChange(src app.Value, e app.Event) {
    h.name = src.Get("value").String()
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

Note that we named the build output `app.wasm`. The reason is that the HTTP handler requires the web assembly app to be named this way in order to be served.

## HTTP handler

Once your app is built, the next step is to serve it.

This package provides an [http.Handler implementation](https://pkg.go.dev/github.com/maxence-charriere/go-app/pkg/app#Handler) ready to serve your PWA and all the required resources to make it work in a web browser.

The handler can be used to create either a web server or a cloud function (AWS Lambda, GCloud function or Azure function).

```go
package main

import (
    "net/http"

    "github.com/maxence-charriere/go-app/pkg/app"
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

Note that **you need to add `app.wasm` to the server location**. The reason is that [app.Handler](https://pkg.go.dev/github.com/maxence-charriere/go-app/pkg/app#Handler) is looking for a file named `app.wasm` in the server directory in order to serve the web assembly binary.

```sh
hello-local        # Server directory
├── app.wasm       # Wasm binary
├── hello-local    # Server binary
└── main.go        # Server code
```

## Works on mainstream browsers

|         | Chrome | Edge | Firefox | Opera | Safari |
| :------ | :----: | :--: | :-----: | :---: | :----: |
| Desktop |   ✔    | ✔\*  |    ✔    |   ✔   |   ✔    |
| Mobile  |   ✔    |  ✔   |    ✔    |   ✔   |   ✔    |

\*_only Chromium based [Edge](https://www.microsoft.com/edge)_

## Demo

### Hello app

The hello example introduced above:

| App sources                                                              | Server sources                                                                               | Description                                                                                     |
| ------------------------------------------------------------------------ | -------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------- |
| [hello](https://github.com/maxence-charriere/go-app/tree/master/demo/hello) | [hello-local](https://github.com/maxence-charriere/go-app/tree/master/demo/hello-local)         | Hello app that runs on a local server.                                                          |
|                                                                          | [hello-docker](https://github.com/maxence-charriere/go-app/tree/master/demo/hello-docker)       | Hello app that run in a Docker container.                                                       |
|                                                                          | [hello-appengine](https://github.com/maxence-charriere/go-app/tree/master/demo/hello-appengine) | Hello app that run on Google Cloud App Engine.<br> [See live](https://goapp-269110.appspot.com) |

### Live apps

<p align="center">
    <a href="https://luck.murlok.io">
        <img alt="luck app"  width="400" src="https://storage.googleapis.com/murlok-github/luck-thumb.png">
    </a>
    <a href="https://demo.murlok.io">
        <img alt="hello app"  width="400" src="https://storage.googleapis.com/murlok-github/hello-thumb.png">
    </a>
    <a href="https://demo.murlok.io/city">
        <img alt="city app"  width="400" src="https://storage.googleapis.com/murlok-github/city-thumb.png">
    </a>
</p>

## Sponsors

Support this project by becoming a sponsor. Your logo/picture will show up here with a link to your website.

<a href="https://opencollective.com/go-app" target="_blank"><img src="https://opencollective.com/go-app/contribute/button@2x.png?color=blue" width=250 /></a>
