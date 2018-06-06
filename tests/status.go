package tests

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/stretchr/testify/require"
)

func testStatusMenu(t *testing.T, d app.Driver) {
	tests := []struct {
		scenario string
		function func(t *testing.T, d app.StatusMenu, driver app.Driver)
	}{
		{
			scenario: "set text success",
			function: testStatusMenuSetTextSuccess,
		},
		{
			scenario: "set icon success",
			function: testStatusMenuSetIconSuccess,
		},
		{
			scenario: "set icon fails",
			function: testStatusMenuSetIconFail,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			menu, err := d.NewStatusMenu(app.StatusMenuConfig{
				Text: "test",
			})
			if app.NotSupported(err) {
				return
			}
			defer menu.Close()

			require.NoError(t, err)
			test.function(t, menu, d)
		})
	}

	testMenu(t, func(c app.MenuConfig) (app.Menu, error) {
		return d.NewStatusMenu(app.StatusMenuConfig{
			Text: "hello",
		})
	})
}

func testStatusMenuSetTextSuccess(t *testing.T, m app.StatusMenu, driver app.Driver) {
	err := m.SetText(uuid.New().String())
	require.NoError(t, err)
}

func testStatusMenuSetIconSuccess(t *testing.T, m app.StatusMenu, driver app.Driver) {
	_, filename, _, _ := runtime.Caller(0)
	filename = filepath.Join(filepath.Dir(filename), "resources", "logo.png")

	err := m.SetIcon(filename)
	if app.NotSupported(err) {
		return
	}
	require.NoError(t, err)
}

func testStatusMenuSetIconFail(t *testing.T, m app.StatusMenu, driver app.Driver) {
	_, filename, _, _ := runtime.Caller(0)
	filename = filepath.Join(filepath.Dir(filename), "resources", "logo.bmp")

	err := m.SetIcon(filename)
	if app.NotSupported(err) {
		return
	}
	require.Error(t, err)
}
