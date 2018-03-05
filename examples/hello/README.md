# hello
The hello world built with this package.

![hello](https://github.com/murlokswarm/app/wiki/assets/hello.gif)

## component
```go
// Hello is a hello world component.
type Hello struct {
	Name string
}

// Render returns a string that describes the component markup.
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
	<input value="{{.Name}}" placeholder="Say something..." onchange="Name" autofocus>
</div>
	`

```


## main (mac)
![hello](https://github.com/murlokswarm/app/wiki/assets/hello-mac.png)

```go
func main() {
	app.Import(&hello.Hello{})

	app.Run(&mac.Driver{
		OnRun: func() {
			newWindow()
		},

		OnReopen: func(hasVisibleWindow bool) {
			if !hasVisibleWindow {
				newWindow()
			}
		},
	})

}

func newWindow() {
	app.NewWindow(app.WindowConfig{
		Title:           "hello world",
		TitlebarHidden:  true,
		Width:           1280,
		Height:          768,
		BackgroundColor: "#21252b",
		DefaultURL:      "/hello.Hello",
	})
}
```

```bash
# Build app
go build

## Launch app
```


## main (web)
![hello](https://github.com/murlokswarm/app/wiki/assets/hello-web.png)

```go
func main() {
	app.Import(&hello.Hello{})

	app.Run(&web.Driver{
		DefaultURL: "/hello.Hello",
	})
}
```

```bash
# Build server and client
goapp web build

# Launch server
./hello

# Launch client (mac)
open http://localhost:7042

# Launch client (windows)
explorer http://localhost:7042

# Launch client (linux)
xdg-open http://localhost:7042
```