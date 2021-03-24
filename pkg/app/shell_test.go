package app

import (
	"testing"
)

func TestShellPreRender(t *testing.T) {
	testSkipWasm(t)

	utests := []struct {
		scenario string
		shell    UI
	}{
		{
			scenario: "empty shell",
			shell:    Shell(),
		},
		{
			scenario: "shell with content",
			shell:    Shell().Content(Div()),
		},
		{
			scenario: "shell with menu",
			shell:    Shell().Menu(Div()),
		},
		{
			scenario: "shell with submenu",
			shell:    Shell().Submenu(Div()),
		},
		{
			scenario: "shell with overlay",
			shell:    Shell().OverlayMenu(Div()),
		},
		{
			scenario: "shell with menu and content",
			shell: Shell().
				Menu(Div()).
				Content(Div()),
		},
		{
			scenario: "shell with menu, submenu and content",
			shell: Shell().
				Menu(Div()).
				Submenu(Div()).
				Content(Div()),
		},
		{
			scenario: "shell with menu, submenu, overlay menu and content",
			shell: Shell().
				Menu(Div()).
				Submenu(Div()).
				OverlayMenu(Div()).
				Content(Div()),
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			d := NewServerTester(u.shell)
			defer d.Close()
			d.PreRender()
		})
	}
}
