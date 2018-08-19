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
	hrefFmt      bool
	compoByID    map[string]*component
	compoByCompo map[app.Compo]*component
	root         *elem
}

// NewDOM creates a dom engine.
func NewDOM(f *app.Factory, hrefFmt bool) *DOM {
	return &DOM{
		factory:      f,
		hrefFmt:      hrefFmt,
		compoByID:    make(map[string]*component),
		compoByCompo: make(map[app.Compo]*component),
		root:         &elem{id: "root:"},
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
	root, err := decodeCompo(c, dom.hrefFmt)
	if err != nil {
		return errors.Wrap(err, "decoding compo failed")
	}

	compoID := "compo-" + uuid.New().String()

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

	new, err := decodeCompo(c, dom.hrefFmt)
	if err != nil {
		return nil, errors.Wrap(err, "decoding compo failed")
	}

	if err = dom.updateNode(old, new); err != nil {
		return nil, errors.Wrapf(err, "updating %s with %s failed", old.ID(), new.ID())
	}

	return p.(node).Flush(), nil
}

func (dom *DOM) updateNode(old, new node) error {
	panic("not implemented")
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
