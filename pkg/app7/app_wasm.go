package app

import "syscall/js"

var (
	window = &browserWindow{value: value{Value: js.Global()}}
)
