# Testing

Testing is an essential step to achieve app reliability. Since go-app is working on 2 different environments (web browser and server), it provides 2 testing [dispatchers](/reference#Dispatcher) to emulate [components lifecycle](/components#lifecycle) behaviors.

## Component server prerendering

[Prerendering](/components#prerender) is a component lifecycle step where a component can be initialized on the server-side before being converted into HTML. The server-side environment can be emulated with a dispatcher created with the [NewServerTester()](/reference#NewServerTester) function.

Here is an example that tests if a component has the expected values after the PreRenderer interface call:

```go
package main

import (
	"testing"

	"github.com/maxence-charriere/go-app/v8/pkg/app"
)

type aTitle struct {
	app.Compo

	title string
}

func (t *aTitle) OnPreRender(ctx app.Context) {
	t.title = "Testing Prerendering"
	t.Update()
}

func (t *aTitle) Render() app.UI {
	return app.H1().
		Class("title").
		Text(t.title)
}

func TestComponentPreRendering(t *testing.T) {
	compo := &aTitle{}

	// Creating the server emulator:
	disp := app.NewServerTester(compo)
	defer disp.Close() // Releases alocated resources.

	if compo.title != "" {
		t.Fatal("component title is not empty")
	}

	// Call OnPreRender() from PreRenderer interface:
	disp.PreRender()

	// When using Update(), Dispatch() ,or Defer(), operation are queued in
	// a go channel. Consume() execute pending operations:
	disp.Consume()

	if compo.title != "Testing Prerendering" {
		t.Fatal("bad component title:", compo.title)
	}
}
```

## Component client lifecycle

Like on the [server-side](#testing-component-server-prerendering), testing a component on the client-side is done by emulating the corresponding environment. On the client-side, it is done with the [NewClientTester()](/reference#NewClientTester) function.

Here is an example that tests if a component has the expected values after mounting and navigation:

```go

import (
	"net/url"
	"testing"

	"github.com/maxence-charriere/go-app/v8/pkg/app"
)

type aTitle struct {
	app.Compo

	title string
}

func (t *aTitle) OnMount(ctx app.Context) {
	t.title = "Testing Mounting"
	t.Update()
}

func (t *aTitle) OnNav(ctx app.Context) {
	t.title = "Testing Nav"
	t.Update()
}

func (t *aTitle) Render() app.UI {
	return app.H1().
		Class("title").
		Text(t.title)
}

func TestComponentLifcycle(t *testing.T) {
	compo := &aTitle{}

	disp := app.NewClientTester(compo)
	defer disp.Close()

	disp.Consume()
	if compo.title != "Testing Mounting" {
		t.Fatal("bad component title:", compo.title)
	}

	disp.Nav(&url.URL{})
	disp.Consume()
	if compo.title != "Testing Nav" {
		t.Fatal("bad component title:", compo.title)
	}
}


```

See [ClientDispatcher](/reference#ClientDispatcher) for other lifecycle and component extension events.

## Asynchronous operations

## UI elements

## Next

- [Understand go-app architecture](/architecture)
- [How to create a component](/components)
- [Handle concurrency](/concurrency)
- [API reference](/reference)
