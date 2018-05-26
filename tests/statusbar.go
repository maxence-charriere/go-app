package tests

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/murlokswarm/app"
	"github.com/stretchr/testify/require"
)

func testStatusBar(t *testing.T, d app.Driver) {
	tests := []struct {
		scenario string
		function func(t *testing.T, d app.StatusBarMenu, driver app.Driver)
	}{
		{
			scenario: "set icon success",
			function: testStatusBarSetIconSuccess,
		},
		{
			scenario: "set icon fails",
			function: testStatusBarSetIconFail,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			menu, err := d.StatusBar()
			if app.NotSupported(err) {
				return
			}

			require.NoError(t, err)
			test.function(t, menu, d)
		})
	}

	testMenu(t, func(c app.MenuConfig) (app.Menu, error) {
		return d.StatusBar()
	})
}

func testStatusBarSetIconSuccess(t *testing.T, m app.StatusBarMenu, driver app.Driver) {
	_, filename, _, _ := runtime.Caller(0)
	filename = filepath.Join(filepath.Dir(filename), "resources", "logo.png")

	err := m.SetIcon(filename)
	if app.NotSupported(err) {
		return
	}
	require.NoError(t, err)
}

func testStatusBarSetIconFail(t *testing.T, m app.StatusBarMenu, driver app.Driver) {
	_, filename, _, _ := runtime.Caller(0)
	filename = filepath.Join(filepath.Dir(filename), "resources", "logo.bmp")

	err := m.SetIcon(filename)
	if app.NotSupported(err) {
		return
	}
	require.Error(t, err)
}
