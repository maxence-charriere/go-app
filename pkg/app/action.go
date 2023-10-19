package app

import (
	"fmt"
	"sync"
)

// Action represents a custom event that can be propagated across the app. It
// can carry a payload and be enriched with additional contextual tags.
type Action struct {
	// Name uniquely identifies the action.
	Name string

	// Value is the data associated with the action and can be nil.
	Value any

	// Tags provide additional context or metadata for the action.
	Tags Tags
}

// ActionHandler defines a callback executed when an action is triggered
// via Context.NewAction().
type ActionHandler func(Context, Action)

// Handle registers the provided handler for a specific action name. When that
// action is triggered, the handler executes in a separate goroutine.
func Handle(actionName string, h ActionHandler) {
	actionHandlers[actionName] = h
}

var actionHandlers = make(map[string]ActionHandler)

type actionHandler struct {
	Source   UI
	Function ActionHandler
	Async    bool
}

// actionManager manages the registration and execution of action handlers. It
// ensures that only actions related to mounted sources are processed.
type actionManager struct {
	mutex    sync.Mutex
	handlers map[string]map[string]actionHandler
}

// TODO: deprecate
func (m *actionManager) post(a Action) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	handlers := m.handlers[a.Name]
	for key, h := range handlers {
		source := h.Source
		if !source.Mounted() {
			delete(handlers, key)
			continue
		}

		ctx := makeContext(source)
		function := h.Function
		if h.Async {
			ctx.Async(func() { function(ctx, a) })
		} else {
			ctx.Dispatch(func(ctx Context) { function(ctx, a) })
		}
	}
}

func (m *actionManager) handle(actionName string, async bool, source UI, h ActionHandler) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.handlers == nil {
		m.handlers = make(map[string]map[string]actionHandler)
	}

	handlers, isRegistered := m.handlers[actionName]
	if !isRegistered {
		handlers = make(map[string]actionHandler)
		m.handlers[actionName] = handlers
	}

	key := fmt.Sprintf("/%T:%p/%p", source, source, h)
	handlers[key] = actionHandler{
		Source:   source,
		Function: h,
		Async:    async,
	}
}

func (m *actionManager) closeUnusedHandlers() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for actionName, handlers := range m.handlers {
		for key, h := range handlers {
			if !h.Source.Mounted() {
				delete(handlers, key)
			}
		}

		if len(handlers) == 0 {
			delete(m.handlers, actionName)
		}
	}
}

// Handle registers an ActionHandler for the given action and source.
func (m *actionManager) Handle(action string, source UI, async bool, handler ActionHandler) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.handlers == nil {
		m.handlers = make(map[string]map[string]actionHandler)
	}

	handlers, ok := m.handlers[action]
	if !ok {
		handlers = make(map[string]actionHandler)
		m.handlers[action] = handlers
	}

	key := actionHandlerKey(source, handler)
	handlers[key] = actionHandler{
		Source:   source,
		Function: handler,
		Async:    async,
	}
}

// Post processes the provided action by executing its associated handlers.
func (m *actionManager) Post(ctx Context, a Action) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for key, handler := range m.handlers[a.Name] {
		source := handler.Source
		if !source.Mounted() {
			delete(m.handlers[a.Name], key)
			continue
		}

		function := handler.Function
		if handler.Async {
			ctx.Async(func() {
				function(ctx, a)
			})
			continue
		}
		ctx.Dispatch(func(ctx Context) {
			function(ctx, a)
		})
	}
}

// Cleanup removes handlers corresponding to unmounted sources.
func (m *actionManager) Cleanup() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for action, handlers := range m.handlers {
		for key, handler := range handlers {
			if !handler.Source.Mounted() {
				delete(handlers, key)
			}

			if len(handlers) == 0 {
				delete(m.handlers, action)
			}
		}
	}
}

func actionHandlerKey(source UI, handler ActionHandler) string {
	return fmt.Sprintf("/%T/%p/%p", source, source, handler)
}
