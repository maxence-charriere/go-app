// +build !wasm

package app

import "net/url"

func run()                                       { panic("no wasm") }
func render(Compo)                               { panic("no wasm") }
func reload()                                    { panic("no wasm") }
func bind(msg string, c Compo) *Binding          { panic("no wasm") }
func windowSize() (w, h int)                     { panic("no wasm") }
func navigate(rawurl string, updateHistory bool) { panic("no wasm") }
func locationURL() *url.URL                      { panic("no wasm") }
