package main

import (
	"github.com/murlokswarm/app"
)

// /!\ Import the component. Required to use a component.
func init() {
	app.Import(&DockMenu{})
}

// DockMenu is the component that represents the dock menu.
type DockMenu struct {
	switchIcon  bool
	switchBadge bool
}

// Render returns the HTML describing the dock menu content.
func (m *DockMenu) Render() string {
	return `
<menu>
	<menuitem label="Change icon" onclick="OnChangeIcon"></menuitem>
	<menuitem label="Change badge" onclick="OnChangeBadge"></menuitem>
</menu>
	`
}

// OnChangeIcon changes the dock icon when the dock menu item named
// "Change icon" is clicked.
func (m *DockMenu) OnChangeIcon() {
	m.switchIcon = !m.switchIcon

	var icon string
	if m.switchIcon {
		icon = app.Resources("like.png")
	} else {
		icon = app.Resources("logo.png")
	}

	dock, err := app.Dock()
	if err != nil {
		app.Log("OnChangeIcon: %s", err)
		return
	}
	if err = dock.SetIcon(icon); err != nil {
		app.Log("OnChangeIcon: %s", err)
	}
}

// OnChangeBadge changes the dock badge when the dock menu item named
// "Change badge" is clicked.
func (m *DockMenu) OnChangeBadge() {
	m.switchBadge = !m.switchBadge

	var badge string
	if m.switchBadge {
		badge = "hello"
	} else {
		badge = "world"
	}

	dock, err := app.Dock()
	if err != nil {
		app.Log("OnChangeBadge: %s", err)
		return
	}
	if err = dock.SetBadge(badge); err != nil {
		app.Log("OnChangeBadge: %s", err)
	}
}
