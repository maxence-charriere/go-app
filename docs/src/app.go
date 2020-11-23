// +build wasm

package main

import (
	"github.com/maxence-charriere/go-app/v7/pkg/app"
)

func main() {
	app.Route("/", newStart())
	app.Route("/start", newStart())
	app.Route("/architecture", newArchitecture())
	app.Route("/reference", newReference())
	app.Route("/components", newCompo())
	app.Route("/concurrency", newConcurrency())
	app.Route("/syntax", newSyntax())
	app.Route("/js", newJS())
	app.Route("/routing", newRouting())
	app.Route("/static-resources", newStaticResources())
	app.Route("/built-with", newBuiltWith())
	app.Run()
}
