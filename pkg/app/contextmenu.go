// +build wasm

package app

import (
	"fmt"
	"strconv"
	"strings"
	"syscall/js"
)

func init() {
	Import(&ContextMenu{})
}

// ContextMenu is a component that displays a context menu.
type ContextMenu struct {
	Items   []MenuItem
	Visible bool
}

func (m *ContextMenu) OnMount() {
	Bind("__app.NewContextMenu", m).Do(m.new)
}

// Render returns the markup that describes the context menu.
func (m *ContextMenu) Render() string {
	return `
<div id="App_ContextMenuBackground" onclick="Hide">
	<div id="App_ContextMenu">
		{{range $idx, $v := .Items}}
	
		{{if .Separator}}
		<div class="App_MenuItemSeparator"></div>
		{{else}}
		<button class="App_MenuItem" onclick="Items.{{$idx}}.OnClick" {{if .Disabled}}disabled{{end}}>
			<div class="App_MenuItemLabel">{{.Label}}</div>
			<div class="App_MenuItemKeys">{{.Keys}}</div>
		</button>
		{{end}}

		{{end}}
	</div>
</div>
	`
}

func (m *ContextMenu) new(items []MenuItem) {
	for i := range items {
		if items[i].OnClick == nil {
			items[i].Disabled = true
		}

		items[i].Keys = convertKeys(items[i].Keys)
	}

	m.Items = items
	m.Visible = true
	Render(m)

	UI(m.Show)
}

func (m *ContextMenu) Show() {
	bg := js.Global().
		Get("document").
		Call("getElementById", "App_ContextMenuBackground")
	bg.Get("style").Set("display", "block")

	menu := js.Global().
		Get("document").
		Call("getElementById", "App_ContextMenu")
	menuWidth := menu.Get("offsetWidth").Int()
	menuHeight := menu.Get("offsetHeight").Int()

	winWidth := getWindowWidth()
	winHeight := getWindowHeight()

	x := cursorX
	if x+menuWidth > winWidth {
		x = winWidth - menuWidth
	}

	y := cursorY
	if y+menuHeight > winHeight {
		y = winHeight - menuHeight
	}

	menu.Get("style").Set("left", strconv.Itoa(x)+"px")
	menu.Get("style").Set("top", strconv.Itoa(y)+"px")
}

func (m *ContextMenu) Hide(s, e js.Value) {
	fmt.Println("gonna hide")
	bg := js.Global().
		Get("document").
		Call("getElementById", "App_ContextMenuBackground")
	bg.Get("style").Set("display", "none")
}

func convertKeys(k string) string {
	k = strings.ToLower(k)

	switch js.Global().Get("navigator").Get("platform").String() {
	case "Macintosh", "MacIntel", "MacPPC", "Mac68K":
		k = strings.Replace(k, "cmdorctrl", "⌘", -1)
		k = strings.Replace(k, "cmd", "⌘", -1)
		k = strings.Replace(k, "command", "⌘", -1)
		k = strings.Replace(k, "ctrl", "⌃", -1)
		k = strings.Replace(k, "control", "⌃", -1)
		k = strings.Replace(k, "alt", "⌥", -1)
		k = strings.Replace(k, "option", "⌥", -1)
		k = strings.Replace(k, "shift", "⇧", -1)
		k = strings.Replace(k, "capslock", "⇪", -1)
		k = strings.Replace(k, "del", "⌫", -1)
		k = strings.Replace(k, "delete", "⌫", -1)
		k = strings.Replace(k, "+", "", -1)

	case "Windows", "Win32":
		k = strings.Replace(k, "cmdorctrl", "ctrl", -1)
		k = strings.Replace(k, "cmd", "win", -1)
		k = strings.Replace(k, "command", "win", -1)
		k = strings.Replace(k, "control", "ctrl", -1)

	default:
		k = ""
	}

	return k
}

func getWindowWidth() int {
	w := js.Global().Get("innerWidth")
	if !w.Truthy() {
		w = js.Global().
			Get("document").
			Get("documentElement").
			Get("clientWidth")
	}
	if !w.Truthy() {
		w = js.Global().
			Get("document").
			Get("body").
			Get("clientWidth")
	}
	if w.Type() != js.TypeNumber {
		return 0
	}
	return w.Int()
}

func getWindowHeight() int {
	h := js.Global().Get("innerHeight")
	if !h.Truthy() {
		h = js.Global().
			Get("document").
			Get("documentElement").
			Get("clientHeight")
	}
	if !h.Truthy() {
		h = js.Global().
			Get("document").
			Get("body").
			Get("clientHeight")
	}
	if h.Type() != js.TypeNumber {
		return 0
	}
	return h.Int()
}
