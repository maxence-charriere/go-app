package main

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type homePage struct {
	app.Compo
}

func newHomePage() *homePage {
	return &homePage{}
}

func (p *homePage) Render() app.UI {
	return newPage()
}
