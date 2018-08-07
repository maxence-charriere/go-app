package core

import (
	"sync"

	"github.com/murlokswarm/app"
)

// ElemWithCompo is the interface that describes an element that hosts
// components.
type ElemWithCompo interface {
	app.Elem

	// Contains reports whether the component is mounted in the element.
	Contains(app.Compo) bool
}

// Elem is a base struct to embed in app.Elem implementations.
type Elem struct {
	err error
}

// ID satisfies the app.Elem interface.
func (e *Elem) ID() string {
	return ""
}

// WhenWindow satisfies the app.Elem interface.
func (e *Elem) WhenWindow(func(app.Window)) {}

// WhenPage satisfies the app.Elem interface.
func (e *Elem) WhenPage(func(app.Page)) {}

// WhenNavigator satisfies the app.Elem interface.
func (e *Elem) WhenNavigator(func(app.Navigator)) {}

// WhenMenu satisfies the app.Elem interface.
func (e *Elem) WhenMenu(func(app.Menu)) {}

// WhenDockTile satisfies the app.Elem interface.
func (e *Elem) WhenDockTile(func(app.DockTile)) {}

// WhenStatusMenu satisfies the app.Elem interface.
func (e *Elem) WhenStatusMenu(func(app.StatusMenu)) {}

// WhenErr satisfies the app.Elem interface.
func (e *Elem) WhenErr(f func(error)) {
	if e.err != nil {
		f(e.err)
	}
}

// Err satisfies the app.Elem interface.
func (e *Elem) Err() error {
	return e.err
}

// SetErr set the element error state with the given error.
func (e *Elem) SetErr(err error) {
	e.err = err
}

// NewElemDB creates an element database.
func NewElemDB() *ElemDB {
	return &ElemDB{
		elems:          make(map[string]app.Elem),
		elemsWithCompo: make([]ElemWithCompo, 0, 32),
	}
}

// ElemDB represents a element database.
type ElemDB struct {
	mutex          sync.RWMutex
	elems          map[string]app.Elem
	elemsWithCompo []ElemWithCompo
}

// Put inserts or update the given element into the database.
func (db *ElemDB) Put(e app.Elem) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	db.elems[e.ID()] = e

	if ec, ok := e.(ElemWithCompo); ok {
		for i := range db.elemsWithCompo {
			if db.elemsWithCompo[i].ID() == ec.ID() {
				db.elemsWithCompo[i] = ec
				return
			}
		}

		db.elemsWithCompo = append(db.elemsWithCompo, ec)
	}
}

// Delete deletes the given element.
func (db *ElemDB) Delete(e app.Elem) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	delete(db.elems, e.ID())

	if ec, ok := e.(ElemWithCompo); ok {
		for i := range db.elemsWithCompo {
			if db.elemsWithCompo[i].ID() == ec.ID() {
				elems := db.elemsWithCompo
				copy(elems[i:], elems[i+1:])
				elems[len(elems)-1] = nil
				elems = elems[:len(elems)-1]

				db.elemsWithCompo = elems
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

	for _, e := range db.elemsWithCompo {
		if e.Contains(c) {
			return e
		}
	}

	return &Elem{err: app.ErrElemNotSet}
}
