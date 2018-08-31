package dom

import (
	"sync"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/pkg/errors"
)

// DOM is a dom (document object model) engine that contains html nodes state.
// It is safe for concurrent operations.
type DOM struct {
	mutex        sync.Mutex
	factory      *app.Factory
	transforms   []Transform
	compoByID    map[string]*component
	compoByCompo map[app.Compo]*component
	root         *elem
}

// NewDOM creates a dom engine.
func NewDOM(f *app.Factory, t ...Transform) *DOM {
	return &DOM{
		factory:      f,
		transforms:   t,
		compoByID:    make(map[string]*component),
		compoByCompo: make(map[app.Compo]*component),
		root: &elem{
			id:      "root:",
			tagName: "body",
		},
	}
}

// CompoByID returns the component with the given identifier.
func (dom *DOM) CompoByID(id string) (app.Compo, error) {
	dom.mutex.Lock()
	defer dom.mutex.Unlock()

	r, ok := dom.compoByID[id]
	if !ok {
		return nil, app.ErrCompoNotMounted
	}

	return r.compo, nil
}

// Contains reports whether the given component is in the dom.
func (dom *DOM) Contains(c app.Compo) bool {
	dom.mutex.Lock()
	defer dom.mutex.Unlock()

	_, ok := dom.compoByCompo[c]
	return ok
}

func (dom *DOM) insertCompo(c *component) {
	if sub, ok := c.compo.(app.Subscriber); ok {
		c.events = sub.Subscribe()
	}

	dom.compoByID[c.id] = c
	dom.compoByCompo[c.compo] = c
}

func (dom *DOM) deleteCompo(c *component) {
	if c.events != nil {
		c.events.Close()
	}

	delete(dom.compoByCompo, c.compo)
	delete(dom.compoByID, c.id)
}

// Len returns the amount of components in the DOM.
func (dom *DOM) Len() int {
	dom.mutex.Lock()
	defer dom.mutex.Unlock()

	return len(dom.compoByID)
}

// New creates the nodes for the given component and use it as its root.
func (dom *DOM) New(c app.Compo) ([]Change, error) {
	dom.mutex.Lock()
	defer dom.mutex.Unlock()

	dom.clean()

	if err := validateCompo(c); err != nil {
		return nil, err
	}

	if err := dom.mountCompo(c, nil); err != nil {
		return nil, errors.Wrap(err, "mounting compo failed")
	}

	compo := dom.compoByCompo[c]
	dom.root.appendChild(compo.root)
	return dom.root.Flush(), nil
}

func (dom *DOM) mountCompo(c app.Compo, parent *compo) error {
	root, err := decodeCompo(c, dom.transforms...)
	if err != nil {
		return errors.Wrap(err, "decoding compo failed")
	}

	compoID := app.CompoName(c) + ":" + uuid.New().String()
	if parent != nil {
		compoID = parent.id
	}

	if err = dom.mountNode(root, compoID); err != nil {
		dom.dismountNode(root)
		return err
	}

	dom.insertCompo(&component{
		id:    compoID,
		root:  root,
		compo: c,
	})

	if parent != nil {
		parent.SetRoot(root)
	}

	if mounter, ok := c.(app.Mounter); ok {
		mounter.OnMount()
	}

	return nil
}

func (dom *DOM) mountNode(n node, compoID string) error {
	switch n := n.(type) {
	case *text:
		n.compoID = compoID

	case *elem:
		n.compoID = compoID
		n.changes = append(n.changes, mountElemChange(n.id, n.compoID))

		for _, c := range n.children {
			if err := dom.mountNode(c, compoID); err != nil {
				return err
			}
		}

	case *compo:
		n.compoID = compoID

		c, err := dom.factory.NewCompo(n.name)
		if err != nil {
			return err
		}

		if err = mapCompoFields(c, n.fields); err != nil {
			return err
		}

		return dom.mountCompo(c, n)
	}

	return nil
}

// Update updates the state of the given component.
func (dom *DOM) Update(c app.Compo) ([]Change, error) {
	dom.mutex.Lock()
	defer dom.mutex.Unlock()

	if err := validateCompo(c); err != nil {
		return nil, err
	}

	compo, ok := dom.compoByCompo[c]
	if !ok {
		return nil, app.ErrCompoNotMounted
	}

	old := compo.root
	p := old.Parent()

	new, err := decodeCompo(c, dom.transforms...)
	if err != nil {
		return nil, errors.Wrap(err, "decoding compo failed")
	}

	if err = dom.updateNode(old, new); err != nil {
		return nil, errors.Wrapf(err, "updating %s with %s failed", old.ID(), new.ID())
	}

	return p.(node).Flush(), nil
}

func (dom *DOM) updateNode(old, new node) error {
	switch old := old.(type) {
	case *text:
		if new, ok := new.(*text); ok {
			return dom.updateText(old, new)
		}

	case *elem:
		if new, ok := new.(*elem); ok {
			return dom.updateElem(old, new)
		}

	case *compo:
		if new, ok := new.(*compo); ok {
			return dom.updateCompo(old, new)
		}
	}

	return dom.replaceNode(old, new)
}

func (dom *DOM) updateText(old, new *text) error {
	if old.text != new.text {
		old.SetText(new.text)
	}

	return nil
}

func (dom *DOM) updateElem(old, new *elem) error {
	if old.tagName != new.tagName {
		return dom.replaceNode(old, new)
	}

	if !attrsEqual(old.attrs, new.attrs) {
		old.SetAttrs(new.attrs)
	}

	oc := old.children
	nc := new.children

	// Sync children.
	for len(oc) != 0 && len(nc) != 0 {
		if err := dom.updateNode(oc[0], nc[0]); err != nil {
			return err
		}

		oc = oc[1:]
		nc = nc[1:]
	}

	// Remove children.
	for len(oc) != 0 {
		c := oc[0]
		dom.dismountNode(c)
		old.removeChild(c)
		oc = oc[1:]
	}

	// Add children.
	for len(nc) != 0 {
		c := nc[0]
		if err := dom.mountNode(c, old.CompoID()); err != nil {
			return err
		}

		old.appendChild(c)
		nc = nc[1:]
	}

	return nil
}

func (dom *DOM) updateCompo(old, new *compo) error {
	if old.name != new.name {
		return dom.replaceNode(old, new)
	}

	if !attrsEqual(old.fields, new.fields) {
		old.fields = new.fields
		c := dom.compoByID[old.id]

		if err := mapCompoFields(c.compo, old.fields); err != nil {
			return err
		}

		newRoot, err := decodeCompo(c.compo, dom.transforms...)
		if err != nil {
			return err
		}

		return dom.updateNode(old.root, newRoot)
	}

	return nil
}

func (dom *DOM) replaceNode(old, new node) error {
	dom.dismountNode(old)

	if err := dom.mountNode(new, old.CompoID()); err != nil {
		return err
	}

	switch p := old.Parent().(type) {
	case *elem:
		p.replaceChild(old, new)

	case *compo:
		p.RemoveRoot()
		p.SetRoot(new)

		c, _ := dom.compoByID[old.CompoID()]
		c.root = new
	}

	return nil
}

// Clean removes all the node from the dom, putting it clean state.
func (dom *DOM) Clean() {
	dom.mutex.Lock()
	defer dom.mutex.Unlock()

	dom.clean()
}

func (dom *DOM) clean() {
	if len(dom.root.children) != 0 {
		dom.dismountCompo(dom.root.children[0])
	}
}

func (dom *DOM) dismountCompo(root node) {
	if c, ok := dom.compoByID[root.CompoID()]; ok {
		dom.dismountNode(root)
		dom.deleteCompo(c)

		if dismounter, ok := c.compo.(app.Dismounter); ok {
			dismounter.OnDismount()
		}
	}
}

func (dom *DOM) dismountNode(n node) {
	switch n := n.(type) {
	case *elem:
		for _, c := range n.children {
			dom.dismountNode(c)
		}

	case *compo:
		if n.root != nil {
			dom.dismountCompo(n.root)
		}
	}
}

type component struct {
	id     string
	root   node
	compo  app.Compo
	events *app.EventSubscriber
}
