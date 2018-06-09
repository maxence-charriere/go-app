package app

import (
	"sync"
)

// Control is the interface that describes an app control.
// eg. window, menu, etc.
type Control interface {
	// ID returns the control identifier.
	ID() string

	// Reports whether it contains the given control.
	Contains(c Component) bool
}

// ControlDB is the interface that describes an control database.
type ControlDB interface {
	// Add adds the control in the database.
	Add(e Control)

	// Remove removes the control from the database.
	Remove(e Control)

	// Control returns the control with the given identifier.
	ControlByID(id string) (Control, error)

	// ControlByComponent returns the control where the component is mounted.
	ControlByComponent(c Component) (Control, error)
}

// NewControlDB creates an control database safe for concurrent use.
// Should be used only in backend implementations.
func NewControlDB() ControlDB {
	return &controlDB{
		controls: make(map[string]Control),
	}
}

type controlDB struct {
	mutex    sync.RWMutex
	controls map[string]Control
}

func (db *controlDB) Add(e Control) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	db.controls[e.ID()] = e
}

func (db *controlDB) Remove(e Control) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	delete(db.controls, e.ID())
}

func (db *controlDB) ControlByID(id string) (Control, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	elem, ok := db.controls[id]
	if !ok {
		return nil, ErrNotFound
	}
	return elem, nil
}

func (db *controlDB) ControlByComponent(c Component) (Control, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	for _, elem := range db.controls {
		if elem.Contains(c) {
			return elem, nil
		}
	}
	return nil, ErrNotFound
}
