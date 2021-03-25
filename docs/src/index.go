package main

import (
	"strings"

	"github.com/maxence-charriere/go-app/v8/pkg/app"
)

type index struct {
	app.Compo

	Ilinks []string
	Iclass string

	currentFragment string
}

func newIndex() *index {
	return &index{}
}

func (i *index) Links(v ...string) *index {
	i.Ilinks = v
	return i
}

func (i *index) Class(v string) *index {
	if v == "" {
		return i
	}
	if i.Iclass != "" {
		i.Iclass += " "
	}
	i.Iclass += v
	return i
}

func (i *index) OnNav(ctx app.Context) {
	i.currentFragment = ctx.Page.URL().Fragment
	i.Update()
}

func (i *index) Render() app.UI {
	return app.Div().
		Class(i.Iclass).
		Body(
			app.Range(i.Ilinks).Slice(func(j int) app.UI {
				link := i.Ilinks[j]
				githubLink := githubIndex(link)

				return newIndexLink().
					Label(link).
					Focus("#"+i.currentFragment == githubLink).
					Href(githubLink)
			}),
		)
}

type indexLink struct {
	app.Compo

	Ilabel string
	Ihref  string
	Ifocus bool
}

func newIndexLink() *indexLink {
	return &indexLink{}
}

func (i *indexLink) Label(v string) *indexLink {
	i.Ilabel = v
	return i
}

func (i *indexLink) Href(v string) *indexLink {
	i.Ihref = v
	return i
}

func (i *indexLink) Focus(v bool) *indexLink {
	i.Ifocus = v
	return i
}

func (i *indexLink) Render() app.UI {
	focus := ""
	if i.Ifocus {
		focus = "focus"
	}

	return app.A().
		Class("index-link").
		Class(focus).
		Href(i.Ihref).
		Text(i.Ilabel)
}

func githubIndex(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "?", "")
	s = strings.ReplaceAll(s, "(", "")
	s = strings.ReplaceAll(s, ")", "")
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, ".", "-")
	s = strings.ReplaceAll(s, "'", "-")
	return "#" + s
}
