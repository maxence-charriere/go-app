# Javascript and DOM

Since WebAssembly is browser-based technology, some scenarios may require DOM access and JavaScript calls.

This is usually done with the help of [syscall/js](https://golang.org/pkg/syscall/js/) but for compatibility and tooling reasons, **go-app wraps the standard package**. Interacting with Javascript is done by using the [Value](/reference#Value) interface.

This article provide examples that show common interactions with Javascript.

## Include JS files

Building UIs can sometime require the need of third party Javascript libraries. Those libraries can either be included when the page is loaded or inlined in a component.

### Handler

This can be done by including a Javascript file URL in the `Scripts` field from the [app.Handler](/reference#Handler):

```go
handler := &app.Handler{
	Name: "My App",
	Scripts: []string{
		"/web/myscript.js",                // Local script
		"https://foo.com/remoteScript.js", // Remote script
	},
}
```

Or by directly putting code in the `RawHeaders` field:

```go
handler := &app.Handler{
	Name: "My App",
	RawHeaders: []string{
		`<!-- Global site tag (gtag.js) - Google Analytics -->
		<script async src="https://www.googletagmanager.com/gtag/js?id=UA-xxxxxxx-x"></script>
		<script>
		  window.dataLayer = window.dataLayer || [];
		  function gtag(){dataLayer.push(arguments);}
		  gtag('js', new Date());

		  gtag('config', 'UA-xxxxxx-x');
		</script>
		`,
	},
}
```

### Inline

Javascript file can also be included in components by using a [Script](/reference#Script) element. Here is an example that asynchronously load Youtube Iframe API script.

```go
type youtubePlayer struct {
	app.Compo
}

func (p *youtubePlayer) Render() app.UI {
	return app.Div().Body(
		app.Script().
			Src("//www.youtube.com/iframe_api").
			Async(true),
		app.IFrame().
			ID("youtube-player").
			Allow("autoplay").
			Allow("accelerometer").
			Allow("encrypted-media").
			Allow("picture-in-picture").
			Sandbox("allow-presentation allow-same-origin allow-scripts allow-popups").
			Src("https://www.youtube.com/embed/LqeRF_0DDCg"),
	)
}
```

## Window

```go
app.Window()
```

[Window()](/reference#Window) returns a global javascript object representing the [browser window](/reference#BrowserWindow) that can be used to call functions with `window` and empty namespaces.

### Get element by ID

`GetElementByID()` allow to get a DOM element from an ID.

```js
// JS version:
let elem = document.getElementById("YOUR_ID");
```

```go
// Go equivalent:
elem := app.Window().GetElementByID("YOUR_ID")
```

It is a helper function equivalent to:

```go
elem := app.Window().
    Get("document").
    Call("getElementById","YOUR_ID")
```

### Create JS object

Creating an object from a library is done by getting its name from the `Window` and call the `New()` function.

Here is an example about how to create a Youtube player:

```js
// JS version:
let player = new YT.Player("player", {
  height: "390",
  width: "640",
  videoId: "M7lc1UVf-VE",
});
```

```go
// Go equivalent:
player := app.Window().
	Get("YT").
	Get("Player").
	New("player", map[string]interface{}{
		"height":  390,
		"width":   640,
		"videoId": "M7lc1UVf-VE",
    })
```

## Cancel an event

When implementing an [event handler](/reference#EventHandler), the event can be canceled by calling [PreventDefault()](/reference#Event.PreventDefault).

```go
type foo struct {
	app.Compo
}

func (f *foo) Render() app.UI {
	return app.Div().
		OnChange(f.onContextMenu).
		Text("Don't copy me!")
}

func (f *foo) onContextMenu(ctx app.Context, e app.Event) {
	e.PreventDefault()
}
```

## Get input value

Input are usually used to get a user inputed value. Here is how to get that value when implementing an [event handler](/reference#EventHandler):

```go
type foo struct {
    app.Compo
}

func (f *foo) Render() app.UI {
    return app.Input().OnChange(f.onInputChange)
}

func (f *foo) onInputChange(ctx app.Context, e app.Event) {
    v := ctx.JSSrc.Get("value").String()
}
```

## Next

- [Customize components with the declarative syntax](/syntax)
- [Deal with static resources](/static-resources)
- [API reference](/reference)
