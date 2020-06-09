// +build !wasm

package app

import "net/url"

type value struct{}

func (v value) Bool() bool {
	panicNoWasm()
	return false
}

func (v value) Call(m string, args ...interface{}) Value {
	panicNoWasm()
	return nil
}

func (v value) Float() float64 {
	panicNoWasm()
	return 0
}

func (v value) Get(p string) Value {
	panicNoWasm()
	return nil
}

func (v value) Index(i int) Value {
	panicNoWasm()
	return nil
}

func (v value) InstanceOf(t Value) bool {
	panicNoWasm()
	return false
}

func (v value) Int() int {
	panicNoWasm()
	return 0
}

func (v value) Invoke(args ...interface{}) Value {
	panicNoWasm()
	return nil
}

func (v value) IsNaN() bool {
	panicNoWasm()
	return false
}

func (v value) IsNull() bool {
	panicNoWasm()
	return false
}

func (v value) IsUndefined() bool {
	panicNoWasm()
	return false
}

func (v value) JSValue() Value {
	panicNoWasm()
	return nil
}

func (v value) Length() int {
	panicNoWasm()
	return 0
}

func (v value) New(args ...interface{}) Value {
	panicNoWasm()
	return nil
}

func (v value) Set(p string, x interface{}) {
	panicNoWasm()
}

func (v value) SetIndex(i int, x interface{}) {
	panicNoWasm()
}

func (v value) String() string {
	panicNoWasm()
	return ""
}

func (v value) Truthy() bool {
	panicNoWasm()
	return false
}

func (v value) Type() Type {
	panicNoWasm()
	return TypeUndefined
}

func null() Value {
	panicNoWasm()
	return nil
}

func undefined() Value {
	panicNoWasm()
	return nil
}

func valueOf(x interface{}) Value {
	panicNoWasm()
	return nil
}

func funcOf(fn func(this Value, args []Value) interface{}) Func {
	panicNoWasm()
	return nil
}

type browserWindow struct {
	value
}

func (w browserWindow) URL() *url.URL {
	panicNoWasm()
	return nil
}

func (w browserWindow) Size() (width, height int) {
	panicNoWasm()
	return 0, 0
}

func (w browserWindow) CursorPosition() (x, y int) {
	panicNoWasm()
	return 0, 0
}

func (w browserWindow) setCursorPosition(x, y int) {
	panicNoWasm()
}

func (w *browserWindow) GetElementByID(id string) Value {
	panicNoWasm()
	return nil
}

func (w *browserWindow) ScrollToID(id string) {
	panicNoWasm()
}

func (w *browserWindow) AddEventListener(event string, h EventHandler) func() {
	panicNoWasm()
	return nil
}

func copyBytesToGo(dst []byte, src Value) int {
	panicNoWasm()
	return 0
}

func copyBytesToJS(dst Value, src []byte) int {
	panicNoWasm()
	return 0
}

func makeEventHandler(h EventHandler) Func {
	panicNoWasm()
	return nil
}
