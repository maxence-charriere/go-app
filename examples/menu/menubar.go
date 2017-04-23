package main

import "github.com/murlokswarm/app"

// MenuBar is the component that define the menu bar.
// It Implements the common menu item present in a menu bar and uses the
// EditMenu component for the Edit section.
type MenuBar struct{}

// Render returns return the HTML describing the menu bar.
func (m *MenuBar) Render() string {
	return `
<menu>
	<menu label="app">
		<menuitem label="About" selector="orderFrontStandardAboutPanel:" separator="true" />
		<menuitem label="Preferences…" shortcut="meta+," separator="true" disabled="true" />
		<menuitem label="Hide" shortcut="meta+h" selector="hide:" />
		<menuitem label="Hide Others" shortcut="meta+alt+h" selector="hideOtherApplications:" />
		<menuitem label="Show All" selector="unhideAllApplications:" separator="true" />
		<menuitem label="Quit" shortcut="meta+q" selector="terminate:" /> 
	</menu>

	<EditMenu />

	<menu label="Window">
		<menuitem label="Minimize" shortcut="meta+m" selector="performMiniaturize:" />
		<menuitem label="Zoom" selector="performZoom:" separator="true" />
		<menuitem label="Bring All to Front" selector="arrangeInFront:" />
		<menuitem label="Close" shortcut="meta+w" selector="performClose:" />
	</menu>
	
	<menu label="Help"></menu>
</menu>
	`
}

// /!\ Register the component. Required to use the component into a context.
func init() {
	app.RegisterComponent(&MenuBar{})
}
