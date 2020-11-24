# Concurrency

![concurrency.png](/web/images/concurrency.png)

## UI goroutine

The UI goroutine is the app's main goroutine. Under the hood, it is an event loop where each event is executed synchronously.

Here are the events that are always executed on the UI goroutine:

- [Component lifecycle events](/components#lifecycle) (OnMount, OnNav, OnDismount)
- [Component updates](/components#update)
- HTML element [event handlers](/syntax#event-handlers)
- [Dispatch()](#dispatch) calls

## Standard goroutines

Those are standard goroutines. They are executed in parallel with the UI goroutine.

```go
go func() {
	// ...
}()
```

### When to use?

Since rendering operations are executed on the [UI goroutine](#ui-goroutine), **blocking and long operations should be executed in another goroutine**. That will prevent the UI to feel unresponsive.

If those operations lead to component field modifications, make sure to perform them back on the UI goroutine by calling [Dispatch()](#dispatch).

## Dispatch()

[Dispatch](/reference#Dispatch) is a call that makes the given function to be executed on the UI goroutine.

```go
func Dispatch(func () {
    // ...
})
```

Here is an example that asynchronously performs an HTTP request and displays the response.

```go
type httpCall struct {
	app.Compo

	response string
}

func (c *httpCall) Render() app.UI {
	return app.Div().Body(
		app.H1().Text("HTTP Call"),

		app.H2().Text("URL:"),
		app.Input().
			Placeholder("Enter an URL").
			OnChange(c.OnURLChange),

		app.H2().Text("Response:"),
		app.P().Text(c.response),
	)
}

func (c *httpCall) OnURLChange(ctx app.Context, e app.Event) {
	// Reseting response value:
	c.response = ""
	c.Update()

	// Launching HTTP request:
	url := ctx.JSSrc.Get("value").String()
	go c.doRequest(url) // Performs blocking operation on a new goroutine.
}

func (c *httpCall) doRequest(url string) {
	r, err := http.Get(url)
	if err != nil {
		c.updateResponse(err.Error())
		return
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.updateResponse(err.Error())
		return
	}

	c.updateResponse(string(b))
}

func (c *httpCall) updateResponse(res string) {
	app.Dispatch(func() { // Ensures response field is updated on UI goroutine.
		c.response = res
		c.Update()
	})
}
```

## Next

- [Interact with Javascript](/js)
- [API reference](/reference)
