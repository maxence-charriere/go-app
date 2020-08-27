package main

import "github.com/maxence-charriere/go-app/v7/pkg/app"

type loader struct {
	app.Compo

	IDescription string
	IState       string
	IErr         error
}

func newLoader() *loader {
	return &loader{
		IState: "none",
	}
}

func (l *loader) Description(v string) *loader {
	l.IDescription = v
	return l
}

func (l *loader) Err(err error) *loader {
	l.IErr = err
	return l
}

func (l *loader) Loading(v bool) *loader {
	if v {
		l.IState = "bloc"
	}

	return l
}

func (l *loader) Render() app.UI {
	return app.Div().
		Style("display", l.IState).
		Body(
			app.Stack().
				Class("loader").
				Center().
				Content(
					app.Div().Class("margin").Body(
						app.Div().Class("icon"),
					),
					app.Div().
						Class("content").
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
		)
}
