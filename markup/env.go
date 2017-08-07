package markup

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Env is the interface that describes an environment that handles components
// lifecycle.
type Env interface {
	// Component returns the component mounted under the identifier id.
	// err should be set if there is no mounted component under id.
	Component(id uuid.UUID) (c Componer, err error)

	// Root returns the root tag of component c.
	Root(c Componer) (root Tag, err error)

	// Mount indexes the component c into the env.
	// The component will live until it is dismounted.
	//
	// Mount should call the Render method from the component and create a tree
	// of Tag.
	// An id should be assigned to each tags.
	// Tags describing other components should trigger their creation and their
	// mount.
	Mount(c Componer) (root Tag, err error)

	// Dismount removes references to a component and its children.
	Dismount(c Componer)
}

// NewEnv creates an environment.
func NewEnv(b CompoBuilder) Env {
	return newEnv(b)
}

func newEnv(b CompoBuilder) *env {
	return &env{
		components:   make(map[uuid.UUID]Componer),
		compoRoots:   make(map[Componer]Tag),
		compoBuilder: b,
	}
}

type env struct {
	components   map[uuid.UUID]Componer
	compoRoots   map[Componer]Tag
	compoBuilder CompoBuilder
}

func (e *env) Component(id uuid.UUID) (c Componer, err error) {
	ok := false
	if c, ok = e.components[id]; !ok {
		err = errors.Errorf("no component with id %v is mounted", id)
	}
	return
}

func (e *env) Root(c Componer) (root Tag, err error) {
	ok := false
	if root, ok = e.compoRoots[c]; !ok {
		err = errors.Errorf("%T is not mounted", c)
	}
	return
}

func (e *env) Mount(c Componer) (root Tag, err error) {
	rootID := uuid.New()
	compoID := uuid.New()
	return e.mount(c, rootID, compoID)
}

func (e *env) mount(c Componer, rootID uuid.UUID, compoID uuid.UUID) (root Tag, err error) {
	if _, ok := e.compoRoots[c]; ok {
		err = errors.Errorf("%T is already mounted", c)
		return
	}

	if err = decodeComponent(c, &root); err != nil {
		err = errors.Wrapf(err, "fail to mount %T", c)
		return
	}

	if err = e.mountTag(&root, rootID, compoID); err != nil {
		err = errors.Wrapf(err, "fail to mount %T", c)
		return
	}

	e.components[compoID] = c
	e.compoRoots[c] = root

	if mounter, ok := c.(Mounter); ok {
		mounter.OnMount()
	}
	return
}

func (e *env) mountTag(t *Tag, id uuid.UUID, compoID uuid.UUID) error {
	t.ID = id
	t.CompoID = compoID

	if t.IsText() {
		return nil
	}

	if t.IsComponent() {
		c, err := e.compoBuilder.New(t.Name)
		if err != nil {
			return errors.Wrapf(err, "fail to mount %s", t.Name)
		}
		if err = mapComponentFields(c, t.Attrs); err != nil {
			return errors.Wrapf(err, "fail to mount %s", t.Name)
		}

		rootID := uuid.New()
		if _, err = e.mount(c, rootID, id); err != nil {
			return errors.Wrapf(err, "fail to mount %s", t.Name)
		}
		return nil
	}

	for i := range t.Children {
		childID := uuid.New()
		if err := e.mountTag(&t.Children[i], childID, compoID); err != nil {
			return errors.Wrapf(err, "fail to mount %s child", t.Name)
		}
	}
	return nil
}

func (e *env) Dismount(c Componer) {
	root, ok := e.compoRoots[c]
	if !ok {
		return
	}

	e.dismountTag(root)
	delete(e.components, root.CompoID)
	delete(e.compoRoots, c)

	if dismounter, ok := c.(Dismounter); ok {
		dismounter.OnDismount()
	}
	return
}

func (e *env) dismountTag(t Tag) {
	if t.IsComponent() {
		c, err := e.Component(t.ID)
		if err != nil {
			return
		}

		e.Dismount(c)
		return
	}

	for i := range t.Children {
		e.dismountTag(t.Children[i])
	}
	return
}

func (e *env) Update(c Componer) (syncs []Sync, err error) {
	syncs, _, err = e.update(c)
	return
}

func (e *env) update(c Componer) (syncs []Sync, syncParent bool, err error) {
	root, ok := e.compoRoots[c]
	if !ok {
		err = errors.Errorf("%T is not mounted", c)
		return
	}

	var newRoot Tag
	if err = decodeComponent(c, &newRoot); err != nil {
		err = errors.Wrapf(err, "fail to update %T", c)
		return
	}

	return e.syncTags(&root, &newRoot)
}

func (e *env) syncTags(l, r *Tag) (syncs []Sync, syncParent bool, err error) {
	if l.Name != r.Name {
		return e.mergeTags(l, r)
	}

	if l.IsText() {
		syncParent = e.syncTextTags(l, r)
		return
	}

	if l.IsComponent() {
		return e.syncComponentTags(l, r)
	}

	var subsyncs []Sync
	var fullsync bool

	if subsyncs, fullsync, err = e.syncTagChildren(l, r); err != nil {
		return
	}

	if !fullsync {
		syncs = append(syncs, subsyncs...)
	}

	if attrEq := AttrEquals(l.Attrs, r.Attrs); !attrEq || fullsync {
		if !attrEq {
			l.Attrs = r.Attrs
		}

		s := Sync{
			Tag:  *l,
			Full: fullsync,
		}
		syncs = append(syncs, s)
	}
	return
}

func (e *env) mergeTags(l, r *Tag) (syncs []Sync, syncParent bool, err error) {
	e.dismountTag(*l)
	if err = e.mountTag(r, l.ID, l.CompoID); err != nil {
		err = errors.Wrapf(err, "fail to merge %s and %s", l.Name, r.Name)
		return
	}

	*l = *r

	if l.IsText() {
		syncParent = true
		return
	}

	s := Sync{
		Tag:  *l,
		Full: true,
	}
	syncs = append(syncs, s)
	return
}

func (e *env) syncTextTags(l, r *Tag) (syncParent bool) {
	if l.Text != r.Text {
		l.Text = r.Text
		syncParent = true
	}
	return
}

func (e *env) syncComponentTags(l, r *Tag) (syncs []Sync, syncParent bool, err error) {
	if AttrEquals(l.Attrs, r.Attrs) {
		return
	}

	l.Attrs = r.Attrs

	c, err := e.Component(l.ID)
	if err != nil {
		err = errors.Wrapf(err, "fail to sync %s", l.Name)
		return
	}
	if err = mapComponentFields(c, l.Attrs); err != nil {
		err = errors.Wrapf(err, "fail to sync %s", l.Name)
		return
	}

	if syncs, syncParent, err = e.update(c); err != nil {
		err = errors.Wrapf(err, "fail to sync %s", l.Name)
	}
	return
}

func (e *env) syncTagChildren(l, r *Tag) (syncs []Sync, fullsync bool, err error) {
	lc := l.Children
	rc := r.Children
	count := 0

	for len(lc) != 0 && len(rc) != 0 {
		var subsyncs []Sync
		var sp bool

		if subsyncs, sp, err = e.syncTags(&lc[0], &rc[0]); err != nil {
			return
		}
		if sp {
			fullsync = true
			syncs = nil
		}
		if !fullsync {
			syncs = append(syncs, subsyncs...)
		}

		lc = lc[1:]
		rc = rc[1:]
		count++
	}

	l.Children = l.Children[:count]

	if len(lc) != len(rc) {
		fullsync = true
		syncs = nil
	}

	for len(lc) != 0 {
		e.dismountTag(lc[0])
		lc = lc[1:]
	}

	for len(rc) != 0 {
		child := &rc[0]
		childID := uuid.New()

		if err = e.mountTag(child, childID, l.CompoID); err != nil {
			return
		}
		l.Children = append(l.Children, *child)

		rc = rc[1:]
	}
	return
}

// Sync represents a sync operatrion.
type Sync struct {
	Tag  Tag
	Full bool
}
