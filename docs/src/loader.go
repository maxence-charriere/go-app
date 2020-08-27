package main

import "github.com/maxence-charriere/go-app/v7/pkg/app"

type loader struct {
	app.Compo

	IDescription string
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

func (l *loader) Err(err error) *loader {
	l.IErr = err
	return l
}

func (l *loader) Loading(v bool) *loader {
	l.ILoading = v
	return l
}

func (l *loader) Render() app.UI {
	display := "none"
	if l.ILoading || l.IErr != nil {
		display = "bloc"
	}

	return app.Div().
		Style("display", display).
		Body(
			app.Stack().
				Class("loader").
				Center().
				Content(
					app.Div().
						Class("margin").
						Body(
							app.If(l.IErr == nil,
								app.Div().Class("icon"),
							),
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
