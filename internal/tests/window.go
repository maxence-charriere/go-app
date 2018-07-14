package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/murlokswarm/app"
)

func testWindow(t *testing.T, d app.Driver) {
	tests := []struct {
		scenario string
		config   app.WindowConfig
		function func(t *testing.T, w app.Window)
	}{
		{
			scenario: "create",
		},
		{
			scenario: "create with a default component",
			config: app.WindowConfig{
				DefaultURL: "tests.hello",
			},
		},
		{
			scenario: "move",
			function: testWindowMove,
		},
		{
			scenario: "resize",
			function: testWindowResize,
		},
		{
			scenario: "focus",
			function: testWindowFocus,
		},
		{
			scenario: "full screen",
			function: testWindowFullScreen,
		},
		{
			scenario: "minimize",
			function: testWindowMinimize,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			w, err := d.NewWindow(test.config)
			if app.NotSupported(err) {
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			defer w.Close()

			if test.function == nil {
				return
			}
			test.function(t, w)
		})
	}

	testElemWithCompo(t, func() (app.ElemWithCompo, error) {
		return d.NewWindow(app.WindowConfig{})
	})

	testElementWithNavigation(t, func() (app.Navigator, error) {
		return d.NewWindow(app.WindowConfig{})
	})
}

func testWindowMove(t *testing.T, w app.Window) {
	w.Move(42, 42)
	x, y := w.Position()
	assert.Equal(t, 42.0, x)
	assert.Equal(t, 42.0, y)

	w.Center()
	cx, cy := w.Position()
	assert.NotEqual(t, x, cx)
	assert.NotEqual(t, y, cy)
}

func testWindowResize(t *testing.T, w app.Window) {
	w.Resize(100, 100)
	width, height := w.Size()
	assert.Equal(t, 100.0, width)
	assert.Equal(t, 100.0, height)
}

func testWindowFocus(t *testing.T, w app.Window) {
	w.Focus()
}

func testWindowFullScreen(t *testing.T, w app.Window) {
	w.ToggleFullScreen()
	w.ToggleFullScreen()
}

func testWindowMinimize(t *testing.T, w app.Window) {
	w.ToggleMinimize()
	w.ToggleMinimize()
}
