package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/maxence-charriere/go-app/v9/pkg/ui"
)

type referenceContent struct {
	app.Compo

	Iid    string
	Iclass string
	Iindex bool

	content         htmlContent
	currentFragment string
}

func newReferenceContent() *referenceContent {
	return &referenceContent{}
}

func (c *referenceContent) ID(v string) *referenceContent {
	c.Iid = v
	return c
}

func (c *referenceContent) Class(v string) *referenceContent {
	c.Iclass = app.AppendClass(c.Iclass, v)
	return c
}

func (c *referenceContent) Index(v bool) *referenceContent {
	c.Iindex = v
	return c
}

func (c *referenceContent) OnPreRender(ctx app.Context) {
	c.load(ctx)
}

func (c *referenceContent) OnMount(ctx app.Context) {
	c.load(ctx)
}

func (c *referenceContent) OnNav(ctx app.Context) {
	c.handleFragment(ctx)
}

func (c *referenceContent) load(ctx app.Context) {
	ctx.ObserveState(referenceState).
		OnChange(func() {
			ctx.Defer(c.handleFragment)
			ctx.Defer(c.scrollTo)
		}).
		Value(&c.content)

	ctx.NewAction(getReference)
}

func (c *referenceContent) Render() app.UI {
	loaderSize := 60
	loaderSpacing := 18
	if c.Iindex {
		loaderSize = 30
		loaderSpacing = 9
	}

	return app.Section().
		ID(c.Iid).
		Class(c.Iclass).
		Body(
			ui.Loader().
				Class("separator").
				Class("fill").
				Class("heading").
				Loading(c.content.Status == loading).
				Err(c.content.Err).
				Size(loaderSize).
				Spacing(loaderSpacing),

			app.If(!c.Iindex && c.content.Content != "",
				app.Raw(c.content.Content),
			).ElseIf(c.Iindex && c.content.Index != "",
				app.Raw(c.content.Index),
			),
			app.Div().Text(c.content.Err),
		)
}

func (c *referenceContent) handleFragment(ctx app.Context) {
	if !c.Iindex {
		return
	}
	if c.currentFragment != "" {
		c.unfocusCurrentIndex(ctx)
	}
	c.focusCurrentIndex(ctx)
}

func (c *referenceContent) unfocusCurrentIndex(ctx app.Context) {
	link := app.Window().GetElementByID(refLinkID(c.currentFragment))
	if !link.Truthy() {
		return
	}
	link.Set("className", "")
}

func (c *referenceContent) focusCurrentIndex(ctx app.Context) {
	fragment := ctx.Page().URL().Fragment
	link := app.Window().GetElementByID(refLinkID(fragment))
	if !link.Truthy() {
		return
	}
	link.Set("className", "focus")
	c.currentFragment = fragment
}

func (c *referenceContent) scrollTo(ctx app.Context) {
	id := ctx.Page().URL().Fragment
	if c.Iindex {
		id = refLinkID(id)
	}
	ctx.ScrollTo(id)
}
