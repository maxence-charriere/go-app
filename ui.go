package app

import (
	"log"
	"runtime"
	"sync"

	"github.com/murlokswarm/markup"
	"github.com/murlokswarm/uid"
	"github.com/pkg/errors"
)

var (
	// UIChan is a channel which take a func as payload.
	// When the app package is initialised, it creates a goroutine dedicated to
	// execute UI related tasks in order to avoid to deal with concurrency when
	// programming a component.
	// UIChan allows to enqueue UI related tasks that should be executed in this
	// goroutine.
	// When implementing a driver, driver and component callbacks should be
	// called through this channel.
	UIChan = make(chan func(), 256)

	// Elements contains the elements in use by the app.
	Elements ElementStorer
)

// Elementer is the interface that describes an app element.
// It wraps the ID method which allow to keep and retrieve a element.
type Elementer interface {
	ID() uid.ID
}

// ElementStorer is the interface that describes a element store.
// It keep a reference of the elements used by the app and allows to retrieve
// them when needed.
// Implementation should be thread safe.
type ElementStorer interface {
	// Add adds an element in the store.
	// Should be called in a driver implementation when a native element is
	// created.
	Add(e Elementer) error

	// Remove removes an element from the store.
	// Should be called in a driver implementation when a native element is
	// closed or removed.
	Remove(e Elementer)

	// Len returns the numbers of element in use.
	Len() int

	// Get returns the element with id.
	// ok will be false if there is no element with id.
	Get(id uid.ID) (e Elementer, ok bool)
}

// Contexter is the interface that describes an element where a component can be
// mounted and rendered. e.g. a window.
type Contexter interface {
	Elementer

	// Mount mounts the component c into the context. The context is displayed
	// with the default appearance of the component.
	Mount(c Componer)

	// Render renders the DOM node targeted into the sync description s.
	// s contains info specifying if the node and its children should be
	// replaced or just updated.
	// Should be called in a driver implementation.
	Render(s markup.Sync)
}

// Context returns the context where c is mounted.
// Panic if there is no context where c is mounted.
func Context(c Componer) Contexter {
	root := markup.Root(c)

	elem, ok := Elements.Get(root.ContextID)
	if !ok {
		log.Panicf("no context with id %v in use", root.ContextID)
	}
	return elem.(Contexter)
}

func init() {
	runtime.LockOSThread()
	go startUIGoroutine()

	Elements = newElementStore()
}

type elementStore struct {
	mutex sync.Mutex
	elems map[uid.ID]Elementer
}

func newElementStore() *elementStore {
	return &elementStore{
		elems: map[uid.ID]Elementer{},
	}
}

func (s *elementStore) Add(e Elementer) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, ok := s.Get(e.ID()); ok {
		return errors.Errorf("[%s %T] is already registered", e.ID(), e)
	}
	s.elems[e.ID()] = e
	return nil
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

func (s *elementStore) Get(id uid.ID) (e Elementer, ok bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	e, ok = s.elems[id]
	return
}
