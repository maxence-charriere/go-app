package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/analytics"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type githubDeployPage struct {
	app.Compo
}

func newGithubDeployPage() *githubDeployPage {
	return &githubDeployPage{}
}

func (p *githubDeployPage) OnPreRender(ctx app.Context) {
	p.initPage(ctx)
}

func (p *githubDeployPage) OnNav(ctx app.Context) {
	p.initPage(ctx)
}

func (p *githubDeployPage) initPage(ctx app.Context) {
	ctx.Page().SetTitle("Deploy a PWA on GitHub Pages")
	ctx.Page().SetDescription("Documentation about how to deploy a PWA created with go-app on GitHub Pages.")
	analytics.Page("github-deploy", nil)
}

func (p *githubDeployPage) Render() app.UI {
	return newPage().
		Title("Deploy on GitHub Pages").
		Icon(githubSVG).
		Index(
			newIndexLink().Title("Intro"),
			newIndexLink().Title("Generate a Static Website"),
			newIndexLink().Title("Deployment"),
			newIndexLink().Title("    Domainless Repository"),

			app.Div().Class("separator"),
		).
		Content(
			newRemoteMarkdownDoc().Src("/web/documents/github-deploy.md"),
		)
}
