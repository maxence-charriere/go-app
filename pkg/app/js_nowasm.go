// +build !wasm

package app

import (
	"net/url"
	"runtime"

	"github.com/maxence-charriere/go-app/v7/pkg/errors"
)

var (
	errNoWasm = errors.New("unsupported instruction").
		Tag("required-architecture", "wasm").
		Tag("current-architecture", runtime.GOARCH)
)

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

func (v value) IsNaN() bool {
	panic(errNoWasm)
}

func (v value) IsNull() bool {
	panic(errNoWasm)
}

func (v value) IsUndefined() bool {
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
	return true
}

func (v value) Type() Type {
	panic(errNoWasm)
}

func (v value) getAttr(k string) string {
	return ""
}

func (v value) setAttr(k, val string) {
}

func (v value) delAttr(k string) {
}

func (v value) firstChild() Value {
	return value{}
}

func (v value) appendChild(c Wrapper) {
}

func (v value) replaceChild(new, old Wrapper) {
}

func (v value) removeChild(c Wrapper) {
}

func (v value) addEventListener(event string, fn Func) {
}

func (v value) removeEventListener(event string, fn Func) {
}

func (v value) setNodeValue(val string) {
}

func (v value) setInnerHTML(val string) {
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

type function struct {
	value
}

func (f function) Release() {
}

func funcOf(fn func(this Value, args []Value) interface{}) Func {
	return function{value: value{}}
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

func (w *browserWindow) GetElementByID(id string) Value {
	panic(errNoWasm)
}

func (w *browserWindow) ScrollToID(id string) {
	panic(errNoWasm)
}

func (w *browserWindow) AddEventListener(event string, h EventHandler) func() {
	panic(errNoWasm)
}

func (w *browserWindow) createElement(tag string) (Value, error) {
	return value{}, nil
}

func (w *browserWindow) createTextNode(v string) Value {
	return value{}
}

func copyBytesToGo(dst []byte, src Value) int {
	panic(errNoWasm)
}

func copyBytesToJS(dst Value, src []byte) int {
	panic(errNoWasm)
}
