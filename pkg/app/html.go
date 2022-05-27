package app

import (
	"context"
	"io"
	"net/url"
	"strings"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

type htmlElement[T any] struct {
	tag           string
	xmlns         string
	isSelfClosing bool
	attributes    map[string]string
	eventHandlers map[string]eventHandler
	parent        UI
	children      []UI

	context       context.Context
	contextCancel func()
	dispatcher    Dispatcher
	jsElement     Value
}

func (e *htmlElement[T]) Attr(k string, v any) T {
	if e.attributes == nil {
		e.attributes = make(map[string]string)
	}

	switch k {
	case "style", "allow":
		var b strings.Builder
		b.WriteString(e.attributes[k])
		b.WriteString(toAttributeValue(v))
		b.WriteByte(';')
		e.attributes[k] = b.String()

	case "class":
		var b strings.Builder
		b.WriteString(e.attributes[k])
		if b.Len() != 0 {
			b.WriteByte(' ')
		}
		b.WriteString(toAttributeValue(v))
		e.attributes[k] = b.String()

	default:
		e.attributes[k] = toAttributeValue(v)
	}

	return e.toHTMLInterface()
}

func (e *htmlElement[T]) On(event string, h EventHandler, scope ...any) T {
	if h == nil {
		return e.toHTMLInterface()
	}

	if e.eventHandlers == nil {
		e.eventHandlers = make(map[string]eventHandler)
	}

	e.eventHandlers[event] = eventHandler{
		event:     event,
		scope:     toPath(scope...),
		goHandler: h,
	}

	return e.toHTMLInterface()
}

func (e *htmlElement[T]) Text(v any) T {
	switch e.tag {
	case "textarea":
		return e.Attr("value", v)

	default:
		return e.Body(Text(v))
	}
}

func (e *htmlElement[T]) Body(v ...UI) T {
	if e.isSelfClosing {
		panic(errors.New("setting html element body failed").
			Tag("reason", "self closing element can't have children").
			Tag("tag", e.tag),
		)
	}

	e.children = FilterUIElems(v...)
	return e.toHTMLInterface()
}

func (e *htmlElement[T]) JSValue() Value {
	return e.jsElement
}

func (e *htmlElement[T]) IsMounted() bool {
	return e.context != nil &&
		e.context.Err() == nil &&
		e.dispatcher != nil &&
		e.jsElement != nil
}

func (e *htmlElement[T]) toHTMLInterface() T {
	var i any = e
	return i.(T)
}

func (e *htmlElement[T]) getContext() context.Context {
	return e.context
}

func (e *htmlElement[T]) getDispatcher() Dispatcher {
	return e.dispatcher
}

func (e *htmlElement[T]) getAttributes() map[string]string {
	return e.attributes
}

func (e *htmlElement[T]) getEventHandlers() map[string]eventHandler {
	return e.eventHandlers
}

func (e *htmlElement[T]) getParent() UI {
	return e.parent
}

func (e *htmlElement[T]) setParent(v UI) {
	e.parent = v
}

func (e *htmlElement[T]) getChildren() []UI {
	return e.children
}

func (e *htmlElement[T]) mount(Dispatcher) error {
	panic("not implemented")
}

func (e *htmlElement[T]) dismount() {
	panic("not implemented")
}

func (e *htmlElement[T]) update(UI) error {
	panic("not implemented")
}

func (e *htmlElement[T]) onNav(*url.URL) {
	panic("not implemented")
}

func (e *htmlElement[T]) onAppUpdate() {
	panic("not implemented")
}

func (e *htmlElement[T]) onAppInstallChange() {
	panic("not implemented")
}

func (e *htmlElement[T]) onResize() {
	panic("not implemented")
}

func (e *htmlElement[T]) preRender(Page) {
	panic("not implemented")
}

func (e *htmlElement[T]) html(w io.Writer) {
	panic("not implemented")
}

func (e *htmlElement[T]) htmlWithIndent(w io.Writer, indent int) {
	panic("not implemented")
}

// -----------------------------------------------------------------------------
// The method below might be removed in later versions.
// -----------------------------------------------------------------------------
func (e *htmlElement[T]) Kind() Kind {
	return HTML
}

func (e *htmlElement[T]) name() string {
	return e.tag
}

func (e *htmlElement[T]) self() UI {
	return e
}

func (e *htmlElement[T]) setSelf(UI) {
}
