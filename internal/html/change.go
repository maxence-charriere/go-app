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
	createNode  = "createNode"
	deleteNode  = "deleteNode"
	appendChild = "appendChild"
	removeChild = "removeChild"
	updateNode  = "updateNode"
	setRoot     = "setRoot"
)

func createNodeChange(n node) Change {
	return Change{
		Type:  createNode,
		Value: n,
	}
}

func deleteNodeChange(id string) Change {
	return Change{
		Type:  deleteNode,
		Value: id,
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

func updateNodeChange(n node) Change {
	return Change{
		Type:  updateNode,
		Value: n,
	}
}

func setRootChange(id string) Change {
	return Change{
		Type:  setRoot,
		Value: id,
	}
}
