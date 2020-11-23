package app

import (
	"time"

	"github.com/google/uuid"
	"github.com/maxence-charriere/go-app/v7/pkg/errors"
)

const (
	shellItemBaseWitdth = 300
	shellRefreshCooldow = time.Millisecond * 100
)

// UIShell is a component that responsively handles the disposition of a side
// menu, a submenu, and a main content.
//
// EXPERIMENTAL WIDGET.
type UIShell interface {
	UI

	// Class adds a CSS class to the layout.
	Class(c string) UIShell

	// Content sets the main content.
	Content(elems ...UI) UIShell

	// Menu sets the side menu.
	Menu(elems ...UI) UIShell

	// Submenu sets the second side menu.
	Submenu(elems ...UI) UIShell

	// MenuButton sets the content of the button that displays an overlayed menu
	// when clicked. The button is displayed only when a menu is set and shrunk.
	// Default is ☰.
	MenuButton(elems ...UI) UIShell

	// OverlayMenu sets the content of the overlay menu. The overlay menu is
	// shown when the menu is shrunk and the menu button is clicked.
	OverlayMenu(elems ...UI) UIShell

	// ItemsBaseWidth set the base width for the menu, submenu and content in
	// px. Default is 300px.
	ItemsBaseWidth(px int) UIShell

	// AlignItemsToCenter vertically aligns menus and content to center.
	AlignItemsToCenter() UIShell
}

// Shell creates a responsive layout that handles the disposition of a side
// menu, a submenu, and a main content.
//
// EXPERIMENTAL WIDGET.
func Shell() UIShell {
	return &shell{
		IitemsBaseWidth: shellItemBaseWitdth,
		Ialignment:      "stretch",
		ImenuButton: []UI{
			Div().
				Class("goapp-shell-menu-button-default").
				Text("☰"),
		},
	}
}

type shell struct {
	Compo

	Icontent          []UI
	Imenu             []UI
	Isubmenu          []UI
	ImenuButton       []UI
	IoverlayMenu      []UI
	IitemsBaseWidth   int
	Ialignment        string
	Iclass            string
	IshowShrunkenMenu bool

	id                  string
	closeResizeListener func()
	refreshCooldown     *time.Timer
	shrunkenMenu        bool
	shrunkenSubmenu     bool
}

func (s *shell) Class(c string) UIShell {
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

func (s *shell) Menu(elems ...UI) UIShell {
	s.Imenu = FilterUIElems(elems...)
	return s
}

func (s *shell) Submenu(elems ...UI) UIShell {
	s.Isubmenu = FilterUIElems(elems...)
	return s
}

func (s *shell) MenuButton(elems ...UI) UIShell {
	s.ImenuButton = elems
	return s
}

func (s *shell) OverlayMenu(elems ...UI) UIShell {
	s.IoverlayMenu = elems
	return s
}

func (s *shell) ItemsBaseWidth(px int) UIShell {
	if px <= 0 {
		px = shellItemBaseWitdth
	}

	s.IitemsBaseWidth = px
	return s
}

func (s *shell) AlignItemsToCenter() UIShell {
	s.Ialignment = "center"
	return s
}

func (s *shell) OnMount(ctx Context) {
	s.id = uuid.New().String()
	s.closeResizeListener = Window().AddEventListener("resize", s.onResize)

	s.Update()
	s.refreshLayout()
}

func (s *shell) OnDismount() {
	if s.refreshCooldown != nil {
		s.refreshCooldown.Stop()
	}

	s.closeResizeListener()
}

func (s *shell) Render() UI {
	showMenu := s.hasMenu() && !s.shrunkenMenu
	showSubmenu := s.hasSubmenu() && !s.shrunkenSubmenu

	visible := func(b bool) string {
		if b {
			return "bloc"
		}
		return "none"
	}

	layoutShrink := "0"
	if !showMenu && !showSubmenu {
		layoutShrink = "1"
	}

	return Div().
		Class("goapp-shell").
		Class(s.Iclass).
		Body(
			Div().
				ID(s.elemID("layout")).
				Class("goapp-shell-layout").
				Style("align-items", s.Ialignment).
				Body(
					Div().
						ID(s.elemID("menu")).
						Class("goapp-shell-item").
						Style("display", visible(showMenu)).
						Style("flex-basis", pxToString(s.IitemsBaseWidth)).
						Style("flex-shrink", "2").
						Body(
							If(showMenu, s.Imenu...),
						),
					Div().
						ID(s.elemID("submenu")).
						Class("goapp-shell-item").
						Style("display", visible(showSubmenu)).
						Style("flex-basis", pxToString(s.IitemsBaseWidth)).
						Style("flex-shrink", "1").
						Body(s.Isubmenu...),
					Div().
						Class("goapp-shell-item").
						Style("flex-basis", pxToString(s.contentItemBaseWidth())).
						Style("flex-grow", "1").
						Style("flex-shrink", layoutShrink).
						Body(s.Icontent...),
				),
			If(s.hasMenu() && s.hasOverlayMenu(),
				Button().
					ID(s.elemID("menu-button")).
					Class("goapp-shell-menu-button").
					Style("display", visible(!showMenu && !s.IshowShrunkenMenu)).
					OnClick(s.onMenuButtonClick).
					Body(s.ImenuButton...),
				Div().
					Class("goapp-shell-overlay-menu").
					Style("display", visible(!showMenu && s.IshowShrunkenMenu)).
					OnClick(s.onMenuOverlayClick).
					Body(s.IoverlayMenu...),
			),
		)
}

func (s *shell) elemID(elem string) string {
	if elem == "" {
		return ""
	}

	return "app-shell-" + elem + "-" + s.id
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

func (s *shell) minItemWidth() int {
	return s.IitemsBaseWidth * 70 / 100
}

func (s *shell) contentItemBaseWidth() int {
	return (s.IitemsBaseWidth * 70 / 100) + s.IitemsBaseWidth
}

func (s *shell) mounted() bool {
	return s.id != ""
}

func (s *shell) onResize(ctx Context, e Event) {
	s.refreshLayout()
}

func (s *shell) refreshLayout() {
	if s.refreshCooldown != nil {
		s.refreshCooldown.Reset(shellRefreshCooldow)
		return
	}

	s.refreshCooldown = time.AfterFunc(shellRefreshCooldow, func() {
		Dispatch(func() {
			s.refreshMenu()
			s.refreshSubmenu()
			s.Update()
		})
	})
}

func (s *shell) refreshMenu() {
	if !s.mounted() || !s.hasMenu() {
		return
	}

	currentWidth := s.minItemWidth() + s.contentItemBaseWidth()
	if len(s.Isubmenu) != 0 {
		currentWidth += s.IitemsBaseWidth
	}

	layoutID := s.elemID("layout")
	layout := Window().GetElementByID(layoutID)
	if !layout.Truthy() {
		Log("%s", errors.New("shell layout not found").Tag("id", layoutID))
		return
	}
	layoutWidth := layout.Get("clientWidth").Int()

	s.shrunkenMenu = currentWidth > layoutWidth
}

func (s *shell) refreshSubmenu() {
	if !s.mounted() || !s.hasSubmenu() {
		return
	}

	layoutID := s.elemID("layout")
	layout := Window().GetElementByID(layoutID)
	if !layout.Truthy() {
		Log("%s", errors.New("shell layout not found").Tag("id", layoutID))
		return
	}

	currentWidth := s.minItemWidth() + s.contentItemBaseWidth()
	layoutWidth := layout.Get("clientWidth").Int()
	s.shrunkenSubmenu = currentWidth > layoutWidth
}

func (s *shell) onMenuButtonClick(ctx Context, e Event) {
	s.IshowShrunkenMenu = true
	s.Update()
}

func (s *shell) onMenuOverlayClick(ctx Context, e Event) {
	s.IshowShrunkenMenu = false
	s.Update()
}
