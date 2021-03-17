package main

import "github.com/maxence-charriere/go-app/v8/pkg/app"

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
	}
	l.Update()
}

func (l *foodList) Render() app.UI {
	return app.Div().Body(
		app.H1().Text("Dealing with List"),
		app.Ul().Body(
			app.Range(l.foods).Slice(func(i int) app.UI {
				f := l.foods[i]

				return app.Li().
					Text("Click to eat "+f.label).
					// Here the eat() function is removing the clicked element.
					//
					// When updating the list of children, event handlers are
					// updated if they have a different pointer address. In that
					// case, the eat function() address will always be the same.
					//
					// Giving scope to the event handler will prevent reusing the
					// former event handler that is associated with a deleted
					// element.
					OnClick(l.eat(f), f.id)
			}),
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
