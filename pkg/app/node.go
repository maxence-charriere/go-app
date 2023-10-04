package app

import (
	"context"
	"io"
	"net/url"
	"reflect"
	"strings"
	"sync"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

// UI is the interface that describes a user interface element such as
// components and HTML elements.
type UI interface {
	// JSValue returns the javascript value linked to the element.
	JSValue() Value

	// Reports whether the element is mounted.
	Mounted() bool

	name() string
	self() UI
	setSelf(UI)
	getDispatcher() Dispatcher
	getAttributes() attributes
	getEventHandlers() eventHandlers
	getParent() UI
	getChildren() []UI
	mount(Dispatcher) error
	dismount()
	canUpdateWith(UI) bool
	updateWith(UI) error
	onComponentEvent(any)
	html(w io.Writer)
	htmlWithIndent(w io.Writer, indent int)

	setParent(UI) UI
}

// FilterUIElems processes and returns a filtered list of the provided UI
// elements.
//
// Specifically, it:
// - Interprets and removes selector elements such as Condition and RangeLoop.
// - Eliminates nil elements and nil pointers.
// - Flattens and includes the children of recognized selector elements.
//
// This function is primarily intended for components that accept ui elements as
// variadic arguments or slice, such as the Body method of HTML elements.
func FilterUIElems(v ...UI) []UI {
	if len(v) == 0 {
		return nil
	}

	removeELemAt := func(i int) {
		copy(v[i:], v[i+1:])
		v[len(v)-1] = nil
		v = v[:len(v)-1]
	}

	var trailing []UI
	replaceElemAt := func(i int, elems ...UI) {
		trailing = append(trailing, v[i+1:]...)
		v = append(v[:i], elems...)
		v = append(v, trailing...)
		trailing = trailing[:0]
	}

	for i := len(v) - 1; i >= 0; i-- {
		elem := v[i]
		if elem == nil {
			removeELemAt(i)
		}
		if elemValue := reflect.ValueOf(elem); elemValue.Kind() == reflect.Pointer && elemValue.IsNil() {
			removeELemAt(i)
		}

		switch elem.(type) {
		case Condition, RangeLoop:
			replaceElemAt(i, elem.getChildren()...)
		}
	}

	return v
}

func mount(d Dispatcher, n UI) error {
	n.setSelf(n)
	return n.mount(d)
}

func dismount(n UI) {
	n.dismount()
	n.setSelf(nil)
}

func canUpdate(a, b UI) bool {
	a.setSelf(a)
	b.setSelf(b)
	return a.canUpdateWith(b)
}

func update(a, b UI) error {
	a.setSelf(a)
	b.setSelf(b)
	return a.updateWith(b)
}

// HTMLString return an HTML string representation of the given UI element.
func HTMLString(ui UI) string {
	var w strings.Builder
	PrintHTML(&w, ui)
	return w.String()
}

// HTMLStringWithIndent return an indented HTML string representation of the
// given UI element.
func HTMLStringWithIndent(ui UI) string {
	var w strings.Builder
	PrintHTMLWithIndent(&w, ui)
	return w.String()
}

// PrintHTML writes an HTML representation of the UI element into the given
// writer.
func PrintHTML(w io.Writer, ui UI) {
	if !ui.Mounted() {
		ui.setSelf(ui)
	}
	ui.html(w)
}

// PrintHTMLWithIndent writes an idented HTML representation of the UI element
// into the given writer.
func PrintHTMLWithIndent(w io.Writer, ui UI) {
	if !ui.Mounted() {
		ui.setSelf(ui)
	}
	ui.htmlWithIndent(w, 0)
}

// nodeManager orchestrates the lifecycle of UI elements, providing specialized
// mechanisms for mounting, dismounting, and updating nodes.
type nodeManager struct {
	// Ctx serves as a base context, utilized to derive child contexts providing
	// mechanisms to interact with the browser, manage pages, handle
	// concurrency, and facilitate component communication.
	Ctx context.Context

	// ResolveURL is used to transform attributes that hold URL values.
	ResolveURL attributeURLResolver

	// Page represents a reference to the current web page.
	Page Page

	// EmitHTMLEvent is called when a specific HTML event occurs on a UI
	// element. 'src' represents the source UI element triggering the event, and
	// 'f' is the callback to be executed in response.
	EmitHTMLEvent func(src UI, f func())

	initOnce sync.Once
}

func (m *nodeManager) init() {
	if m.Ctx == nil {
		m.Ctx = context.Background()
	}

	if m.ResolveURL == nil {
		m.ResolveURL = func(s string) string {
			return s
		}
	}

	if m.Page == nil {
		url, _ := url.Parse("https://goapp.dev")
		if IsServer {
			m.Page = &requestPage{
				url:                   url,
				resolveStaticResource: m.ResolveURL,
			}
		} else {
			m.Page = browserPage{
				url:                   url,
				resolveStaticResource: m.ResolveURL,
			}
		}
	}

	if m.EmitHTMLEvent == nil {
		m.EmitHTMLEvent = func(u UI, f func()) {
			f()
		}
	}
}

// Mount mounts a UI element based on its type and the specified depth. It
// returns the mounted UI element and any potential error during the process.
func (m *nodeManager) Mount(depth uint, v UI) (UI, error) {
	m.initOnce.Do(m.init)

	switch v := v.(type) {
	case *text:
		return m.mountText(depth, v)

	case HTML:
		return m.mountHTML(depth, v)

	case Composer:
		return m.mountComponent(depth, v)

	case *raw:
		return m.mountRawHTML(depth, v)

	default:
		return nil, errors.New("unsupported element").
			WithTag("type", reflect.TypeOf(v)).
			WithTag("depth", depth)
	}
}

func (m *nodeManager) mountText(depth uint, v *text) (UI, error) {
	if v.Mounted() {
		return nil, errors.New("text is already mounted").
			WithTag("parent-type", reflect.TypeOf(v.getParent())).
			WithTag("preview-value", previewText(v.value))
	}

	v.jsvalue = Window().createTextNode(v.value)
	return v, nil
}

func (m *nodeManager) mountHTML(depth uint, v HTML) (UI, error) {
	if v.Mounted() {
		return nil, errors.New("html element is already mounted").
			WithTag("parent-type", reflect.TypeOf(v.getParent())).
			WithTag("type", reflect.TypeOf(v)).
			WithTag("tag", v.Tag()).
			WithTag("depth", v.depth())
	}

	var jsElement Value
	switch v.(type) {
	case *htmlBody:
		jsElement = Window().Get("document").Get("body")

	default:
		jsElement, _ = Window().createElement(v.Tag(), v.XMLNamespace())
	}
	if IsClient && !jsElement.Truthy() {
		return nil, errors.New("creating js element failed").
			WithTag("type", reflect.TypeOf(v)).
			WithTag("tag", v.Tag()).
			WithTag("xmlns", v.XMLNamespace()).
			WithTag("depth", depth)
	}
	v = v.setJSElement(jsElement)
	m.mountHTMLAttributes(v)
	m.mountHTMLEventHandlers(v)

	v = v.setDepth(depth).(HTML)
	children := v.body()
	for i, child := range children {
		var err error
		if child, err = m.Mount(depth+1, child); err != nil {
			return nil, errors.New("mounting child failed").
				WithTag("type", reflect.TypeOf(v)).
				WithTag("tag", v.Tag()).
				WithTag("depth", depth).
				WithTag("index", i).
				Wrap(err)
		}
		child = child.setParent(v)
		children[i] = child
		v.JSValue().appendChild(child)
	}

	return v, nil
}

func (m *nodeManager) mountHTMLAttributes(v HTML) {
	for name, value := range v.attrs() {
		setJSAttribute(v.JSValue(), name, resolveAttributeURLValue(
			name,
			value,
			m.ResolveURL,
		))
	}
}

func (m *nodeManager) mountHTMLEventHandlers(v HTML) {
	events := v.events()
	for event, handler := range events {
		events[event] = m.mountHTMLEventHandler(v, handler)
	}
}

func (m *nodeManager) mountHTMLEventHandler(v HTML, handler eventHandler) eventHandler {
	event := handler.event

	jsHandler := FuncOf(func(this Value, args []Value) any {
		if len(args) != 0 {
			event := Event{Value: args[0]}
			trackMousePosition(event)
			handler.goHandler(m.MakeContext(v), event)
		}
		return nil
	})
	v.JSValue().addEventListener(event, jsHandler)

	return eventHandler{
		event:     event,
		scope:     handler.scope,
		goHandler: handler.goHandler,
		jsHandler: jsHandler,
		close: func() {
			v.JSValue().removeEventListener(event, jsHandler)
			jsHandler.Release()
		},
	}
}

func (m *nodeManager) mountComponent(depth uint, v Composer) (UI, error) {
	if !v.Mounted() {
		return nil, errors.New("component is already mounted").
			WithTag("parent-type", reflect.TypeOf(v.getParent())).
			WithTag("type", reflect.TypeOf(v)).
			WithTag("depth", v.depth())
	}

	// set self or alternative
	// change Compo.Mounted to not check for a mounted root mounted
	// change JSValue to return an empty js value
	// make a list of deprecated things.

	if mounter, ok := v.(Mounter); ok {
		mounter.OnMount(m.MakeContext(v))
	}

	root, err := m.renderComponent(v)
	if err != nil {
		return nil, errors.New("rendering component failed").
			WithTag("type", reflect.TypeOf(v)).
			WithTag("depth", v.depth()).
			Wrap(err)
	}

	panic("not completely implemented")
	return v, nil
}

func (m *nodeManager) renderComponent(v Composer) (UI, error) {
	rendering := FilterUIElems(v.Render())
	if len(rendering) == 0 {
		return nil, errors.New("render method does not returns a text, html element, or component")
	}
	return rendering[0], nil
}

func (m *nodeManager) mountRawHTML(depth uint, v *raw) (UI, error) {
	panic("not implemented")
}

// Dismount removes a UI element based on its type.
func (m *nodeManager) Dismount(v UI) {
	switch v := v.(type) {
	case *text:

	case HTML:
		m.dismountHTML(v)

	case Composer:
		m.dismountComponent(v)

	case *raw:
		m.dismountRawHTML(v)
	}
}

func (m *nodeManager) dismountHTML(v HTML) {
	for _, child := range v.body() {
		m.Dismount(child)
	}

	for _, handler := range v.events() {
		m.dismountHTMLEventHandler(handler)
	}

	v.setJSElement(nil)
}

func (m *nodeManager) dismountHTMLEventHandler(handler eventHandler) {
	if handler.close != nil {
		handler.close()
	}
}

func (m *nodeManager) dismountComponent(v Composer) error {
	panic("not implemented")
}

func (m *nodeManager) dismountRawHTML(v *raw) error {
	panic("not implemented")
}

// CanUpdate determines whether a given UI element 'v' can be updated with a new
// UI element 'new'. It returns false if the types of the two elements are
// different.
//
// For HTML elements, it ensures that the tag names match. Otherwise, it returns
// true indicating that an update is feasible.
func (m *nodeManager) CanUpdate(v, new UI) bool {
	if vType, newType := reflect.TypeOf(v), reflect.TypeOf(new); vType != newType {
		return false
	}

	switch v.(type) {
	case *htmlElem, *htmlElemSelfClosing:
		return v.(HTML).Tag() == new.(HTML).Tag()

	default:
		return true
	}
}

// Update updates the existing UI element 'v' with a new UI element 'new'. It
// returns the updated UI element and any error encountered during the update
// process.
func (m *nodeManager) Update(v, new UI) (UI, error) {
	if !v.Mounted() {
		return nil, errors.New("element not mounted").WithTag("type", reflect.TypeOf(v))
	}

	switch v := v.(type) {
	case *text:
		return m.updateText(v, new.(*text))

	case HTML:
		return m.updateHTML(v, new.(HTML))

	case Composer:
		return m.updateComponent(v, new.(Composer))

	case *raw:
		return m.updateRawHTML(v, new.(*raw))

	default:
		return nil, errors.New("unsupported element").WithTag("type", reflect.TypeOf(v))
	}
}

func (m *nodeManager) updateText(v, new *text) (UI, error) {
	if v.value == new.value {
		return v, nil
	}

	v.value = new.value
	v.JSValue().setNodeValue(v.value)
	return v, nil
}

func (m *nodeManager) updateHTML(v, new HTML) (UI, error) {
	attrs := v.attrs()
	newAttrs := new.attrs()
	if attrs == nil && len(newAttrs) != 0 {
		v = v.setAttrs(newAttrs)
		m.mountHTMLAttributes(v)
	} else if attrs != nil {
		m.updateHTMLAttributes(v, newAttrs)
	}

	events := v.events()
	newEvents := new.events()
	if events == nil && len(newEvents) != 0 {
		v = v.setEvents(newEvents)
		m.mountHTMLEventHandlers(v)
	} else if events != nil {
		m.updateHTMLEventHandlers(v, newEvents)
	}

	children := v.body()
	newChildren := new.body()
	sharedLen := min(len(children), len(newChildren))
	for i := 0; i < min(len(children), len(newChildren)); i++ {
		child := children[i]
		newChild := newChildren[i]
		if m.CanUpdate(child, newChild) {
			child, err := m.Update(child, newChild)
			if err != nil {
				return nil, errors.New("updating child failed").
					WithTag("type", reflect.TypeOf(v)).
					WithTag("tag", v.Tag()).
					WithTag("depth", v.depth()).
					WithTag("index", i).
					Wrap(err)
			}
			children[i] = child
			continue
		}

		newChild, err := m.Mount(v.depth()+1, newChildren[i])
		if err != nil {
			return nil, errors.New("mounting child failed").
				WithTag("type", reflect.TypeOf(v)).
				WithTag("tag", v.Tag()).
				WithTag("depth", v.depth()).
				WithTag("index", i).
				Wrap(err)
		}
		v.JSValue().replaceChild(newChild, child)
		newChild = newChild.setParent(v)
		children[i] = newChild
		m.Dismount(child)
	}

	for i := sharedLen; i < len(children); i++ {
		child := children[i]
		v.JSValue().removeChild(child)
		m.Dismount(child)
		children[i] = nil
	}
	children = children[:sharedLen]

	for i := sharedLen; i < len(newChildren); i++ {
		newChild, err := m.Mount(v.depth()+1, newChildren[i])
		if err != nil {
			return nil, errors.New("mounting child failed").
				WithTag("type", reflect.TypeOf(v)).
				WithTag("tag", v.Tag()).
				WithTag("depth", v.depth()).
				WithTag("index", i).
				Wrap(err)
		}
		v.JSValue().appendChild(newChild)
		newChild = newChild.setParent(v)
		children = append(children, newChild)
	}

	v = v.setBody(children)
	return v, nil
}

func (m *nodeManager) updateHTMLAttributes(v HTML, newAttrs attributes) {
	attrs := v.attrs()
	for name := range attrs {
		if _, remains := newAttrs[name]; !remains {
			deleteJSAttribute(v.JSValue(), name)
			delete(attrs, name)
		}
	}

	for name, value := range newAttrs {
		if attrs[name] == value {
			continue
		}

		attrs[name] = value
		setJSAttribute(v.JSValue(), name, resolveAttributeURLValue(
			name,
			value,
			m.ResolveURL,
		))
	}
}

func (m *nodeManager) updateHTMLEventHandlers(v HTML, newEvents eventHandlers) {
	events := v.events()
	for event, handler := range events {
		if _, remains := newEvents[event]; !remains {
			m.dismountHTMLEventHandler(handler)
			delete(events, event)
		}
	}

	for event, newHandler := range newEvents {
		handler, exists := events[event]
		if !exists {
			events[event] = m.mountHTMLEventHandler(v, newHandler)
			continue
		}

		if handler.Equal(newHandler) {
			continue
		}

		m.dismountHTMLEventHandler(handler)
		events[event] = m.mountHTMLEventHandler(v, newHandler)
	}
}

func (m *nodeManager) updateComponent(v, new Composer) (UI, error) {
	panic("not implemented")
}

func (m *nodeManager) updateRawHTML(v, new *raw) (UI, error) {
	panic("not implemented")
}

// MakeContext creates and returns a new context derived from the nodeManager's
// base context (Ctx). The derived context is configured and tailored for the
// provided UI element 'v'.
func (m *nodeManager) MakeContext(v UI) Context {
	m.initOnce.Do(m.init)

	return uiContext{
		Context:            m.Ctx,
		src:                v,
		jsSrc:              v.JSValue(),
		appUpdateAvailable: appUpdateAvailable,
		page:               m.Page,
		emit:               m.EmitHTMLEvent,
	}
}
