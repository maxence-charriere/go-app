<div class="header-separator">
    <header id="go-app" class="title fit center">
        <img alt="go-app logo." src="https://storage.googleapis.com/murlok-github/icon-192.png">
        <h1>go-app</h1>
    </header>
    <aside class="subtitle">
	    <a href="https://circleci.com/gh/maxence-charriere/go-app"><img src="https://circleci.com/gh/maxence-charriere/go-app.svg?style=svg" alt="Circle CI Go build"></a>
        <a href="https://goreportcard.com/report/github.com/maxence-charriere/go-app"><img src="https://goreportcard.com/badge/github.com/maxence-charriere/go-app" alt="Go Report Card"></a>
	    <a href="https://GitHub.com/maxence-charriere/go-app/releases/"><img src="https://img.shields.io/github/release/maxence-charriere/go-app.svg" alt="GitHub release"></a>
        <a href="https://twitter.com/jonhymaxoo"><img alt="Twitter URL" src="https://img.shields.io/badge/twitter-@jonhymaxoo-35A9F8?logo=twitter&style=flat"></a>
        <a href="https://opencollective.com/go-app" alt="Financial Contributors on Open Collective"><img src="https://opencollective.com/go-app/all/badge.svg?label=open+collective&color=4FB9F6" /></a>
    </aside>
</div>

Go-app is a package for **building progressive web apps (PWA)** with the [Go programming language (Golang)](https://golang.org) and [WebAssembly (Wasm)](https://webassembly.org).

Shaping a UI is done by using a **[declarative syntax](/syntax) that creates and compose HTML elements only by using the Go programing language**.

It **uses [Go HTTP standard](https://golang.org/pkg/net/http) model**.

An app created with go-app can out of the box **run in its own window**, **supports offline mode**, and are **SEO friendly**.

## Declarative syntax

Go-app uses a [declarative syntax](/syntax) so you can **write reusable component-based UI elements** just by using the Go programming language.

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
			app.If(h.name != "",
				app.Text(h.name),
			).Else(
				app.Text("World!"),
			),
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
	app.Route("/", &hello{})
	app.Route("/hello", &hello{})
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

## Other features

- [SEO friendly](/seo)
- [Installable](/install)
- Offline mode

## Next

- [Getting started](/start)
- [Understand go-app architecture](/architecture)
- [API reference](/reference)
