package test

import (
	"net/url"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
)

func init() {
	app.Import(&Menu{})
}

// Menu is a component to test menu based elements.
type Menu struct {
	Name             string
	DisableAll       bool
	Checked          bool
	RandomTitle      string
	Separator        bool
	RenderRootToggle bool
	RenderToggle     bool
}

// Render statisfies the app.Compo interface.
func (m *Menu) Render() string {
	return `
<menu label="{{.Name}}">
	<menuitem label="button"
			  title="A button."
			  onclick="OnButtonClick"
			  {{if .DisableAll}}disabled{{end}}>
	</menuitem>
	<menuitem label="button with icon"
			  icon="logo.png"
			  onclick="OnButtonWithIconClick" 
			  {{if .DisableAll}}disabled{{end}}>
	</menuitem>
	<menuitem label="checkable button"
			  onclick="OnCheckableButton"
			  {{if .Checked}}checked{{end}}
	          {{if .DisableAll}}disabled{{end}}>
	</menuitem>

	<menuitem separator></menuitem>

	<menuitem label="set dock badge" 
			  onclick="OnSetDockBadge"
			  {{if .DisableAll}}disabled{{end}}>
	</menuitem>
	<menuitem label="unset dock badge" 
			  onclick="OnUnsetDockBadge"
			  {{if .DisableAll}}disabled{{end}}>
	</menuitem>
	
	<menuitem separator></menuitem>

	<menuitem label="set dock icon" 
			  onclick="OnSetDockIcon"
			  {{if .DisableAll}}disabled{{end}}>
	</menuitem>
	<menuitem label="unset dock icon" 
			  onclick="OnUnsetDockIcon"
			  {{if .DisableAll}}disabled{{end}}>
	</menuitem>
	
	<menuitem separator></menuitem>

	<menu label="submenu" {{if .DisableAll}}disabled{{end}}>
		<menuitem label="sub button" onclick="OnSubButtonClick"></menuitem>
		<menuitem label="sub button without action"></menuitem>	
	</menu>
	<menu label="disabled submenu" disabled></menu>

	<menuitem separator></menuitem>

	<menuitem label="random button: {{.RandomTitle}}" 
			  onclick="OnButtonWithRandomTitleClicked"
			  {{if .DisableAll}}disabled{{end}}>
	</menuitem>
	<menu label="random menu: {{.RandomTitle}}"
				onclick="OnButtonWithRandomTitleClicked"
				{{if .DisableAll}}disabled{{end}}>
		<menuitem label="disabled item" disabled></menuitem>
	</menu>
	
	<menuitem separator></menuitem>

	<menu label="separator test" {{if .DisableAll}}disabled{{end}}>
		<menuitem label="switch separator" 
				  onclick="OnSwitchSeparatorClick">
		</menuitem>

		{{if .Separator}}
			<menuitem separator></menuitem>		
		{{else}}
			<menuitem label="----- Separator -----"></menuitem>	
		{{end}}

		<menuitem label="disabled button" disabled></menuitem>		
	</menu>

	<menuitem separator></menuitem>

	<menu label="full render test" {{if .DisableAll}}disabled{{end}}>
		<menuitem label="render root" onclick="OnRenderRootClicked"></menuitem>
		<menuitem label="render root attribute" onclick="OnRenderRootAttributeClicked"></menuitem>
		{{if .RenderToggle}}
			<menuitem label="render menu" onclick="OnRenderTest"></menuitem>
		{{else}}
			<menu label="render item">
				<menuitem label="render" onclick="OnRenderTest"></menuitem>
			</menu>
		{{end}}
	</menu>

	{{if .RenderRootToggle}}
		<menuitem label="blank button to test root render" disabled></menuitem>
	{{end}}

	<menuitem separator></menuitem>
	
	<menuitem label="enable all" 
			onclick="OnEnableAllClick"
			{{if not .DisableAll}}disabled{{end}}>
	</menuitem>
	<menuitem label="disable all" 
			  onclick="OnDisableAllClick"
			  {{if .DisableAll}}disabled{{end}}>
	</menuitem>	
</menu>
	`
}

// OnNavigate is the function that is called when the component is navigated on.
func (m *Menu) OnNavigate(u *url.URL) {
	m.RandomTitle = uuid.New().String()
	app.Render(m)
}

// OnButtonClick is the function that is called when the button labelled
// "button" is clicked.
func (m *Menu) OnButtonClick() {
	app.Logf("button clicked")
}

// OnButtonWithIconClick is the function that is called when the button labelled
// "button with icon" is clicked.
func (m *Menu) OnButtonWithIconClick() {
	app.Logf("button with icon clicked")
}

// OnCheckableButton is the function that is called when the button labelled
// "checkable button" is clicked.
func (m *Menu) OnCheckableButton() {
	app.Logf("checkable button clicked")
	m.Checked = !m.Checked
	app.Render(m)
}

// OnSetDockBadge is the function that is called when the button labelled "set
// dock badge" is clicked.
func (m *Menu) OnSetDockBadge() {
	app.Logf("button set dock badge clicked")
	app.Dock().SetBadge(uuid.New().String())
}

// OnUnsetDockBadge is the function that is called when the button labelled
// "unset dock badge" is clicked.
func (m *Menu) OnUnsetDockBadge() {
	app.Logf("button unset dock badge clicked")
	app.Dock().SetBadge(nil)
}

// OnSetDockIcon is the function that is called when the button labelled "set
// dock icon" is clicked.
func (m *Menu) OnSetDockIcon() {
	app.Logf("button set dock icon clicked")
	app.Dock().SetIcon(app.Resources("logo.png"))
}

// OnUnsetDockIcon is the function that is called when the button labelled
// "unset dock icon" is clicked.
func (m *Menu) OnUnsetDockIcon() {
	app.Logf("button unset dock icon clicked")
	app.Dock().SetIcon("")
}

// OnButtonWithRandomTitleClicked is the function that is called when the button
// with randow title is clicked.
func (m *Menu) OnButtonWithRandomTitleClicked() {
	app.Logf("button with random title clicked")
	m.RandomTitle = uuid.New().String()
	app.Render(m)
}

// OnSubButtonClick is the function that is called when the button labelled "sub
// button" is clicked.
func (m *Menu) OnSubButtonClick() {
	app.Logf("sub button clicked")
}

// OnEnableAllClick is the function that is called when the button labelled
// "enable all" is clicked.
func (m *Menu) OnEnableAllClick() {
	app.Logf("button enable all clicked")
	m.DisableAll = false
	app.Render(m)
}

// OnDisableAllClick is the function that is called when the button labelled
// "disable all" is clicked.
func (m *Menu) OnDisableAllClick() {
	app.Logf("button disable all clicked")
	m.DisableAll = true
	app.Render(m)
}

// OnSwitchSeparatorClick is the function that is called when the button
// labelled "switch separator" is clicked.
func (m *Menu) OnSwitchSeparatorClick() {
	app.Logf("button switch separator clicked")
	m.Separator = !m.Separator
	app.Render(m)
}

// OnRenderRootClicked is the function that is called when the button labelled
// "render root" is clicked.
func (m *Menu) OnRenderRootClicked() {
	app.Logf("button render root clicked")
	m.RenderRootToggle = !m.RenderRootToggle
	app.Render(m)
}

// OnRenderRootAttributeClicked is the function that is called when the button
// labelled "render root attribute" is clicked.
func (m *Menu) OnRenderRootAttributeClicked() {
	app.Logf("button render root attribute clicked")
	m.Name = uuid.New().String()
	app.Render(m)
}

// OnRenderTest is the function that is called when a button form the render
// test menu is clicked.
func (m *Menu) OnRenderTest() {
	m.RenderToggle = !m.RenderToggle
	app.Render(m)
}
