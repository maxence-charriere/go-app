package main

import (
	"github.com/murlokswarm/app"
)

// Menu is a component that contains menu related examples.
type Menu struct {
	SupportsMenuBar bool
	SupportsDock    bool
	DockBadge       bool
	DockCustomIcon  bool
}

// OnMount is the func called when the component is mounted.
func (m *Menu) OnMount() {
	m.SupportsMenuBar = app.MenuBar().Err() != app.ErrNotSupported
	m.SupportsDock = app.Dock().Err() != app.ErrNotSupported
	app.Render(m)
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
					Take a look in the menu bar at the top of the screen and click
					on the "Test menu" section to show the testing menu.
				</p>
			</div>
			{{end}}

			{{if .SupportsDock}}
			<div class="Menu-OthersItem">
				<h2>Dock</h2>
				<p>
					Right click on the app dock icon to show the testing menu. Other
					docks actions are available below.
				</p>
				<ul>
					{{if .DockBadge}}
					<li><a onclick="OnRemoveDockBadge">Remove bagde</a></li>
					{{else}}
					<li><a onclick="OnSetDockBadge">Set badge</a></li>
					{{end}}

					{{if .DockCustomIcon}}
					<li><a onclick="OnRemoveDockCustomIcon">Remove custom icon</a></li>
					{{else}}
					<li><a onclick="OnSetDockCustomIcon">Set custom icon</a></li>
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

// OnSetDockBadge is the function called when "Set badge" is clicked.
func (m *Menu) OnSetDockBadge() {
	app.Dock().SetBadge("hello")
	m.DockBadge = true
	app.Render(m)
}

// OnRemoveDockBadge is the function called when "Remove bagde" is clicked.
func (m *Menu) OnRemoveDockBadge() {
	app.Dock().SetBadge(nil)
	m.DockBadge = false
	app.Render(m)
}

// OnSetDockCustomIcon is the function called when "Set custom icon" is clicked.
func (m *Menu) OnSetDockCustomIcon() {
	app.Dock().SetIcon(app.Resources("like.png"))
	m.DockCustomIcon = true
	app.Render(m)
}

// OnRemoveDockCustomIcon is the function called when "Remove custom icon" is
// clicked.
func (m *Menu) OnRemoveDockCustomIcon() {
	app.Dock().SetIcon("")
	m.DockCustomIcon = false
	app.Render(m)
}

// ContextMenu is a component that describes a context menu.
type ContextMenu app.ZeroCompo

// Render returns a html string that describes the component.
func (m *ContextMenu) Render() string {
	return `
<menu>
	<menuitem label="Cut" keys="cmdorctrl+x" selector="cut:"></menuitem>
	<menuitem label="Copy" keys="cmdorctrl+c" selector="copy:"></menuitem>
	<menuitem label="Paste" keys="cmdorctrl+v" selector="paste:"></menuitem>
	<menuitem separator></menuitem>
	<menuitem label="Select All" keys="cmdorctrl+a" selector="selectAll:"></menuitem>
</menu>
	`
}

// TestMenu is a component that describes a menu for testing.
type TestMenu struct {
	ShowSeparator bool
	Disable       string
}

// Render returns a html string that describes the component.
func (m *TestMenu) Render() string {
	return `
<menu label="Test menu">
	<menuitem label="Hello" onclick="OnHelloClick" {{.Disable}}></menuitem>
	<menuitem label="Hello with Icon" 
			  icon="{{resources "logo.png"}}" 
			  onclick="OnHelloClick"
			  {{.Disable}}></menuitem>
	<menuitem label="Hello with bad onclick" onclick="unknown" {{.Disable}}></menuitem>
	<menuitem label="Hello without onclick" {{.Disable}}></menuitem>
	<menuitem label="Disabled Hello" onclick="OnHelloClick" disabled></menuitem>
	<menuitem separator></menuitem>

	{{if .Disable}}
	<menuitem label="Enable all" onclick="OnEnableAll"></menuitem>
	{{else}}
	<menuitem label="Disable all" onclick="OnDisableAll"></menuitem>
	{{end}}
</menu>
	`
}

// OnHelloClick is the function called when a hello button is clicked.
func (m *TestMenu) OnHelloClick() {
	app.Log("hello clicked")
}

// OnEnableAll is the function called when the "Enable all" button is clicked.
func (m *TestMenu) OnEnableAll() {
	m.Disable = ""
	app.Render(m)
}

// OnDisableAll is the function called when the "Disable all" button is clicked.
func (m *TestMenu) OnDisableAll() {
	m.Disable = "disabled"
	app.Render(m)
}
