package main

import (
	"net/url"
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

func TestComponentLifcycle(t *testing.T) {
	compo := &aTitle{}

	disp := app.NewClientTester(compo)
	defer disp.Close() // Releases alocated resources.

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
