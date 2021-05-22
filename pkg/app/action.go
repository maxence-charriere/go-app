package app

import (
	"fmt"
	"sync"
)

type Action struct {
	// The name that identifies the action..
	Name string

	// The value passed along with the action. Can be nil.
	Value interface{}

	// Tags that provide some context to the action.
	Tags Tags
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
	handlers map[string]map[string]actionHandler
}

func (m *actionManager) init() {
	m.handlers = make(map[string]map[string]actionHandler)
}

func (m *actionManager) post(a Action) {
	m.once.Do(m.init)
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

// ActionBuilder is the interface that describes a builder that builds and posts
// action.
type ActionBuilder interface {
	// Sets the action value.
	Value(v interface{}) ActionBuilder

	// Gives the action a tag with the given name and value. The value is
	// converted to a string.
	Tag(name string, v interface{}) ActionBuilder

	// Posts the action built. The action is then handled by handlers registered
	// with Handle() and Context.Handle().
	Post()
}

type actionBuilder struct {
	disp  Dispatcher
	name  string
	value interface{}
	tags  Tags
}

func newActionBuilder(d Dispatcher, actionName string) ActionBuilder {
	return &actionBuilder{
		disp: d,
		name: actionName,
	}
}

func (b *actionBuilder) Value(v interface{}) ActionBuilder {
	b.value = v
	return b
}

func (b *actionBuilder) Tag(name string, v interface{}) ActionBuilder {
	if b.tags == nil {
		b.tags = make(Tags)
	}
	b.tags.Set(name, v)
	return b
}

func (b *actionBuilder) Post() {
	b.disp.Post(Action{
		Name:  b.name,
		Value: b.value,
		Tags:  b.tags,
	})
}
