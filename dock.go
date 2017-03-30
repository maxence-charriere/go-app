package app

// Docker represents a context with dock specific interactions.
type Docker interface {
	Contexter
	SetIcon(path string)
	SetBadge(v interface{})
}
