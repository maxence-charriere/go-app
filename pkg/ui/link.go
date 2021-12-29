package ui

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// ILink is the interface that describes a clickable link.
type ILink interface {
	app.UI

	// Sets the ID.
	ID(v string) ILink

	// Sets the class. Multiple classes can be defined by successive calls.
	Class(v string) ILink

	// Sets the style. Multiple styles can be defined by successive calls.
	Style(k, v string) ILink

	// Sets the icon SVG code or its source path.
	Icon(v string) ILink

	// Sets the icon size.
	IconSize(px int) ILink

	// Sets the space between the icon and the label.
	IconSpace(px int) ILink

	// Sets the content vertical padding.
	Padding(v int) ILink

	// Sets the label.
	Label(v string) ILink

	// Set the tooltip.
	Help(v string) ILink

	// Sets the href to go when the link is clicked
	Href(v string) ILink

	// Sets the event handler called when the link is clicked.
	OnClick(v app.EventHandler) ILink
}

// Link create a clickable link.
func Link() ILink {
	return &link{
		IiconSize:  DefaultIconSize,
		IiconSpace: DefaultIconSpace,
		Ipadding:   3,
	}
}

type link struct {
	app.Compo

	Iid        string
	Iclass     string
	Istyles    []style
	Iicon      string
	IiconSize  int
	IiconSpace int
	Ipadding   int
	Ilabel     string
	Ihelp      string
	Ihref      string
	IonClick   app.EventHandler
}

func (l *link) ID(v string) ILink {
	l.Iid = v
	return l
}

func (l *link) Class(v string) ILink {
	l.Iclass = app.AppendClass(l.Iclass, v)
	return l
}

func (l *link) Style(k, v string) ILink {
	if v == "" {
		return l
	}
	l.Istyles = append(l.Istyles, style{
		key:   k,
		value: v,
	})
	return l
}

func (l *link) Icon(v string) ILink {
	l.Iicon = v
	return l
}

func (l *link) IconSize(px int) ILink {
	l.IiconSize = px
	return l
}

func (l *link) IconSpace(px int) ILink {
	l.IiconSpace = px
	return l
}

func (l *link) Padding(px int) ILink {
	l.Ipadding = px
	return l
}

func (l *link) Label(v string) ILink {
	l.Ilabel = v
	return l
}

func (l *link) Help(v string) ILink {
	l.Ihelp = v
	return l
}

func (l *link) Href(v string) ILink {
	l.Ihref = v
	return l
}

func (l *link) OnClick(v app.EventHandler) ILink {
	l.IonClick = v
	return l
}

func (l *link) Render() app.UI {
	link := app.A().
		ID(l.Iid).
		Class(l.Iclass).
		Title(l.Ihelp).
		Href(l.Ihref).
		OnClick(l.onClick).
		Body(
			Stack().
				Style("padding", fmt.Sprintf("%vpx 0", l.Ipadding)).
				Middle().
				Content(
					app.If(l.Iicon != "",
						Icon().
							Style("margin-right", pxToString(l.IiconSpace)).
							Size(l.IiconSize).
							Src(l.Iicon),
					),
					app.Div().Text(l.Ilabel),
				),
		)
	for _, s := range l.Istyles {
		link.Style(s.key, s.value)
	}
	return link
}

func (l *link) onClick(ctx app.Context, e app.Event) {
	if l.IonClick != nil {
		e.PreventDefault()
		l.IonClick(ctx, e)
	}
}
