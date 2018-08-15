package html

import (
	"sync"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
)

// DOM describes a document object model store that manages node states.
type DOM struct {
	mutex           sync.Mutex
	factory         *app.Factory
	controlID       string
	root            *elemNode
	compoRowByID    map[string]compoRow
	compoRowByCompo map[app.Compo]compoRow
	hrefFmt         bool
}

// NewDOM create a document object model store.
func NewDOM(f *app.Factory, controlID string, hrefFmt bool) *DOM {
	return &DOM{
		factory:   f,
		controlID: controlID,
		root: &elemNode{
			id: "goapp-root",
		},
		compoRowByID:    make(map[string]compoRow),
		compoRowByCompo: make(map[app.Compo]compoRow),
		hrefFmt:         hrefFmt,
	}
}

// CompoByID returns the component with the given identifier.
func (d *DOM) CompoByID(id string) (app.Compo, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	r, ok := d.compoRowByID[id]
	if !ok {
		return nil, app.ErrCompoNotMounted
	}
	return r.compo, nil
}

// Contains reports whether the given component is in the dom.
func (d *DOM) Contains(c app.Compo) bool {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	_, ok := d.compoRowByCompo[c]
	return ok
}

func (d *DOM) insertCompoRow(r compoRow) {
	if sub, ok := r.compo.(app.Subscriber); ok {
		r.events = sub.Subscribe()
	}

	d.compoRowByID[r.id] = r
	d.compoRowByCompo[r.compo] = r
}

func (d *DOM) deleteCompoRow(id string) {
	if r, ok := d.compoRowByID[id]; ok {
		if r.events != nil {
			r.events.Close()
		}

		delete(d.compoRowByCompo, r.compo)
		delete(d.compoRowByID, id)
	}
}

// Len returns the amount of components in the DOM.
func (d *DOM) Len() int {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return len(d.compoRowByCompo)
}

// Render create or update the nodes of the given component.
// Component must be a pointer and not based on an empty struct.
func (d *DOM) Render(c app.Compo) ([]Change, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if err := validateCompo(c); err != nil {
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

	new, err := decodeCompo(c, d.hrefFmt)
	if err != nil {
		return nil, err
	}

	if err := d.syncNodes(old, new); err != nil {
		return nil, err
	}
	return parent.(node).ConsumeChanges(), nil
}

func (d *DOM) mountCompo(c app.Compo, parent *compoNode) error {
	root, err := decodeCompo(c, d.hrefFmt)
	if err != nil {
		return err
	}

	compoID := uuid.New().String()
	if err := d.mountNode(root, compoID); err != nil {
		return err
	}
	d.insertCompoRow(compoRow{
		id:    compoID,
		compo: c,
		root:  root,
	})

	if parent != nil {
		parent.SetRoot(root)
	}

	if mounter, ok := c.(app.Mounter); ok {
		mounter.OnMount()
	}
	return nil
}

func (d *DOM) mountNode(n node, compoID string) error {
	switch n := n.(type) {
	case *textNode:
		n.compoID = compoID
		n.controlID = d.controlID

	case *elemNode:
		n.compoID = compoID
		n.controlID = d.controlID
		n.changes = append(n.changes, mountElemChange(
			n.id,
			n.compoID,
		))

		for _, c := range n.children {
			if err := d.mountNode(c, compoID); err != nil {
				return err
			}
		}

	case *compoNode:
		n.compoID = compoID
		n.controlID = d.controlID

		c, err := d.factory.NewCompo(n.Name())
		if err != nil {
			return err
		}

		if err = mapCompoFields(c, n.fields); err != nil {
			return err
		}

		n.compo = c
		return d.mountCompo(c, n)
	}
	return nil
}

func (d *DOM) dismountCompo(c app.Compo) {
	row, _ := d.compoRowByCompo[c]
	d.dismountNode(row.root)
	d.deleteCompoRow(row.id)

	if dismounter, ok := c.(app.Dismounter); ok {
		dismounter.OnDismount()
	}
}

func (d *DOM) dismountNode(n node) {
	switch n := n.(type) {
	case *elemNode:
		for _, c := range n.children {
			d.dismountNode(c)
		}

	case *compoNode:
		d.dismountCompo(n.compo)
	}
}

func (d *DOM) syncNodes(old, new node) error {
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

func (d *DOM) syncTextNodes(old, new *textNode) error {
	if old.Text() != new.Text() {
		old.SetText(new.Text())
	}
	return nil
}

func (d *DOM) syncElemNodes(old, new *elemNode) error {
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

func (d *DOM) syncCompoNodes(old, new *compoNode) error {
	if old.Name() != new.Name() {
		return d.replaceNode(old, new)
	}

	if !attrsEqual(old.fields, new.fields) {
		old.fields = new.fields
		if err := mapCompoFields(old.compo, new.fields); err != nil {
			return err
		}

		newRoot, err := decodeCompo(old.compo, d.hrefFmt)
		if err != nil {
			return err
		}
		return d.syncNodes(old.root, newRoot)
	}
	return nil
}

func (d *DOM) replaceNode(old, new node) error {
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
