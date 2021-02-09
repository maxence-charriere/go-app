package app

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/maxence-charriere/go-app/v7/pkg/errors"
)

const (
	flowItemBaseWidth = 300
)

// UIFlow is the interface that describes a container that displays its items as
// a flow.
//
// EXPERIMENTAL WIDGET.
type UIFlow interface {
	UI

	// The HTML Class.
	Class(c string) UIFlow

	// Content sets the content with the given UI elements.
	Content(elems ...UI) UIFlow

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
	}
}

type flow struct {
	Compo

	IitemsBaseWitdh    int
	Iclass             string
	Icontent           []UI
	IstrechOnSingleRow bool

	id         string
	width      int
	itemWidth  int
	contentLen int
}

func (f *flow) Class(c string) UIFlow {
	if f.Iclass != "" {
		f.Iclass += " "
	}

	f.Iclass += c
	return f
}

func (f *flow) Content(elems ...UI) UIFlow {
	f.Icontent = FilterUIElems(elems...)
	return f
}

func (f *flow) ItemsBaseWidth(px int) UIFlow {
	f.IitemsBaseWitdh = px
	return f
}

func (f *flow) StrechtOnSingleRow() UIFlow {
	f.IstrechOnSingleRow = true
	return f
}

func (f *flow) OnMount(ctx Context) {
	f.id = "app-flow-" + uuid.New().String()

	f.Update()
	f.refreshLayout()
}

func (f *flow) OnNav(ctx Context) {
	f.refreshLayout()
}

func (f *flow) OnAppResize(ctx Context) {
	f.refreshLayout()
}

func (f *flow) Render() UI {
	if contentLen := len(f.Icontent); f.Mounted() && contentLen != f.contentLen {
		f.contentLen = contentLen
		f.refreshLayout()
	}

	return Div().
		ID(f.id).
		Class("goapp-flow").
		Class(f.Iclass).
		Body(
			Range(f.Icontent).Slice(func(i int) UI {
				item := f.Icontent[i]
				baseWidth := strconv.Itoa(f.itemWidth) + "px"

				return Div().
					Class("goapp-flow-item").
					Style("flex-basis", baseWidth).
					Body(item)
			}),
		)
}

func (f *flow) mounted() bool {
	return f.id != ""
}

func (f *flow) refreshLayout() {
	f.Dispatcher().Dispatch(f.adjustItemSizes)
}

func (f *flow) adjustItemSizes() {
	if !f.mounted() || f.IitemsBaseWitdh == 0 || len(f.Icontent) == 0 {
		return
	}

	elem := Window().GetElementByID(f.id)
	if !elem.Truthy() {
		Log("%s", errors.New("flow not found").Tag("id", f.id))
		return
	}

	width := elem.Get("clientWidth").Int()
	if width == 0 {
		return
	}

	f.width = width
	defer f.Update()

	itemWidth := f.IitemsBaseWitdh
	itemsPerRow := width / itemWidth
	if itemsPerRow <= 1 {
		f.itemWidth = width
		return
	}

	itemWidth = width / itemsPerRow
	if l := len(f.Icontent); l < itemsPerRow && f.IstrechOnSingleRow {
		itemWidth = width / l
	}
	f.itemWidth = itemWidth
}
