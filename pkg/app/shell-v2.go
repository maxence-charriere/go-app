package app

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// UIShell is a layout that responsively displays a content with a menu pane,
// an index page, and a hamburger menu.
//
// EXPERIMENTAL - Subject to change.
type UIShell2 interface {
	UI

	// Sets the ID.
	ID(v string) UIShell2

	// Sets the class. Multiple classes can be defined by successive calls.
	Class(v string) UIShell2

	// Sets the width in px for the menu and index panes.
	PaneWidth(px int) UIShell2

	// Customizes the hamburger menu button with the given element.
	// Default is ☰.
	HamburgerButton(v UI) UIShell2

	// Sets the hamburger menu content.
	HamburgerMenu(v ...UI) UIShell2

	// Sets the menu pane content.
	Menu(v ...UI) UIShell2

	// Sets the index pane content.
	Index(v ...UI) UIShell2

	// Sets the content.
	Content(v ...UI) UIShell2
}

// Shell returns a layout that responsively displays a content with a menu pane,
// an index pane, and a hamburger menu.
//
// EXPERIMENTAL - Subject to change.
func Shell2() UIShell2 {
	return &shell2{
		IpaneWidth:      270,
		id:              uuid.NewString(),
		refreshInterval: time.Millisecond * 50,
	}
}

type shell2 struct {
	Compo

	Iid              string
	Iclass           string
	IpaneWidth       int
	IhamburgerButton UI
	IhamburgerMenu   []UI
	Imenu            []UI
	Iindex           []UI
	Icontent         []UI

	id                string
	hideMenu          bool
	hideIndex         bool
	showHamburgerMenu bool
	refreshInterval   time.Duration
	refreshTimer      *time.Timer
}

func (s *shell2) ID(v string) UIShell2 {
	s.Iid = v
	return s
}

func (s *shell2) Class(v string) UIShell2 {
	if v == "" {
		return s
	}
	if s.Iclass != "" {
		s.Iclass += " "
	}
	s.Iclass += v
	return s
}

func (s *shell2) PaneWidth(px int) UIShell2 {
	s.IpaneWidth = px
	return s
}

func (s *shell2) HamburgerButton(v UI) UIShell2 {
	b := FilterUIElems(v)
	if len(b) != 0 {
		s.IhamburgerButton = b[0]
	}
	return s
}

func (s *shell2) HamburgerMenu(v ...UI) UIShell2 {
	s.IhamburgerMenu = FilterUIElems(v...)
	return s
}

func (s *shell2) Menu(v ...UI) UIShell2 {
	s.Imenu = FilterUIElems(v...)
	return s
}

func (s *shell2) Index(v ...UI) UIShell2 {
	s.Iindex = FilterUIElems(v...)
	return s
}

func (s *shell2) Content(v ...UI) UIShell2 {
	s.Icontent = FilterUIElems(v...)
	return s
}

func (s *shell2) OnPreRender(ctx Context) {
	s.refresh(ctx)
}

func (s *shell2) OnMount(ctx Context) {
	s.refresh(ctx)
}

func (s *shell2) OnResize(ctx Context) {
	s.scheduleRefresh(ctx)
}

func (s *shell2) OnUpdate(ctx Context) {
	s.scheduleRefresh(ctx)
}

func (s *shell2) Render() UI {
	visible := func(v bool) string {
		if v {
			return "block"
		}
		return "none"
	}

	return Div().
		DataSet("goapp", "shell").
		ID(s.Iid).
		Class(s.Iclass).
		Body(
			Div().
				ID("goapp-shell-"+s.id).
				Style("display", "flex").
				Style("width", "100%").
				Style("height", "100%").
				Style("overflow", "hidden").
				Body(
					Div().
						Style("display", visible(!s.hideMenu)).
						Style("flex-shrink", "0").
						Style("flex-basis", pxToString(s.IpaneWidth)).
						Body(s.Imenu...),
					Div().
						Style("display", visible(!s.hideIndex)).
						Style("flex-shrink", "0").
						Style("flex-basis", pxToString(s.IpaneWidth)).
						Body(s.Iindex...),
					Div().
						Style("flex-grow", "1").
						Body(s.Icontent...),
				),
			Div().
				Style("display", visible(s.hideMenu && len(s.IhamburgerMenu) != 0)).
				Style("position", "absolute").
				Style("top", "0").
				Style("left", "0").
				Style("cursor", "pointer").
				OnClick(s.onHamburgerButtonClick).
				Body(
					If(s.IhamburgerButton == nil,
						Div().
							Class("goapp-shell-hamburger-button-default").
							Text("☰"),
					),
				),
			Div().
				Style("display", visible(s.hideMenu && s.showHamburgerMenu)).
				Style("position", "absolute").
				Style("top", "0").
				Style("left", "0").
				Style("width", "100%").
				Style("height", "100%").
				Style("overflow", "hidden").
				OnClick(s.hideHamburgerMenu).
				Body(s.IhamburgerMenu...),
		)
}

func (s *shell2) scheduleRefresh(ctx Context) {
	if s.refreshTimer != nil {
		s.refreshTimer.Stop()
		s.refreshTimer.Reset(s.refreshInterval)
		return
	}

	if IsClient {
		s.refreshTimer = time.AfterFunc(s.refreshInterval, func() {
			ctx.Dispatch(s.refresh)
		})
	}
}

func (s *shell2) refresh(ctx Context) {
	w, h := s.layoutSize()
	fmt.Println("adust - layout:", w, h)

	hideIndex := len(s.Iindex) == 0 || 3*s.IpaneWidth > w
	hideMenu := false
	if hideIndex {
		hideMenu = len(s.Imenu) == 0 || 3*s.IpaneWidth > w
	} else {
		hideMenu = len(s.Imenu) == 0 || 5*s.IpaneWidth > w
	}

	if hideMenu != s.hideMenu || hideIndex != s.hideIndex {
		s.ResizeContent()
	}

	s.hideMenu = hideMenu
	s.hideIndex = hideIndex
}

func (s *shell2) layoutSize() (int, int) {
	layout := Window().GetElementByID("goapp-shell-" + s.id)
	if !layout.Truthy() {
		return 320, 568
	}
	return layout.Get("clientWidth").Int(), layout.Get("clientHeight").Int()
}

func (s *shell2) onHamburgerButtonClick(ctx Context, e Event) {
	s.showHamburgerMenu = true
}

func (s *shell2) hideHamburgerMenu(ctx Context, e Event) {
	s.showHamburgerMenu = false
}
