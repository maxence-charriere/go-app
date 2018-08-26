package app

// Node is the interface that describes a DOM node.
// Node have a limited support.
type Node interface {
	// The node identifier.
	ID() string

	// The parent node.
	Parent() Node
}
