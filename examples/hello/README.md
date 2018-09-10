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
		URL: "/hello.Hello",
	})
}
```

```bash
# In $GOPATH/src/github.com/murlokswarm/app/examples/hello/bin/hello-mac
# Build and launch the app
goapp mac run
```


## main (web)
![hello](https://github.com/murlokswarm/app/wiki/assets/hello-web.png)

```go
func main() {
	app.Import(&hello.Hello{})

	app.Run(&web.Driver{
		URL: "/hello.Hello",
	})
}
```

```bash
# In $GOPATH/src/github.com/murlokswarm/app/examples/hello/bin/hello-web
# Build/launch server and client
goapp web run -b
```
