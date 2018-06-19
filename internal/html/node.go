package html

import (
	"github.com/murlokswarm/app"
)

type node interface {
	app.DOMNode
	SetParent(node)
	ConsumeChanges() []Change
}

func attrsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}

	for k, va := range a {
		if vb, ok := b[k]; !ok || va != vb {
			return false
		}
	}
	return true
}

// Change represents a change to perform in order to render a component within
// a control.
type Change struct {
	// The change type.
	Type string

	// A value that describes how to make the change.
	Value interface{}
}

const (
	createText = "createText"
	setText    = "setText"

	createElem  = "createElem"
	setAttrs    = "setAttrs"
	appendChild = "appendChild"
	removeChild = "removeChild"

	setRoot    = "setRoot"
	deleteNode = "deleteNode"
)

func createTextChange(n *textNode) Change {
	return Change{
		Type:  createText,
		Value: n,
	}
}

type textValue struct {
	ID   string
	Text string
}

func setTextChange(id, text string) Change {
	return Change{
		Type: setText,
		Value: textValue{
			ID:   id,
			Text: text,
		},
	}
}

func createElemChange(n *elemNode) Change {
	return Change{
		Type:  createElem,
		Value: n,
	}
}

func setAttrsChange(a map[string]string) Change {
	return Change{
		Type:  setAttrs,
		Value: a,
	}
}

type childValue struct {
	ParentID string
	ChildID  string
}

func appendChildChange(parentID, childID string) Change {
	return Change{
		Type: appendChild,
		Value: childValue{
			ParentID: parentID,
			ChildID:  childID,
		},
	}
}

func removeChildChange(parentID, childID string) Change {
	return Change{
		Type: removeChild,
		Value: childValue{
			ParentID: parentID,
			ChildID:  childID,
		},
	}
}

func deleteNodeChange(id string) Change {
	return Change{
		Type:  deleteNode,
		Value: id,
	}
}
