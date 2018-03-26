package app

import "sync"

var (
	// DefaultActionRegistry is the default action registry.
	DefaultActionRegistry ActionRegistry
)

// Handle handles the named action with the given handler.
//
// It is a helper function that call DefaultActionRegistry.Handle.
func Handle(name string, h ActionHandler) {
	DefaultActionRegistry.Handle(name, h)
}

// NewAction creates and posts the named action with the given arg.
// The action is then handled in a separate goroutine.
//
// It is a helper function that call DefaultActionRegistry.Post.
func NewAction(name string, arg interface{}) {
	DefaultActionRegistry.Post(name, arg)
}

// NewActions creates and posts a batch of actions.
// All the actions are handled sequentially in a separate goroutine.
//
// It is a helper function that call DefaultActionRegistry.PostBatch.
func NewActions(a ...Action) {
	DefaultActionRegistry.PostBatch(a...)
}

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

// NewActionRegistry creates an action registry.
// Returned registry is safe for concurrent operations.
func NewActionRegistry(dispatcher EventDispatcher) ActionRegistry {
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

	if !ok {
		Error("no handler for action", a.Name)
		return
	}
	h(r.dispatcher, a)
}

func (r *actionRegistry) PostBatch(a ...Action) {
	go func() {
		for _, action := range a {
			r.exec(action)
		}
	}()
}

// ActionRegistryWithLogs returns a decorated version of the given action
// registry that logs its operations.
func ActionRegistryWithLogs(r ActionRegistry) ActionRegistry {
	return &actionRegistryWithLogs{
		base: r,
	}
}

type actionRegistryWithLogs struct {
	base ActionRegistry
}

func (r *actionRegistryWithLogs) Handle(name string, h ActionHandler) {
	Logf("action %s is handled", name)
	r.base.Handle(name, h)
}

func (r *actionRegistryWithLogs) Post(name string, arg interface{}) {
	Logf("posting action %s %+v", name, arg)
	r.base.Post(name, arg)
}

func (r *actionRegistryWithLogs) PostBatch(a ...Action) {
	Logf("posting batch of actions %+v", a)
	r.base.PostBatch(a...)
}
