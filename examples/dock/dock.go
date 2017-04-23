package main

import (
	"path/filepath"

	"github.com/murlokswarm/app"
)

// DockMenu is the component representing the dock menu.
type DockMenu struct {
	switchIcon  bool
	switchBadge bool
}

// Render returns the HTML describing the dock menu content.
func (m *DockMenu) Render() string {
	return `
<menu>
	<menuitem label="Change icon" onclick="OnChangeIcon" />
	<menuitem label="Change badge" onclick="OnChangeBadge" />
</menu>
	`
}

// OnChangeIcon changes the dock icon when the dock menu item named
// "Change icon" is clicked.
func (m *DockMenu) OnChangeIcon() {
	dock, _ := app.Dock()
	m.switchIcon = !m.switchIcon

	if m.switchIcon {
		icon := filepath.Join(app.Resources(), "like.png")
		dock.SetIcon(icon)
		return
	}
	icon := filepath.Join(app.Resources(), "logo.png")
	dock.SetIcon(icon)
}

// OnChangeBadge changes the dock badge when the dock menu item named
// "Change badge" is clicked.
func (m *DockMenu) OnChangeBadge() {
	dock, _ := app.Dock()
	m.switchBadge = !m.switchBadge

	if m.switchBadge {
		dock.SetBadge("Hello")
		return
	}
	dock.SetBadge("World!")
}

// /!\ Register the component. Required to use the component into a context.
func init() {
	app.RegisterComponent(&DockMenu{})
}
