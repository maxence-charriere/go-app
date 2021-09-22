package app

import (
	"net/url"
	"reflect"
	"syscall/js"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

var (
	window = &browserWindow{value: value{Value: js.Global()}}
)

type value struct {
	js.Value
}

func (v value) Call(m string, args ...interface{}) Value {
	args = cleanArgs(args...)
	return val(v.Value.Call(m, args...))
}

func (v value) Get(p string) Value {
	return val(v.Value.Get(p))
}

func (v value) Set(p string, x interface{}) {
	if wrapper, ok := x.(Wrapper); ok {
		x = jsval(wrapper.JSValue())
	}
	v.Value.Set(p, x)
}

func (v value) Index(i int) Value {
	return val(v.Value.Index(i))
}

func (v value) InstanceOf(t Value) bool {
	return v.Value.InstanceOf(jsval(t))
}

func (v value) Invoke(args ...interface{}) Value {
	return val(v.Value.Invoke(args...))
}

func (v value) JSValue() Value {
	return v
}

func (v value) New(args ...interface{}) Value {
	args = cleanArgs(args...)
	return val(v.Value.New(args...))
}

func (v value) Type() Type {
	return Type(v.Value.Type())
}

func (v value) getAttr(k string) string {
	return v.Call("getAttribute", k).String()
}

func (v value) setAttr(k, val string) {
	v.Call("setAttribute", k, val)
}

func (v value) delAttr(k string) {
	v.Call("removeAttribute", k)
}

func (v value) firstChild() Value {
	return v.Get("firstChild")
}

func (v value) appendChild(c Wrapper) {
	v.Call("appendChild", c)
}

func (v value) replaceChild(new, old Wrapper) {
	v.Call("replaceChild", new, old)
}

func (v value) removeChild(c Wrapper) {
	v.Call("removeChild", c)
}

func (v value) firstElementChild() Value {
	return v.Get("firstElementChild")
}

func (v value) addEventListener(event string, fn Func) {
	v.Call("addEventListener", event, fn)
}

func (v value) removeEventListener(event string, fn Func) {
	v.Call("removeEventListener", event, fn)
}

func (v value) setNodeValue(val string) {
	v.Set("nodeValue", val)
}

func (v value) setInnerHTML(val string) {
	v.Set("innerHTML", val)
}

func (v value) setInnerText(val string) {
	v.Set("innerText", val)
}

func null() Value {
	return val(js.Null())
}

func undefined() Value {
	return val(js.Undefined())
}

func valueOf(x interface{}) Value {
	switch t := x.(type) {
	case value:
		x = t.Value

	case function:
		x = t.fn

	case *browserWindow:
		x = t.Value

	case Event:
		return valueOf(t.Value)
	}

	return val(js.ValueOf(x))
}

type function struct {
	value
	fn js.Func
}

func (f function) Release() {
	f.fn.Release()
}

func funcOf(fn func(this Value, args []Value) interface{}) Func {
	f := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		wargs := make([]Value, len(args))
		for i, a := range args {
			wargs[i] = val(a)
		}

		return fn(val(this), wargs)
	})

	return function{
		value: value{Value: f.Value},
		fn:    f,
	}
}

type browserWindow struct {
	value

	body    UI
	cursorX int
	cursorY int
}

func (w *browserWindow) URL() *url.URL {
	rawurl := w.
		Get("location").
		Get("href").
		String()

	u, _ := url.Parse(rawurl)
	return u
}

func (w *browserWindow) Size() (width int, height int) {
	getSize := func(axis string) int {
		size := w.Get("inner" + axis)
		if !size.Truthy() {
			size = w.
				Get("document").
				Get("documentElement").
				Get("client" + axis)
		}
		if !size.Truthy() {
			size = w.
				Get("document").
				Get("body").
				Get("client" + axis)
		}
		if size.Type() != TypeNumber {
			return 0
		}
		return size.Int()
	}

	return getSize("Width"), getSize("Height")
}

func (w *browserWindow) CursorPosition() (x, y int) {
	return w.cursorX, w.cursorY
}

func (w *browserWindow) setCursorPosition(x, y int) {
	w.cursorX = x
	w.cursorY = y
}

func (w *browserWindow) GetElementByID(id string) Value {
	return w.Get("document").Call("getElementById", id)
}

func (w *browserWindow) ScrollToID(id string) {
	if elem := w.GetElementByID(id); elem.Truthy() {
		elem.Call("scrollIntoView")
	}
}

func (w *browserWindow) AddEventListener(event string, h EventHandler) func() {
	callback := makeJsEventHandler(w.body, func(ctx Context, e Event) {
		h(ctx, e)

		// Trigger children components updates:
		if len(w.body.children()) == 0 {
			return
		}
		compo, ok := w.body.children()[0].(Composer)
		if !ok {
			return
		}
		ctx.Dispatcher().Dispatch(Dispatch{
			Mode:   Update,
			Source: compo,
		})
	})
	w.addEventListener(event, callback)

	return func() {
		w.removeEventListener(event, callback)
		callback.Release()
	}
}

func (w *browserWindow) setBody(body UI) {
	w.body = body
}

func (w *browserWindow) createElement(tag string) (Value, error) {
	v := w.Get("document").Call("createElement", tag)
	if !v.Truthy() {
		return nil, errors.New("creating element failed").
			Tag("reason", "create javascript element returned nil").
			Tag("tag", tag)
	}
	return v, nil
}

func (w *browserWindow) createTextNode(v string) Value {
	return w.Get("document").Call("createTextNode", v)
}

func (w *browserWindow) addHistory(u *url.URL) {
	w.Get("history").Call("pushState", nil, "", u.String())
	lastURLVisited = u
}

func (w *browserWindow) replaceHistory(u *url.URL) {
	w.Get("history").Call("replaceState", nil, "", u.String())
	lastURLVisited = u
}

func val(v js.Value) Value {
	return value{Value: v}
}

func jsval(v Value) js.Value {
	switch v := v.(type) {
	case value:
		return v.Value

	case function:
		return v.Value

	case *browserWindow:
		return v.Value

	case Event:
		return jsval(v.Value)

	default:
		Log("%s", errors.New("syscall/js value conversion failed").
			Tag("type", reflect.TypeOf(v)),
		)
		return js.Undefined()
	}
}

// JSValue returns the underlying syscall/js value of the given Javascript
// value.
func JSValue(v Value) js.Value {
	return jsval(v)
}

func copyBytesToGo(dst []byte, src Value) int {
	return js.CopyBytesToGo(dst, jsval(src))
}

func copyBytesToJS(dst Value, src []byte) int {
	return js.CopyBytesToJS(jsval(dst), src)
}

func cleanArgs(args ...interface{}) []interface{} {
	for i, a := range args {

		args[i] = cleanArg(a)
	}

	return args
}

func cleanArg(v interface{}) interface{} {
	switch v := v.(type) {
	case map[string]interface{}:
		m := make(map[string]interface{}, len(v))
		for key, val := range v {
			m[key] = cleanArg(val)
		}
		return m

	case []interface{}:
		s := make([]interface{}, len(v))
		for i, val := range v {
			s[i] = cleanArgs(val)
		}

	case Wrapper:
		return jsval(v.JSValue())
	}

	return v

}
