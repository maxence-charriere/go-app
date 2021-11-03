## Intro

Since WebAssembly is browser-based technology, some scenarios may require DOM access and JavaScript calls.

This is usually done with the help of [syscall/js](https://golang.org/pkg/syscall/js/) but for compatibility and tooling reasons, **go-app wraps the JS standard package**. Interacting with JavaScript is done by using the [Value](/reference#Value) interface.

This article provides examples that show common interactions with JavaScript.

## Include JS files

Building UIs can sometimes require the need of third-party JavaScript libraries. Those libraries can either be included at the [page](/architecture#html-pages) level or inlined in a [component](/components).

### Page's scope

JS files can be included on a page by using the [Handler](/reference#Handler) `Scripts` field:

```go
handler := &app.Handler{
	Name: "My App",
	Scripts: []string{
		"/web/myscript.js",                // Local script
		"https://foo.com/remoteScript.js", // Remote script
	},
}
```

Or by directly putting JS markup in the `RawHeaders` field:

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

### Inlined in Components

JS files can also be included directly inlined into [components](/components) in the `Render()` method by using the [\<script\>](/reference#Script) HTML element.

The following example asynchronously loads a YouTube video into an `<iframe>`, using a YouTube JavaScript file:

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
			ID("yt-container").
			Allow("autoplay").
			Allow("accelerometer").
			Allow("encrypted-media").
			Allow("picture-in-picture").
			Sandbox("allow-presentation allow-same-origin allow-scripts allow-popups").
			Src("https://www.youtube.com/embed/LqeRF_0DDCg"),
	)
}
```

## Using window global object

The `window` JS global object is usable from the [Window](/reference#Window) function.

```go
app.Window()
```

### Get element by ID

`GetElementByID()` is to get a DOM element from an ID.

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
let player = new YT.Player("yt-container", {
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
	New("yt-container", map[string]interface{}{
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
    v := ctx.JSSrc().Get("value").String()
}
```

## Next

- [Concurrency](/concurrency)
- [Reference](/reference)
