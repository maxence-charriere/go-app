package app

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// UIFlow is the interface that describes a container that displays its items as
// a flow.
//
// EXPERIMENTAL - Subject to change.
type UIFlow interface {
	UI

	// Sets the ID.
	ID(v string) UIFlow

	// Sets the class. Multiple classes can be defined by successive calls.
	Class(v string) UIFlow

	// Sets the width in px for the content items.
	// Default is 300px.
	ItemWidth(px int) UIFlow

	// Sets the space between content elements in px.
	ItemSpacing(px int) UIFlow

	// Makes the items occupy all the available space when the content has only
	// one row.
	StretchItems() UIFlow

	// Sets the content.
	Content(elems ...UI) UIFlow
}

// Flow creates a container that displays its items as a flow.
//
// EXPERIMENTAL - Subject to change.
func Flow() UIFlow {
	return &flow{
		IitemWidth:      300,
		id:              "goapp-flow-" + uuid.NewString(),
		itemsPerRow:     1,
		refreshInterval: time.Millisecond * 50,
	}
}

type flow struct {
	Compo

	Iid           string
	Iclass        string
	IitemWidth    int
	IitemSpacing  int
	IstretchItems bool
	Icontent      []UI

	id              string
	itemsPerRow     int
	itemWidth       int
	refreshInterval time.Duration
	refreshTimer    *time.Timer
}

func (f *flow) ID(v string) UIFlow {
	f.Iid = v
	return f
}

func (f *flow) Class(v string) UIFlow {
	if v == "" {
		return f
	}
	if f.Iclass != "" {
		f.Iclass += " "
	}
	f.Iclass += v
	return f
}

func (f *flow) ItemWidth(px int) UIFlow {
	if px > 0 {
		f.IitemWidth = px
	}
	return f
}

func (f *flow) ItemSpacing(px int) UIFlow {
	if px > 0 {
		f.IitemSpacing = px
	}
	return f
}

func (f *flow) StretchItems() UIFlow {
	f.IstretchItems = true
	return f
}

func (f *flow) Content(elems ...UI) UIFlow {
	f.Icontent = FilterUIElems(elems...)
	return f
}

func (f *flow) OnPreRender(ctx Context) {
	f.refresh(ctx)
}

func (f *flow) OnMount(ctx Context) {
	f.refresh(ctx)
}

func (f *flow) OnResize(ctx Context) {
	f.scheduleRefresh(ctx)
}

func (f *flow) OnUpdate(ctx Context) {
	f.scheduleRefresh(ctx)
}

func (f *flow) OnDismount() {
	if f.refreshTimer != nil {
		f.refreshTimer.Stop()
	}
}

func (f *flow) Render() UI {
	return Div().
		DataSet("goapp-kit", "flow").
		ID(f.Iid).
		Class(f.Iclass).
		Body(
			Div().
				ID(f.id).
				Style("display", "flex").
				Style("flex-wrap", "wrap").
				Style("width", "100%").
				Style("overflow", "hidden").
				Body(
					Range(f.Icontent).Slice(func(i int) UI {
						marginTop := "0"
						if i >= f.itemsPerRow {
							marginTop = pxToString(f.IitemSpacing)
						}

						marginLeft := "0"
						if i%f.itemsPerRow != 0 {
							marginLeft = pxToString(f.IitemSpacing)
						}

						return Div().
							Style("position", "relative").
							Style("flex-shrink", "0").
							Style("flex-basis", pxToString(f.itemWidth)).
							Style("margin-top", marginTop).
							Style("margin-left", marginLeft).
							Style("overflow", "hidden").
							Body(f.Icontent[i])
					}),
				),
		)
}

func (f *flow) scheduleRefresh(ctx Context) {
	if f.refreshTimer != nil {
		f.refreshTimer.Stop()
		f.refreshTimer.Reset(f.refreshInterval)
		return
	}

	if IsClient {
		f.refreshTimer = time.AfterFunc(f.refreshInterval, func() {
			ctx.Dispatch(f.refresh)
		})
	}
}

func (f *flow) refresh(ctx Context) {
	w, _ := f.layoutSize()
	fmt.Println("layout width:", w)
	w += f.IitemSpacing
	fmt.Println("layout v-width:", w)

	itemWidth := f.IitemWidth + f.IitemSpacing
	itemsPerRow := w / itemWidth
	if f.IstretchItems && len(f.Icontent) < itemsPerRow {
		itemsPerRow = len(f.Icontent)
	}
	if itemsPerRow == 0 {
		itemsPerRow = 1
	}
	itemWidth = (w - f.IitemSpacing*itemsPerRow) / itemsPerRow

	if itemsPerRow != f.itemsPerRow || itemWidth != f.itemWidth {
		f.ResizeContent()
	}

	f.itemsPerRow = itemsPerRow
	f.itemWidth = itemWidth
}

func (f *flow) layoutSize() (int, int) {
	layout := Window().GetElementByID(f.id)
	if !layout.Truthy() {
		return 320 - 24 - 24, 568
	}
	return layout.Get("clientWidth").Int(), layout.Get("clientHeight").Int()
}
