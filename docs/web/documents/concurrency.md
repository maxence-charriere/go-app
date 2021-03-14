# Concurrency

In go-app, every event and user interaction are handled on a single goroutine. Because some scenarios can have a long execution time, like performing an HTTP request, there is a risk that the UI feels slow or unresponsive.

This document describes how it works and what tools go-app provides to solve this problem.

![concurrency.png](/web/images/concurrency.svg)

## UI goroutine

The UI goroutine is the app's main goroutine. Under the hood, it is an event loop where each event is executed synchronously.

Here are the events that are always executed on the UI goroutine:

- [Component lifecycle events](/components#lifecycle) (OnMount, OnNav, OnDismount)
- [Component updates](/components#update)
- HTML element [event handlers](/syntax#event-handlers)
- [Dispatch()](#dispatch) calls
- [Defer()](/reference#Compo.Defer) calls

## Async

```go
func (ctx Context) Async(fn func())
```

[Async()](/reference#Context.Async) is a [Context](/reference#Context) method that executes a given function on a new goroutine. It is usually used to perform long or blocking operations.

Here is an example where an HTTP request is performed when a page is loaded.

```go
type foo struct {
	app.Compo
}

func (f *foo) OnNav(ctx app.Context) {
	// Launching a new goroutine:
	ctx.Async(func() {
		r, err := http.Get("/bar")
		if err != nil {
			app.Log(err)
			return
		}
		defer r.Body.Close()

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			app.Log(err)
			return
		}

		app.Logf("request response: %s", b)
	})
}
```

The difference with manually launching a goroutine is that go-app has no insights about when a manually launched goroutine ceases its execution. It's not a problem on the client-side but when prerendering on the server-side, go-app has to wait for all launched goroutines to finish their jobs in order to properly generate HTML markup. Therefore, manually launching a goroutine for UI-related purposes introduces reliability issues on the server-side.

**Prefer the use of [Async()](/reference#Context.Async) rather than manually launching a goroutine when dealing with UI**.

## Dispatch

```go
func (ctx Context) Dispatch(fn func())
```

[Dispatch()](reference#Context.Dispatch) is a [Context](/reference#Context) method that executes a given function on the [UI goroutine](#ui-goroutine). It is used to update the UI after an [Async()](#async) call, in order to avoid concurrent calls when updating a component field.

Here is an example where an HTTP request is performed when a page is loaded, and its result is stored in a component field:

```go
type foo struct {
	app.Compo

	response []byte
}

func (f *foo) OnNav(ctx app.Context) {
	// Launching a new goroutine:
	ctx.Async(func() {
		r, err := http.Get("/bar")
		if err != nil {
			app.Log(err)
			return
		}
		defer r.Body.Close()

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			app.Log(err)
			return
		}

		// Storing HTTP response in component field:
		ctx.Dispatch(func() {
			f.response = b
			f.Update()
		})
	})
}
```

**Always update a component fields on the [UI goroutine](#ui-goroutine)!**

## Defer

```go
func (c *Compo) Defer(fn func(Context))
```

[Defer()](/reference#Compo.Defer) is a [Compo](/reference#Compo) method that like [Dispatch()](#dispatch), executes a given function on the [UI goroutine](#ui-goroutine). The difference is that the given function is executed only if the component is still mounted when it is called.

Here is an example where an HTTP request is performed when a page is loaded, and its result is stored in a component field with `Defer()`:

```go
type foo struct {
	app.Compo

	response []byte
}

func (f *foo) OnNav(ctx app.Context) {
	// Launching a new goroutine:
	ctx.Async(func() {
		r, err := http.Get("/bar")
		if err != nil {
			app.Log(err)
			return
		}
		defer r.Body.Close()

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			app.Log(err)
			return
		}

		// Storing HTTP response in component field:
		f.Defer(func(app.Context) {
			f.response = b
			f.Update()
		})
	})
}
```

## Next

- [Interact with Javascript](/js)
- [API reference](/reference)
