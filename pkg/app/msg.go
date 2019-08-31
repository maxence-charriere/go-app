package app

import (
	"context"
	"reflect"
	"sync"

	"github.com/maxence-charriere/app/pkg/log"
)

// Binding represents a series of functions that are executed successively when
// a message is emitted.
type Binding struct {
	msg   string
	funcs []reflect.Value
}

// Do adds a function to be executed when a message is emitted.
//
// Function first argument must implement the context.Context interface.
// https://golang.org/pkg/context/#Context
//
// Other functions arguments are mapped to the emitted arguments or the return
// values of the previous function added with Do.
//
// It panics if f is not a function or if f first argument does not implements
// context.Context.
func (b *Binding) Do(f interface{}) *Binding {
	v := reflect.ValueOf(f)
	typ := v.Type()

	if v.Kind() != reflect.Func {
		log.Error("adding function to binding failed").
			T("reason", "argument passed in not a function").
			T("argument type", typ).
			T("message", b.msg).
			T("index", len(b.funcs)).
			Panic()
	}
	if typ.NumIn() < 1 {
		log.Error("adding function to binding failed").
			T("reason", "function does not have a context argument").
			T("function type", typ).
			T("message", b.msg).
			T("index", len(b.funcs)).
			Panic()
	}

	if !typ.In(0).Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
		log.Error("adding function to binding failed").
			T("reason", "first argument does not implement context.Context").
			T("function type", typ).
			T("message", b.msg).
			T("index", len(b.funcs)).
			Panic()
	}

	b.funcs = append(b.funcs, v)
	return b
}

func (b *Binding) exec(ctx context.Context, args ...interface{}) {
	for _, f := range b.funcs {
		f.Call([]reflect.Value{reflect.ValueOf(ctx)})
	}
}

type messenger struct {
	mutex    sync.RWMutex
	bindings map[string][]*Binding
}

func (m *messenger) emit(ctx context.Context, msg string, args ...interface{}) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, b := range m.bindings[msg] {
		b.exec(ctx, args...)
	}
}

func (m *messenger) bind(msg string, c Compo) (*Binding, func()) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.bindings == nil {
		m.bindings = make(map[string][]*Binding)
	}

	b := &Binding{msg: msg}
	m.bindings[msg] = append(m.bindings[msg], b)

	close := func() {
		m.removeBinding(b)
	}

	return b, close
}

func (m *messenger) removeBinding(b *Binding) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	bindings := m.bindings[b.msg]
	for i, bind := range bindings {
		if bind == b {
			copy(bindings[i:], bindings[i+1:])
			bindings[len(bindings)-1] = nil
			bindings = bindings[:len(bindings)-1]
			break
		}
	}
	m.bindings[b.msg] = bindings
}
