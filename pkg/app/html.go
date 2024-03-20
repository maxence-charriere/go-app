package app

// HTML provides an interface for representing HTML elements within the
// application.
type HTML interface {
	UI

	// Tag retrieves the name of the HTML tag that the element represents.
	Tag() string

	// XMLNamespace fetches the XML namespace associated with the HTML element.
	// This is relevant for elements like SVG which might have a different
	// namespace.
	XMLNamespace() string

	// SelfClosing determines whether the HTML element is self-closing.
	// For elements like <img> or <br> which don't have closing tags, this
	// method returns true. Otherwise, it returns false.
	SelfClosing() bool

	depth() uint
	attrs() attributes
	setAttrs(attributes) HTML
	events() eventHandlers
	setEvents(eventHandlers) HTML
	setDepth(uint) UI
	setJSElement(Value) HTML
	parent() UI
	body() []UI
	setBody([]UI) HTML
}

type htmlElement struct {
	tag           string
	xmlns         string
	treeDepth     uint
	isSelfClosing bool
	jsElement     Value
	attributes    attributes
	eventHandlers eventHandlers
	parentElement UI
	children      []UI
}

func (e *htmlElement) JSValue() Value {
	return e.jsElement
}

func (e *htmlElement) Mounted() bool {
	return e.jsElement != nil
}

func (e *htmlElement) Tag() string {
	return e.tag
}

func (e *htmlElement) XMLNamespace() string {
	return e.xmlns
}

func (e *htmlElement) SelfClosing() bool {
	return e.isSelfClosing
}

func (e *htmlElement) depth() uint {
	return e.treeDepth
}

func (e *htmlElement) attrs() attributes {
	return e.attributes
}

func (e *htmlElement) setAttr(name string, value any) {
	if e.attributes == nil {
		e.attributes = make(attributes)
	}
	e.attributes.Set(name, value)
}

func (e *htmlElement) events() eventHandlers {
	return e.eventHandlers
}

func (e *htmlElement) setEventHandler(event string, h EventHandler, options ...EventOption) {
	if e.eventHandlers == nil {
		e.eventHandlers = make(eventHandlers)
	}
	e.eventHandlers.Set(event, h, options...)
}

func (e *htmlElement) parent() UI {
	return e.parentElement
}

func (e *htmlElement) body() []UI {
	return e.children
}
