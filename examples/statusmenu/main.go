package main

import (
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/mac"
)

func main() {
	app.Import(&Menu{})

	app.Run(&mac.Driver{
		Bundle: mac.Bundle{
			Background: true,
		},

		OnRun: func() {
			app.NewStatusMenu(app.StatusMenuConfig{
				Icon: app.Resources("logo.png"),
				Text: "Background app",
				URL:  "/Menu",
			})
		},
	}, app.Logs())
}

// Menu is a component that describes a status able to change its text and icon.
type Menu struct {
	IconHidden bool
	TextHidden bool
}

// Render returns the HTML describing the status menu.
func (m *Menu) Render() string {
	return `
<menu>
	<menuitem label="Show icon" onclick="OnShowIcon" {{if not .IconHidden}}disabled{{end}}></menuitem>
	<menuitem label="Hide icon" onclick="OnHideIcon" {{if .IconHidden}}disabled{{end}}></menuitem>
	<menuitem separator></menuitem>

	<menuitem label="Show text" onclick="OnShowText" {{if not .TextHidden}}disabled{{end}}></menuitem>
	<menuitem label="Hide text" onclick="OnHideText" {{if .TextHidden}}disabled{{end}}></menuitem>
	<menuitem separator></menuitem>

	<menuitem label="Open a window" onclick="OnOpenWindow"></menuitem>
	<menuitem separator></menuitem>

	<menuitem label="Close" selector="terminate:"></menuitem>
</menu>
	`
}

// OnShowIcon is the function called when the show icon button is clicked.
func (m *Menu) OnShowIcon() {
	app.ElemByCompo(m).WhenStatusMenu(func(s app.StatusMenu) {
		s.SetIcon(app.Resources("logo.png"))
		m.IconHidden = false
		app.Render(m)
	})
}

// OnHideIcon is the function called when the hide icon button is clicked.
func (m *Menu) OnHideIcon() {
	app.ElemByCompo(m).WhenStatusMenu(func(s app.StatusMenu) {
		s.SetIcon("")
		m.IconHidden = true
		app.Render(m)

		if m.TextHidden {
			m.OnShowText()
		}
	})

}

// OnShowText is the function called when the show text button is clicked.
func (m *Menu) OnShowText() {
	app.ElemByCompo(m).WhenStatusMenu(func(s app.StatusMenu) {
		s.SetText("Background app")
		m.TextHidden = false
		app.Render(m)
	})
}

// OnHideText is the function called when the hide text button is clicked.
func (m *Menu) OnHideText() {
	app.ElemByCompo(m).WhenStatusMenu(func(s app.StatusMenu) {
		s.SetText("")
		m.TextHidden = true
		app.Render(m)

		if m.IconHidden {
			m.OnShowIcon()
		}
	})
}

// OnOpenWindow is the function called when the open a window button is clicked.
func (m *Menu) OnOpenWindow() {
	app.NewWindow(app.WindowConfig{
		Width:          1024,
		Height:         720,
		TitlebarHidden: true,
	})
}
