package db

import (
	"sort"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/markup"
	"github.com/pkg/errors"
)

// ElementDB is an element database that implements app.ElementDB.
type ElementDB struct {
	capacity               int
	elements               map[uuid.UUID]app.Element
	elementsWithComponents elementWithComponentList
}

// NewElementDB creates an element database with given capacity.
func NewElementDB(capacity int) *ElementDB {
	return &ElementDB{
		capacity:               capacity,
		elements:               make(map[uuid.UUID]app.Element, capacity),
		elementsWithComponents: make(elementWithComponentList, 0, capacity),
	}
}

// Add satisfies the app.ElementDB interface.
func (db *ElementDB) Add(e app.Element) error {
	if len(db.elements) == db.capacity {
		return errors.Errorf("can't handle more than %d elements simultaneously", db.capacity)
	}

	if _, ok := db.elements[e.ID()]; ok {
		return errors.Errorf("element with id %s is already added", e.ID())
	}

	db.elements[e.ID()] = e

	if elemWithComp, ok := e.(app.ElementWithComponent); ok {
		db.elementsWithComponents = append(db.elementsWithComponents, elemWithComp)
		sort.Sort(db.elementsWithComponents)
	}
	return nil
}

// Remove satisfies the app.ElementDB interface.
func (db *ElementDB) Remove(e app.Element) {
	delete(db.elements, e.ID())

	if _, ok := e.(app.ElementWithComponent); ok {
		elements := db.elementsWithComponents
		for i, elem := range elements {
			if elem == e {
				copy(elements[i:], elements[i+1:])
				elements[len(elements)-1] = nil
				elements = elements[:len(elements)-1]
				db.elementsWithComponents = elements
				return
			}
		}
	}
}

// Element satisfies the app.ElementDB interface.
func (db *ElementDB) Element(id uuid.UUID) (e app.Element, ok bool) {
	e, ok = db.elements[id]
	return
}

// ElementByComponent satisfies the app.ElementDB interface.
func (db *ElementDB) ElementByComponent(c markup.Component) (e app.ElementWithComponent, err error) {
	for _, elem := range db.elementsWithComponents {
		if elem.Contains(c) {
			e = elem
			return
		}
	}

	err = errors.Errorf("component %+v is not mounted in any elements", c)
	return
}

// Sort satisfies the app.ElementDB interface.
func (db *ElementDB) Sort() {
	sort.Sort(db.elementsWithComponents)
}

// Len satisfies the app.ElementDB interface.
func (db *ElementDB) Len() int {
	return len(db.elements)
}

// ConcurentElemDB is a concurent element database that implements
// app.ElementDB.
// It is safe for multiple goroutines to call its methods concurrently.
type ConcurentElemDB struct{}

// TO DO:
// Implementation.

// Slice of []ElementWithComponent that implements sort.Interface.
type elementWithComponentList []app.ElementWithComponent

func (l elementWithComponentList) Len() int {
	return len(l)
}

func (l elementWithComponentList) Less(i, j int) bool {
	return l[i].LastFocus().After(l[j].LastFocus())
}

func (l elementWithComponentList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
