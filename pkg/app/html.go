package app

import (
	"context"
	"io"
	"net/url"
	"reflect"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

type htmlElement[T any] struct {
	tag           string
	xmlns         string
	isSelfClosing bool
	attributes    attributes
	eventHandlers map[string]eventHandler
	parent        UI
	children      []UI

	context       context.Context
	contextCancel func()
	dispatcher    Dispatcher
	jsElement     Value
}

func (e *htmlElement[T]) Attr(name string, value any) T {
	if e.attributes == nil {
		e.attributes = make(attributes)
	}
	e.attributes.Set(name, value)

	return e.toHTMLInterface()
}

func (e *htmlElement[T]) On(event string, h EventHandler, scope ...any) T {
	if h == nil {
		return e.toHTMLInterface()
	}

	if e.eventHandlers == nil {
		e.eventHandlers = make(map[string]eventHandler)
	}
	e.eventHandlers[event] = makeEventHandler(event, h, scope...)

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
	return e.context != nil && e.context.Err() == nil
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

func (e *htmlElement[T]) mount(d Dispatcher) error {
	if e.IsMounted() {
		return errors.New("html element is already mounted").Tag("tag", e.tag)
	}

	e.context, e.contextCancel = context.WithCancel(context.Background())
	e.dispatcher = d

	jsElement, err := Window().createElement(e.tag, e.xmlns)
	if err != nil {
		return errors.New("creating javascript element failed").
			Tag("tag", e.tag).
			Tag("xmlns", e.xmlns).
			Wrap(err)
	}
	e.jsElement = jsElement

	e.attributes.Mount(jsElement, d.resolveStaticResource)

	for event, eh := range e.eventHandlers {
		e.eventHandlers[event] = eh.Mount(e)
	}

	for i, c := range e.children {
		if err := mount(d, c); err != nil {
			return errors.New("mounting child failed").
				Tag("index", i).
				Tag("child-type", reflect.TypeOf(c)).
				Wrap(err)
		}

		c.setParent(e)
		e.JSValue().appendChild(c)
	}

	return nil
}

func (e *htmlElement[T]) dismount() {
	for _, c := range e.children {
		dismount(c)
	}

	for _, eh := range e.eventHandlers {
		eh.Dismount()
	}

	e.contextCancel()
}

func (e *htmlElement[T]) canUpdateWith(v UI) bool {
	if v, ok := v.(*htmlElement[T]); ok {
		return ok && e.tag == v.tag
	}
	return false
}

func (e *htmlElement[T]) updateWith(v UI) error {
	if e.IsMounted() {
		return nil
	}

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
