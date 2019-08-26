// +build js

package app

import "syscall/js"

// MenuItem represents a menu item.
type MenuItem struct {
	Disabled  bool
	Keys      string
	Icon      string
	Label     string
	OnClick   func(s, e js.Value)
	Separator bool
}
