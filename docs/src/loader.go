package main

import "github.com/maxence-charriere/go-app/v7/pkg/app"

type loader struct {
	app.Compo

	IDescription string
	IClass       string
	ILoading     bool
	IErr         error
}

func newLoader() *loader {
	return &loader{}
}

func (l *loader) Description(v string) *loader {
	l.IDescription = v
	return l
}

func (l *loader) Class(v string) *loader {
	if l.IClass != "" {
		l.IClass += " "
	}
	l.IClass += v

	return l
}

func (l *loader) Err(err error) *loader {
	l.IErr = err
	return l
}

func (l *loader) Loading(v bool) *loader {
	l.ILoading = v
	return l
}

func (l *loader) Render() app.UI {
	visibility := "hidden"
	if l.ILoading || l.IErr != nil {
		visibility = ""
	}

	return app.Stack().
		Class("loader").
		Class(visibility).
		Center().
		Content(
			app.Div().
				Class("content").
				Body(
					app.Stack().
						Center().
						Content(
							app.Div().
								Class("icon-brackground").
								Body(
									app.If(l.IErr == nil,
										app.Div().Class("icon"),
									),
								),
							app.Div().
								Class("info").
								Body(
									app.Div().
										Class("label").
										Text("Loading"),
									app.If(l.IErr != nil,
										app.Div().
											Class("error").
											Text(l.IErr),
									).ElseIf(l.IDescription != "",
										app.Div().Text(l.IDescription),
									),
								),
							app.Div().Class("margin"),
						),
				),
		)
}
