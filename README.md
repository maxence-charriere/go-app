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

Go-app is a package for **building progressive web apps (PWAs)** with [Go](https://go.dev) and [WebAssembly (Wasm)](https://webassembly.org).

Write your UI in Go using **reusable components** and a [declarative syntax](https://go-app.dev/declarative-syntax). Handle user interactions in Go and serve your app with the standard [`net/http`](https://pkg.go.dev/net/http) package.

Apps built with go-app are **SEO-friendly** and support **offline use**. On supported browsers and devices, users can **install your app and run it in its own window**.

## Documentation (built with go-app)

[![go-app documentation](https://raw.githubusercontent.com/maxence-charriere/go-app-docs/main/docs/web/images/go-app.png)](https://go-app.dev)

## Install

Requires [Go 1.26 or newer](https://go.dev/dl/).

From your project directory, initialize a Go module if you do not already have one. Replace `example.com/myapp` with your module path:

```sh
go mod init example.com/myapp
```

Add go-app to your project:

```sh
go get github.com/maxence-charriere/go-app/v11/pkg/app
```

## Components

Describe your UI with Go functions using go-app's [declarative syntax](https://go-app.dev/declarative-syntax). Each component embeds `app.Compo` and defines a `Render` method that returns its UI.

This component displays a greeting that updates as you type your name:

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
		app.Input().
			Type("text").
			Value(h.name).
			Placeholder("What is your name?").
			AutoFocus(true).
			OnInput(h.ValueTo(&h.name)),
	)
}
```

## Serving the app

Serve your app with Go's standard [`net/http`](https://pkg.go.dev/net/http) package. The `app.Handler` type implements `http.Handler` and can be used with an existing Go server.

This example registers the component above at `/` and serves the app on port `8000`:

```go
func main() {
	// Register the component.
	app.Route("/", func() app.Composer { return &hello{} })

	// Start the app in the browser; on the server, continue below.
	app.RunWhenOnBrowser()

	// Serve the app over HTTP.
	http.Handle("/", &app.Handler{
		Name:        "Hello",
		Description: "A Hello World example",
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
