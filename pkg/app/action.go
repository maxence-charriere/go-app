package app

import (
	"context"
	"fmt"
	"sync"
)

type Action struct {
	// The name that identifies the action..
	Name string

	// The value passed along with the action. Can be nil.
	Value interface{}
}

// ActionHandler represents a handler that is executed when an action is created
// with Context.NewAction().
type ActionHandler func(Context, Action)

// Handle registers the handler for the given action name. When an action
// occurs, the handler is executed on its own goroutine.
func Handle(actionName string, h ActionHandler) {
	actionHandlers[actionName] = h
}

var actionHandlers = make(map[string]ActionHandler)

type actionHandler struct {
	async    bool
	source   UI
	function ActionHandler
}

type actionManager struct {
	once     sync.Once
	mutex    sync.Mutex
	queue    chan Action
	stop     func()
	handlers map[string]map[string]actionHandler
}

func (m *actionManager) init() {
	m.queue = make(chan Action, 128)
	m.handlers = make(map[string]map[string]actionHandler)

	ctx, cancel := context.WithCancel(context.Background())
	m.stop = cancel

	go func() {
		defer close(m.queue)
		defer cancel()

		for {
			select {
			case <-ctx.Done():
				return

			case action := <-m.queue:
				m.execute(action)
			}
		}
	}()
}

func (m *actionManager) post(actionName string, v interface{}) {
	m.once.Do(m.init)
	m.queue <- Action{
		Name:  actionName,
		Value: v,
	}
}

func (m *actionManager) execute(a Action) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	handlers := m.handlers[a.Name]
	for key, h := range handlers {
		source := h.source
		if !source.Mounted() {
			delete(handlers, key)
			continue
		}

		ctx := makeContext(source)
		function := h.function

		if h.async {
			ctx.Async(func() { function(ctx, a) })
		} else {
			ctx.Dispatch(func(ctx Context) { function(ctx, a) })
		}
	}
}

func (m *actionManager) handle(actionName string, async bool, source UI, h ActionHandler) {
	m.once.Do(m.init)
	m.mutex.Lock()
	defer m.mutex.Unlock()

	handlers, isRegistered := m.handlers[actionName]
	if !isRegistered {
		handlers = make(map[string]actionHandler)
		m.handlers[actionName] = handlers
	}

	key := fmt.Sprintf("/%T:%p/%p", source, source, h)
	handlers[key] = actionHandler{
		async:    async,
		source:   source,
		function: h,
	}
}

func (m *actionManager) closeUnusedHandlers() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for actionName, handlers := range m.handlers {
		for key, h := range handlers {
			if !h.source.Mounted() {
				delete(handlers, key)
			}
		}

		if len(handlers) == 0 {
			delete(m.handlers, actionName)
		}
	}
}

func (m *actionManager) close() {
	if m.stop != nil {
		m.stop()
	}
	m.handlers = nil
}
