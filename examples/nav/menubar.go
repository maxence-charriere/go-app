package main

import "github.com/murlokswarm/app"

// MenuBar is the component that define the menu bar.
type MenuBar struct{}

// Render returns return the HTML describing the menu bar.
func (m *MenuBar) Render() string {
	return `
<menu>
	<menu label="app">
		<menuitem label="Close" selector="performClose:" shortcut="meta+w" />
		<menuitem label="Quit" selector="terminate:" shortcut="meta+q" /> 
	</menu>
</menu>
	`
}

func init() {
	app.RegisterComponent(&MenuBar{})
}
