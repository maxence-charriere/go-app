package ui

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// IScroll is the interface that describes scrollable content surrounded by a
// fixed header and footer.
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

// Scroll creates scrollable content surrounded by a fixed header and footer.
func Scroll() IScroll {
	return &scroll{
		IheaderHeight: 90,
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

func (s *scroll) Render() app.UI {
	return app.Div().
		DataSet("goapp-ui", "scroll").
		ID(s.Iid).
		Class(s.Iclass).
		Body(
			app.Div().
				Style("width", "100%").
				Style("height", "100%").
				Body(
					app.Div().
						Style("height", pxToString(s.IheaderHeight)).
						Body(s.Iheader...),
					app.Div().
						Style("height", fmt.Sprintf("calc(100%s - %vpx)", "%", s.IheaderHeight+s.IfooterHeight)).
						Body(s.Icontent...),
					app.Div().
						Style("height", pxToString(s.IfooterHeight)).
						Body(s.Ifooter...),
				),
		)
}
