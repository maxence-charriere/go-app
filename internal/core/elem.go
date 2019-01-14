package core

import (
	"sync"

	"github.com/murlokswarm/app"
)

// Elem is a base struct to embed in app.Elem implementations.
type Elem struct {
	id  string
	err error
}

// ID satisfies the app.Elem interface.
func (e *Elem) ID() string {
	return e.id
}

// Contains satisfies the app.Elem interface.
func (e *Elem) Contains(app.Compo) bool {
	return false
}

// WhenView satisfies the app.Elem interface.
func (e *Elem) WhenView(func(app.View)) {}

// WhenWindow satisfies the app.Elem interface.
func (e *Elem) WhenWindow(func(app.Window)) {}

// WhenMenu satisfies the app.Elem interface.
func (e *Elem) WhenMenu(func(app.Menu)) {}

// WhenDockTile satisfies the app.Elem interface.
func (e *Elem) WhenDockTile(func(app.DockTile)) {}

// WhenStatusMenu satisfies the app.Elem interface.
func (e *Elem) WhenStatusMenu(func(app.StatusMenu)) {}

// Err satisfies the app.Elem interface.
func (e *Elem) Err() error {
	return e.err
}

// SetErr set the element error state with the given error.
// TODO: remove it.
func (e *Elem) SetErr(err error) {
	e.err = err
}

// NewElemDB creates an element database.
func NewElemDB() *ElemDB {
	return &ElemDB{
		elems:    make(map[string]app.Elem),
		elemList: make([]app.Elem, 0, 32),
	}
}

// ElemDB represents a element database.
type ElemDB struct {
	mutex    sync.RWMutex
	elems    map[string]app.Elem
	elemList []app.Elem
}

// Put inserts or update the given element into the database.
func (db *ElemDB) Put(e app.Elem) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if _, ok := db.elems[e.ID()]; !ok {
		db.elems[e.ID()] = e
		db.elemList = append(db.elemList, e)
	}
}

// Delete deletes the given element.
func (db *ElemDB) Delete(e app.Elem) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if _, ok := db.elems[e.ID()]; ok {
		delete(db.elems, e.ID())

		for i := range db.elemList {
			if db.elemList[i].ID() == e.ID() {
				elems := db.elemList
				copy(elems[i:], elems[i+1:])
				elems[len(elems)-1] = nil
				elems = elems[:len(elems)-1]

				db.elemList = elems
				return
			}
		}
	}
}

// GetByID returns the element with the given id.
func (db *ElemDB) GetByID(id string) app.Elem {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	if e, ok := db.elems[id]; ok {
		return e
	}

	return &Elem{err: app.ErrElemNotSet}
}

// GetByCompo returns the element where the given component is mounted.
func (db *ElemDB) GetByCompo(c app.Compo) app.Elem {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	for _, e := range db.elemList {
		if e.Contains(c) {
			return e
		}
	}

	return &Elem{err: app.ErrElemNotSet}
}
