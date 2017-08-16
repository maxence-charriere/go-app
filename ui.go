package app

type Window interface{}

type WindowConfig struct{}

type MenuBar interface{}

type Dock interface {
	SetIcon(name string)

	SetBadge(v interface{})
}
