package app

import (
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/maxence-charriere/go-app/v7/pkg/errors"
)

const (
	flowItemBaseWidth   = 300
	flowRefreshCooldown = time.Millisecond * 100
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

	id                  string
	width               int
	itemWidth           int
	closeResizeListener func()
	refreshCooldown     *time.Timer
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
	f.closeResizeListener = Window().AddEventListener("resize", f.onResize)

	f.Update()
	f.refreshLayout()
}

func (f *flow) OnDismount() {
	if f.refreshCooldown != nil {
		f.refreshCooldown.Stop()
	}

	f.closeResizeListener()
}

func (f *flow) Render() UI {
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

func (f *flow) onResize(ctx Context, e Event) {
	f.refreshLayout()
}

func (f *flow) refreshLayout() {
	if f.refreshCooldown != nil {
		f.refreshCooldown.Reset(flowRefreshCooldown)
		return
	}

	f.refreshCooldown = time.AfterFunc(flowRefreshCooldown, func() {
		Dispatch(f.adjustItemSizes)
	})
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

	itemsPerRow := width / f.IitemsBaseWitdh
	if itemsPerRow == 0 {
		f.itemWidth = width
		return
	}

	if l := len(f.Icontent); l < itemsPerRow && f.IstrechOnSingleRow {
		f.itemWidth = width / l
		return
	}

	remainingSpace := width - f.IitemsBaseWitdh*itemsPerRow
	spaceToAddPerItem := remainingSpace / itemsPerRow
	f.itemWidth = f.IitemsBaseWitdh + spaceToAddPerItem
}
