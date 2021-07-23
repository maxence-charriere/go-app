package ui

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// IScroll is the interface that describes a base with a scrollable content
// surrounded by a fixed header and footer.
type IScroll interface {
	app.UI

	// Sets the ID.
	ID(v string) IScroll

	// Sets the class. Multiple classes can be defined by successive calls.
	Class(v string) IScroll

	// Sets the header height in px. Default is 90px.
	HeaderHeight(px int) IScroll

	// Sets the header.
	Header(v ...app.UI) IScroll

	// Sets the content.
	Content(v ...app.UI) IScroll

	// Sets the footer height in px. Default is 0.
	FooterHeight(px int) IScroll

	// Sets the footer.
	Footer(v ...app.UI) IScroll
}

// Scroll creates base with a scrollable content surrounded by a fixed header
// and footer.
func Scroll() IScroll {
	return &scroll{
		IheaderHeight: defaultHeaderHeight,
		hpadding:      BaseHPadding,
		vpadding:      BaseVPadding,
	}
}

type scroll struct {
	app.Compo

	Iid           string
	Iclass        string
	IheaderHeight int
	IfooterHeight int
	Iheader       []app.UI
	Icontent      []app.UI
	Ifooter       []app.UI

	hpadding int
	vpadding int
	width    int
}

func (s *scroll) ID(v string) IScroll {
	s.Iid = v
	return s
}

func (s *scroll) Class(v string) IScroll {
	if v == "" {
		return s
	}
	if s.Iclass != "" {
		s.Iclass += " "
	}
	s.Iclass += v
	return s
}

func (s *scroll) HeaderHeight(px int) IScroll {
	s.IheaderHeight = px
	return s
}

func (s *scroll) Header(v ...app.UI) IScroll {
	s.Iheader = app.FilterUIElems(v...)
	return s
}

func (s *scroll) Content(v ...app.UI) IScroll {
	s.Icontent = app.FilterUIElems(v...)
	return s
}

func (s *scroll) FooterHeight(px int) IScroll {
	s.IfooterHeight = px
	return s
}

func (s *scroll) Footer(v ...app.UI) IScroll {
	s.Ifooter = app.FilterUIElems(v...)
	return s
}

func (s *scroll) OnMount(ctx app.Context) {
	s.resize(ctx)
}

func (s *scroll) OnResize(ctx app.Context) {
	s.resize(ctx)
}

func (s *scroll) OnUpdate(ctx app.Context) {
	s.resize(ctx)
}

func (s *scroll) Render() app.UI {
	return app.Div().
		DataSet("goapp-ui", "scroll").
		ID(s.Iid).
		Class(s.Iclass).
		Body(
			app.Div().
				Style("width", "100%").
				Style("height", fmt.Sprintf("calc(100%s - %vpx)", "%", s.vpadding*2)).
				Style("padding", fmt.Sprintf("%vpx 0", s.vpadding)).
				Body(
					app.Div().
						Style("width", fmt.Sprintf("calc(100%s - %vpx)", "%", s.hpadding*2)).
						Style("padding", fmt.Sprintf("0 %vpx", s.hpadding)).
						Style("height", pxToString(s.IheaderHeight)).
						Body(s.Iheader...),
					app.Div().
						Style("width", fmt.Sprintf("calc(100%s - %vpx)", "%", s.hpadding*2)).
						Style("height", fmt.Sprintf("calc(100%s - %vpx)", "%", s.IheaderHeight+s.IfooterHeight)).
						Style("padding", fmt.Sprintf("0 %vpx", s.hpadding)).
						Style("overflow-x", "hidden").
						Style("overflow-y", "scroll").
						Body(s.Icontent...),
					app.Div().
						Style("width", fmt.Sprintf("calc(100%s - %vpx)", "%", s.hpadding*2)).
						Style("padding", fmt.Sprintf("0 %vpx", s.hpadding)).
						Style("height", pxToString(s.IfooterHeight)).
						Body(s.Ifooter...),
				),
		)
}

func (s *scroll) resize(ctx app.Context) {
	w, _ := ctx.Page().Size()
	if w <= 480 {
		s.hpadding = BaseMobileHPadding
	} else {
		s.hpadding = BaseHPadding
	}

	if w != s.width {
		s.width = w
	}
}
