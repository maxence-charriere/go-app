package dom

// Change represents a node change that can be read by native platform to update
// their display.
type Change struct {
	// The change type.
	Type string

	// A value that describes how to make the change.
	Value interface{}
}

const (
	createText = "createText"
	setText    = "setText"

	createElem   = "createElem"
	setAttrs     = "setAttrs"
	appendChild  = "appendChild"
	removeChild  = "removeChild"
	replaceChild = "replaceChild"
	mountElem    = "mountElem"

	createCompo  = "createCompo"
	setCompoRoot = "setCompoRoot"

	deleteNode = "deleteNode"
)

type textValue struct {
	ID   string
	Text string `json:",omitempty"`
}

func createTextChange(id string) Change {
	return Change{
		Type: createText,
		Value: textValue{
			ID: id,
		},
	}
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

type elemValue struct {
	ID        string
	CompoID   string            `json:",omitempty"`
	TagName   string            `json:",omitempty"`
	Namespace string            `json:",omitempty"`
	Attrs     map[string]string `json:",omitempty"`
}

type childValue struct {
	ParentID string
	ChildID  string
	OldID    string `json:",omitempty"`
}

func createElemChange(id, tagName, namespace string) Change {
	return Change{
		Type: createElem,
		Value: elemValue{
			ID:        id,
			TagName:   tagName,
			Namespace: namespace,
		},
	}
}

func setAttrsChange(id string, a map[string]string) Change {
	return Change{
		Type: setAttrs,
		Value: elemValue{
			ID:    id,
			Attrs: a,
		},
	}
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

func replaceChildChange(parentID, oldID, newID string) Change {
	return Change{
		Type: replaceChild,
		Value: childValue{
			ParentID: parentID,
			ChildID:  newID,
			OldID:    oldID,
		},
	}
}

func mountElemChange(id, compoID string) Change {
	return Change{
		Type: mountElem,
		Value: elemValue{
			ID:      id,
			CompoID: compoID,
		},
	}
}

type compoValue struct {
	ID     string
	Name   string `json:",omitempty"`
	RootID string `json:",omitempty"`
}

func createCompoChange(id string, name string) Change {
	return Change{
		Type: createCompo,
		Value: compoValue{
			ID:   id,
			Name: name,
		},
	}
}

func setCompoRootChange(id, rootID string) Change {
	return Change{
		Type: setCompoRoot,
		Value: compoValue{
			ID:     id,
			RootID: rootID,
		},
	}
}

type deleteValue struct {
	ID string
}

func deleteNodeChange(id string) Change {
	return Change{
		Type: deleteNode,
		Value: deleteValue{
			ID: id,
		},
	}
}
