package app

// Node is the interface that describes an element that is used to represent a
// user interface.
type Node interface {
	Kind() Kind
}

// UI is the interface that describes a user interface element such as
// components and HTML elements.
type UI interface {
	Node

	// JSValue returns the javascript value linked to the element.
	JSValue() Value

	// Reports whether the element is mounted.
	Mounted() bool

	parent() UI
	setParent(UI)
	children() []UI
	appendChild(UI)
	removeChild(UI)
	mount() error
	update(UI) error
	dismount()
}

// Kind represents the specific kind of a user interface element.
type Kind uint

func (k Kind) String() string {
	switch k {
	case SimpleText:
		return "text"

	case HTML:
		return "html"

	case Component:
		return "component"

	case Selector:
		return "selector"

	default:
		return "undefined"
	}
}

const (
	// UndefinedElem represents an undefined UI element.
	UndefinedElem Kind = iota

	// SimpleText represents a simple text element.
	SimpleText

	// HTML represents an HTML element.
	HTML

	// Component represents a customized, independent and reusable UI element.
	Component

	// Selector represents an element that is used to select a subset of
	// elements within a given list.
	Selector
)
