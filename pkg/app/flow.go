package app

import (
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

const (
	flowDefaultItemsWidth = 300
	flowResizeSizeDelay   = time.Millisecond * 100
)

// UIFlow is the interface that describes a container that displays its items as
// a flow.
//
// EXPERIMENTAL WIDGET.
type UIFlow interface {
	UI

	// Class adds a CSS class to the flow root HTML element class property.
	Class(v string) UIFlow

	// Content sets the content with the given UI elements.
	Content(elems ...UI) UIFlow

	// ID sets the flow root HTML element id property.
	ID(v string) UIFlow

	// ItemsWidth sets the items base width in px. Items size is adjusted to
	// fit the space in the container. Default is 300px.
	ItemsWidth(px int) UIFlow

	// StrechtOnSingleRow makes the items to occupy all the available space when
	// the flow spreads on a single row.
	StrechtOnSingleRow() UIFlow
}

// Flow creates a container that displays its items as a flow.
//
// EXPERIMENTAL WIDGET.
func Flow() UIFlow {
	return &flow{
		IitemsWidth: flowDefaultItemsWidth,
		id:          "goapp-flow-" + uuid.New().String(),
	}
}

type flow struct {
	Compo

	IitemsWidth        int
	Iclass             string
	Iid                string
	Icontent           []UI
	IstrechOnSingleRow bool

	id              string
	contentLen      int
	width           int
	itemWidth       int
	adjustSizeTimer *time.Timer
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

func (f *flow) Content(elems ...UI) UIFlow {
	f.Icontent = FilterUIElems(elems...)
	return f
}

func (f *flow) ID(v string) UIFlow {
	f.Iid = v
	return f
}

func (f *flow) ItemsWidth(px int) UIFlow {
	if px > 0 {
		f.IitemsWidth = px
	}
	return f
}

func (f *flow) StrechtOnSingleRow() UIFlow {
	f.IstrechOnSingleRow = true
	return f
}

func (f *flow) OnMount(ctx Context) {
	if f.requiresLayoutUpdate() {
		f.refreshLayout(ctx)
	}
}

func (f *flow) OnNav(ctx Context) {
	if f.requiresLayoutUpdate() {
		f.refreshLayout(ctx)
	}
}

func (f *flow) OnUpdate(ctx Context) {
	if f.requiresLayoutUpdate() {
		f.refreshLayout(ctx)
	}
}

func (f *flow) OnResize(ctx Context) {
	f.refreshLayout(ctx)
}

func (f *flow) OnDismount() {
	f.cancelAdjustItemSizes()
}

func (f *flow) Render() UI {
	return Div().
		DataSet("goapp", "Flow").
		ID(f.id).
		Class(f.Iclass).
		Style("display", "flex").
		Style("flex-direction", "row").
		Style("flex-wrap", "wrap").
		Style("align-items", "stretch").
		Body(
			Range(f.Icontent).Slice(func(i int) UI {
				itemWidth := strconv.Itoa(f.itemWidth) + "px"
				return Div().
					Style("flex-grow", "0").
					Style("flex-shrink", "1").
					Style("max-width", itemWidth).
					Style("width", itemWidth).
					Body(f.Icontent[i])
			}),
		)
}

func (f *flow) requiresLayoutUpdate() bool {
	return (f.Iid != "" && f.Iid != f.id) ||
		len(f.Icontent) != f.contentLen
}

func (f *flow) refreshLayout(ctx Context) {
	if f.Iid != "" && f.Iid != f.id {
		f.id = f.Iid
		return
	}

	f.contentLen = len(f.Icontent)

	if IsServer {
		return
	}

	f.cancelAdjustItemSizes()
	if f.adjustSizeTimer != nil {
		f.adjustSizeTimer.Reset(flowResizeSizeDelay)
		return
	}

	f.adjustSizeTimer = time.AfterFunc(flowResizeSizeDelay, func() {
		f.adjustItemSizes(ctx)
	})
}

func (f *flow) adjustItemSizes(ctx Context) {
	if f.IitemsWidth == 0 || len(f.Icontent) == 0 {
		return
	}

	elem := Window().GetElementByID(f.id)
	if !elem.Truthy() {
		Log(errors.New("flow root element found").Tag("id", f.id))
		return
	}

	width := elem.Get("clientWidth").Int()
	if width == 0 {
		return
	}

	defer f.ResizeContent()

	itemWidth := f.IitemsWidth
	itemsPerRow := width / itemWidth

	if itemsPerRow <= 1 {
		f.itemWidth = width
		return
	}

	itemWidth = width / itemsPerRow
	if l := len(f.Icontent); l <= itemsPerRow && f.IstrechOnSingleRow {
		itemWidth = width / l
	}
	f.itemWidth = itemWidth
}

func (f *flow) cancelAdjustItemSizes() {
	if f.adjustSizeTimer != nil {
		f.adjustSizeTimer.Stop()
	}
}
