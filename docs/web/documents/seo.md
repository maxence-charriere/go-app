## Intro

[SEO](https://en.wikipedia.org/wiki/Search_engine_optimization) (Search engine optimization) is to make web content readable and indexable by search engines such as [Google](https://google.com).

Because app built with go-app are Go binaries loaded in the web browser and their content is generated dynamically on the client-side, it is tricky for search engines to crawl and index apps content.

Go-app is solving that issue by prerendering components on the server-side before including the generated HTML markup into requested pages.

## Prerendering

Prerendering is about converting a [component](/components) to its plain HTML representation.

On the server-side, when the [Handler](/reference#Handler) receives a page request, it creates an instance of the component associated with the requested URL, generates its HTML representation, and includes it on the page.

For a simple Hello world component:

```go
type hello struct {
	app.Compo
}

func (h *hello) Render() app.UI {
	return app.H1().Text("Hello World!")
}
```

The generated page will look like the following code:

```html
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8" />
    <meta httpequiv="Content-Type" content="text/html; charset=utf-8" />
    <link rel="manifest" href="/manifest.webmanifest" />
    <link type="text/css" rel="stylesheet" href="/app.css" />
    <script defer src="/wasm_exec.js"></script>
    <script src="/app.js" defer></script>
  </head>

  <body>
    <div>
      <!-- Prerendered component is included in the div below: -->
      <div id="app-pre-render">
        <h1>Hello World!</h1>
      </div>

      <div id="app-wasm-loader" class="goapp-app-info">
        <img
          class="goapp-logo goapp-spin"
          src="https://storage.googleapis.com/murlok-github/icon-192.png"
          id="app-wasm-loader-icon"
        />
        <p id="app-wasm-loader-label" class="goapp-label">Loading...</p>
      </div>
    </div>
    <div id="app-end"></div>
  </body>
</html>
```

### Customizing prerendering

Like on the client-side, a component might require further initializations. Launching additional instructions is done by implementing the [PreRender](/reference#PreRenderer) interface.

Here is a Hello example that gets the username from an URL parameter:

```go
type hello struct {
	app.Compo

	name string
}

func (h *hello) OnPreRender(ctx app.Context) {
	username := ctx.Page.URL().Query().Get("username")
	h.name = username
}

func (h *hello) Render() app.UI {
	return app.H1().Text("Hello " + h.name)
}
```

### Customizing page metadata

An essential step for a good SEO is to have meta tags, such as page title, well-formed.

Page metadata can be set from the `OnPreRender` [Context](/reference#Context) argument. Context contains a [page](/reference#Page) field to set page metadata. Here is an example that sets the page title and author.

```go
func (h *hello) OnPreRender(ctx app.Context) {
	ctx.Page.SetTitle("A Hello World written with go-app")
	ctx.Page.SetAuthor("Maxence")
}
```

**See the [page reference](/reference#Page) for the detail of customizable metadata**.

### Caching

For performance reasons, prerendered pages are cached once generated.

The [Handler](/reference#Handler) provides a `PreRenderCache` field that allows customizing the cache behavior. By default, it uses an in-memory [LRU cache](<https://en.wikipedia.org/wiki/Cache_replacement_policies#Least_recently_used_(LRU)>) that keeps cached data for 24 hours with a maximum size of 8MB.

Cache behavior can be customized by setting the PreRendering cache to [another LRU cache](/reference#NewPreRenderLRUCache) with different values:

```go
h := app.Handler{
		Name:           "Hello world",
		PreRenderCache: app.NewPreRenderLRUCache(100*100000, time.Hour), // 10MB/1hour
}
```

Or any other cache that satisfies the [PreRenderCache](/reference#PreRenderCache) interface:

```go
type PreRenderCache interface {
    // Get returns the item at the given path.
    Get(ctx context.Context, path string) (PreRenderedItem, bool)

    // Set stored the item at the given path.
    Set(ctx context.Context, i PreRenderedItem)
}
```

You could implement a cache based on Redis or any other datastore.

## Next

- [Lifecycle and Updates](/lifecycle)
- [Reference](/reference)
