package app

import (
	"net/url"
	"syscall/js"

	"github.com/maxence-charriere/go-app/v10/pkg/errors"
)

type value struct {
	jsValue
}

func (v value) Call(m string, args ...any) Value {
	res := v.jsValue.Call(m, syscallJSArgs(args)...)
	return ValueOf(res)
}

func (v value) Delete(p string) {
	v.jsValue.Delete(p)
}

func (v value) Equal(w Value) bool {
	return v.jsValue.Equal(syscalJSValueOf(w))
}

func (v value) Get(p string) Value {
	return ValueOf(v.jsValue.Get(p))
}

func (v value) Set(p string, x any) {
	v.jsValue.Set(p, syscalJSValueOf(x))
}

func (v value) Index(i int) Value {
	return ValueOf(v.jsValue.Index(i))
}

func (v value) InstanceOf(t Value) bool {
	return v.jsValue.InstanceOf(syscalJSValueOf(t))
}

func (v value) Invoke(args ...any) Value {
	res := v.jsValue.Invoke(syscallJSArgs(args)...)
	return ValueOf(res)
}

func (v value) JSValue() Value {
	return v
}

func (v value) New(args ...any) Value {
	res := v.jsValue.New(syscallJSArgs(args)...)
	return ValueOf(res)
}

func (v value) Type() Type {
	return Type(v.jsValue.Type())
}

func (v value) Release() {
	if function, ok := v.jsValue.(jsFunc); ok {
		function.Release()
	}
}

func (v value) Then(f func(Value)) {
	release := func() {}

	then := FuncOf(func(this Value, args []Value) any {
		var arg Value
		if len(args) > 0 {
			arg = args[0]
		}

		f(arg)
		release()
		return nil
	})

	release = then.Release
	v.Call("then", then)
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

func (v value) addEventListener(event string, fn Func, options map[string]any) {
	if len(options) == 0 {
		v.Call("addEventListener", event, fn)
		return
	}
	v.Call("addEventListener", event, fn, options)
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
	return valueOf(js.Null())
}

func undefined() Value {
	return valueOf(js.Undefined())
}

func valueOf(v any) Value {
	switch v := v.(type) {
	case jsValue:
		return value{jsValue: v}

	case Wrapper:
		return v.JSValue()

	default:
		return value{jsValue: js.ValueOf(v)}
	}
}

func funcOf(function func(this Value, args []Value) any) Func {
	return value{
		jsValue: js.FuncOf(func(this js.Value, args []js.Value) any {
			goArgs := make([]Value, len(args))
			for i, arg := range args {
				goArgs[i] = ValueOf(arg)
			}

			return syscalJSValueOf(function(ValueOf(this), goArgs))
		}),
	}
}

type jsValue interface {
	Bool() bool
	Call(string, ...any) js.Value
	Delete(string)
	Equal(js.Value) bool
	Float() float64
	Get(string) js.Value
	Index(int) js.Value
	InstanceOf(js.Value) bool
	Int() int
	Invoke(...any) js.Value
	IsNaN() bool
	IsNull() bool
	IsUndefined() bool
	Length() int
	New(...any) js.Value
	Set(string, any)
	SetIndex(int, any)
	String() string
	Truthy() bool
	Type() js.Type
}

type jsFunc interface {
	jsValue
	Release()
}

type jsError interface {
	jsValue
	Error() string
}

func syscallJSArgs(v []any) []any {
	if len(v) == 0 {
		return nil
	}

	s := make([]any, len(v))
	for i, value := range v {
		s[i] = syscalJSValueOf(value)
	}
	return s
}

func syscalJSValueOf(v any) js.Value {
	switch v := v.(type) {
	case value:
		return syscalJSValueOf(v.jsValue)

	case Wrapper:
		return syscalJSValueOf(v.JSValue())

	case map[string]any:
		m := make(map[string]any, len(v))
		for key, value := range v {
			m[key] = syscalJSValueOf(value)
		}
		return js.ValueOf(m)

	case []any:
		s := make([]any, len(v))
		for i, value := range v {
			s[i] = syscalJSValueOf(value)
		}
		return js.ValueOf(s)

	default:
		return js.ValueOf(v)
	}
}

type browserWindow struct {
	value

	body    UI
	cursorX int
	cursorY int
}

func newBrowserWindow() *browserWindow {
	return &browserWindow{value: value{jsValue: js.Global()}}
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

func (w *browserWindow) setBody(body UI) {
	w.body = body
}

func (w *browserWindow) createElement(tag, xmlns string) (Value, error) {
	var element Value
	if xmlns == "" {
		element = w.Get("document").Call("createElement", tag)
	} else {
		element = w.Get("document").Call("createElementNS", xmlns, tag)
	}

	if !element.Truthy() {
		return nil, errors.New("creating javascript element failed").
			WithTag("tag", tag).
			WithTag("xmlns", xmlns)
	}
	return element, nil
}

func (w *browserWindow) createTextNode(v string) Value {
	return w.Get("document").Call("createTextNode", v)
}

func (w *browserWindow) addHistory(u *url.URL) {
	u.Scheme = w.URL().Scheme
	u.Host = w.URL().Host
	w.Get("history").Call("pushState", nil, "", u.String())
}

func (w *browserWindow) replaceHistory(u *url.URL) {
	u.Scheme = w.URL().Scheme
	u.Host = w.URL().Host
	w.Get("history").Call("replaceState", nil, "", u.String())
}

func copyBytesToGo(dst []byte, src Value) int {
	return js.CopyBytesToGo(dst, syscalJSValueOf(src))
}

func copyBytesToJS(dst Value, src []byte) int {
	return js.CopyBytesToJS(syscalJSValueOf(dst), src)
}
