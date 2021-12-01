package ui

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// IFlow is the interface that describes a container that displays its items as
// a flow.
type IFlow interface {
	app.UI

	// Sets the ID.
	ID(v string) IFlow

	// Sets the class. Multiple classes can be defined by successive calls.
	Class(v string) IFlow

	// Sets the width in px for the content items.
	// Default is 300px.
	ItemWidth(px int) IFlow

	// Sets the space between content elements in px.
	Spacing(px int) IFlow

	// Makes the items occupy all the available space when the content has only
	// one row.
	StretchItems() IFlow

	// Sets the content.
	Content(elems ...app.UI) IFlow
}

// Flow creates a container that displays its items as a flow.
func Flow() IFlow {
	return &flow{
		IitemWidth:  300,
		id:          "goapp-flow-" + uuid.NewString(),
		itemsPerRow: 1,
	}
}

type flow struct {
	app.Compo

	Iid           string
	Iclass        string
	IitemWidth    int
	Ispacing      int
	IstretchItems bool
	Icontent      []app.UI

	id          string
	itemsPerRow int
	itemWidth   float64
}

func (f *flow) ID(v string) IFlow {
	f.Iid = v
	return f
}

func (f *flow) Class(v string) IFlow {
	f.Iclass = app.AppendClass(f.Iclass, v)
	return f
}

func (f *flow) ItemWidth(px int) IFlow {
	if px > 0 {
		f.IitemWidth = px
	}
	return f
}

func (f *flow) Spacing(px int) IFlow {
	if px > 0 {
		f.Ispacing = px
	}
	return f
}

func (f *flow) StretchItems() IFlow {
	f.IstretchItems = true
	return f
}

func (f *flow) Content(elems ...app.UI) IFlow {
	f.Icontent = app.FilterUIElems(elems...)
	return f
}

func (f *flow) OnPreRender(ctx app.Context) {
	f.refresh(ctx)
}

func (f *flow) OnMount(ctx app.Context) {
	f.refresh(ctx)
}

func (f *flow) OnResize(ctx app.Context) {
	f.refresh(ctx)
}

func (f *flow) OnUpdate(ctx app.Context) {
	f.refresh(ctx)
}

func (f *flow) Render() app.UI {
	return app.Div().
		DataSet("goapp-ui", "flow").
		ID(f.Iid).
		Class(f.Iclass).
		Body(
			app.Div().
				ID(f.id).
				Style("display", "flex").
				Style("flex-wrap", "wrap").
				Style("width", "100%").
				Style("max-width", "100%").
				Body(
					app.Range(f.Icontent).Slice(func(i int) app.UI {
						marginTop := "0"
						if i >= f.itemsPerRow {
							marginTop = pxToString(f.Ispacing)
						}

						marginLeft := "0"
						if i%f.itemsPerRow != 0 {
							marginLeft = pxToString(f.Ispacing)
						}

						width := fmt.Sprintf("%.6fpx", f.itemWidth)

						return app.Div().
							Style("position", "relative").
							Style("flex-shrink", "0").
							Style("flex-basis", width).
							Style("max-width", width).
							Style("margin-top", marginTop).
							Style("margin-left", marginLeft).
							Body(f.Icontent[i])
					}),
				),
		)
}

func (f *flow) refresh(ctx app.Context) {
	w, _ := f.layoutSize()
	w += f.Ispacing

	itemWidth := f.IitemWidth + f.Ispacing
	itemsPerRow := w / itemWidth
	if f.IstretchItems && len(f.Icontent) < itemsPerRow {
		itemsPerRow = len(f.Icontent)
	}
	if itemsPerRow == 0 {
		itemsPerRow = 1
	}
	itemWidthFloat := float64(w-f.Ispacing*itemsPerRow) / float64(itemsPerRow)

	if itemsPerRow != f.itemsPerRow || itemWidthFloat != f.itemWidth {
		f.itemsPerRow = itemsPerRow
		f.itemWidth = itemWidthFloat

		ctx.Defer(func(app.Context) {
			f.ResizeContent()
		})
	}
}

func (f *flow) layoutSize() (int, int) {
	layout := app.Window().GetElementByID(f.id)
	if !layout.Truthy() {
		return 320 - 24 - 24, 568
	}
	return layout.Get("clientWidth").Int(), layout.Get("clientHeight").Int()
}
