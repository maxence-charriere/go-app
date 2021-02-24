package main

import "github.com/maxence-charriere/go-app/v8/pkg/app"

type link struct {
	app.Compo

	Iclass   string
	Ilabel   string
	Ihref    string
	Ifocus   bool
	IonClick func()
	Iicon    app.UI
}

func newLink() *link {
	return &link{}
}

func (l *link) Class(v string) *link {
	if v == "" {
		return l
	}
	if l.Iclass != "" {
		l.Iclass += " "
	}
	l.Iclass += v
	return l
}

func (l *link) Label(v string) *link {
	l.Ilabel = v
	return l
}

func (l *link) Href(v string) *link {
	l.Ihref = v
	return l
}

func (l *link) Focus(v bool) *link {
	l.Ifocus = v
	return l
}

func (l *link) OnClick(v func()) *link {
	l.IonClick = v
	return l
}

func (l *link) Icon(v app.UI) *link {
	l.Iicon = v
	return l
}

func (l *link) Render() app.UI {
	iconVisibility := ""
	if l.Iicon == nil {
		iconVisibility = "hide"
	}

	focus := ""
	if l.Ifocus {
		focus = "focus"
	}

	return app.A().
		Class("link").
		Class("heading").
		Class("fit").
		Class(l.Iclass).
		Class(focus).
		Href(l.Ihref).
		OnClick(l.onClick).
		Body(
			app.Stack().
				Center().
				Content(
					app.Div().
						Class(iconVisibility).
						Class("link-icon").
						Body(l.Iicon),
					app.Div().Text(l.Ilabel),
				),
		)
}

func (l *link) onClick(ctx app.Context, e app.Event) {
	if l.IonClick == nil {
		return
	}

	e.PreventDefault()
	l.IonClick()
}
