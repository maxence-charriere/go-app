package app

import (
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/maxence-charriere/go-app/v8/pkg/errors"
)

const (
	flowItemBaseWidth   = 300
	flowResizeSizeDelay = time.Millisecond * 100
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

	// ItemsBaseWidth sets the items base width in px. Items size is adjusted to
	// fit the space in the container. Default is 300px.
	ItemsBaseWidth(px int) UIFlow

	// StrechtOnSingleRow makes the items to occupy all the available space when
	// the flow spreads on a single row.
	StrechtOnSingleRow() UIFlow
}

// Flow creates a container that displays its items as a flow.
//
// EXPERIMENTAL WIDGET.
func Flow() UIFlow {
	return &flow{
		IitemsBaseWitdh: flowItemBaseWidth,
		id:              "goapp-flow-" + uuid.New().String(),
	}
}

type flow struct {
	Compo

	IitemsBaseWitdh    int
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

func (f *flow) ItemsBaseWidth(px int) UIFlow {
	if px > 0 {
		f.IitemsBaseWitdh = px
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

func (f *flow) OnResize(ctx Context) {
	f.refreshLayout(ctx)
}

func (f *flow) OnDismount() {
	f.cancelAdjustItemSizes()
}

func (f *flow) Render() UI {
	if f.requiresLayoutUpdate() {
		f.Defer(f.refreshLayout)
	}

	return Div().
		ID(f.id).
		Class("goapp-flow").
		Class(f.Iclass).
		Body(
			Range(f.Icontent).Slice(func(i int) UI {
				itemWidth := strconv.Itoa(f.itemWidth) + "px"

				return Div().
					Class("goapp-flow-item").
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
		f.Update()
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
		f.Defer(f.adjustItemSizes)
	})
}

func (f *flow) adjustItemSizes(ctx Context) {
	if f.IitemsBaseWitdh == 0 || len(f.Icontent) == 0 {
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
	defer f.Update()

	itemWidth := f.IitemsBaseWitdh
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
