# Getting Started

GoApp is a package for building fast and modern Progressive Web Apps (PWAs) using Go and WebAssembly.

It allows you to build web user interfaces using Go syntax by composing HTML-based components. Here's a minimal "Hello World" example:

```go
func newHello() app.Composer {
    return &hello{}
}

type hello struct {
    app.Compo
}

func (h *hello) Render() app.UI {
    return app.H1().Text("Hello World!")
}
```

The component is then associated with a path using an HTTP handler that follows the Go standard library model:

```go
func main() {
    app.Route("/", newHello)
    app.RunWhenOnBrowser()

    http.Handle("/", &app.Handler{
        Name:        "Hello",
        Description: "A Hello World example",
    })

    if err := http.ListenAndServe(":8000", nil); err != nil {
        log.Fatal(err)
    }
}
```

The program is compiled into two binaries:

- A WebAssembly (WASM) binary, executed in the web browser:

  ```sh
  GOARCH=wasm GOOS=js go build -o web/app.wasm
  ```

- A standard Go binary, which serves the WebAssembly file, support files, and static resources:
  ```sh
  go build
  ```
