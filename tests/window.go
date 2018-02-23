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
			scenario: "load a component",
			function: testWindowLoadSuccess,
		},
		{
			scenario: "load a component fails",
			function: testWindowLoadFail,
		},
		{
			scenario: "render a component",
			function: testWindowRenderSuccess,
		},
		{
			scenario: "render a component fails",
			function: testWindowRenderFail,
		},
		{
			scenario: "reload a component",
			function: testWindowReloadSuccess,
		},
		{
			scenario: "reload a component fails",
			function: testWindowReloadFail,
		},
		{
			scenario: "load previous component",
			function: testWindowPreviousSuccess,
		},
		{
			scenario: "load previous component fails",
			function: testWindowPreviousFail,
		},
		{
			scenario: "load next component",
			function: testWindowNextSuccess,
		},
		{
			scenario: "load next component fails",
			function: testWindowNextFail,
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
}

func testWindowIsDecorated(t *testing.T, w app.Window) {
	if base := w.Base(); base == w {
		t.Error("window is not decorated")
	}
}

func testWindowLoadSuccess(t *testing.T, w app.Window) {
	if err := w.Load("tests.hello"); err != nil {
		t.Fatal(err)
	}
}

func testWindowLoadFail(t *testing.T, w app.Window) {
	err := w.Load("tests.abracadabra")
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testWindowRenderSuccess(t *testing.T, w app.Window) {
	if err := w.Load("tests.hello"); err != nil {
		t.Fatal(err)
	}

	compo := w.Component()
	if compo == nil {
		t.Fatal("component is nil")
	}

	hello := compo.(*Hello)
	hello.Name = "Maxence"

	if err := w.Render(hello); err != nil {
		t.Fatal(err)
	}
}

func testWindowRenderFail(t *testing.T, w app.Window) {
	if err := w.Load("tests.hello"); err != nil {
		t.Fatal(err)
	}

	compo := w.Component()
	if compo == nil {
		t.Fatal("component is nil")
	}

	hello := compo.(*Hello)
	hello.TmplErr = true

	err := w.Render(hello)
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testWindowReloadSuccess(t *testing.T, w app.Window) {
	if err := w.Load("tests.hello"); err != nil {
		t.Fatal(err)
	}

	if err := w.Reload(); err != nil {
		t.Fatal(err)
	}
}

func testWindowReloadFail(t *testing.T, w app.Window) {
	err := w.Reload()
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testWindowPreviousSuccess(t *testing.T, w app.Window) {
	if err := w.Load("tests.hello"); err != nil {
		t.Fatal(err)
	}

	if w.CanPrevious() {
		t.Fatal("can previous is true")
	}

	if err := w.Load("tests.world"); err != nil {
		t.Fatal(err)
	}

	if !w.CanPrevious() {
		t.Fatal("can previous is false")
	}

	if err := w.Previous(); err != nil {
		t.Fatal(err)
	}
}

func testWindowPreviousFail(t *testing.T, w app.Window) {
	if err := w.Load("tests.hello"); err != nil {
		t.Fatal(err)
	}

	if w.CanPrevious() {
		t.Fatal("can previous is true")
	}

	err := w.Previous()
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testWindowNextSuccess(t *testing.T, w app.Window) {
	if err := w.Load("tests.hello"); err != nil {
		t.Fatal(err)
	}

	if w.CanNext() {
		t.Fatal("can next is true")
	}

	if err := w.Load("tests.world"); err != nil {
		t.Fatal(err)
	}

	if w.CanNext() {
		t.Fatal("can next is true")
	}

	if err := w.Previous(); err != nil {
		t.Fatal(err)
	}

	if !w.CanNext() {
		t.Fatal("can next is false")
	}

	if err := w.Next(); err != nil {
		t.Fatal(err)
	}
}

func testWindowNextFail(t *testing.T, w app.Window) {
	if err := w.Load("tests.hello"); err != nil {
		t.Fatal(err)
	}

	if w.CanNext() {
		t.Fatal("can next is true")
	}

	err := w.Next()
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
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
