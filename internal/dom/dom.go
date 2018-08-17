package dom

import (
	"sync"

	"github.com/murlokswarm/app"
)

// DOM is a dom (document object model) engine that contains html nodes state.
// It is safe for concurrent operations.
type DOM struct {
	mutex        sync.Mutex
	factory      *app.Factory
	hrefFmt      bool
	compoByID    map[string]*component
	compoByCompo map[app.Compo]*component
	root         *elem
}

// NewDOM creates a dom engine.
func NewDOM(f *app.Factory, hrefFmt bool) *DOM {
	return &DOM{
		factory:      f,
		hrefFmt:      hrefFmt,
		compoByID:    make(map[string]*component),
		compoByCompo: make(map[app.Compo]*component),
		root:         &elem{id: "goapp-root"},
	}
}

// New creates the nodes for the given component and use it as its root.
func (dom *DOM) New(c app.Compo) ([]Change, error) {
	dom.mutex.Lock()
	defer dom.mutex.Unlock()

	panic("not implemented")
}

// Update updates the state of the given component.
func (dom *DOM) Update([]Change, error) {
	dom.mutex.Lock()
	defer dom.mutex.Unlock()

	panic("not implemented")
}

// Clean removes all the node from the dom, putting it clean state.
func (dom *DOM) Clean() error {
	dom.mutex.Lock()
	defer dom.mutex.Unlock()

	return dom.clean()
}

func (dom *DOM) clean() error {
	panic("not implemented")
}

type component struct {
	id     string
	root   node
	compo  app.Compo
	events *app.EventSubscriber
}
