package tests

import (
	"testing"

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
			scenario: "window is decorated with logs",
			function: testWindowIsDecorated,
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

	testElementWithComponent(t, func() (app.ElementWithComponent, error) {
		return d.NewWindow(app.WindowConfig{})
	})

	testElementWithNavigation(t, func() (app.Navigator, error) {
		return d.NewWindow(app.WindowConfig{})
	})
}

func testWindowIsDecorated(t *testing.T, w app.Window) {
	if base := w.Base(); base == w {
		t.Error("window is not decorated")
	}
}

func testWindowMove(t *testing.T, w app.Window) {
	w.Move(420, 420)
	x, y := w.Position()
	if x != 420 {
		t.Error("window x is not 420:", x)
	}
	if y != 420 {
		t.Error("window y is not 420:", y)
	}

	w.Center()
	cx, cy := w.Position()
	if cx == x {
		t.Error("window was not centered on x axis")
	}
	if cy == y {
		t.Error("window was not centered on y axis")
	}
}

func testWindowResize(t *testing.T, w app.Window) {
	w.Resize(100, 100)
	width, height := w.Size()
	if width != 100 {
		t.Error("window width is not 100:", width)
	}
	if height != 100 {
		t.Error("window height is not 100:", height)
	}
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
