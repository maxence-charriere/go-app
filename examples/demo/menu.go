package main

import (
	"github.com/murlokswarm/app"
)

// Menu is a component that contains menu related examples.
type Menu app.ZeroCompo

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
	</div>
</div>
	`
}

// OnContextMenu is the function called to display a context menu.
func (m *Menu) OnContextMenu() {
	app.NewContextMenu("contextmenu")
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
