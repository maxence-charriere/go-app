// +build darwin,amd64

package mac

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/core"
	"github.com/murlokswarm/app/internal/dom"
)

func newWindow(d *core.Driver) *core.Window {
	return &core.Window{
		DefaultWidth:  1280,
		DefaultHeight: 720,
		DOM: dom.Engine{
			Resources: d.Resources,
			AttrTransforms: []dom.Transform{
				dom.JsToGoHandler,
				dom.HrefCompoFmt,
			},
		},
		Driver: d,
	}
}

func handleWindow(h func(w *core.Window, in map[string]interface{})) core.GoHandler {
	return func(in map[string]interface{}) {
		e := driver.Elems.GetByID(in["ID"].(string))
		if e.Err() == app.ErrElemNotSet {
			return
		}

		h(e.(*core.Window), in)
	}
}

func onWindowNavigate(w *core.Window, in map[string]interface{}) {
	e := app.ElemByCompo(w.Compo())

	e.WhenWindow(func(w app.Window) {
		w.Load(in["URL"].(string))
	})
}

func onWindowCallback(w *core.Window, in map[string]interface{}) {
	mappingStr := in["Mapping"].(string)

	var m dom.Mapping
	if err := json.Unmarshal([]byte(mappingStr), &m); err != nil {
		app.Logf("window callback failed: %s", err)
		return
	}

	if m.Override == "Files" {
		data, _ := json.Marshal(driver.droppedFiles)
		driver.droppedFiles = nil

		m.JSONValue = strings.Replace(
			m.JSONValue,
			`"FileOverride":"xxx"`,
			fmt.Sprintf(`"Files":%s`, data),
			1,
		)
	}

	c, err := w.DOM.CompoByID(m.CompoID)
	if err != nil {
		app.Logf("window callback failed: %s", err)
		return
	}

	var f func()
	if f, err = m.Map(c); err != nil {
		app.Logf("window callback failed: %s", err)
		return
	}

	if f != nil {
		f()
		return
	}

	app.Render(c)
}

func onWindowAlert(w *core.Window, in map[string]interface{}) {
	app.Logf("%s", in["Alert"])
}

func onWindowMove(w *core.Window, in map[string]interface{}) {
	x := in["X"].(float64)
	y := in["Y"].(float64)

	w.SetPosition(x, y)
}

func onWindowResize(w *core.Window, in map[string]interface{}) {
	width := in["Width"].(float64)
	height := in["Height"].(float64)

	w.SetSize(width, height)
}

func onWindowFocus(w *core.Window, in map[string]interface{}) {
	w.SetIsFocus(true)
}

func onWindowBlur(w *core.Window, in map[string]interface{}) {
	w.SetIsFocus(false)
}

func onWindowFullScreen(w *core.Window, in map[string]interface{}) {
	w.SetIsFullScreen(true)
}

func onWindowExitFullScreen(w *core.Window, in map[string]interface{}) {
	w.SetIsFullScreen(false)
}

func onWindowMinimize(w *core.Window, in map[string]interface{}) {
	w.SetIsMinimized(true)
}

func onWindowDeminimize(w *core.Window, in map[string]interface{}) {
	w.SetIsMinimized(false)
}

func onWindowClose(w *core.Window, in map[string]interface{}) {
	w.Release()
}
