package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/maxence-charriere/go-app/v7/pkg/app"
)

type document struct {
	app.Compo

	path  string
	title string
}

func newDocument(path string) *document {
	return &document{path: path}
}

func (d *document) Title(t string) *document {
	d.title = t
	return d
}

func (d *document) Render() app.UI {
	return app.Main().Body()
}

type complexScenario struct {
	app.Compo

	Total int
	OK    int
	KO    int
	Logs  []string
}

func (c *complexScenario) Render() app.UI {
	return app.Div().
		Body(
			app.H1().Text("Complex scenario"),
			app.P().Text("Example for a complex scenario"),
			app.P().Body(
				app.Button().
					Disabled(c.OK+c.KO != c.Total). // This could also be done with css class.
					OnClick(c.onButtonClick).
					Text("Start"),
			),

			app.H2().Text("Progress"),
			app.P().Body(
				app.Progress().
					Value(c.OK+c.KO).
					Max(c.Total),
			),

			app.H2().Text("Logs"),
			app.P().Body(
				app.Range(c.Logs).Slice(func(i int) app.UI {
					return app.Div().Text(c.Logs[i])
				}),
			),
		)
}

func (c *complexScenario) onButtonClick(ctx app.Context, e app.Event) {
	c.Total = 5
	c.OK = 0
	c.KO = 0
	c.Logs = nil
	c.Update()

	go c.fakeJob(1, true)
	go c.fakeJob(2, true)
	go c.fakeJob(3, true)
	go c.fakeJob(4, true)
	go c.fakeJob(5, false)
}

func (c *complexScenario) fakeJob(id int, result bool) {
	c.log("launching job %v", id)

	d := time.Duration(rand.Intn(5)+1) * time.Second
	time.Sleep(d)

	if !result {
		c.incKO()
		c.log("job %v failed", id)
		return
	}

	c.incOK()
	c.log("job %v succeeded", id)
}

// Functions below are the one that update the UI. Since they are called from
// new goroutines, Dispatch() function is used to ensure the component fieldd
// are updated on the UI goroutine.

func (c *complexScenario) incOK() {
	app.Dispatch(func() {
		c.OK++
		c.Update()
	})
}

func (c *complexScenario) incKO() {
	app.Dispatch(func() {
		c.KO++
		c.Update()
	})
}

func (c *complexScenario) log(format string, v ...interface{}) {
	app.Dispatch(func() {
		c.Logs = append(c.Logs, fmt.Sprintf(format, v...))
		c.Update()
	})
}
