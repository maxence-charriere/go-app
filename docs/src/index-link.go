package main

import (
	"strings"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type indexLink struct {
	app.Compo

	Iclass string
	Ititle string
	Ihref  string
}

func newIndexLink() *indexLink {
	return &indexLink{}
}

func (l *indexLink) Class(v string) *indexLink {
	l.Iclass = app.AppendClass(l.Iclass, v)
	return l
}

func (l *indexLink) Title(v string) *indexLink {
	l.Ititle = v
	return l
}

func (l *indexLink) Href(v string) *indexLink {
	l.Ihref = v
	return l
}

func (l *indexLink) OnNav(ctx app.Context) {}

func (l *indexLink) Render() app.UI {
	fragment := titleToFragment(l.Ititle)

	href := l.Ihref
	if href == "" {
		href = "#" + fragment
	}

	return app.A().
		Class("index-link").
		Class(l.Iclass).
		Class(fragmentFocus(fragment)).
		Href(href).
		Text(l.Ititle).
		Title(l.Ititle)
}

func titleToFragment(v string) string {
	v = strings.TrimSpace(v)
	v = strings.ToLower(v)
	v = strings.ReplaceAll(v, " ", "-")
	v = strings.ReplaceAll(v, ".", "-")
	v = strings.ReplaceAll(v, "?", "")
	return v
}
