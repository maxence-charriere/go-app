package main

import (
	"fmt"

	"github.com/murlokswarm/app"
)

// CustomMenu is a component to demonstrate how to work with menu and menuitems.
// It contains different customization for menuitems.
type CustomMenu struct{}

// Render returns return the HTML describing the custom menu.
func (c *CustomMenu) Render() string {
	return `
<menu>
	<menuitem label="MenuItem with Go callback" onclick="OnClick" />
	<menuitem label="MenuItem with shortcut" shortcut="meta+e" onclick="OnMenuItemWithShorcutClick" />
	<menuitem label="MenuItem with separator" separator="true" onclick="OnMenuItemWithSeparatorClick" />
	<menuitem label="MenuItem with icon" icon="logo.png" onclick="OnMenuItemWithIconClick" />
	<menuitem label="MenuItem disabled" onclick="OnClick" disabled="true" />
</menu>
	`
}

// OnClick is called when MenuItem with Go callback is clicked.
func (c *CustomMenu) OnClick() {
	fmt.Println("MenuItem with Go callback clicked")
}

// OnMenuItemWithShorcutClick is called when MenuItem with shortcutis clicked.
func (c *CustomMenu) OnMenuItemWithShorcutClick() {
	fmt.Println("MenuItem with shortcut clicked")
}

// OnMenuItemWithSeparatorClick is called when MenuItem with separator is
// clicked.
func (c *CustomMenu) OnMenuItemWithSeparatorClick() {
	fmt.Println("MenuItem with separator clicked")
}

// OnMenuItemWithIconClick is called when MenuItem with icon is clicked.
func (c *CustomMenu) OnMenuItemWithIconClick() {
	fmt.Println("MenuItem with icon clicked")
}

// /!\ Register the component. Required to use the component into a context.
func init() {
	app.RegisterComponent(&CustomMenu{})
}
