package app

import (
	"net/url"
	"reflect"
	"syscall/js"

	"github.com/maxence-charriere/go-app/v7/pkg/errors"
)

type value struct {
	js.Value
}

func (v value) Call(m string, args ...interface{}) Value {
	for i, a := range args {
		switch a := a.(type) {
		case Wrapper:
			args[i] = jsval(a.JSValue())
		}
	}

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
	return val(v.Value.New(args...))
}

func (v value) Type() Type {
	return Type(v.Value.Type())
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
	callback := makeJsEventHandler(body, h)
	w.Call("addEventListener", event, callback)

	return func() {
		w.Call("removeEventListener", event, callback)
		callback.Release()
	}
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
