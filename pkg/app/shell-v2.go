package app

import (
	"time"

	"github.com/google/uuid"
)

// UIShell is a layout that responsively displays a content with a menu pane,
// an index page, and a hamburger menu.
//
// EXPERIMENTAL - Subject to change.
type UIShell interface {
	UI

	// Sets the ID.
	ID(v string) UIShell

	// Sets the class. Multiple classes can be defined by successive calls.
	Class(v string) UIShell

	// Sets the width in px for the menu and index panes.
	PaneWidth(px int) UIShell

	// Customizes the hamburger menu button with the given element.
	// Default is ☰.
	HamburgerButton(v UI) UIShell

	// Sets the hamburger menu content.
	HamburgerMenu(v ...UI) UIShell

	// Sets the menu pane content.
	Menu(v ...UI) UIShell

	// Sets the index pane content.
	Index(v ...UI) UIShell

	// Sets the content.
	Content(v ...UI) UIShell
}

// Shell returns a layout that responsively displays a content with a menu pane,
// an index pane, and a hamburger menu.
//
// EXPERIMENTAL - Subject to change.
func Shell() UIShell {
	return &shell{
		IpaneWidth:      270,
		id:              uuid.NewString(),
		refreshInterval: time.Millisecond * 50,
	}
}

type shell struct {
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

func (s *shell) ID(v string) UIShell {
	s.Iid = v
	return s
}

func (s *shell) Class(v string) UIShell {
	if v == "" {
		return s
	}
	if s.Iclass != "" {
		s.Iclass += " "
	}
	s.Iclass += v
	return s
}

func (s *shell) PaneWidth(px int) UIShell {
	s.IpaneWidth = px
	return s
}

func (s *shell) HamburgerButton(v UI) UIShell {
	b := FilterUIElems(v)
	if len(b) != 0 {
		s.IhamburgerButton = b[0]
	}
	return s
}

func (s *shell) HamburgerMenu(v ...UI) UIShell {
	s.IhamburgerMenu = FilterUIElems(v...)
	return s
}

func (s *shell) Menu(v ...UI) UIShell {
	s.Imenu = FilterUIElems(v...)
	return s
}

func (s *shell) Index(v ...UI) UIShell {
	s.Iindex = FilterUIElems(v...)
	return s
}

func (s *shell) Content(v ...UI) UIShell {
	s.Icontent = FilterUIElems(v...)
	return s
}

func (s *shell) OnPreRender(ctx Context) {
	s.refresh(ctx)
}

func (s *shell) OnMount(ctx Context) {
	s.refresh(ctx)
}

func (s *shell) OnResize(ctx Context) {
	s.scheduleRefresh(ctx)
}

func (s *shell) OnUpdate(ctx Context) {
	s.scheduleRefresh(ctx)
}

func (s *shell) Render() UI {
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
						Style("overflow", "hidden").
						Body(s.Imenu...),
					Div().
						Style("display", visible(!s.hideIndex)).
						Style("flex-shrink", "0").
						Style("flex-basis", pxToString(s.IpaneWidth)).
						Style("overflow", "hidden").
						Body(s.Iindex...),
					Div().
						Style("flex-grow", "1").
						Style("overflow", "hidden").
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

func (s *shell) scheduleRefresh(ctx Context) {
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

func (s *shell) refresh(ctx Context) {
	w, _ := s.layoutSize()

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

func (s *shell) layoutSize() (int, int) {
	layout := Window().GetElementByID("goapp-shell-" + s.id)
	if !layout.Truthy() {
		return 320, 568
	}
	return layout.Get("clientWidth").Int(), layout.Get("clientHeight").Int()
}

func (s *shell) onHamburgerButtonClick(ctx Context, e Event) {
	s.showHamburgerMenu = true
}

func (s *shell) hideHamburgerMenu(ctx Context, e Event) {
	s.showHamburgerMenu = false
}
