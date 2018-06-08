package app

import (
	"sync"
)

// Action represents an action to handle.
type Action struct {
	Name string
	Arg  interface{}
}

// ActionHandler represent an action handler.
type ActionHandler func(e EventDispatcher, a Action)

// ActionRegistry is the interface that describes an action registry.
type ActionRegistry interface {
	// Handler handles the named action with the given handler.
	Handle(name string, h ActionHandler)

	// Post posts the named action with the given arg.
	// The action is then handled in a separate goroutine.
	Post(name string, arg interface{})

	// PostBatch posts a batch of actions.
	// All the actions are handled sequentially in a separate goroutine.
	PostBatch(a ...Action)
}

func newActionRegistry(dispatcher EventDispatcher) ActionRegistry {
	return &actionRegistry{
		actions:    make(map[string]ActionHandler),
		dispatcher: dispatcher,
	}
}

type actionRegistry struct {
	mutex      sync.RWMutex
	actions    map[string]ActionHandler
	dispatcher EventDispatcher
}

func (r *actionRegistry) Handle(name string, h ActionHandler) {
	r.mutex.Lock()
	r.actions[name] = h
	r.mutex.Unlock()
}

func (r *actionRegistry) Post(name string, arg interface{}) {
	go func() {
		r.exec(Action{
			Name: name,
			Arg:  arg,
		})
	}()
}

func (r *actionRegistry) exec(a Action) {
	r.mutex.RLock()
	h, ok := r.actions[a.Name]
	r.mutex.RUnlock()

	if ok {
		h(r.dispatcher, a)
	}
}

func (r *actionRegistry) PostBatch(a ...Action) {
	go func() {
		for _, action := range a {
			r.exec(action)
		}
	}()
}
