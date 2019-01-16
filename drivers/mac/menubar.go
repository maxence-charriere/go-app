// +build darwin,amd64

package mac

import (
	"net/url"
	"os/exec"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/core"
)

func init() {
	app.Import(&MenuBar{})
	app.Import(&AppMenu{})
	app.Import(&EditMenu{})
	app.Import(&WindowMenu{})
	app.Import(&HelpMenu{})
}

// MenuBar is a component that describes a menu bar.
type MenuBar struct {
	AppURL    string
	EditURL   string
	WindowURL string
	HelpURL   string
	CutomURLs []string
}

// OnNavigate setup the menu bar sections.
func (m *MenuBar) OnNavigate(u *url.URL) {
	m.AppURL = u.Query().Get("appurl")
	m.AppURL = core.CompoNameFromURLString(m.AppURL)
	if len(m.AppURL) == 0 {
		m.AppURL = "mac.appmenu"
	}

	m.EditURL = u.Query().Get("editurl")
	m.EditURL = core.CompoNameFromURLString(m.EditURL)
	if len(m.EditURL) == 0 {
		m.EditURL = "mac.editmenu"
	}

	m.WindowURL = u.Query().Get("windowurl")
	m.WindowURL = core.CompoNameFromURLString(m.WindowURL)
	if len(m.WindowURL) == 0 {
		m.WindowURL = "mac.windowmenu"
	}

	m.HelpURL = u.Query().Get("helpurl")
	m.HelpURL = core.CompoNameFromURLString(m.HelpURL)
	if len(m.HelpURL) == 0 {
		m.HelpURL = "mac.helpmenu"
	}

	for _, u := range u.Query()["custom"] {
		customURL := core.CompoNameFromURLString(u)
		if len(customURL) != 0 {
			m.CutomURLs = append(m.CutomURLs, customURL)
		}
	}

	app.Render(m)
}

// Render returns the markup that describes the menu bar.
func (m *MenuBar) Render() string {
	return `
<menu>
	{{if .AppURL}}
		{{compo .AppURL}}
	{{else}}
		<!-- prevent cocoa to generate a non modifiable menu -->
		<mac.appmenu>
	{{end}}
	{{if .EditURL}}
		{{compo .EditURL}}
	{{end}}
	{{if .WindowURL}}
		{{compo .WindowURL}}
	{{end}}
	{{range .CutomURLs}}
		{{compo .}}
	{{end}}
	{{if .HelpURL}}
		{{compo .HelpURL}}
	{{end}}
</menu>
	`
}

// AppMenu is a component that describes the default app menu.
type AppMenu struct {
	AppName string
}

// OnMount initializes the menu application name.
func (m *AppMenu) OnMount() {
	m.AppName = app.Name()
	app.Render(m)
}

// Render returns the markup that describes the app menu.
func (m *AppMenu) Render() string {
	return `
<menu>
	<menuitem label="About {{.AppName}}" role="about">
	<menuitem separator>

	<menuitem label="Preferences…" keys="cmdorctrl+," onclick="OnPreferences">
	<menuitem separator>

	<menuitem label="Hide {{.AppName}}" keys="cmdorctrl+h" role="hide">
	<menuitem label="Hide Others" keys="cmdorctrl+alt+h" role="hideOthers">
	<menuitem label="Show All" role="unhide">
	<menuitem separator>

	<menuitem label="Quit {{.AppName}}" keys="cmdorctrl+q" role="quit">
</menu>
	`
}

// OnPreferences is the function called when the Preferences button is clicked.
func (m *AppMenu) OnPreferences() {
	app.Emit(PreferencesRequested)
}

// EditMenu is a component that describes the default edit menu.
type EditMenu app.ZeroCompo

// Render returns the markup that describes the edit menu.
func (m *EditMenu) Render() string {
	return `
<menu label="Edit">
	<menuitem label="Undo" keys="cmdorctrl+z" role="undo">
	<menuitem label="Redo" keys="cmdorctrl+shift+z" role="redo">
	<menuitem separator>
	<menuitem label="Cut" keys="cmdorctrl+x" role="cut">
	<menuitem label="Copy" keys="cmdorctrl+c" role="copy">
	<menuitem label="Paste" keys="cmdorctrl+v" role="paste">
	<menuitem label="Paste and Match Style" keys="shift+alt+cmdorctrl+v" role="pasteAndMatchStyle">
	<menuitem label="Delete" role="delete">
	<menuitem label="Select All" keys="cmdorctrl+a" role="selectAll">
</menu>
	`
}

// WindowMenu is a component that describes the default window menu.
type WindowMenu app.ZeroCompo

// Render returns the markup that describes the window menu.
func (m *WindowMenu) Render() string {
	return `
<menu label="Window">
	<menuitem label="Minimize" keys="cmdorctrl+m" role="minimize">
	<menuitem label="Zoom" role="zoom">
	<menuitem separator>
	<menuitem label="Bring All to Front" role="arrangeInFront">
	<menuitem label="Close" keys="cmdorctrl+w" role="close">
</menu>
	`
}

// HelpMenu is a component that describes the default help menu.
type HelpMenu app.ZeroCompo

// Render returns the markup that describes the help menu.
func (m *HelpMenu) Render() string {
	return `
<menu label="Help">
	<menuitem label="Built with github.com/murlokswarm/app" onclick="OnBuiltWith">
</menu>
	`
}

// OnBuiltWith is called when the On Built with button is clicked.
// It opens the app repository on the default browser.
func (m *HelpMenu) OnBuiltWith() {
	cmd := exec.Command("open", "https://github.com/murlokswarm/app")
	if err := cmd.Run(); err != nil {
		app.Logf("opening https://github.com/murlokswarm/app failed: %s", err)
	}
}
