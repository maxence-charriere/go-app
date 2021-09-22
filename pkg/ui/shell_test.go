package ui

import (
	"testing"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func TestShellPreRender(t *testing.T) {
	utests := []struct {
		scenario string
		shell    app.UI
	}{
		{
			scenario: "empty shell",
			shell:    Shell(),
		},
		{
			scenario: "shell with content",
			shell:    Shell().Content(app.Div()),
		},
		{
			scenario: "shell with menu",
			shell:    Shell().Menu(app.Div()),
		},
		{
			scenario: "shell with submenu",
			shell:    Shell().Index(app.Div()),
		},
		{
			scenario: "shell with overlay",
			shell:    Shell().HamburgerMenu(app.Div()),
		},
		{
			scenario: "shell with menu and content",
			shell: Shell().
				Menu(app.Div()).
				Content(app.Div()),
		},
		{
			scenario: "shell with menu, submenu and content",
			shell: Shell().
				Menu(app.Div()).
				Index(app.Div()).
				Content(app.Div()),
		},
		{
			scenario: "shell with menu, submenu, overlay menu and content",
			shell: Shell().
				Menu(app.Div()).
				Index(app.Div()).
				HamburgerMenu(app.Div()).
				Content(app.Div()),
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			d := app.NewServerTester(u.shell)
			defer d.Close()
			d.PreRender()
		})
	}
}
