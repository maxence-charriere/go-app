package app

import (
	"fmt"
	"io"
	"net/url"
	"reflect"
)

// Node is the interface that describes an UI node.
type Node interface {
	nodeType() reflect.Type
}

// UI is the interface that describes a node that is a user interface element.
// eg. HTML elements and components.
type UI interface {
	Node
	Wrapper

	parent() UI
	setParent(p UI)
	dismount()
}

type nodeWithChildren interface {
	replaceChild(old, new UI)
}

type standardNode interface {
	UI
	nodeWithChildren

	attributes() map[string]string
	setAttribute(k string, v interface{})
	setAttributeValue(k, v string)
	removeAttributeValue(k string)
	eventHandlers() map[string]eventHandler
	setEventHandler(k string, h EventHandler)
	setEventHandlerValue(k string, h eventHandler)
	removeEventHandlerValue(k string, h eventHandler)
	mount() error
	children() []UI
	appendChild(child UI)
	appendChildValue(child UI)
	removeChild(child UI)
	removeChildValue(child UI)
	replaceChildValue(old, new UI)
	update(n standardNode)
}

type textNode interface {
	UI

	text() string
	mount() error
	update(t textNode)
}

type rawNode interface {
	UI

	raw() string
	mount() error
}

// Composer is the interface that describes a component that embeds other nodes.
//
// Satisfying this interface is done by embedding app.Compo into a struct and
// implementing the Render function.
//
// Example:
//  type Hello struct {
//      app.Compo
//  }
//
//  func (c *Hello) Render() app.UI {
//      return app.Text("hello")
//  }
type Composer interface {
	UI
	nodeWithChildren

	// Render returns the node tree that define how the component is desplayed.
	Render() UI

	// Update update the component appearance. It should be called when a field
	// used to render the component has been modified.
	Update()

	setCompo(n Composer)
	mount(c Composer) error
	update(n Composer)
}

// Mounter is the interface that describes a component that can perform
// additional actions when mounted.
type Mounter interface {
	Composer

	// The function that is called when the component is mounted.
	OnMount()
}

// Dismounter is the interface that describes a component that can perform
// additional actions when dismounted.
type Dismounter interface {
	Composer

	// The function that is called when the component is dismounted.
	OnDismount()
}

// Navigator is the interface that describes a component that can perform
// additional actions when navigated on.
type Navigator interface {
	Composer

	// The function that is called when the component is navigated on.
	OnNav(u *url.URL)
}

type writableNode interface {
	html(w io.Writer)
	htmlWithIndent(w io.Writer, indent int)
}

type conditionNode interface {
	Node

	nodes() []UI
}

type eventHandler struct {
	function   EventHandler
	jsFunction Func
}

func (h eventHandler) equals(o eventHandler) bool {
	return fmt.Sprintf("%p", h.function) == fmt.Sprintf("%p", o.function)
}

func indirect(nodes ...Node) []UI {
	inodes := make([]UI, 0, len(nodes))

	for _, n := range nodes {
		switch t := n.(type) {
		case conditionNode:
			inodes = append(inodes, t.nodes()...)

		case Composer:
			t.setCompo(t)
			inodes = append(inodes, t)

		case UI:
			inodes = append(inodes, t)
		}
	}

	return inodes
}

func mount(n Node) error {
	switch t := n.(type) {
	case textNode:
		return t.mount()

	case standardNode:
		return t.mount()

	case rawNode:
		return t.mount()

	case Composer:
		return t.mount(t)

	default:
		return fmt.Errorf("%T is not mountable", n)
	}
}

func update(a, b UI) error {
	if a.nodeType() != b.nodeType() {
		return replace(a, b)
	}

	switch t := a.(type) {
	case textNode:
		t.update(b.(textNode))

	case standardNode:
		return updateStandardNode(t, b.(standardNode))

	case rawNode:
		return updateRawNode(t, b.(rawNode))

	case Composer:
		t.update(b.(Composer))
		t.Update()

	default:
		return fmt.Errorf("%T: node can't be updated", t)
	}

	return nil
}

func replace(a, b UI) error {
	if err := mount(b); err != nil {
		return err
	}

	parent := a.parent()
	b.setParent(parent)
	parent.(nodeWithChildren).replaceChild(a, b)

	for {
		parentValue, ok := parent.(standardNode)
		if ok {
			parentValue.replaceChildValue(a, b)
			break
		}
		parent = parent.parent()
	}

	a.dismount()
	return nil
}

func updateStandardNode(a, b standardNode) error {
	a.update(b)

	achildren := a.children()
	bchildren := b.children()

	for len(achildren) != 0 && len(bchildren) != 0 {
		if err := update(achildren[0], bchildren[0]); err != nil {
			return err
		}
		achildren = achildren[1:]
		bchildren = bchildren[1:]
	}

	for len(achildren) != 0 {
		c := achildren[0]
		a.removeChildValue(c)
		a.removeChild(c)
		c.dismount()
		achildren = achildren[:len(achildren)-1]
	}

	for len(bchildren) != 0 {
		c := bchildren[0]
		if err := mount(c); err != nil {
			return err
		}
		a.appendChild(c)
		a.appendChildValue(c)
		bchildren = bchildren[1:]
	}

	return nil
}

func updateRawNode(a, b rawNode) error {
	if a.raw() != b.raw() {
		return replace(a, b)
	}
	return nil
}
