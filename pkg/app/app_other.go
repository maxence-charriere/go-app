// +build !wasm

package app

func run()                              {}
func render(Compo)                      {}
func reload()                           {}
func bind(msg string, c Compo) *Binding { panic("no wasm") }
