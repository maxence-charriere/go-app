package main

import (
	"github.com/murlokswarm/app"
)

// Menu is a component that contains menu related examples.
type Menu struct {
	SupportsMenuBar bool

	SupportsDock   bool
	DockBadge      bool
	DockCustomIcon bool

	SupportsStatusMenu bool
	StatusText         bool
	StatusIcon         bool
	statusMenu         app.StatusMenu
}

// OnMount is the func called when the component is mounted.
func (m *Menu) OnMount() {
	m.SupportsMenuBar = app.MenuBar().Err() != app.ErrNotSupported
	m.SupportsDock = app.Dock().Err() != app.ErrNotSupported

	m.statusMenu = app.NewStatusMenu(app.StatusMenuConfig{
		Icon: app.Resources("logo.png"),
		URL:  "testmenu",
	})
	m.StatusIcon = true
	m.SupportsStatusMenu = m.statusMenu.Err() != app.ErrNotSupported

	app.Render(m)
}

// OnDismount is the func called when the component is dismounted.
func (m *Menu) OnDismount() {
	m.statusMenu.Close()
}

// Render returns a html string that describes the component.
func (m *Menu) Render() string {
	return `
<div class="Layout">
	<navpane current="menu">
	<div class="Menu-ContextMenu">
		<h1>Context menu</h1>

		<h2>Copy</h2>
		<p>Select and right click the text below.</p>
		<textarea class="Menu-CopyPasteArea" readonly oncontextmenu="OnContextMenu">
			Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.
		</textarea>

		<h2>Paste</h2>
		<p>Right click and paste in the area below.</p>
		<textarea class="Menu-CopyPasteArea" contenteditable oncontextmenu="OnContextMenu"></textarea>
	</div>
	<div class="Menu-Others">
		<h1>Other menus</h1>
		<div class="Menu-OthersContent">
			{{if .SupportsMenuBar}}
			<div class="Menu-OthersItem">
				<h2>Menu Bar</h2>
				<p>
					Take a look in the menu bar at the top of the screen and 
					click on the "Test menu" section to show the testing menu.
				</p>
			</div>
			{{end}}

			{{if .SupportsDock}}
			<div class="Menu-OthersItem">
				<h2>Dock</h2>
				<p>
					Right click on the app dock icon to show the testing menu. 
					Other docks actions are available below.
				</p>
				<ul>
					{{if .DockBadge}}
					<li><a onclick="ToggleDockBadge">Remove bagde</a></li>
					{{else}}
					<li><a onclick="ToggleDockBadge">Set badge</a></li>
					{{end}}

					{{if .DockCustomIcon}}
					<li><a onclick="ToggleDockCustomIcon">Remove custom icon</a></li>
					{{else}}
					<li><a onclick="ToggleDockCustomIcon">Set custom icon</a></li>
					{{end}}
				</ul>
			</div>
			{{end}}

			{{if .SupportsStatusMenu}}
			<div class="Menu-OthersItem">
				<h2>Status menu</h2>
				<p>
					Click on the app icon in the status menu to show the testing 
					menu. Other status menu actions are available below.
				</p>
				<ul>
					{{if .DockBadge}}
					<li><a onclick="ToggleStatusText">Remove text</a></li>
					{{else}}
					<li><a onclick="ToggleStatusText">Set text</a></li>
					{{end}}

					{{if .DockCustomIcon}}
					<li><a onclick="ToggleStatusIcon">Remove icon</a></li>
					{{else}}
					<li><a onclick="ToggleStatusIcon">Set icon</a></li>
					{{end}}
				</ul>
			</div>
			{{end}}
		</div>
	</div>
</div>
	`
}

// OnContextMenu is the function called to display a context menu.
func (m *Menu) OnContextMenu() {
	app.NewContextMenu("contextmenu")
}

// ToggleDockBadge is the function called the show/hide the dock badge.
func (m *Menu) ToggleDockBadge() {
	if m.DockBadge {
		app.Dock().SetBadge(nil)
	} else {
		app.Dock().SetBadge("hello")
	}

	m.DockBadge = !m.DockBadge
	app.Render(m)
}

// ToggleDockCustomIcon is the function called the show/hide the dock custom
// icon.
func (m *Menu) ToggleDockCustomIcon() {
	if m.DockCustomIcon {
		app.Dock().SetIcon("")
	} else {
		app.Dock().SetIcon(app.Resources("like.png"))
	}

	m.DockCustomIcon = true
	app.Render(m)
}

// ToggleStatusText is the function called the show/hide the status menu text.
func (m *Menu) ToggleStatusText() {
	if m.StatusText {
		m.statusMenu.SetText("")

		if !m.StatusIcon {
			m.ToggleStatusIcon()
		}
	} else {
		m.statusMenu.SetText("goapp-demo")
	}

	m.StatusText = !m.StatusText
	app.Render(m)
}

// ToggleStatusIcon is the function called the show/hide the status menu icon.
func (m *Menu) ToggleStatusIcon() {
	if m.StatusIcon {
		m.statusMenu.SetIcon("")

		if !m.StatusText {
			m.ToggleStatusText()
		}
	} else {
		m.statusMenu.SetIcon(app.Resources("logo.png"))
	}

	m.StatusIcon = !m.StatusIcon
	app.Render(m)
}

// ContextMenu is a component that describes a context menu.
type ContextMenu app.ZeroCompo

// Render returns a html string that describes the component.
func (m *ContextMenu) Render() string {
	return `
<menu>
	<menuitem label="Cut" keys="cmdorctrl+x" selector="cut:">
	<menuitem label="Copy" keys="cmdorctrl+c" selector="copy:">
	<menuitem label="Paste" keys="cmdorctrl+v" selector="paste:">
	<menuitem separator>
	<menuitem label="Select All" keys="cmdorctrl+a" selector="selectAll:">
</menu>
	`
}

// TestMenu is a component that describes a menu for testing.
type TestMenu struct {
	ShowSeparator bool
	ShowSubMenu   bool

	Disable string
}

// Render returns a html string that describes the component.
func (m *TestMenu) Render() string {
	return `
<menu label="Test menu">
	<menuitem label="Hello" onclick="OnHelloClick" {{.Disable}}>
	<menuitem label="Hello with Icon" 
			  icon="{{resources "logo.png"}}" 
			  onclick="OnHelloClick"
			  {{.Disable}}>
	<menuitem label="Hello with bad onclick" onclick="unknown" {{.Disable}}>
	<menuitem label="Hello without onclick" {{.Disable}}>
	<menuitem label="Disabled Hello" onclick="OnHelloClick" disabled>
	<menuitem separator>

	<menu label="Sub hello" {{.Disable}}>
		<menuitem label="World">

		{{if .ShowSeparator}}
		<menuitem separator>
		<menuitem label="Remove separator above" onclick="ToggleSeparator">
		{{else}}
		<menuitem label="Add separator below" onclick="ToggleSeparator">
		{{end}}

		<menuitem separator>
		{{if .ShowSubMenu}}
		<menu label="Sub menu">
			<menuitem label="Transform to item" onclick="ToggleSubMenu">
		</menu>
		{{else}}
		<menuitem label="Transform to menu" onclick="ToggleSubMenu">
		{{end}}
		
	</menu>
	<menuitem separator>

	{{if .Disable}}
	<menuitem label="Enable all" onclick="ToggleEnableAll">
	{{else}}
	<menuitem label="Disable all" onclick="ToggleEnableAll">
	{{end}}
</menu>
	`
}

// OnHelloClick is the function called when a hello button is clicked.
func (m *TestMenu) OnHelloClick() {
	app.Log("hello clicked")
}

// ToggleEnableAll is the function called to enable or disable all the menu
// buttons.
func (m *TestMenu) ToggleEnableAll() {
	if len(m.Disable) != 0 {
		m.Disable = ""
	} else {
		m.Disable = "disabled"
	}

	app.Render(m)
}

// ToggleSeparator is the function called to show or hide the Sub hello
// separator.
func (m *TestMenu) ToggleSeparator() {
	m.ShowSeparator = !m.ShowSeparator
	app.Render(m)
}

// ToggleSubMenu is the function called to perform item/submenu transforms.
func (m *TestMenu) ToggleSubMenu() {
	m.ShowSubMenu = !m.ShowSubMenu
	app.Render(m)
}
