package html

import (
	"github.com/google/uuid"
	"github.com/murlokswarm/app"
)

// DOM is the interface that describes a document object model store that
// manages node states.
type DOM interface {
	// Returns the component with the given identifier.
	ComponentByID(id string) (app.Component, error)

	// Reports whether the given component is in the dom.
	ContainsComponent(c app.Component) bool

	// Create or update the nodes of the given component.
	Render(c app.Component) ([]Change, error)
}

// NewDOM create a document object model store.
func NewDOM(f app.Factory, controlID string) DOM {
	return &dom{
		factory:   f,
		controlID: controlID,
	}
}

type dom struct {
	factory         app.Factory
	controlID       string
	compoRowByID    map[string]compoRow
	compoRowByCompo map[app.Component]compoRow
}

func (d *dom) ComponentByID(id string) (app.Component, error) {
	r, ok := d.compoRowByID[id]
	if !ok {
		return nil, app.ErrNotFound
	}
	return r.component, nil
}

func (d *dom) ContainsComponent(c app.Component) bool {
	_, ok := d.compoRowByCompo[c]
	return ok
}

func (d *dom) insertCompoRow(r compoRow) {
	if sub, ok := r.component.(app.Subscriber); ok {
		r.events = sub.Subscribe()
	}

	d.compoRowByID[r.id] = r
	d.compoRowByCompo[r.component] = r
}

func (d *dom) deleteCompoRow(id string) {
	if r, ok := d.compoRowByID[id]; ok {
		if r.events != nil {
			r.events.Close()
		}

		delete(d.compoRowByCompo, r.component)
		delete(d.compoRowByID, id)
	}
}

// Render create or update the given component.
// It satisfies the app.DOM interface.
func (d *dom) Render(c app.Component) ([]Change, error) {
	panic("not implemented")
}

func (d *dom) mountComponent(c app.Component, isRoot bool) ([]Change, error) {
	id := uuid.New().String()

	root, err := decodeComponent(c)
	if err != nil {
		return nil, err
	}

	d.insertCompoRow(compoRow{
		id:        id,
		component: c,
		root:      root,
	})

	var changes []Change
	if changes, err = d.mountNode(root, id); err != nil {
		return nil, err
	}

	if isRoot {
		changes = append(changes, setRootChange(root.ID()))
	}

	if mounter, ok := c.(app.Mounter); ok {
		mounter.OnMount()
	}
	return changes, nil
}

func (d *dom) mountNode(n node, compoID string) ([]Change, error) {
	switch n := n.(type) {
	case *textNode:
		n.compoID = compoID
		n.controlID = d.controlID
		return []Change{createNodeChange(n)}, nil

	case *compoNode:
		n.compoID = compoID
		n.controlID = d.controlID

		compo, err := d.factory.New(n.Name())
		if err != nil {
			return nil, err
		}
		if err = mapComponentFields(compo, n.fields); err != nil {
			return nil, err
		}

		var changes []Change
		if changes, err = d.mountComponent(compo, false); err != nil {
			return nil, err
		}
		n.component = compo

		row, _ := d.compoRowByCompo[compo]
		n.setRoot(row.root)
		return changes, err

	case *elemNode:
		n.compoID = compoID
		n.controlID = d.controlID

		changes := make([]Change, 0, len(n.children)+1)
		changes = append(changes, createNodeChange(n))

		for _, c := range n.children {
			childChanges, err := d.mountNode(c, compoID)
			if err != nil {
				return nil, err
			}

			childID := c.ID()
			if n, ok := c.(*compoNode); ok {
				childID = n.root.ID()
			}

			changes = append(changes, childChanges...)
			changes = append(changes, appendChildChange(n.ID(), childID))
		}
		return changes, nil
	}
	return nil, nil
}

func (d *dom) dismountComponent(c app.Component) ([]Change, error) {
	r, ok := d.compoRowByCompo[c]
	if !ok {
		return nil, nil
	}

	changes, _ := d.dismountNode(r.root)

	d.deleteCompoRow(r.id)

	if dismounter, ok := c.(app.Dismounter); ok {
		dismounter.Render()
	}
	return changes, nil
}

func (d *dom) dismountNode(n node) ([]Change, error) {
	switch n := n.(type) {
	case *textNode:
		return []Change{deleteNodeChange(n.ID())}, nil

	case *compoNode:
		n.removeRoot()
		return d.dismountComponent(n.component)

	case *elemNode:
		changes := make([]Change, 0, 2*len(n.children)+1)
		for _, c := range n.children {
			childID := c.ID()
			if n, ok := c.(*compoNode); ok {
				childID = n.root.ID()
			}
			changes = append(changes, removeChildChange(n.ID(), childID))

			childChanges, _ := d.dismountNode(c)
			changes = append(changes, childChanges...)
		}

		changes = append(changes, deleteNodeChange(n.ID()))
		return changes, nil
	}
	return nil, nil
}

func (d *dom) syncNodes(current, new node) ([]Change, error) {
	switch current := current.(type) {
	case *textNode:
		if new, ok := new.(*textNode); ok {
			return d.syncTextNodes(current, new)
		}
		return d.replaceNode(current, new)

	case *compoNode:
		if new, ok := new.(*compoNode); ok {
			return d.syncCompoNodes(current, new)
		}
		return d.replaceNode(current, new)

	case *elemNode:
		if new, ok := new.(*elemNode); ok {
			return d.syncElemNodes(current, new)
		}
		return d.replaceNode(current, new)
	}
	return nil, nil
}

func (d *dom) syncTextNodes(current, new *textNode) ([]Change, error) {
	if current.text == new.text {
		return nil, nil
	}
	current.text = new.text
	return []Change{updateNodeChange(current)}, nil
}

func (d *dom) syncCompoNodes(current, new *compoNode) ([]Change, error) {
	if current.name != new.name {
		return d.replaceNode(current, new)
	}
	if attrsEqual(current.fields, new.fields) {
		return nil, nil
	}

	current.fields = new.fields
	if err := mapComponentFields(current.component, current.fields); err != nil {
		return nil, err
	}
	return d.Render(current.component)
}

func (d *dom) syncElemNodes(current, new *elemNode) ([]Change, error) {
	if current.tagName != new.tagName {
		return d.replaceNode(current, new)
	}

	var changes []Change

	curChildren := current.children
	newChildren := new.children

	// Sync children.
	for len(curChildren) != 0 && len(newChildren) != 0 {
		childChange, err := d.syncNodes(curChildren[0], newChildren[0])
		if err != nil {
			return nil, err
		}
		changes = append(changes, childChange...)

		newChildren = newChildren[:1]
	}

	// Remove children.
	for len(curChildren) != 0 {
		c := curChildren[0]

		childID := c.ID()
		if n, ok := c.(*compoNode); ok {
			childID = n.root.ID()
		}
		changes = append(changes, removeChildChange(current.ID(), childID))

		childChange, _ := d.dismountNode(c)
		changes = append(changes, childChange...)
		current.removeChild(c)

		curChildren = curChildren[:1]
	}

	// Append new children.
	for len(newChildren) != 0 {
		c := newChildren[0]
		current.appendChild(c)

		childChanges, err := d.mountNode(c, current.compoID)
		if err != nil {
			return nil, err
		}

		childID := c.ID()
		if n, ok := c.(*compoNode); ok {
			childID = n.root.ID()
		}

		changes = append(changes, childChanges...)
		changes = append(changes, appendChildChange(current.ID(), childID))
	}
	return changes, nil
}

func (d *dom) replaceNode(old, new node) ([]Change, error) {
	panic("not implemented")
}
