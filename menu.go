package app

// MenuItem represents a menu item.
type MenuItem struct {
	Disabled  bool
	Keys      string
	Icon      string
	Label     string
	OnClick   func()
	Separator bool
}
