package main

import "github.com/murlokswarm/app"

// EditMenu is the component that define the an edit menu.
type EditMenu struct{}

// Render returns return the HTML describing the edit menu.
func (e *EditMenu) Render() string {
	return `
<menu label="Edit">
	<menuitem label="Undo" shortcut="meta+z" selector="undo:" />
	<menuitem label="Redo" shortcut="meta+shift+z" selector="redo:" separator="true" />
	<menuitem label="Cut" shortcut="meta+x" selector="cut:" />
	<menuitem label="Copy" shortcut="meta+c" selector="copy:" />
	<menuitem label="Paste" shortcut="meta+v" selector="paste:" />
	<menuitem label="Paste and Match Style" shortcut="shift+alt+meta+v" selector="pasteAsPlainText:" />
	<menuitem label="Delete" selector="delete:" />
	<menuitem label="Select All" shortcut="meta+a" selector="selectAll:" />
</menu>
	`
}

// /!\ Register the component. Required to use the component into a context.
func init() {
	app.RegisterComponent(&EditMenu{})
}
