package tests

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/murlokswarm/app"
	"github.com/stretchr/testify/require"
)

func testDockTile(t *testing.T, d app.Driver) {
	tests := []struct {
		scenario string
		function func(t *testing.T, d app.DockTile, driver app.Driver)
	}{
		{
			scenario: "set icon success",
			function: testDockSetIconSuccess,
		},
		{
			scenario: "set icon fails",
			function: testDockSetIconFail,
		},
		{
			scenario: "set badge success",
			function: testDockSetBadgeSuccess,
		},
		{
			scenario: "set badge fails",
			function: testDockSetBadgeFails,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			dock, err := d.Dock()
			if err == app.ErrNotSupported {
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			test.function(t, dock, d)
		})
	}

	testMenu(t, func(c app.MenuConfig) (app.Menu, error) {
		return d.Dock()
	})
}

func testDockSetIconSuccess(t *testing.T, d app.DockTile, driver app.Driver) {
	_, filename, _, _ := runtime.Caller(0)
	filename = filepath.Join(filepath.Dir(filename), "resources", "logo.png")

	err := d.SetIcon(filename)
	if err == app.ErrNotSupported {
		return
	}
	require.NoError(t, err)
}

func testDockSetIconFail(t *testing.T, d app.DockTile, driver app.Driver) {
	_, filename, _, _ := runtime.Caller(0)
	filename = filepath.Join(filepath.Dir(filename), "resources", "logo.bmp")

	err := d.SetIcon(filename)
	if err == app.ErrNotSupported {
		return
	}
	require.Error(t, err)
}

func testDockSetBadgeSuccess(t *testing.T, d app.DockTile, driver app.Driver) {
	err := d.SetBadge("Hello")
	if err == app.ErrNotSupported {
		return
	}
	if err != nil {
		t.Fatal(err)
	}
}

func testDockSetBadgeFails(t *testing.T, d app.DockTile, driver app.Driver) {
	err := d.SetBadge(func() {})
	if err == app.ErrNotSupported {
		return
	}
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}
