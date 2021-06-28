package ui

import (
	"time"

	"github.com/google/uuid"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// IShell is the interface that describes a layout that responsively displays a
// content with a menu pane, an index page, and a hamburger menu.
type IShell interface {
	app.UI

	// Sets the ID.
	ID(v string) IShell

	// Sets the class. Multiple classes can be defined by successive calls.
	Class(v string) IShell

	// Sets the width in px for the menu and index panes.
	// Default is 270px.
	PaneWidth(px int) IShell

	// Customizes the hamburger menu button with the given element.
	// Default is ☰.
	HamburgerButton(v app.UI) IShell

	// Sets the hamburger menu content.
	HamburgerMenu(v ...app.UI) IShell

	// Sets the menu pane content.
	Menu(v ...app.UI) IShell

	// Sets the index pane content.
	Index(v ...app.UI) IShell

	// Sets the content.
	Content(v ...app.UI) IShell
}

// Shell returns a layout that responsively displays a content with a menu pane,
// an index pane, and a hamburger menu.
func Shell() IShell {
	return &shell{
		IpaneWidth:      270,
		id:              "goapp-shell-" + uuid.NewString(),
		refreshInterval: time.Millisecond * 50,
	}
}

type shell struct {
	app.Compo

	Iid              string
	Iclass           string
	IpaneWidth       int
	IhamburgerButton app.UI
	IhamburgerMenu   []app.UI
	Imenu            []app.UI
	Iindex           []app.UI
	Icontent         []app.UI

	id                string
	hideMenu          bool
	hideIndex         bool
	showHamburgerMenu bool
	refreshInterval   time.Duration
	refreshTimer      *time.Timer
	width             int
}

func (s *shell) ID(v string) IShell {
	s.Iid = v
	return s
}

func (s *shell) Class(v string) IShell {
	if v == "" {
		return s
	}
	if s.Iclass != "" {
		s.Iclass += " "
	}
	s.Iclass += v
	return s
}

func (s *shell) PaneWidth(px int) IShell {
	if px > 0 {
		s.IpaneWidth = px
	}
	return s
}

func (s *shell) HamburgerButton(v app.UI) IShell {
	b := app.FilterUIElems(v)
	if len(b) != 0 {
		s.IhamburgerButton = b[0]
	}
	return s
}

func (s *shell) HamburgerMenu(v ...app.UI) IShell {
	s.IhamburgerMenu = app.FilterUIElems(v...)
	return s
}

func (s *shell) Menu(v ...app.UI) IShell {
	s.Imenu = app.FilterUIElems(v...)
	return s
}

func (s *shell) Index(v ...app.UI) IShell {
	s.Iindex = app.FilterUIElems(v...)
	return s
}

func (s *shell) Content(v ...app.UI) IShell {
	s.Icontent = app.FilterUIElems(v...)
	return s
}

func (s *shell) OnPreRender(ctx app.Context) {
	s.refresh(ctx)
}

func (s *shell) OnMount(ctx app.Context) {
	s.refresh(ctx)
}

func (s *shell) OnResize(ctx app.Context) {
	s.scheduleRefresh(ctx)
}

func (s *shell) OnUpdate(ctx app.Context) {
	s.scheduleRefresh(ctx)
}

func (s *shell) OnDismount() {
	if s.refreshTimer != nil {
		s.refreshTimer.Stop()
	}
}

func (s *shell) Render() app.UI {
	visible := func(v bool) string {
		if v {
			return "block"
		}
		return "none"
	}

	return app.Div().
		DataSet("goapp-ui", "shell").
		ID(s.Iid).
		Class(s.Iclass).
		Body(
			app.Div().
				ID(s.id).
				Style("display", "flex").
				Style("width", "100%").
				Style("height", "100%").
				Style("overflow", "hidden").
				Body(
					app.Div().
						Style("position", "relative").
						Style("display", visible(!s.hideMenu)).
						Style("flex-shrink", "0").
						Style("flex-basis", pxToString(s.IpaneWidth)).
						Style("overflow", "hidden").
						Body(s.Imenu...),
					app.Div().
						Style("position", "relative").
						Style("display", visible(!s.hideIndex)).
						Style("flex-shrink", "0").
						Style("flex-basis", pxToString(s.IpaneWidth)).
						Style("overflow", "hidden").
						Body(s.Iindex...),
					app.Div().
						Style("position", "relative").
						Style("flex-grow", "1").
						Style("overflow", "hidden").
						Body(s.Icontent...),
				),
			app.Div().
				Style("display", visible(s.hideMenu && len(s.IhamburgerMenu) != 0)).
				Style("position", "absolute").
				Style("top", "0").
				Style("left", "0").
				Style("cursor", "pointer").
				OnClick(s.onHamburgerButtonClick).
				Body(
					app.If(s.IhamburgerButton == nil,
						app.Div().
							Class("goapp-shell-hamburger-button-default").
							Text("☰"),
					),
				),
			app.Div().
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

func (s *shell) scheduleRefresh(ctx app.Context) {
	if s.refreshTimer != nil {
		s.refreshTimer.Stop()
		s.refreshTimer.Reset(s.refreshInterval)
		return
	}

	if app.IsClient {
		s.refreshTimer = time.AfterFunc(s.refreshInterval, func() {
			ctx.Dispatch(s.refresh)
		})
	}
}

func (s *shell) refresh(ctx app.Context) {
	w, _ := s.layoutSize()

	hideIndex := len(s.Iindex) == 0 || 3*s.IpaneWidth > w
	hideMenu := false
	if hideIndex {
		hideMenu = len(s.Imenu) == 0 || 3*s.IpaneWidth > w
	} else {
		hideMenu = len(s.Imenu) == 0 || 5*s.IpaneWidth > w
	}

	if hideMenu != s.hideMenu || hideIndex != s.hideIndex || w != s.width {
		s.hideMenu = hideMenu
		s.hideIndex = hideIndex
		s.width = w
		s.ResizeContent()
	}
}

func (s *shell) layoutSize() (int, int) {
	layout := app.Window().GetElementByID(s.id)
	if !layout.Truthy() {
		return 320, 568
	}
	return layout.Get("clientWidth").Int(), layout.Get("clientHeight").Int()
}

func (s *shell) onHamburgerButtonClick(ctx app.Context, e app.Event) {
	s.showHamburgerMenu = true
}

func (s *shell) hideHamburgerMenu(ctx app.Context, e app.Event) {
	s.showHamburgerMenu = false
}
