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
	Contains(c app.Component) bool

	// The amount of components in the DOM.
	Len() int

	// Create or update the nodes of the given component.
	// Component must be a pointer and not based on an empty struct.
	Render(c app.Component) ([]Change, error)
}

// NewDOM create a document object model store.
func NewDOM(f app.Factory, controlID string) DOM {
	return newDOM(f, controlID)
}

func newDOM(f app.Factory, controlID string) *dom {
	return &dom{
		factory:   f,
		controlID: controlID,
		root: &elemNode{
			id: "goapp-root",
		},
		compoRowByID:    make(map[string]compoRow),
		compoRowByCompo: make(map[app.Component]compoRow),
	}
}

type dom struct {
	factory         app.Factory
	controlID       string
	root            *elemNode
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

func (d *dom) Contains(c app.Component) bool {
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

func (d *dom) Len() int {
	return len(d.compoRowByCompo)
}

func (d *dom) Render(c app.Component) ([]Change, error) {
	if err := validateComponent(c); err != nil {
		return nil, err
	}

	row, ok := d.compoRowByCompo[c]
	if !ok {
		if len(d.root.children) != 0 {
			c := d.root.children[0]
			d.dismountNode(c)
			d.root.removeChild(c)
		}

		// Mounting root component.
		if err := d.mountCompo(c, nil); err != nil {
			return nil, err
		}
		row, _ = d.compoRowByCompo[c]

		d.root.appendChild(row.root)
		return d.root.ConsumeChanges(), nil
	}

	old := row.root
	parent := old.Parent()

	new, err := decodeComponent(c)
	if err != nil {
		return nil, err
	}

	if err := d.syncNodes(old, new); err != nil {
		return nil, err
	}
	return parent.(node).ConsumeChanges(), nil
}

func (d *dom) mountCompo(c app.Component, parent *compoNode) error {
	root, err := decodeComponent(c)
	if err != nil {
		return err
	}

	compoID := uuid.New().String()
	if err := d.mountNode(root, compoID); err != nil {
		return err
	}
	d.insertCompoRow(compoRow{
		id:        compoID,
		component: c,
		root:      root,
	})

	if parent != nil {
		parent.SetRoot(root)
	}

	if mounter, ok := c.(app.Mounter); ok {
		mounter.OnMount()
	}
	return nil
}

func (d *dom) mountNode(n node, compoID string) error {
	switch n := n.(type) {
	case *textNode:
		n.compoID = compoID
		n.controlID = d.controlID

	case *elemNode:
		n.compoID = compoID
		n.controlID = d.controlID

		for _, c := range n.children {
			if err := d.mountNode(c, compoID); err != nil {
				return err
			}
		}

	case *compoNode:
		n.compoID = compoID
		n.controlID = d.controlID

		c, err := d.factory.New(n.Name())
		if err != nil {
			return err
		}

		if err = mapComponentFields(c, n.fields); err != nil {
			return err
		}

		n.component = c
		return d.mountCompo(c, n)
	}
	return nil
}

func (d *dom) dismountCompo(c app.Component) {
	row, _ := d.compoRowByCompo[c]
	d.dismountNode(row.root)
	d.deleteCompoRow(row.id)

	if dismounter, ok := c.(app.Dismounter); ok {
		dismounter.OnDismount()
	}
}

func (d *dom) dismountNode(n node) {
	switch n := n.(type) {
	case *elemNode:
		for _, c := range n.children {
			d.dismountNode(c)
		}

	case *compoNode:
		d.dismountCompo(n.component)
	}
}

func (d *dom) syncNodes(old, new node) error {
	switch old := old.(type) {
	case *textNode:
		if new, ok := new.(*textNode); ok {
			return d.syncTextNodes(old, new)
		}
		return d.replaceNode(old, new)

	case *compoNode:
		if new, ok := new.(*compoNode); ok {
			return d.syncCompoNodes(old, new)
		}
		return d.replaceNode(old, new)

	default:
		if new, ok := new.(*elemNode); ok {
			return d.syncElemNodes(old.(*elemNode), new)
		}
		return d.replaceNode(old, new)
	}
}

func (d *dom) syncTextNodes(old, new *textNode) error {
	if old.Text() != new.Text() {
		old.SetText(new.Text())
	}
	return nil
}

func (d *dom) syncElemNodes(old, new *elemNode) error {
	if old.TagName() != new.TagName() {
		return d.replaceNode(old, new)
	}

	if !attrsEqual(old.attrs, new.attrs) {
		old.SetAttrs(new.attrs)
	}

	oc := old.children
	nc := new.children

	// Sync children.
	for len(oc) != 0 && len(nc) != 0 {
		if err := d.syncNodes(oc[0], nc[0]); err != nil {
			return err
		}
		oc = oc[1:]
		nc = nc[1:]
	}

	// Remove children.
	for len(oc) != 0 {
		c := oc[0]
		d.dismountNode(c)
		old.removeChild(c)
		oc = oc[1:]
	}

	// Add children.
	for len(nc) != 0 {
		c := nc[0]
		if err := d.mountNode(c, old.CompoID()); err != nil {
			return err
		}
		old.appendChild(c)
		nc = nc[1:]
	}
	return nil
}

func (d *dom) syncCompoNodes(old, new *compoNode) error {
	if old.Name() != new.Name() {
		return d.replaceNode(old, new)
	}

	if !attrsEqual(old.fields, new.fields) {
		old.fields = new.fields
		if err := mapComponentFields(old.component, new.fields); err != nil {
			return err
		}

		newRoot, err := decodeComponent(old.component)
		if err != nil {
			return err
		}
		return d.syncNodes(old.root, newRoot)
	}
	return nil
}

func (d *dom) replaceNode(old, new node) error {
	d.dismountNode(old)

	if err := d.mountNode(new, old.CompoID()); err != nil {
		return err
	}

	switch p := old.Parent().(type) {
	case *elemNode:
		p.replaceChild(old, new)

	case *compoNode:
		p.RemoveRoot()
		p.SetRoot(new)
	}
	return nil
}
