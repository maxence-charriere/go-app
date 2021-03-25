package app

import (
	"time"

	"github.com/google/uuid"
	"github.com/maxence-charriere/go-app/v8/pkg/errors"
)

const (
	shellMenuDefaultWidth    = 300
	shellContentDefaultWidth = 480
	shellAdjustLayoutDelay   = time.Millisecond * 50
)

// UIShell is a component that responsively handles the disposition of a side
// menu, a submenu, and a main content.
//
// EXPERIMENTAL WIDGET.
type UIShell interface {
	UI

	// AlignItemsToCenter vertically aligns menus and content to center.
	AlignItemsToCenter() UIShell

	// Class adds a CSS class to the shell root HTML element class property.
	Class(c string) UIShell

	// Content sets the main content.
	Content(elems ...UI) UIShell

	// ID sets the shell root HTML element id property.
	ID(v string) UIShell

	// Menu sets the side menu.
	Menu(elems ...UI) UIShell

	// MenuButton sets the content of the button that displays an overlayed menu
	// when clicked. The button is displayed only when a menu is set and shrunk.
	// Default is ☰.
	MenuButton(elems ...UI) UIShell

	// MenuWidth set the base width for the menu and submenu in px. Default is
	// 300px.
	MenuWidth(px int) UIShell

	// OverlayMenu sets the content of the overlay menu. The overlay menu is
	// shown when the menu is shrunk and the menu button is clicked.
	OverlayMenu(elems ...UI) UIShell

	// Submenu sets the second side menu.
	Submenu(elems ...UI) UIShell
}

// Shell creates a responsive layout that handles the disposition of a side
// menu, a submenu, and a main content.
//
// EXPERIMENTAL WIDGET.
func Shell() UIShell {
	return &shell{
		ImenusWidth: shellMenuDefaultWidth,
		Ialignment:  "stretch",
		ImenuButton: []UI{
			Div().
				Class("goapp-shell-menu-button-default").
				Text("☰"),
		},
		id:               "goapp-shell-" + uuid.New().String(),
		isMenuVisible:    true,
		isSubmenuVisible: true,
	}
}

type shell struct {
	Compo

	Icontent     []UI
	Imenu        []UI
	Isubmenu     []UI
	ImenuButton  []UI
	IoverlayMenu []UI
	ImenusWidth  int
	Ialignment   string
	Iclass       string
	Iid          string

	id                         string
	isMenuVisible              bool
	isSubmenuVisible           bool
	isOverlayMenuButtonVisible bool
	isOverlayMenuVisible       bool
	adjustLayoutTimer          *time.Timer
}

func (s *shell) AlignItemsToCenter() UIShell {
	s.Ialignment = "center"
	return s
}

func (s *shell) Class(c string) UIShell {
	if c == "" {
		return s
	}
	if s.Iclass != "" {
		s.Iclass += " "
	}
	s.Iclass += c
	return s
}

func (s *shell) Content(elems ...UI) UIShell {
	s.Icontent = FilterUIElems(elems...)
	return s
}

func (s *shell) ID(v string) UIShell {
	s.Iid = v
	return s
}

func (s *shell) Menu(elems ...UI) UIShell {
	s.Imenu = FilterUIElems(elems...)
	return s
}

func (s *shell) MenuButton(elems ...UI) UIShell {
	s.ImenuButton = elems
	return s
}

func (s *shell) MenuWidth(px int) UIShell {
	if px > 0 {
		s.ImenusWidth = px
	}
	return s
}

func (s *shell) OverlayMenu(elems ...UI) UIShell {
	s.IoverlayMenu = elems
	return s
}

func (s *shell) Submenu(elems ...UI) UIShell {
	s.Isubmenu = FilterUIElems(elems...)
	return s
}

func (s *shell) OnMount(ctx Context) {
	s.refreshLayout(ctx)
}

func (s *shell) OnNav(ctx Context) {
	s.refreshLayout(ctx)
}

func (s *shell) OnResize(ctx Context) {
	s.refreshLayout(ctx)
}

func (s *shell) Dismount() {
	s.cancelAdjustLayout()
}

func (s *shell) Render() UI {
	if s.requiresLayoutUpdate() {
		s.Defer(s.refreshLayout)
	}

	visible := func(b bool) string {
		if b {
			return "bloc"
		}
		return "none"
	}

	menuWidth := pxToString(s.ImenusWidth)

	return Div().
		ID(s.id).
		Class("goapp-shell").
		Class(s.Iclass).
		Body(
			Div().
				Class("goapp-shell-layout").
				Style("align-items", s.Ialignment).
				Body(
					Div().
						Class("goapp-shell-item").
						Style("display", visible(s.isMenuVisible)).
						Style("width", menuWidth).
						Style("max-width", menuWidth).
						Body(s.Imenu...),
					Div().
						Class("goapp-shell-item").
						Style("display", visible(s.isSubmenuVisible)).
						Style("width", menuWidth).
						Style("max-width", menuWidth).
						Body(s.Isubmenu...),
					Div().
						Class("goapp-shell-item").
						Style("flex-basis", pxToString(shellContentDefaultWidth)).
						Style("flex-grow", "1").
						Body(s.Icontent...),
				),
			If(s.isOverlayMenuButtonVisible,
				Button().
					Class("goapp-shell-menu-button").
					Style("display", visible(!s.isOverlayMenuVisible)).
					OnClick(s.onMenuButtonClick).
					Body(s.ImenuButton...),
				Div().
					Class("goapp-shell-overlay-menu").
					Style("display", visible(s.isOverlayMenuVisible)).
					OnClick(s.onMenuOverlayClick).
					Body(s.IoverlayMenu...),
			),
		)
}

func (s *shell) hasMenu() bool {
	return len(s.Imenu) != 0
}

func (s *shell) hasSubmenu() bool {
	return len(s.Isubmenu) != 0
}

func (s *shell) hasOverlayMenu() bool {
	return len(s.IoverlayMenu) != 0
}

func (s *shell) requiresLayoutUpdate() bool {
	return s.Iid != "" && s.Iid != s.id
}

func (s *shell) refreshLayout(ctx Context) {
	if s.Iid != "" && s.Iid != s.id {
		s.id = s.Iid
		s.Update()
		return
	}

	if IsServer {
		return
	}

	s.cancelAdjustLayout()
	if s.adjustLayoutTimer != nil {
		s.adjustLayoutTimer.Reset(shellAdjustLayoutDelay)
		return
	}

	s.adjustLayoutTimer = time.AfterFunc(0, func() {
		s.Defer(s.adjustLayout)
	})
}

func (s *shell) adjustLayout(ctx Context) {
	root := Window().GetElementByID(s.id)
	if !root.Truthy() {
		Log(errors.New("shell root element not found").Tag("id", s.id))
		return
	}

	width := root.Get("clientWidth").Int()
	if width == 0 {
		return
	}

	s.isSubmenuVisible = len(s.Isubmenu) != 0 && shellContentDefaultWidth+s.ImenusWidth <= width

	if s.isSubmenuVisible {
		s.isMenuVisible = len(s.Imenu) != 0 && shellContentDefaultWidth+2*s.ImenusWidth <= width
	} else {
		s.isMenuVisible = len(s.Imenu) != 0 && shellContentDefaultWidth+s.ImenusWidth <= width

	}

	s.isOverlayMenuButtonVisible = len(s.IoverlayMenu) != 0 && !s.isMenuVisible
	s.Update()
	s.ResizeContent()
}

func (s *shell) cancelAdjustLayout() {
	if s.adjustLayoutTimer != nil {
		s.adjustLayoutTimer.Stop()
	}
}

func (s *shell) onMenuButtonClick(ctx Context, e Event) {
	s.isOverlayMenuVisible = true
	s.Update()
}

func (s *shell) onMenuOverlayClick(ctx Context, e Event) {
	s.isOverlayMenuVisible = false
	s.Update()
}
