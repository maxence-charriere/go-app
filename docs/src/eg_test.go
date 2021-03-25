package main

import (
	"net/url"
	"testing"
	"time"

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

func (t *aTitle) setAsyncTitle(ctx app.Context) {
	ctx.Async(func() {
		time.Sleep(time.Millisecond * 100)
		t.Defer(func(ctx app.Context) {
			t.title = "Testing Async"
		})
	})
}

func TestComponentPreRendering(t *testing.T) {
	compo := &aTitle{}

	// Creating the server emulator:
	disp := app.NewServerTester(compo)
	defer disp.Close() // Releases alocated resources.

	if compo.title == "Testing Prerendering" {
		t.Fatal("bad component title:", compo.title)
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

func TestComponentLifcycle(t *testing.T) {
	compo := &aTitle{}

	disp := app.NewClientTester(compo)
	defer disp.Close()

	disp.Nav(&url.URL{})
	disp.Consume()
	if compo.title != "Testing Nav" {
		t.Fatal("bad component title:", compo.title)
	}

}

func TestComponentAsync(t *testing.T) {
	compo := &aTitle{}

	disp := app.NewClientTester(compo)
	defer disp.Close()

	compo.setAsyncTitle(disp.Context()) // Async operation queued.
	disp.Consume()                      // Async operation launched but not completed.
	if compo.title == "Testing Async" {
		t.Fatal("bad component title:", compo.title)
	}

	disp.Wait()    // Wait for the async operations do complete.
	disp.Consume() // Apply changes.
	if compo.title != "Testing Async" {
		t.Fatal("bad component title:", compo.title)
	}
}

func TestUIElement(t *testing.T) {
	compo := &aTitle{}
	disp := app.NewClientTester(compo)
	defer disp.Close()

	app.TestMatch(compo, app.TestUIDescriptor{
		Path:     app.TestPath(0), // Component root.
		Expected: app.H2().Text("Testing Mounting"),
	})
}
