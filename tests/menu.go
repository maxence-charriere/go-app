package tests

import (
	"testing"

	"github.com/murlokswarm/app"
)

func testContextMenu(t *testing.T, d app.Driver) {
	testMenu(t, d.NewContextMenu)
}

func testMenu(t *testing.T, setup func(c app.MenuConfig) (app.Menu, error)) {
	tests := []struct {
		scenario string
		config   app.MenuConfig
		function func(t *testing.T, w app.Menu)
	}{
		{
			scenario: "create",
		},
		{
			scenario: "create with a default component",
			config: app.MenuConfig{
				DefaultURL: "tests.menu",
			},
		},
		{
			scenario: "load a component",
			function: testMenuLoadSuccess,
		},
		{
			scenario: "load a component fails",
			function: testMenuLoadFail,
		},
		{
			scenario: "render a component",
			function: testMenuRenderSuccess,
		},
		{
			scenario: "render a component fails",
			function: testMenuRenderFail,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			m, err := setup(test.config)
			if app.NotSupported(err) {
				return
			}
			if err != nil {
				t.Fatal(err)
			}

			if test.function == nil {
				return
			}
			test.function(t, m)
		})
	}
}

func testMenuLoadSuccess(t *testing.T, m app.Menu) {
	if err := m.Load("tests.menu"); err != nil {
		t.Fatal(err)
	}
}

func testMenuLoadFail(t *testing.T, m app.Menu) {
	err := m.Load("tests.tralala")
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testMenuRenderSuccess(t *testing.T, m app.Menu) {
	if err := m.Load("tests.menu"); err != nil {
		t.Fatal(err)
	}

	compo := m.Component()
	if compo == nil {
		t.Fatal("component is nil")
	}

	menu := compo.(*Menu)
	menu.Label = "a menu for test"

	if err := m.Render(menu); err != nil {
		t.Fatal(err)
	}
}

func testMenuRenderFail(t *testing.T, m app.Menu) {
	if err := m.Load("tests.menu"); err != nil {
		t.Fatal(err)
	}

	compo := m.Component()
	if compo == nil {
		t.Fatal("component is nil")
	}

	menu := compo.(*Menu)
	menu.SimulateErr = true

	err := m.Render(menu)
	if err == nil {
		t.Fatal(err)
	}
	t.Log(err)
}
