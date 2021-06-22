package app

import (
	"context"
	"io"
	"net/url"
	"strconv"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

type elemWithChildren interface {
	UI

	replaceChildAt(idx int, new UI) error
}

type elem struct {
	attrs       map[string]string
	body        []UI
	disp        Dispatcher
	ctx         context.Context
	ctxCancel   func()
	events      map[string]eventHandler
	jsvalue     Value
	parentElem  UI
	selfClosing bool
	tag         string
	this        UI
}

func (e *elem) Kind() Kind {
	return HTML
}

func (e *elem) JSValue() Value {
	return e.jsvalue
}

func (e *elem) Mounted() bool {
	return e.dispatcher() != nil &&
		e.ctx != nil &&
		e.ctx.Err() == nil &&
		e.self() != nil &&
		e.jsvalue != nil
}

func (e *elem) name() string {
	return e.tag
}

func (e *elem) self() UI {
	return e.this
}

func (e *elem) setSelf(n UI) {
	e.this = n
}

func (e *elem) context() context.Context {
	return e.ctx
}

func (e *elem) dispatcher() Dispatcher {
	return e.disp
}

func (e *elem) attributes() map[string]string {
	return e.attrs
}

func (e *elem) eventHandlers() map[string]eventHandler {
	return e.events
}

func (e *elem) parent() UI {
	return e.parentElem
}

func (e *elem) setParent(p UI) {
	e.parentElem = p
}

func (e *elem) children() []UI {
	return e.body
}

func (e *elem) mount(d Dispatcher) error {
	if e.Mounted() {
		return errors.New("mounting ui element failed").
			Tag("reason", "already mounted").
			Tag("name", e.name()).
			Tag("kind", e.Kind())
	}

	e.disp = d
	e.ctx, e.ctxCancel = context.WithCancel(context.Background())

	v, err := Window().createElement(e.tag)
	if err != nil {
		return errors.New("mounting ui element failed").
			Tag("name", e.name()).
			Tag("kind", e.Kind()).
			Wrap(err)
	}
	e.jsvalue = v

	for k, v := range e.attrs {
		v = e.resolveURLAttr(k, v)
		e.attrs[k] = v
		e.setJsAttr(k, v)
	}

	for k, v := range e.events {
		e.setJsEventHandler(k, v)
	}

	for _, c := range e.children() {
		if err := e.appendChild(c, true); err != nil {
			return errors.New("mounting ui element failed").
				Tag("name", e.name()).
				Tag("kind", e.Kind()).
				Wrap(err)
		}
	}

	return nil
}

func (e *elem) dismount() {
	for _, c := range e.children() {
		dismount(c)
	}

	for k, v := range e.events {
		e.delJsEventHandler(k, v)
	}

	e.ctxCancel()
	e.jsvalue = nil
}

func (e *elem) update(n UI) error {
	if !e.Mounted() {
		return nil
	}

	if n.Kind() != e.Kind() || n.name() != e.name() {
		return errors.New("updating ui element failed").
			Tag("replace", true).
			Tag("reason", "different element types").
			Tag("current-kind", e.Kind()).
			Tag("current-name", e.name()).
			Tag("updated-kind", n.Kind()).
			Tag("updated-name", n.name())
	}

	e.updateAttrs(n.attributes())
	e.updateEventHandler(n.eventHandlers())

	achildren := e.children()
	bchildren := n.children()
	i := 0

	// Update children:
	for len(achildren) != 0 && len(bchildren) != 0 {
		a := achildren[0]
		b := bchildren[0]

		err := update(a, b)
		if isErrReplace(err) {
			err = e.replaceChildAt(i, b)
		}

		if err != nil {
			return errors.New("updating ui element failed").
				Tag("kind", e.Kind()).
				Tag("name", e.name()).
				Wrap(err)
		}

		achildren = achildren[1:]
		bchildren = bchildren[1:]
		i++
	}

	// Remove children:
	for len(achildren) != 0 {
		if err := e.removeChildAt(i); err != nil {
			return errors.New("updating ui element failed").
				Tag("kind", e.Kind()).
				Tag("name", e.name()).
				Wrap(err)
		}

		achildren = achildren[1:]
	}

	// Add children:
	for len(bchildren) != 0 {
		c := bchildren[0]

		if err := e.appendChild(c, false); err != nil {
			return errors.New("updating ui element failed").
				Tag("kind", e.Kind()).
				Tag("name", e.name()).
				Wrap(err)
		}

		bchildren = bchildren[1:]
	}

	return nil
}

func (e *elem) appendChild(c UI, onlyJsValue bool) error {
	if err := mount(e.dispatcher(), c); err != nil {
		return errors.New("appending child failed").
			Tag("name", e.name()).
			Tag("kind", e.Kind()).
			Tag("child-name", c.name()).
			Tag("child-kind", c.Kind()).
			Wrap(err)
	}

	if !onlyJsValue {
		e.body = append(e.body, c)
	}

	c.setParent(e.self())
	e.JSValue().appendChild(c)
	return nil
}

func (e *elem) replaceChildAt(idx int, new UI) error {
	old := e.body[idx]

	if err := mount(e.dispatcher(), new); err != nil {
		return errors.New("replacing child failed").
			Tag("name", e.name()).
			Tag("kind", e.Kind()).
			Tag("index", idx).
			Tag("old-name", old.name()).
			Tag("old-kind", old.Kind()).
			Tag("new-name", new.name()).
			Tag("new-kind", new.Kind()).
			Wrap(err)
	}

	e.body[idx] = new
	new.setParent(e.self())
	e.JSValue().replaceChild(new, old)

	dismount(old)
	return nil
}

func (e *elem) removeChildAt(idx int) error {
	body := e.body
	if idx < 0 || idx >= len(body) {
		return errors.New("removing child failed").
			Tag("reason", "index out of range").
			Tag("index", idx).
			Tag("name", e.name()).
			Tag("kind", e.Kind())
	}

	c := body[idx]

	copy(body[idx:], body[idx+1:])
	body[len(body)-1] = nil
	body = body[:len(body)-1]
	e.body = body

	e.JSValue().removeChild(c)
	dismount(c)
	return nil
}

func (e *elem) updateAttrs(attrs map[string]string) {
	for k := range e.attrs {
		if _, exists := attrs[k]; !exists {
			e.delAttr(k)
		}
	}

	if e.attrs == nil && len(attrs) != 0 {
		e.attrs = make(map[string]string, len(attrs))
	}

	for k, v := range attrs {
		v = e.resolveURLAttr(k, v)
		if curval, exists := e.attrs[k]; !exists || curval != v {
			e.attrs[k] = v
			e.setJsAttr(k, v)
		}
	}
}

func (e *elem) setAttr(k string, v interface{}) {
	if e.attrs == nil {
		e.attrs = make(map[string]string)
	}

	switch k {
	case "style", "allow":
		s := e.attrs[k] + toString(v) + ";"
		e.attrs[k] = s

	case "class":
		s := e.attrs[k]
		if s != "" {
			s += " "
		}
		s += toString(v)
		e.attrs[k] = s

	default:
		e.attrs[k] = toString(v)
	}
}

func (e *elem) resolveURLAttr(k, v string) string {
	if !isURLAttrValue(k) {
		return v
	}
	return e.disp.resolveStaticResource(v)
}

func (e *elem) setJsAttr(k, v string) {
	switch k {
	case "value":
		e.JSValue().Set("value", v)

	case "class":
		e.JSValue().Set("className", v)

	case "contenteditable":
		e.JSValue().Set("contentEditable", v)

	case "async",
		"autofocus",
		"autoplay",
		"checked",
		"default",
		"defer",
		"disabled",
		"hidden",
		"ismap",
		"loop",
		"multiple",
		"muted",
		"open",
		"readonly",
		"required",
		"reversed",
		"selected":
		switch k {
		case "ismap":
			k = "isMap"
		case "readonly":
			k = "readOnly"
		}
		v, _ := strconv.ParseBool(v)
		e.JSValue().Set(k, v)

	default:
		if isURLAttrValue(k) {
			v = e.dispatcher().resolveStaticResource(v)
		}
		e.JSValue().setAttr(k, v)
	}
}

func (e *elem) delAttr(k string) {
	e.JSValue().delAttr(k)
	delete(e.attrs, k)
}

func (e *elem) updateEventHandler(handlers map[string]eventHandler) {
	for k, current := range e.events {
		if _, exists := handlers[k]; !exists {
			e.delJsEventHandler(k, current)
		}
	}

	if e.events == nil && len(handlers) != 0 {
		e.events = make(map[string]eventHandler, len(handlers))
	}

	for k, new := range handlers {
		if current, exists := e.events[k]; !current.equal(new) {
			if exists {
				e.delJsEventHandler(k, current)
			}

			e.events[k] = new
			e.setJsEventHandler(k, new)
		}
	}
}

func (e *elem) setEventHandler(k string, h EventHandler, scope ...interface{}) {
	if e.events == nil {
		e.events = make(map[string]eventHandler)
	}

	e.events[k] = eventHandler{
		event: k,
		scope: toPath(scope...),
		value: h,
	}
}

func (e *elem) setJsEventHandler(k string, h eventHandler) {
	jshandler := makeJsEventHandler(e.self(), h.value)
	h.jsvalue = jshandler
	e.events[k] = h
	e.JSValue().addEventListener(k, jshandler)
}

func (e *elem) delJsEventHandler(k string, h eventHandler) {
	e.jsvalue.removeEventListener(k, h.jsvalue)
	h.jsvalue.Release()
	delete(e.events, k)
}

func (e *elem) setBody(body ...UI) {
	if e.selfClosing {
		panic(errors.New("setting html element body failed").
			Tag("reason", "self closing element can't have children").
			Tag("name", e.name()),
		)
	}

	e.body = FilterUIElems(body...)
}

func (e *elem) onNav(u *url.URL) {
	for _, c := range e.children() {
		c.onNav(u)
	}
}

func (e *elem) onAppUpdate() {
	for _, c := range e.children() {
		c.onAppUpdate()
	}
}

func (e *elem) onAppInstallChange() {
	for _, c := range e.children() {
		c.onAppInstallChange()
	}
}

func (e *elem) onResize() {
	for _, c := range e.children() {
		c.onResize()
	}
}

func (e *elem) preRender(p Page) {
	for _, c := range e.children() {
		c.preRender(p)
	}
}

func (e *elem) html(w io.Writer) {
	w.Write(stob("<"))
	w.Write(stob(e.tag))

	for k, v := range e.attrs {
		w.Write(stob(" "))
		w.Write(stob(k))

		if v != "" {
			w.Write(stob(`="`))
			w.Write(stob(v))
			w.Write(stob(`"`))
		}
	}

	w.Write(stob(">"))

	if e.selfClosing {
		return
	}

	for _, c := range e.body {
		w.Write(ln())
		if c.self() == nil {
			c.setSelf(c)
		}
		c.html(w)
	}

	if len(e.body) != 0 {
		w.Write(ln())
	}

	w.Write(stob("</"))
	w.Write(stob(e.tag))
	w.Write(stob(">"))
}

func (e *elem) htmlWithIndent(w io.Writer, indent int) {
	writeIndent(w, indent)
	w.Write(stob("<"))
	w.Write(stob(e.tag))

	for k, v := range e.attrs {
		w.Write(stob(" "))
		w.Write(stob(k))

		if v != "" {
			w.Write(stob(`="`))
			w.Write(stob(v))
			w.Write(stob(`"`))
		}
	}

	w.Write(stob(">"))

	if e.selfClosing {
		return
	}

	for _, c := range e.body {
		w.Write(ln())
		if c.self() == nil {
			c.setSelf(c)
		}
		c.htmlWithIndent(w, indent+1)
	}

	if len(e.body) != 0 {
		w.Write(ln())
		writeIndent(w, indent)
	}

	w.Write(stob("</"))
	w.Write(stob(e.tag))
	w.Write(stob(">"))
}

func isURLAttrValue(k string) bool {
	switch k {
	case "cite",
		"data",
		"href",
		"src",
		"srcset":
		return true

	default:
		return false
	}
}
