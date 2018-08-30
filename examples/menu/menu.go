package main

import (
	"github.com/murlokswarm/app"
)

// /!\ Import the component. Required to use a component.
func init() {
	app.Import(&EditMenu{})
	app.Import(&CustomMenu{})
}

// EditMenu is the component that describes the edit menu.
type EditMenu app.ZeroCompo

// Render returns return the HTML describing the edit menu.
func (e *EditMenu) Render() string {
	return `
<menu label="Edit">
	<menuitem label="Cut" keys="cmdorctrl+x" selector="cut:"></menuitem>
	<menuitem label="Copy" keys="cmdorctrl+c" selector="copy:"></menuitem>
	<menuitem label="Paste" keys="cmdorctrl+v" selector="paste:"></menuitem>
	<menuitem separator></menuitem>
	<menuitem label="Select All" keys="cmdorctrl+a" selector="selectAll:"></menuitem>
</menu>
	`
}

// CustomMenu is a component to demonstrate how to work with menu and menuitems.
// It contains different customization for menuitems.
type CustomMenu app.ZeroCompo

// Render returns return the HTML describing the custom menu.
func (c *CustomMenu) Render() string {
	return `
<menu label="Custom">
	<menuitem label="MenuItem with Go callback" onclick="OnClick"></menuitem>
	<menuitem label="MenuItem with icon" onclick="OnClickWithIcon" icon="logo.png" checked></menuitem>
	<menuitem label="MenuItem with keys" keys="cmdorctrl+e" onclick="OnMenuItemWithShorcutClick" checked></menuitem>
	<menuitem separator></menuitem>
	<menuitem label="MenuItem disabled" onclick="OnClick" disabled="true"></menuitem>
</menu>
	`
}

// OnClick is called when MenuItem with Go callback is clicked.
func (c *CustomMenu) OnClick() {
	app.Log("MenuItem with Go callback clicked")
}

// OnClickWithIcon is called when MenuItem with Go callback is clicked.
func (c *CustomMenu) OnClickWithIcon() {
	app.Log("MenuItem with icon clicked")
}

// OnMenuItemWithShorcutClick is called when MenuItem with keysis clicked.
func (c *CustomMenu) OnMenuItemWithShorcutClick() {
	app.Log("MenuItem with keys clicked")
}
