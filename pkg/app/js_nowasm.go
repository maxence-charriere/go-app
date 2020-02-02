// +build !wasm

package app

import "net/url"

type value struct{}

func (v value) Bool() bool {
	panic(errNoWasm)
}

func (v value) Call(m string, args ...interface{}) Value {
	panic(errNoWasm)
}

func (v value) Float() float64 {
	panic(errNoWasm)
}

func (v value) Get(p string) Value {
	panic(errNoWasm)
}

func (v value) Index(i int) Value {
	panic(errNoWasm)
}

func (v value) InstanceOf(t Value) bool {
	panic(errNoWasm)
}

func (v value) Int() int {
	panic(errNoWasm)
}

func (v value) Invoke(args ...interface{}) Value {
	panic(errNoWasm)
}

func (v value) JSValue() Value {
	panic(errNoWasm)
}

func (v value) Length() int {
	panic(errNoWasm)
}

func (v value) New(args ...interface{}) Value {
	panic(errNoWasm)
}

func (v value) Set(p string, x interface{}) {
	panic(errNoWasm)
}

func (v value) SetIndex(i int, x interface{}) {
	panic(errNoWasm)
}

func (v value) String() string {
	panic(errNoWasm)
}

func (v value) Truthy() bool {
	panic(errNoWasm)
}

func (v value) Type() Type {
	panic(errNoWasm)
}

func null() Value {
	panic(errNoWasm)
}

func undefined() Value {
	panic(errNoWasm)
}

func valueOf(x interface{}) Value {
	panic(errNoWasm)
}

func funcOf(fn func(this Value, args []Value) interface{}) Func {
	panic(errNoWasm)
}

type browserWindow struct {
	value
}

func (w browserWindow) URL() *url.URL {
	panic(errNoWasm)
}

func (w browserWindow) Size() (width, height int) {
	panic(errNoWasm)
}

func (w browserWindow) CursorPosition() (x, y int) {
	panic(errNoWasm)
}

func (w browserWindow) setCursorPosition(x, y int) {
	panic(errNoWasm)
}
