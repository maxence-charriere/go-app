package app

// DOM is the interface that describes a document object model store that
// manages node states.
type DOM interface {
	// Returns the component with the given identifier.
	ComponentByID(id string) (Component, error)

	// Create or update the nodes of the given component.
	Render(Component) (changes []DOMChange, err error)
}

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

// DOMChange represents a change to perform on a DOM node.
// It is used by backends to synchronize a node with its local representation.
type DOMChange struct {
	Type  DOMChangeType
	Value interface{}
}

// DOMChangeType represents a DOM change type.
type DOMChangeType int

// Constants that enumerates DOM change types.
const (
	DOMNoChanges DOMChangeType = iota
	DOMUpdateAttrs
	DOMAppendChild
	DOMInsertBefore
	DOMReplaceChild
	DOMRemoveChild
)
