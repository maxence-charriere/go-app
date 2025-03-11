package app

import (
	"bytes"
	"html"
	"io"
	"reflect"
	"strconv"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/errors"
)

// UI is the interface that describes a user interface element such as
// components and HTML elements.
type UI interface {
	// JSValue returns the javascript value linked to the element.
	JSValue() Value

	// Reports whether the element is mounted.
	Mounted() bool

	parent() UI
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

		switch elem := elem.(type) {
		case Condition:
			replaceElemAt(i, elem.body()...)

		case RangeLoop:
			replaceElemAt(i, elem.body()...)
		}
	}

	return v
}

// HTMLString returns a string that represents the HTML markup for the provided
// UI element.
func HTMLString(ui UI) string {
	engine := NewTestEngine().(*engineX)
	var b bytes.Buffer
	engine.nodes.Encode(engine.baseContext(), &b, ui)
	return b.String()
}

// PrintHTML writes the HTML representation of the given UI element into the
// specified writer.
func PrintHTML(w io.Writer, ui UI) {
	w.Write([]byte(HTMLString(ui)))
}

// Component events.
type nav struct{}
type appUpdate struct{}
type appInstallChange struct{}
type resize struct{}

// nodeManager orchestrates the lifecycle of UI elements, providing specialized
// mechanisms for mounting, dismounting, and updating nodes.
type nodeManager struct {
}

// Mount mounts a UI element based on its type and the specified depth. It
// returns the mounted UI element and any potential error during the process.
func (m nodeManager) Mount(ctx Context, depth uint, v UI) (UI, error) {
	ctx = m.context(ctx, v)

	switch v := v.(type) {
	case *text:
		return m.mountText(v)

	case HTML:
		return m.mountHTML(ctx, depth, v)

	case Composer:
		return m.mountComponent(ctx, depth, v)

	case *raw:
		return m.mountRawHTML(depth, v)

	default:
		return nil, errors.New("unsupported element").
			WithTag("type", reflect.TypeOf(v)).
			WithTag("depth", depth)
	}
}

func (m nodeManager) mountText(v *text) (UI, error) {
	if v.Mounted() {
		return nil, errors.New("text is already mounted").
			WithTag("parent-type", reflect.TypeOf(v.parent())).
			WithTag("preview-value", previewText(v.value))
	}

	v.jsvalue = Window().createTextNode(v.value)
	return v, nil
}

func (m nodeManager) mountHTML(ctx Context, depth uint, v HTML) (UI, error) {
	if v.Mounted() {
		return nil, errors.New("html element is already mounted").
			WithTag("parent-type", reflect.TypeOf(v.parent())).
			WithTag("type", reflect.TypeOf(v)).
			WithTag("tag", v.Tag()).
			WithTag("depth", v.depth())
	}

	jsElement, _ := Window().createElement(v.Tag(), v.XMLNamespace())
	v = v.setJSElement(jsElement)
	m.mountHTMLAttributes(ctx, v)
	m.mountHTMLEventHandlers(ctx, v)

	v = v.setDepth(depth).(HTML)
	children := v.body()
	for i, child := range children {
		var err error
		if child, err = m.Mount(ctx, depth+1, child); err != nil {
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

func (m nodeManager) mountHTMLAttributes(ctx Context, v HTML) {
	for name, value := range v.attrs() {
		setJSAttribute(v.JSValue(), name, resolveAttributeURLValue(
			name,
			value,
			ctx.ResolveStaticResource,
		))
	}
}

func (m nodeManager) mountHTMLEventHandlers(ctx Context, v HTML) {
	events := v.events()
	for event, handler := range events {
		events[event] = m.mountHTMLEventHandler(ctx, v, handler)
	}
}

func (m nodeManager) mountHTMLEventHandler(ctx Context, v HTML, handler eventHandler) eventHandler {
	event := handler.event

	jsHandler := FuncOf(func(this Value, args []Value) any {
		if len(args) != 0 {
			ctx.Dispatch(func(ctx Context) {
				event := Event{Value: args[0]}
				trackMousePosition(event)
				handler.goHandler(ctx, event)
			})
		}
		return nil
	})
	v.JSValue().addEventListener(event, jsHandler, handler.options())

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

func (m nodeManager) mountComponent(ctx Context, depth uint, v Composer) (UI, error) {
	if v.Mounted() {
		return nil, errors.New("component is already mounted").
			WithTag("parent-type", reflect.TypeOf(v.parent())).
			WithTag("type", reflect.TypeOf(v)).
			WithTag("depth", v.depth())
	}

	v = v.setRef(v)
	v = v.setDepth(depth)

	if initializer, ok := v.(Initializer); ok {
		initializer.OnInit()
	}

	if preRenderer, ok := v.(PreRenderer); ok && IsServer {
		ctx.Dispatch(preRenderer.OnPreRender)
	}

	if mounter, ok := v.(Mounter); ok && IsClient {
		ctx.Dispatch(mounter.OnMount)
	}

	root, err := m.renderComponent(v)
	if err != nil {
		return nil, errors.New("rendering component failed").
			WithTag("type", reflect.TypeOf(v)).
			WithTag("depth", v.depth()).
			Wrap(err)
	}
	if root, err = m.Mount(ctx, depth+1, root); err != nil {
		return nil, errors.New("mounting component root failed").
			WithTag("type", reflect.TypeOf(v)).
			WithTag("depth", v.depth()).
			Wrap(err)
	}
	root = root.setParent(v)
	v = v.setRoot(root)

	return v, nil
}

func (m nodeManager) renderComponent(v Composer) (UI, error) {
	rendering := FilterUIElems(v.Render())
	if len(rendering) == 0 {
		return nil, errors.New("render method does not returns a text, html element, or component")
	}
	return rendering[0], nil
}

func (m nodeManager) mountRawHTML(depth uint, v *raw) (UI, error) {
	if v.Mounted() {
		return nil, errors.New("raw html is already mounted").
			WithTag("parent-type", reflect.TypeOf(v.parent())).
			WithTag("type", reflect.TypeOf(v)).
			WithTag("depth", v.depth()).
			WithTag("raw-preview", previewText(v.value))
	}

	wrapper, _ := Window().createElement("div", "")
	wrapper.setInnerHTML(v.value)
	v.jsElement = wrapper.firstChild()
	wrapper.removeChild(v.jsElement)
	return v, nil
}

// Dismount removes a UI element based on its type.
func (m nodeManager) Dismount(v UI) {
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

func (m nodeManager) dismountHTML(v HTML) {
	for _, child := range v.body() {
		m.Dismount(child)
	}

	for _, handler := range v.events() {
		m.dismountHTMLEventHandler(handler)
	}

	v.setJSElement(nil)
}

func (m nodeManager) dismountHTMLEventHandler(handler eventHandler) {
	if handler.close != nil {
		handler.close()
	}
}

func (m nodeManager) dismountComponent(v Composer) {
	m.Dismount(v.root())
	v.setRef(nil)

	if dismounter, ok := v.(Dismounter); ok {
		dismounter.OnDismount()
	}
}

func (m nodeManager) dismountRawHTML(v *raw) {
	v.jsElement = nil
}

// CanUpdate determines whether a given UI element 'v' can be updated with a new
// UI element 'new'. It returns false if the types of the two elements are
// different.
//
// For HTML elements, it ensures that the tag names match. Otherwise, it returns
// true indicating that an update is feasible.
func (m nodeManager) CanUpdate(v, new UI) bool {
	if vType, newType := reflect.TypeOf(v), reflect.TypeOf(new); vType != newType {
		return false
	}

	switch v.(type) {
	case DismountEnforcer:
		return v.(DismountEnforcer).CompoID() == new.(DismountEnforcer).CompoID()

	case *htmlElem, *htmlElemSelfClosing:
		return v.(HTML).Tag() == new.(HTML).Tag()

	default:
		return true
	}
}

// Update updates the existing UI element 'v' with a new UI element 'new'. It
// returns the updated UI element and any error encountered during the update
// process.
func (m nodeManager) Update(ctx Context, v, new UI) (UI, error) {
	if !v.Mounted() {
		return nil, errors.New("element not mounted").WithTag("type", reflect.TypeOf(v))
	}

	ctx = m.context(ctx, v)
	switch v := v.(type) {
	case *text:
		return m.updateText(v, new.(*text))

	case HTML:
		return m.updateHTML(ctx, v, new.(HTML))

	case Composer:
		return m.updateComponent(ctx, v, new.(Composer))

	case *raw:
		return m.updateRawHTML(ctx, v, new.(*raw))

	default:
		return nil, errors.New("unsupported element").WithTag("type", reflect.TypeOf(v))
	}
}

func (m nodeManager) updateText(v, new *text) (UI, error) {
	if v.value == new.value {
		return v, nil
	}

	v.value = new.value
	v.JSValue().setNodeValue(v.value)
	return v, nil
}

func (m nodeManager) updateHTML(ctx Context, v, new HTML) (UI, error) {
	attrs := v.attrs()
	newAttrs := new.attrs()
	if attrs == nil && len(newAttrs) != 0 {
		v = v.setAttrs(newAttrs)
		m.mountHTMLAttributes(ctx, v)
	} else if attrs != nil {
		m.updateHTMLAttributes(ctx, v, newAttrs)
	}

	events := v.events()
	newEvents := new.events()
	if events == nil && len(newEvents) != 0 {
		v = v.setEvents(newEvents)
		m.mountHTMLEventHandlers(ctx, v)
	} else if events != nil {
		m.updateHTMLEventHandlers(ctx, v, newEvents)
	}

	children := v.body()
	newChildren := new.body()
	sharedLen := min(len(children), len(newChildren))
	for i := 0; i < min(len(children), len(newChildren)); i++ {
		child := children[i]
		newChild := newChildren[i]
		if m.CanUpdate(child, newChild) {
			child, err := m.Update(ctx, child, newChild)
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

		newChild, err := m.Mount(ctx, v.depth()+1, newChildren[i])
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
		newChild, err := m.Mount(ctx, v.depth()+1, newChildren[i])
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

func (m nodeManager) updateHTMLAttributes(ctx Context, v HTML, newAttrs attributes) {
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
			ctx.ResolveStaticResource,
		))
	}
}

func (m nodeManager) updateHTMLEventHandlers(ctx Context, v HTML, newEvents eventHandlers) {
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
			events[event] = m.mountHTMLEventHandler(ctx, v, newHandler)
			continue
		}

		if handler.Equal(newHandler) {
			continue
		}

		m.dismountHTMLEventHandler(handler)
		events[event] = m.mountHTMLEventHandler(ctx, v, newHandler)
	}
}

func (m nodeManager) updateComponent(ctx Context, v, new Composer) (UI, error) {
	value := reflect.Indirect(reflect.ValueOf(v))
	newValue := reflect.Indirect(reflect.ValueOf(new))

	var modifiedFields bool
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		newField := newValue.Field(i)
		if !field.CanSet() {
			continue
		}
		if _, compoStruct := field.Interface().(Compo); compoStruct {
			continue
		}
		if !canUpdateValue(field, newField) {
			continue
		}
		field.Set(newField)
		modifiedFields = true
	}
	if !modifiedFields {
		return v, nil
	}

	if updater, ok := v.(Updater); ok {
		updater.OnUpdate(ctx)
	}

	ctx.removeComponentUpdate(v)
	return m.UpdateComponentRoot(ctx, v)
}

// UpdateComponentRoot updates the root element of the given component.
func (m nodeManager) UpdateComponentRoot(ctx Context, v Composer) (UI, error) {
	ctx = m.context(ctx, v)

	root := v.root()
	newRoot, err := m.renderComponent(v)
	if err != nil {
		return nil, errors.New("rendering component failed").
			WithTag("type", reflect.TypeOf(v)).
			WithTag("depth", v.depth()).
			Wrap(err)
	}

	if m.CanUpdate(root, newRoot) {
		if root, err = m.Update(ctx, root, newRoot); err != nil {
			return nil, errors.New("updating component root failed").
				WithTag("type", reflect.TypeOf(v)).
				WithTag("depth", v.depth()).
				Wrap(err)
		}
		v.setRoot(root)
	} else {
		if newRoot, err = m.Mount(ctx, v.depth()+1, newRoot); err != nil {
			return nil, errors.New("mounting component root failed").
				WithTag("type", reflect.TypeOf(v)).
				WithTag("depth", v.depth()).
				Wrap(err)
		}

		for parent := v.parent(); parent != nil; parent = parent.parent() {
			if parent, isHTML := parent.(HTML); isHTML {
				parent.JSValue().replaceChild(newRoot, root)
				break
			}
		}
		newRoot.setParent(v)
		v.setRoot(newRoot)
		m.Dismount(root)
	}

	return v, nil
}

func (m nodeManager) updateRawHTML(ctx Context, v, new *raw) (UI, error) {
	if v.value == new.value {
		return v, nil
	}

	newMount, err := m.Mount(ctx, v.depth(), new)
	if err != nil {
		return nil, errors.New("mounting updated raw html failed").
			WithTag("type", reflect.TypeOf(v)).
			WithTag("depth", v.depth()).
			WithTag("raw-preview", previewText(v.value)).
			Wrap(err)
	}

	for parent := v.parent(); parent != nil; parent = parent.parent() {
		if parent, isHTML := parent.(HTML); isHTML {
			parent.JSValue().replaceChild(newMount, v)
			newMount.setParent(parent)
			break
		}
	}
	m.Dismount(v)
	return newMount, nil
}

func (m nodeManager) context(ctx Context, v UI) Context {
	ctx.sourceElement = v
	ctx.notifyComponentEvent = m.NotifyComponentEvent
	return ctx
}

// NotifyComponentEvent traverses a UI element tree to propagate a component
// event, activating pertinent component handlers and potentially enqueuing
// component updates as needed.
func (m nodeManager) NotifyComponentEvent(ctx Context, root UI, event any) {
	ctx = m.context(ctx, root)

	switch element := root.(type) {
	case HTML:
		for _, child := range element.body() {
			m.NotifyComponentEvent(ctx, child, event)
		}

	case Composer:
		switch event.(type) {
		case nav:
			if navigator, ok := element.(Navigator); ok {
				ctx.Dispatch(navigator.OnNav)
			}

		case appUpdate:
			if appUpdater, ok := element.(AppUpdater); ok {
				ctx.Dispatch(appUpdater.OnAppUpdate)
			}

		case appInstallChange:
			if appInstaller, ok := element.(AppInstaller); ok {
				ctx.Dispatch(appInstaller.OnAppInstallChange)
			}

		case resize:
			if resizer, ok := element.(Resizer); ok {
				ctx.Dispatch(resizer.OnResize)
			}
		}
		m.NotifyComponentEvent(ctx, element.root(), event)
	}
}

// Encode transforms the provided UI element into its HTML byte slice
// representation. This allows for the conversion of in-memory UI structures
// into a format suitable for server rendering.
func (m nodeManager) Encode(ctx Context, w *bytes.Buffer, v UI) {
	m.encode(ctx, w, 0, v)
}

func (m nodeManager) encode(ctx Context, w *bytes.Buffer, depth int, v UI) {
	switch v := v.(type) {
	case *text:
		m.encodeText(w, depth, v)

	case HTML:
		m.encodeHTML(ctx, w, depth, v)

	case Composer:
		m.encodeComponent(ctx, w, depth, v)

	case *raw:
		m.encodeRawHTML(w, depth, v)
	}
}

func (m nodeManager) encodeText(w *bytes.Buffer, depth int, v *text) {
	if v.value != "" {
		m.encodeIndent(w, depth)
		w.WriteString(html.EscapeString(v.value))
	}
}

func (m nodeManager) encodeIndent(w *bytes.Buffer, depth int) {
	for i := 0; i < depth*2; i++ {
		w.WriteByte(' ')
	}
}

func (m nodeManager) encodeHTML(ctx Context, w *bytes.Buffer, depth int, v HTML) {
	m.encodeIndent(w, depth)
	w.WriteByte('<')
	w.WriteString(v.Tag())
	for name, value := range v.attrs() {
		m.encodeHTMLAttribute(ctx, w, name, value)
	}
	w.WriteByte('>')

	if v.SelfClosing() {
		return
	}

	children := v.body()
	switch {
	case len(children) > 1:
		w.WriteByte('\n')
		for _, child := range children {
			m.encode(ctx, w, depth+1, child)
			w.WriteByte('\n')
		}
		m.encodeIndent(w, depth)

	case len(children) == 1:
		child := children[0]
		if text, ok := child.(*text); ok {
			m.encodeText(w, 0, text)
		} else {
			w.WriteByte('\n')
			m.encode(ctx, w, depth+1, child)
			w.WriteByte('\n')
			m.encodeIndent(w, depth)
		}
	}

	w.WriteString("</")
	w.WriteString(v.Tag())
	w.WriteByte('>')
}

func (m nodeManager) encodeHTMLAttribute(ctx Context, w *bytes.Buffer, name, value string) {
	if value == "" {
		switch name {
		case "id", "class", "title":
			return
		}
	}

	w.WriteString(" ")
	w.WriteString(name)
	if value != "" && value != "true" {
		w.WriteString("=")
		w.WriteString(strconv.Quote(resolveAttributeURLValue(name, value, ctx.ResolveStaticResource)))
	}
}

func (m nodeManager) encodeComponent(ctx Context, w *bytes.Buffer, depth int, v Composer) {
	root := v.root()
	if root == nil {
		root, _ = m.renderComponent(v)
	}
	if root != nil {
		m.encode(ctx, w, depth, root)
	}
}

func (m nodeManager) encodeRawHTML(w *bytes.Buffer, depth int, v *raw) {
	if v.value != "" {
		m.encodeIndent(w, depth)
		w.WriteString(v.value)
	}
}

func canUpdateValue(v, new reflect.Value) bool {
	switch v.Kind() {
	case reflect.String,
		reflect.Bool,
		reflect.Int,
		reflect.Int64,
		reflect.Int32,
		reflect.Int16,
		reflect.Int8,
		reflect.Uint,
		reflect.Uint64,
		reflect.Uint32,
		reflect.Uint16,
		reflect.Uint8,
		reflect.Float64,
		reflect.Float32:
		return !v.Equal(new)

	default:
		switch v.Interface().(type) {
		case time.Time:
			return !v.Equal(new)

		default:
			return !reflect.DeepEqual(v.Interface(), new.Interface())
		}
	}
}

func component(v UI) (Composer, bool) {
	for element := v; element != nil; element = element.parent() {
		if component, ok := element.(Composer); ok {
			return component, true
		}
	}
	return nil, false
}
