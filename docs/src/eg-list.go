package main

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v8/pkg/app"
)

type food struct {
	id    int
	label string
}

type foodList struct {
	app.Compo
	foods []food
}

func newFoodList() *foodList {
	return &foodList{}
}

func (l *foodList) OnMount(ctx app.Context) {
	l.initFood(ctx)
}

func (l *foodList) initFood(ctx app.Context) {
	l.foods = []food{
		{
			id:    1,
			label: "French fries",
		},
		{
			id:    2,
			label: "Pasta",
		},
		{
			id:    3,
			label: "Rice",
		},
		{
			id:    4,
			label: "Steak",
		},
		{
			id:    5,
			label: "Fish",
		},
	}
	l.Update()
}

func (l *foodList) Render() app.UI {
	return app.Div().Body(
		app.H1().Text("Dealing with List"),
		app.P().
			Class("center-content").
			Body(
				app.Range(l.foods).Slice(func(i int) app.UI {
					f := l.foods[i]

					return app.Div().
						Class("button").
						Text("Remove "+f.label).
						// The eat() function is removing the clicked element.
						//
						// When updating the list of children, event handlers are
						// updated if they have a different pointer address. In that
						// case, the eat function() address will always be the same.
						//
						// Here, a scope is given to the event handler to
						// prevent reusing the former event handler that is
						// associated with a deleted element.
						OnClick(l.eat(f), f.id)
				}),
			),
		app.P().
			Class("center-content").
			Body(
				app.Span().
					Class("button").
					Text("Reset").
					OnClick(l.reset),
			),
	)
}

func (l *foodList) eat(f food) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		for i, fi := range l.foods {
			if fi.id == f.id {
				// Removing food item:
				copy(l.foods[i:], l.foods[i+1:])
				l.foods = l.foods[:len(l.foods)-1]
				l.Update()
				return
			}
		}
	}
}

func (l *foodList) reset(ctx app.Context, e app.Event) {
	l.initFood(ctx)
}

// This will be replaced later by Go 1.16 embed once it will be more widely s
// upported.
var egListCode = fmt.Sprintf("```go\n%s\n```", `
package main

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v8/pkg/app"
)

type food struct {
	id    int
	label string
}

type foodList struct {
	app.Compo
	foods []food
}

func newFoodList() *foodList {
	return &foodList{}
}

func (l *foodList) OnMount(ctx app.Context) {
	l.initFood(ctx)
}

func (l *foodList) initFood(ctx app.Context) {
	l.foods = []food{
		{
			id:    1,
			label: "French fries",
		},
		{
			id:    2,
			label: "Pasta",
		},
		{
			id:    3,
			label: "Rice",
		},
		{
			id:    4,
			label: "Steak",
		},
		{
			id:    5,
			label: "Fish",
		},
	}
	l.Update()
}

func (l *foodList) Render() app.UI {
	return app.Div().Body(
		app.H1().Text("Dealing with List"),
		app.P().
			Class("center-content").
			Body(
				app.Range(l.foods).Slice(func(i int) app.UI {
					f := l.foods[i]

					return app.Div().
						Class("button").
						Text("Remove "+f.label).
						// The eat() function is removing the clicked element.
						//
						// When updating the list of children, event handlers are
						// updated if they have a different pointer address. In that
						// case, the eat function() address will always be the same.
						//
						// Here, a scope is given to the event handler to
						// prevent reusing the former event handler that is
						// associated with a deleted element.
						OnClick(l.eat(f), f.id)
				}),
			),
		app.P().
			Class("center-content").
			Body(
				app.Span().
					Class("button").
					Text("Reset").
					OnClick(l.reset),
			),
	)
}

func (l *foodList) eat(f food) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		for i, fi := range l.foods {
			if fi.id == f.id {
				// Removing food item:
				copy(l.foods[i:], l.foods[i+1:])
				l.foods = l.foods[:len(l.foods)-1]
				l.Update()
				return
			}
		}
	}
}

func (l *foodList) reset(ctx app.Context, e app.Event) {
	l.initFood(ctx)
}
`)
