package core

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
)

// ElemWithCompo is the interface that describes an element that hosts
// components.
type ElemWithCompo interface {
	app.Elem

	// Contains reports whether the component is mounted in the element.
	Contains(app.Compo) bool

	// Render the given component.
	Render(app.Compo) error
}

// Elem is a base struct to embed in an app.Elem implementations.
type Elem struct {
	notSet bool
}

// ID satisfies the app.Elem interface.
func (e *Elem) ID() uuid.UUID {
	return uuid.UUID{}
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

// WhenNotSet satisfies the app.Elem interface.
func (e *Elem) WhenNotSet(f func()) {
	if e.notSet {
		f()
	}
}

// IsNotSet satisfies the app.Elem interface.
func (e *Elem) IsNotSet() bool {
	return e.notSet
}

// Contains satisfies the ElemWithCompo interface.
func (e *Elem) Contains(c app.Compo) bool {
	return false
}

// Render satisfies the ElemWithCompo interface.
func (e *Elem) Render(app.Compo) error {
	if e.notSet {
		return errors.New("not set")
	}

	return nil
}

// NewElemDB creates an element database.
func NewElemDB() *ElemDB {
	return &ElemDB{
		elems:          make(map[uuid.UUID]app.Elem),
		elemsWithCompo: make([]ElemWithCompo, 0, 32),
	}
}

// ElemDB represents a element database.
type ElemDB struct {
	mutex          sync.RWMutex
	elems          map[uuid.UUID]app.Elem
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
func (db *ElemDB) GetByID(id uuid.UUID) app.Elem {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	if e, ok := db.elems[id]; ok {
		return e
	}

	return &Elem{notSet: true}
}

// GetByCompo returns the element where the given component is mounted.
func (db *ElemDB) GetByCompo(c app.Compo) ElemWithCompo {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	for _, e := range db.elemsWithCompo {
		if e.Contains(c) {
			return e
		}
	}

	return &Elem{notSet: true}
}
