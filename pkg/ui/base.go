package ui

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// IBase is the interface that describes a component that serves as a base for a
// content.
type IBase interface {
	app.UI

	// Sets the ID.
	ID(v string) IBase

	// Sets the class. Multiple classes can be defined by successive calls.
	Class(v string) IBase

	// The content.
	Content(v ...app.UI) IBase
}

// Base creates a base for content.
func Base() IBase {
	return &base{
		hpadding: BaseHPadding,
		vpadding: 12,
	}
}

type base struct {
	app.Compo

	Iid      string
	Iclass   string
	Icontent []app.UI

	hpadding int
	vpadding int
	width    int
}

func (b *base) ID(v string) IBase {
	b.Iid = v
	return b
}

func (b *base) Class(v string) IBase {
	b.Iclass = app.AppendClass(b.Iclass, v)
	return b
}

func (b *base) Content(v ...app.UI) IBase {
	b.Icontent = app.FilterUIElems(v...)
	return b
}

func (b *base) OnMount(ctx app.Context) {
	b.resize(ctx)
}

func (b *base) OnResize(ctx app.Context) {
	b.resize(ctx)
}

func (b *base) OnUpdate(ctx app.Context) {
	b.resize(ctx)
}

func (b *base) Render() app.UI {
	return app.Div().
		ID(b.Iid).
		Class(b.Iclass).
		Body(
			app.Div().
				Style("position", "relative").
				Style("top", "0").
				Style("left", "0").
				Style("height", fmt.Sprintf("calc(100%s - %vpx)", "%", b.vpadding*2)).
				Style("width", fmt.Sprintf("calc(100%s - %vpx)", "%", b.hpadding*2)).
				Style("padding", fmt.Sprintf("%vpx %vpx", b.vpadding, b.hpadding)).
				Style("overflow", "hidden").
				Body(b.Icontent...),
		)
}

func (b *base) resize(ctx app.Context) {
	w, _ := ctx.Page().Size()
	if w <= 480 {
		b.hpadding = BaseMobileHPadding
	} else {
		b.hpadding = BaseHPadding
	}

	if w != b.width {
		b.width = w
	}
}
