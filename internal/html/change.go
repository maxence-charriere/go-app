package html

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

func setTextChange(text string) Change {
	return Change{
		Type:  setText,
		Value: text,
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

func appendChildChange(parentID, childID string) Change {
	return Change{
		Type: appendChild,
		Value: struct {
			ParentID string
			ChildID  string
		}{
			ParentID: parentID,
			ChildID:  childID,
		},
	}
}

func removeChildChange(parentID, childID string) Change {
	return Change{
		Type: removeChild,
		Value: struct {
			ParentID string
			ChildID  string
		}{
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
