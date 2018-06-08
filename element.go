package app

import (
	"sync"
)

// Element is the interface that describes an app element.
type Element interface {
	// ID returns the element identifier.
	ID() string

	// Reports whether it contains the given element.
	Contains(c Component) bool
}

// ElementDB is the interface that describes an element database.
type ElementDB interface {
	// Add adds the element in the database.
	Add(e Element)

	// Remove removes the element from the database.
	Remove(e Element)

	// Element returns the element with the given identifier.
	ElementByID(id string) (Element, error)

	// ElementByComponent returns the element where the component is mounted.
	ElementByComponent(c Component) (Element, error)
}

// NewElementDB creates an element database safe for concurrent use.
// Should be used only in backend implementations.
func NewElementDB() ElementDB {
	return &elementDB{
		elements: make(map[string]Element),
	}
}

type elementDB struct {
	mutex    sync.RWMutex
	elements map[string]Element
}

func (db *elementDB) Add(e Element) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	db.elements[e.ID()] = e
}

func (db *elementDB) Remove(e Element) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	delete(db.elements, e.ID())
}

func (db *elementDB) ElementByID(id string) (Element, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	elem, ok := db.elements[id]
	if !ok {
		return nil, ErrNotFound
	}
	return elem, nil
}

func (db *elementDB) ElementByComponent(c Component) (Element, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	for _, elem := range db.elements {
		if elem.Contains(c) {
			return elem, nil
		}
	}
	return nil, ErrNotFound
}
