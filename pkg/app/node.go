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

type writableNode interface {
	html(w io.Writer)
	htmlWithIndent(w io.Writer, indent int)
}

func indirect(nodes ...Node) []UI {
	inodes := make([]UI, 0, len(nodes))

	for _, n := range nodes {
		if v := reflect.ValueOf(n); v.Kind() == reflect.Ptr && v.IsNil() {
			continue
		}

		switch t := n.(type) {
		case Condition:
			inodes = append(inodes, t.nodes()...)

		case RangeLoop:
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

func nav(n UI, u *url.URL) {
	switch t := n.(type) {
	case standardNode:
		t.nav(u)

	case Composer:
		t.mount(t)
	}
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
		c.setParent(a)
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
