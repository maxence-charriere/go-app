package main

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v8/pkg/app"
)

type issue499Data struct {
	ID    int
	Value string
}

type issue499 struct {
	app.Compo

	data []issue499Data
}

func newIssue499Data() *issue499 {
	return &issue499{}
}

func (c *issue499) OnMount(app.Context) {
	c.data = []issue499Data{
		{11, "one"},
		{22, "two"},
		{33, "three"},
		{44, "four"},
		{55, "five"},
		{66, "six"},
		{77, "sever"},
		{88, "eight"},
		{99, "nine"},
	}
	c.Update()
}

func (c *issue499) Render() app.UI {
	return app.Div().Body(
		app.H1().Text("Issue 499"),
		app.Div().
			Body(
				app.Range(c.data).Slice(func(i int) app.UI {
					d := c.data[i]
					return app.Button().
						ID(fmt.Sprintf("elem-%v", d.ID)).
						OnClick(c.newListener(d.ID), d.ID).
						Text(d.Value)
				}),
			),
	)
}

func (c *issue499) newListener(id int) app.EventHandler {
	return func(app.Context, app.Event) {
		for i, d := range c.data {
			if id == d.ID {
				c.data = append(c.data[:i], c.data[i+1:]...)
				c.Update()
				return
			}
		}
	}
}
