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

// ValueNode is the interface that describes a node that holds a value. HTML
// nodes and component nodes are value nodes.
type ValueNode interface {
	Node
	Wrapper

	parent() nodeWithChildren
	setParent(p nodeWithChildren)
	dismount()
}

type nodeWithChildren interface {
	ValueNode

	replaceChild(old, new ValueNode)
}

type standardNode interface {
	ValueNode

	attributes() map[string]string
	setAttribute(k string, v interface{})
	setAttributeValue(k, v string)
	removeAttributeValue(k string)
	eventHandlers() map[string]eventHandler
	setEventHandler(k string, h EventHandler)
	setEventHandlerValue(k string, h eventHandler)
	removeEventHandlerValue(k string, h eventHandler)
	mount() error
	children() []ValueNode
	appendChild(child ValueNode)
	appendChildValue(child ValueNode)
	removeChild(child ValueNode)
	removeChildValue(child ValueNode)
	replaceChildValue(old, new ValueNode)
	update(n standardNode)
}

type textNode interface {
	ValueNode

	text() string
	mount() error
	update(t textNode)
}

type rawNode interface {
	ValueNode

	raw() string
	mount() error
}

// CompoNode is the interface that describes a component node that is built on
// the top of a Compo.
//
// Example:
//  type Hello struct {
//      app.Compo
//  }
//
//  func (c *Hello) Render() app.Node {
//      return app.Text("hello")
//  }
type CompoNode interface {
	ValueNode

	// Render returns the node tree that define how the component is desplayed.
	//
	// Returned node must be a standard HTML node or another component node.
	// Mounting a component node that returns a condition node as root triggers
	// a panic.
	Render() ValueNode

	// Update update the component appearance. It should be called when a field
	// used to render the component has been modified.
	Update()

	setCompo(n CompoNode)
	mount(c CompoNode) error
	update(n CompoNode)
}

// Mounter is the interface that describes a component node that can perform
// additional actions when mounted.
type Mounter interface {
	CompoNode

	// The function that is called when the component is mounted.
	OnMount()
}

// Dismounter is the interface that describes a component node that can perform
// additional actions when dismounted.
type Dismounter interface {
	CompoNode

	// The function that is called when the component is dismounted.
	OnDismount()
}

// Navigator is the interface that describes a component node that can perform
// additional actions when navigated on.
type Navigator interface {
	CompoNode

	// The function that is called when the component is navigated on.
	OnNav(u *url.URL)
}

type writableNode interface {
	html(w io.Writer)
	htmlWithIndent(w io.Writer, indent int)
}

type conditionNode interface {
	Node

	nodes() []ValueNode
}

type eventHandler struct {
	function   EventHandler
	jsFunction Func
}

func (h eventHandler) equals(o eventHandler) bool {
	return fmt.Sprintf("%p", h.function) == fmt.Sprintf("%p", o.function)
}

func indirect(nodes ...Node) []ValueNode {
	inodes := make([]ValueNode, 0, len(nodes))

	for _, n := range nodes {
		switch t := n.(type) {
		case conditionNode:
			inodes = append(inodes, t.nodes()...)

		case CompoNode:
			t.setCompo(t)
			inodes = append(inodes, t)

		case ValueNode:
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

	case CompoNode:
		return t.mount(t)

	default:
		return fmt.Errorf("%T is not mountable", n)
	}
}

func update(a, b ValueNode) error {
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

	case CompoNode:
		t.update(b.(CompoNode))
		t.Update()

	default:
		return fmt.Errorf("%T: node can't be updated", t)
	}

	return nil
}

func replace(a, b ValueNode) error {
	if err := mount(b); err != nil {
		return err
	}

	parent := a.parent()
	b.setParent(parent)
	parent.replaceChild(a, b)

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
