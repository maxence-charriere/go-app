package dom

import (
	"sync"

	"github.com/murlokswarm/app"
)

// DOM is a dom (document object model) engine that contains html nodes state.
// It is safe for concurrent operations.
type DOM struct {
	mutex sync.Mutex
}

// NewDOM creates a dom engine.
func NewDOM() *DOM {
	return &DOM{}
}

// New creates the nodes for the given component and use it as its root.
func (dom *DOM) New(c app.Compo) error {
	panic("not implemented")
}

// Update updates the state of the given component.
func (dom *DOM) Update(c app.Compo) error {
	panic("not implemented")
}

// Clean removes all the node from the dom, putting it clean state.
func (dom *DOM) Clean() error {
	dom.mutex.Lock()
	err := dom.clean()
	dom.mutex.Unlock()

	return err
}

func (dom *DOM) clean() error {
	panic("not implemented")
}
