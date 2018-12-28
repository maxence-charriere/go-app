// +build darwin,amd64

package mac

import (
	"net/url"
	"os/exec"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/core"
	"github.com/pkg/errors"
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
	<menuitem label="About {{.AppName}}" selector="orderFrontStandardAboutPanel:"></menuitem>
	<menuitem separator></menuitem>

	<menuitem label="Preferences…" keys="cmdorctrl+," onclick="OnPreferences"></menuitem>
	<menuitem separator></menuitem>

	<menuitem label="Hide {{.AppName}}" keys="cmdorctrl+h" selector="hide:"></menuitem>
	<menuitem label="Hide Others" keys="cmdorctrl+alt+h" selector="hideOtherApplications:"></menuitem>
	<menuitem label="Show All" selector="unhideAllApplications:"></menuitem>
	<menuitem separator></menuitem>

	<menuitem label="Quit {{.AppName}}" keys="cmdorctrl+q" selector="terminate:"></menuitem>
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
	<menuitem label="Undo" keys="cmdorctrl+z" selector="undo:"></menuitem>
	<menuitem label="Redo" keys="cmdorctrl+shift+z" selector="redo:"></menuitem>
	<menuitem separator></menuitem>
	<menuitem label="Cut" keys="cmdorctrl+x" selector="cut:"></menuitem>
	<menuitem label="Copy" keys="cmdorctrl+c" selector="copy:"></menuitem>
	<menuitem label="Paste" keys="cmdorctrl+v" selector="paste:"></menuitem>
	<menuitem label="Paste and Match Style" keys="shift+alt+cmdorctrl+v" selector="pasteAsPlainText:"></menuitem>
	<menuitem label="Delete" selector="delete:"></menuitem>
	<menuitem label="Select All" keys="cmdorctrl+a" selector="selectAll:"></menuitem>
</menu>
	`
}

// WindowMenu is a component that describes the default window menu.
type WindowMenu app.ZeroCompo

// Render returns the markup that describes the window menu.
func (m *WindowMenu) Render() string {
	return `
<menu label="Window">
	<menuitem label="Minimize" keys="cmdorctrl+m" selector="performMiniaturize:"></menuitem>
	<menuitem label="Zoom" selector="performZoom:"></menuitem>
	<menuitem separator></menuitem>
	<menuitem label="Bring All to Front" selector="arrangeInFront:"></menuitem>
	<menuitem label="Close" keys="cmdorctrl+w" selector="performClose:"></menuitem>
</menu>
	`
}

// HelpMenu is a component that describes the default help menu.
type HelpMenu app.ZeroCompo

// Render returns the markup that describes the help menu.
func (m *HelpMenu) Render() string {
	return `
<menu label="Help">
	<menuitem label="Built with github.com/murlokswarm/app" onclick="OnBuiltWith"></menuitem>
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

func newMenuBar(c MenuBarConfig) *Menu {
	m := newMenu(app.MenuConfig{}, "menu bar")
	if m.Err() != nil {
		return m
	}

	if len(c.URL) == 0 {
		format := "mac.menubar?appurl=%s&editurl=%s&windowurl=%s&helpurl=%s"
		for _, u := range c.CustomURLs {
			format += "&custom=" + u
		}

		m.Load(
			format,
			c.AppURL,
			c.EditURL,
			c.WindowURL,
			c.HelpURL,
		)
	} else {
		m.Load(c.URL)
	}

	if m.Err() != nil {
		return m
	}

	if err := driver.macRPC.Call("driver.SetMenubar", nil, m.id); err != nil {
		m.SetErr(errors.Wrap(err, "set menu bar"))
	}

	return m
}
