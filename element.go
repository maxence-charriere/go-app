package app

import (
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Element is the interface that describes an app element.
type Element interface {
	// ID returns the element identifier.
	ID() uuid.UUID
}

// ElementWithComponent is the interface that describes an element that hosts
// components.
type ElementWithComponent interface {
	Element

	// Load loads the page specified by the URL.
	// URL can be formated as fmt package functions.
	// Calls with an URL which contains a component name will load the named
	// component.
	// e.g. hello will load the component named hello.
	// It returns an error if the component is not imported.
	Load(url string, v ...interface{}) error

	// Component returns the loaded component.
	Component() Component

	// Contains reports whether the component is mounted in the element.
	Contains(c Component) bool

	// Render renders the component.
	Render(c Component) error

	// LastFocus returns the last time when the element was focused.
	LastFocus() time.Time
}

// ElementWithNavigation is the interface that describes an element that
// supports navigation.
type ElementWithNavigation interface {
	ElementWithComponent

	// Reload reloads the current page.
	Reload() error

	// CanPrevious reports whether load the previous page is possible.
	CanPrevious() bool

	// Previous loads the previous page.
	// It returns an error if there is no previous page to load.
	Previous() error

	// CanNext indicates if loading next page is possible.
	CanNext() bool

	// Next loads the next page.
	// It returns an error if there is no next page to load.
	Next() error
}

// Closer is the interface that describes an element that can be closed.
type Closer interface {
	// Close closes the element and free its allocated resources.
	Close()
}

// NotificationConfig is a struct that describes a notification.
type NotificationConfig struct {
	Title     string `json:"title"`
	Subtitle  string `json:"subtitle"`
	Text      string `json:"text"`
	ImageName string `json:"image-name"`
	Sound     bool   `json:"sound"`

	OnReply func(reply string) `json:"-"`
}

// ElemDB is the interface that describes an element database.
type ElemDB interface {
	// Add adds the element in the database.
	Add(e Element) error

	// Remove removes the element from the database.
	Remove(e Element)

	// Element returns the element with the given identifier.
	Element(id uuid.UUID) (Element, error)

	// ElementByComponent returns the element where the component is mounted.
	// It returns an error if the component is not mounted in any element.
	ElementByComponent(c Component) (ElementWithComponent, error)

	// ElementsWithComponents returns the elements that contains components.
	ElementsWithComponents() []ElementWithComponent

	// Sort sorts the elements that hosts components. Latest focused elements
	// will be at the beginning.
	Sort()

	// Len returns the number of element.
	Len() int
}

// NewElemDB creates an element database.
func NewElemDB() ElemDB {
	return &elementDB{
		elements:               make(map[uuid.UUID]Element),
		elementsWithComponents: make(elementWithComponentList, 0, 64),
	}
}

// elementDB is an element database that implements ElemDB.
type elementDB struct {
	elements               map[uuid.UUID]Element
	elementsWithComponents elementWithComponentList
}

func (db *elementDB) Add(e Element) error {
	if _, ok := db.elements[e.ID()]; ok {
		return errors.Errorf("element with id %s is already added", e.ID())
	}

	db.elements[e.ID()] = e

	if elemWithComp, ok := e.(ElementWithComponent); ok {
		db.elementsWithComponents = append(db.elementsWithComponents, elemWithComp)
		sort.Sort(db.elementsWithComponents)
	}
	return nil
}

func (db *elementDB) Remove(e Element) {
	delete(db.elements, e.ID())

	if _, ok := e.(ElementWithComponent); ok {
		elements := db.elementsWithComponents
		for i, elem := range elements {
			if elem.ID() == e.ID() {
				copy(elements[i:], elements[i+1:])
				elements[len(elements)-1] = nil
				elements = elements[:len(elements)-1]
				db.elementsWithComponents = elements
				return
			}
		}
	}
}

func (db *elementDB) Element(id uuid.UUID) (Element, error) {
	e, ok := db.elements[id]
	if !ok {
		return nil, NewErrNotFound("element")
	}
	return e, nil
}

func (db *elementDB) ElementByComponent(c Component) (ElementWithComponent, error) {
	for _, elem := range db.elementsWithComponents {
		if elem.Contains(c) {
			return elem, nil
		}
	}

	return nil, errors.Errorf("component %+v is not mounted in any elements", c)
}

func (db *elementDB) ElementsWithComponents() []ElementWithComponent {
	elems := make([]ElementWithComponent, len(db.elementsWithComponents))
	copy(elems, db.elementsWithComponents)
	return elems
}

func (db *elementDB) Sort() {
	sort.Sort(db.elementsWithComponents)
}

func (db *elementDB) Len() int {
	return len(db.elements)
}

// NewConcurrentElemDB decorates the given element database to ensure concurrent
// access safety.
func NewConcurrentElemDB(db ElemDB) ElemDB {
	return &concurrentElemDB{
		base: db,
	}
}

// concurrentElemDB is a concurrent element database that implements
// ElemDB.
// It is safe for concurrent access.
type concurrentElemDB struct {
	mutex sync.Mutex
	base  ElemDB
}

func (db *concurrentElemDB) Add(e Element) error {
	db.mutex.Lock()
	err := db.base.Add(e)
	db.mutex.Unlock()
	return err
}

func (db *concurrentElemDB) Remove(e Element) {
	db.mutex.Lock()
	db.base.Remove(e)
	db.mutex.Unlock()
}

func (db *concurrentElemDB) Element(id uuid.UUID) (Element, error) {
	db.mutex.Lock()
	e, err := db.base.Element(id)
	db.mutex.Unlock()
	return e, err
}

func (db *concurrentElemDB) ElementByComponent(c Component) (ElementWithComponent, error) {
	db.mutex.Lock()
	e, err := db.base.ElementByComponent(c)
	db.mutex.Unlock()
	return e, err
}

func (db *concurrentElemDB) ElementsWithComponents() []ElementWithComponent {
	db.mutex.Lock()
	elems := db.base.ElementsWithComponents()
	db.mutex.Unlock()
	return elems
}

func (db *concurrentElemDB) Sort() {
	db.mutex.Lock()
	db.base.Sort()
	db.mutex.Unlock()
}

func (db *concurrentElemDB) Len() int {
	db.mutex.Lock()
	l := db.base.Len()
	db.mutex.Unlock()
	return l
}

// Slice of []ElementWithComponent that implements sort.Interface.
type elementWithComponentList []ElementWithComponent

func (l elementWithComponentList) Len() int {
	return len(l)
}

func (l elementWithComponentList) Less(i, j int) bool {
	return l[i].LastFocus().After(l[j].LastFocus())
}

func (l elementWithComponentList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
