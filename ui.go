package app

import (
	"log"
	"runtime"
	"sync"

	"github.com/murlokswarm/markup"
	"github.com/satori/go.uuid"
)

var (
	// UIChan is a channel which take a func as payload.
	// When the app package is initialized, it creates a goroutine dedicated to
	// execute UI related tasks in order to avoid to deal with concurrency when
	// programming a component.
	// UIChan allows to enqueue UI related tasks that should be executed in this
	// goroutine.
	// When implementing a driver, driver and component callbacks should be
	// called through this channel.
	UIChan = make(chan func(), 256)

	elements ElementStorer
)

// Elementer is the interface that describes an app element.
// It wraps the ID method which allow to keep and retrieve a element.
//
// Driver implementation:
// - Elements().Add() should be called when an element is created.
// - Elements().Remove() should be called when an element is closed.
type Elementer interface {
	// ID returns the identifier of the element.
	ID() uuid.UUID
}

// ElementStorer is the interface that describes a element store.
// It keep a reference of the elements used by the app and allows to retrieve
// them when needed.
// Implementation should be thread safe.
type ElementStorer interface {
	// Add adds an element in the store.
	// Should be called in a driver implementation when a native element is
	// created.
	Add(e Elementer)

	// Remove removes an element from the store.
	// Should be called in a driver implementation when a native element is
	// closed or removed.
	Remove(e Elementer)

	// Len returns the numbers of element in use.
	Len() int

	// Get returns the element with id.
	// ok will be false if there is no element with id.
	Get(id uuid.UUID) (e Elementer, ok bool)
}

// Contexter is the interface that describes an element where a component can be
// mounted and rendered. e.g. a window.
type Contexter interface {
	Elementer

	// Mount mounts the component c into the context. The context is displayed
	// with the default appearance of the component.
	//
	// Driver implementation:
	// - Should call markup.Mount()
	Mount(c Componer)

	// Component returns the component mounted with Mount().
	// Returns nil if there is no component mounted.
	Component() Componer

	// Render renders the DOM node targeted into the sync description s.
	// s contains info specifying if the node and its children should be
	// replaced or just updated.
	// Should be called in a driver implementation.
	Render(s markup.Sync)
}

// Elements returns the element store containing all the elements in use by the
// app.
func Elements() ElementStorer {
	return elements
}

// Context returns the context where c is mounted.
// Panic if there is no context where c is mounted.
func Context(c Componer) Contexter {
	root := markup.Root(c)

	elem, ok := Elements().Get(root.ContextID)
	if !ok {
		log.Panicf("no context with id %v in use", root.ContextID)
	}
	return elem.(Contexter)
}

func init() {
	runtime.LockOSThread()
	go startUIGoroutine()

	elements = newElementStore()
}

func startUIGoroutine() {
	for f := range UIChan {
		f()
	}
}

type elementStore struct {
	mutex sync.Mutex
	elems map[uuid.UUID]Elementer
}

func newElementStore() *elementStore {
	return &elementStore{
		elems: map[uuid.UUID]Elementer{},
	}
}

func (s *elementStore) Add(e Elementer) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, ok := s.elems[e.ID()]; ok {
		log.Panicf("[%s %T] is already registered", e.ID(), e)
	}
	s.elems[e.ID()] = e
}

func (s *elementStore) Remove(e Elementer) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.elems, e.ID())
}

func (s *elementStore) Len() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return len(s.elems)
}

func (s *elementStore) Get(id uuid.UUID) (e Elementer, ok bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	e, ok = s.elems[id]
	return
}
