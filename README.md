<h1 align="center">
    <a href="https://go-app.dev">
        <img alt="go-app"  width="150" height="150" src="https://raw.githubusercontent.com/maxence-charriere/go-app-docs/main/docs/web/icon.svg">
    </a>
</h1>

<p align="center">
	<a href="https://github.com/maxence-charriere/go-app/actions"><img src="https://github.com/maxence-charriere/go-app/actions/workflows/build.yml/badge.svg?branch=master" alt="GitHub actions"></a>
	<a href="https://GitHub.com/maxence-charriere/go-app/releases/"><img src="https://img.shields.io/github/release/maxence-charriere/go-app.svg" alt="GitHub release"></a>
	<a href="https://pkg.go.dev/github.com/maxence-charriere/go-app/v11/pkg/app"><img src="https://img.shields.io/badge/dev-reference-007d9c?logo=go&logoColor=white&style=flat" alt="pkg.go.dev docs"></a>
    <a href="https://twitter.com/jonhymaxoo"><img alt="Twitter URL" src="https://img.shields.io/badge/twitter-@jonhymaxoo-35A9F8?logo=twitter&style=flat"></a>
</p>

Go-app is a package for **building progressive web apps (PWA)** with the [Go programming language (Golang)](https://golang.org) and [WebAssembly (Wasm)](https://webassembly.org).

Shaping a UI is done by using a **[declarative syntax](https://go-app.dev/declarative-syntax) that creates and compose HTML elements only by using the Go programing language**.

It **uses [Go HTTP standard](https://golang.org/pkg/net/http) model**.

An app created with go-app can out of the box **runs in its own window**, **supports offline mode**, and is **SEO friendly**.

## Documentation (built with go-app)

[![go-app documentation](https://raw.githubusercontent.com/maxence-charriere/go-app-docs/main/docs/web/images/go-app.png)](https://go-app.dev)

## Install

**go-app** requirements:

- [Go 1.18](https://golang.org/doc/go1.17) or newer
- [Go module](https://github.com/golang/go/wiki/Modules)

```sh
go mod init
go get -u github.com/maxence-charriere/go-app/v11/pkg/app
```

## Declarative syntax

Go-app uses a [declarative syntax](https://go-app.dev/declarative-syntax) so you can **write reusable component-based UI elements** just by using the Go programming language.

Here is a Hello World component that takes an input and displays its value in its title:

```go
type hello struct {
	app.Compo

	name string
}

func (h *hello) Render() app.UI {
	return app.Div().Body(
		app.H1().Body(
			app.Text("Hello, "),
			app.If(h.name != "", func() app.UI {
				return app.Text(h.name)
			}).Else(func() app.UI {
				return app.Text("World!")
			}),
		),
		app.P().Body(
			app.Input().
				Type("text").
				Value(h.name).
				Placeholder("What is your name?").
				AutoFocus(true).
				OnChange(h.ValueTo(&h.name)),
		),
	)
}
```

## Standard HTTP

Apps created with go-app complies with [Go standard HTTP](https://golang.org/pkg/net/http) package interfaces.

```go
func main() {
	// Components routing:
	app.Route("/", func() app.Composer { return &hello{} })
	app.Route("/hello", func() app.Composer { return &hello{} })
	app.RunWhenOnBrowser()

	// HTTP routing:
	http.Handle("/", &app.Handler{
		Name:        "Hello",
		Description: "An Hello World! example",
	})

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
```

## Getting started

**Read the [Getting Started](https://go-app.dev/getting-started) document**.

## Built with go-app

- [Murlok.io](https://murlok.io/)
- [liwasc](https://pojntfx.github.io/liwasc/)
- [go-app Docs](https://go-app.dev/)

## Support the project

If go-app is useful to you, consider [sponsoring its development](https://github.com/sponsors/maxence-charriere).
