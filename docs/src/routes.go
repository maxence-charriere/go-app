package main

import "github.com/maxence-charriere/go-app/v7/pkg/app"

func init() {
	for path, new := range pages() {
		app.Route("/"+path, new())
	}
}
