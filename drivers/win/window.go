// +build windows

package win

import (
	"encoding/json"
	"path"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/core"
	"github.com/murlokswarm/app/internal/dom"
)

func newWindow(d *core.Driver) *core.Window {
	return &core.Window{
		ConvertHTMLPaths: convertHTMLPaths,
		DefaultWidth:     1280,
		DefaultHeight:    720,
		DOM: dom.Engine{
			Resources: resourcesDir,
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

func onWindowClose(w *core.Window, in map[string]interface{}) {
	w.Release()
}

func resourcesDir(p ...string) string {
	r := path.Join(p...)
	r = strings.TrimLeft(r, "/")
	return "ms-appx-web:///Resources/" + r
}

func convertHTMLPaths(paths []string) []string {
	convs := make([]string, len(paths))
	resources := driver.Resources()

	for i, p := range paths {
		p = strings.TrimPrefix(p, resources)
		p = strings.Replace(p, "\\", "/", -1)
		convs[i] = resourcesDir(p)
	}

	spew.Dump(paths)
	spew.Dump(convs)

	return convs
}
