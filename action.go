package app

import (
	"sync"
)

// HandleAction handles the named action with the given handler.
func HandleAction(name string, h ActionHandler) {
	actions.Handle(name, h)
}

// PostAction creates and posts the named action with the given arg.
// The action is handled in its own goroutine.
func PostAction(name string, arg interface{}) {
	actions.Post(name, arg)
}

// PostActions creates and posts a batch of actions.
// All the actions are handled sequentially in a separate goroutine.
func PostActions(a ...Action) {
	actions.PostBatch(a...)
}

// Action represents an action to handle.
type Action struct {
	Name string
	Arg  interface{}
}

// ActionHandler represent an action handler.
type ActionHandler func(e EventDispatcher, a Action)

func newActionRegistry(dispatcher EventDispatcher) *actionRegistry {
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
