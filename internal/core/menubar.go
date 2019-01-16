package core

import (
	"encoding/json"
	"net/url"

	"github.com/murlokswarm/app"
)

func init() {
	app.Import(&MenuBar{})
}

// MenuBar is a component that represents a menu bar.
type MenuBar struct {
	AppName    string
	AppURL     string
	CustomURLs []string
	EditURL    string
	FileURL    string
	HelpURL    string
	WindowURL  string
}

// OnNavigate satisfies the app.Navigable interface.
func (m *MenuBar) OnNavigate(u *url.URL) {
	m.AppName = app.Name()
	app.Render(m)
}

// Render satisfies the app.Compo interface.
func (m *MenuBar) Render() string {
	return `
<menu>
	{{if .AppURL}}
	{{compo .AppURL}}
	{{else}}
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
	{{end}}
</menu>
		`
}

func menuBarConfigToAddr(c app.MenuBarConfig) string {
	u, _ := url.Parse(app.CompoName(&MenuBar{}))
	u.Query().Set("AppURL", c.AppURL)
	u.Query().Set("EditURL", c.EditURL)
	u.Query().Set("FileURL", c.FileURL)
	u.Query().Set("HelpURL", c.HelpURL)
	u.Query().Set("WindowURL", c.WindowURL)

	customURLs, _ := json.Marshal(c.CustomURLs)
	u.Query().Set("CustomURLs", string(customURLs))

	return u.String()
}
