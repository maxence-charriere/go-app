package tests

import (
	"testing"

	"github.com/murlokswarm/app"
	"github.com/stretchr/testify/require"
)

func testContextMenu(t *testing.T, d app.Driver) {
	testMenu(t, d.NewContextMenu)
}

func testMenubar(t *testing.T, d app.Driver) {
	testMenu(t, func(c app.MenuConfig) (app.Menu, error) {
		return d.MenuBar()
	})
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
			require.NoError(t, err)

			if test.function == nil {
				return
			}
			test.function(t, m)
		})
	}
}

func testMenuLoadSuccess(t *testing.T, m app.Menu) {
	err := m.Load("tests.menu")
	require.NoError(t, err)
}

func testMenuLoadFail(t *testing.T, m app.Menu) {
	err := m.Load("tests.tralala")
	require.Error(t, err)
}

func testMenuRenderSuccess(t *testing.T, m app.Menu) {
	err := m.Load("tests.menu")
	require.NoError(t, err)

	compo := m.Component()
	require.NotNil(t, compo)

	menu := compo.(*Menu)
	menu.Label = "a menu for test"

	err = m.Render(menu)
	require.NoError(t, err)
}

func testMenuRenderFail(t *testing.T, m app.Menu) {
	err := m.Load("tests.menu")
	require.NoError(t, err)

	compo := m.Component()
	require.NotNil(t, compo)

	menu := compo.(*Menu)
	menu.SimulateErr = true

	err = m.Render(menu)
	require.Error(t, err)
}
