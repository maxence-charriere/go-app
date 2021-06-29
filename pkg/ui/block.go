package ui

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// Block is the interface that describes a block of content.
type IBlock interface {
	app.UI

	// Sets the ID.
	ID(v string) IBlock

	// Sets the class. Multiple classes can be defined by successive calls.
	Class(v string) IBlock

	// Aligns content to the top.
	Top() IBlock

	// Aligns content to the middle.
	Middle() IBlock

	// The maximum content width. Default is 540px.
	MaxContentWidth(px int) IBlock

	// Sets the content.
	Content(v ...app.UI) IBlock
}

// Block creates a block of content.
func Block() IBlock {
	return &block{
		Ialignment:       stretch,
		ImaxContentWidth: 540,
		padding:          BlockPadding,
	}
}

type block struct {
	app.Compo

	Iid              string
	Iclass           string
	Ialignment       alignment
	ImaxContentWidth int
	Icontent         []app.UI

	padding int
	width   int
}

func (b *block) ID(v string) IBlock {
	b.Iid = v
	return b
}

func (b *block) Class(v string) IBlock {
	if v == "" {
		return b
	}
	if b.Iclass != "" {
		b.Iclass += " "
	}
	b.Iclass += v
	return b
}

func (b *block) Top() IBlock {
	b.Ialignment = top
	return b
}

func (b *block) Middle() IBlock {
	b.Ialignment = middle
	return b
}

func (b *block) MaxContentWidth(px int) IBlock {
	b.ImaxContentWidth = px
	return b
}

func (b *block) Content(v ...app.UI) IBlock {
	b.Icontent = app.FilterUIElems(v...)
	return b
}

func (b *block) OnMount(ctx app.Context) {
	b.resize(ctx)
}

func (b *block) OnResize(ctx app.Context) {
	b.resize(ctx)
}

func (b *block) OnUpdate(ctx app.Context) {
	b.resize(ctx)
}

func (b *block) Render() app.UI {
	layout := Stack().
		Style("width", "100%").
		Style("height", "100%").
		Center().
		Content(
			app.Div().
				Style("padding", pxToString(b.padding)).
				Style("width", fmt.Sprintf("calc(100%s - %vpx)", "%", b.padding*2)).
				Style("max-width", pxToString(b.ImaxContentWidth)).
				Body(b.Icontent...),
		)

	switch b.Ialignment {
	case stretch:
		layout.Stretch()

	case top:
		layout.Top()

	case middle:
		layout.Middle()
	}

	return app.Div().
		DataSet("goapp-ui", "block").
		ID(b.Iid).
		Class(b.Iclass).
		Body(layout)
}

func (b *block) resize(ctx app.Context) {
	w, _ := ctx.Page().Size()
	if w <= 480 {
		b.padding = BlockMobilePadding
	} else {
		b.padding = BlockPadding
	}

	if w != b.width {
		b.width = w
		b.ResizeContent()
	}
}
