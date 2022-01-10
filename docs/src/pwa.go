//go:generate go run gen/godoc.go
//go:generate go fmt

package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/analytics"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/maxence-charriere/go-app/v9/pkg/ui"
)

const (
	defaultTitle       = "A Go package for building Progressive Web Apps"
	defaultDescription = "A package for building progressive web apps (PWA) with the Go programming language (Golang) and WebAssembly (Wasm). It uses a declarative syntax that allows creating and dealing with HTML elements only by using Go, and without writing any HTML markup."
	backgroundColor    = "#2e343a"

	buyMeACoffeeURL     = "https://www.buymeacoffee.com/maxence"
	openCollectiveURL   = "https://opencollective.com/go-app"
	githubURL           = "https://github.com/maxence-charriere/go-app"
	githubSponsorURL    = "https://github.com/sponsors/maxence-charriere"
	twitterURL          = "https://twitter.com/jonhymaxoo"
	coinbaseBusinessURL = "https://commerce.coinbase.com/checkout/851320a4-35b5-41f1-897b-74dd5ee207ae"
)

func setupPWA() {
	ui.BaseHPadding = 42
	ui.BlockPadding = 18
	analytics.Add(analytics.NewGoogleAnalytics())

	app.Route("/", newHomePage())
	app.Route("/getting-started", newGettingStartedPage())
	app.Route("/architecture", newArchitecturePage())
	app.Route("/reference", newReferencePage())

	app.Route("/components", newComponentsPage())
	app.Route("/declarative-syntax", newDeclarativeSyntaxPage())
	app.Route("/routing", newRoutingPage())
	app.Route("/static-resources", newStaticResourcePage())
	app.Route("/js", newJSPage())
	app.Route("/concurrency", newConcurrencyPage())
	app.Route("/seo", newSEOPage())
	app.Route("/lifecycle", newLifecyclePage())
	app.Route("/install", newInstallPage())
	app.Route("/testing", newTestingPage())
	app.Route("/actions", newActionPage())
	app.Route("/states", newStatesPage())

	app.Route("/migrate", newMigratePage())
	app.Route("/github-deploy", newGithubDeployPage())

	app.Route("/privacy-policy", newPrivacyPolicyPage())

	app.Handle(installApp, handleAppInstall)
	app.Handle(updateApp, handleAppUpdate)
	app.Handle(getMarkdown, handleGetMarkdown)
	app.Handle(getReference, handleGetReference)
}
