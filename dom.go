package app

// DOMNode is the interface that describes a DOM node.
// DOM have a limited support.
type DOMNode interface {
	// The node identifier.
	ID() string

	// The indentifier of the component where the node is mounted.
	// It is set only when the node is mounted.
	CompoID() string

	// The identifier of the control where the node is mounted.
	// It is set only when the node is mounted.
	ControlID() string

	// The parent node.
	Parent() DOMNode
}

// DOMElem is the interface that describes a DOM element node.
type DOMElem interface {
	DOMNode

	// The element tag name.
	TagName() string

	// The element attributes.
	Attrs() map[string]string

	// The children nodes.
	Children() []DOMNode
}

// DOMText is the interface that describes a DOM text node.
type DOMText interface {
	DOMNode

	// The text contained in the node.
	Text() string
}

// DOMCompo is the interface that describes a DOM component node.
type DOMCompo interface {
	DOMNode

	// The component name.
	Name() string

	// The compenent fields.
	Fields() map[string]string
}
