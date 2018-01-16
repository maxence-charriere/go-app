package main

import (
	"net/url"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
)

func init() {
	app.Import(&Menu{})
}

// Menu is a component to test menu based elements.
type Menu struct {
	DisableAll  bool
	RandomTitle uuid.UUID
}

// Render statisfies the app.Component interface.
func (m *Menu) Render() string {
	return `
<menu>
	<menuitem label="button" onclick="OnButtonClick" {{if .DisableAll}}disabled{{end}}></menuitem>
	<menuitem label="button with icon" onclick="OnButtonWithIconClick" {{if .DisableAll}}disabled{{end}}></menuitem>
	<menuitem label="{{.RandomTitle}}" onclick="OnButtonWithRandomTitleClicked"></menuitem>

	<menuitem separator></menuitem>

	<menuitem label="set dock badge" onclick="OnSetDockBadge"></menuitem>
	<menuitem label="unset dock badge" onclick="OnUnsetDockBadge"></menuitem>
	
	<menuitem separator></menuitem>

	<menuitem label="set dock icon" onclick="OnSetDockIcon"></menuitem>
	<menuitem label="unset dock icon" onclick="OnUnsetDockIcon"></menuitem>
	
	<menuitem separator></menuitem>

	<menu label="submenu">
		<menuitem label="sub button" onclick="OnSubButtonClick" {{if .DisableAll}}disabled{{end}}></menuitem>
		<menuitem label="sub button without action"></menuitem>	
	</menu>
	
	<menuitem separator></menuitem>
	
	<menuitem label="enable all" onclick="OnEnableAllClick" {{if not .DisableAll}}disabled{{end}}></menuitem>
	<menuitem label="disable all" onclick="OnDisableAllClick" {{if .DisableAll}}disabled{{end}}></menuitem>	
</menu>
	`
}

// OnNavigate is the function that is called when the component is navigated on.
func (m *Menu) OnNavigate(u *url.URL) {
	m.RandomTitle = uuid.New()
	app.Render(m)
}

// OnButtonClick is the function that is called when the button labelled
// "button" is clicked.
func (m *Menu) OnButtonClick() {
	app.DefaultLogger.Log("button clicked")
}

// OnButtonWithIconClick is the function that is called when the button labelled
// "button with icon" is clicked.
func (m *Menu) OnButtonWithIconClick() {
	app.DefaultLogger.Log("button with icon clicked")
}

// OnSetDockBadge is the function that is called when the button labelled "set
// dock badge" is clicked.
func (m *Menu) OnSetDockBadge() {
	app.DefaultLogger.Log("button set dock badge clicked")

	if app.SupportsDock() {
		app.Dock().SetBadge(uuid.New())
	}
}

// OnUnsetDockBadge is the function that is called when the button labelled
// "unset dock badge" is clicked.
func (m *Menu) OnUnsetDockBadge() {
	app.DefaultLogger.Log("button unset dock badge clicked")

	if app.SupportsDock() {
		app.Dock().SetBadge(nil)
	}
}

// OnSetDockIcon is the function that is called when the button labelled "set
// dock icon" is clicked.
func (m *Menu) OnSetDockIcon() {
	app.DefaultLogger.Log("button set dock icon clicked")

	if app.SupportsDock() {
		app.Dock().SetIcon(filepath.Join(app.Resources(), "logo.png"))
	}
}

// OnUnsetDockIcon is the function that is called when the button labelled
// "unset dock icon" is clicked.
func (m *Menu) OnUnsetDockIcon() {
	app.DefaultLogger.Log("button unset dock icon clicked")

	if app.SupportsDock() {
		app.Dock().SetIcon("")
	}
}

// OnButtonWithRandomTitleClicked is the function that is called when the button
// with randow title is clicked.
func (m *Menu) OnButtonWithRandomTitleClicked() {
	app.DefaultLogger.Log("button with random title clicked")
	m.RandomTitle = uuid.New()
	app.Render(m)
}

// OnSubButtonClick is the function that is called when the button labelled "sub
// button" is clicked.
func (m *Menu) OnSubButtonClick() {
	app.DefaultLogger.Log("sub button clicked")
}

// OnEnableAllClick is the function that is called when the button labelled
// "enable all" is clicked.
func (m *Menu) OnEnableAllClick() {
	app.DefaultLogger.Log("button enable all clicked")
	m.DisableAll = false
	app.Render(m)
}

// OnDisableAllClick is the function that is called when the button labelled
// "disable all" is clicked.
func (m *Menu) OnDisableAllClick() {
	app.DefaultLogger.Log("button disable all clicked")
	m.DisableAll = true
	app.Render(m)
}
