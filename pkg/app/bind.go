package app

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/maxence-charriere/go-app/v5/pkg/log"
)

// Binding represents a serie of actions that are executed successively when a
// message is emitted.
type Binding struct {
	msg           string
	compo         Compo
	actions       []interface{}
	callOnUI      func(func())
	deferDuration time.Duration
	deferEnd      time.Time
	whenCancel    func(BindContext)

	mutex      sync.Mutex
	args       []interface{}
	deferTimer *time.Timer
}

// Do adds the given function to the actions of the binding. The function will
// be executed on the ui goroutine.
//
// If the function is the first one to be added to the binding, its arguments
// are mapped with the ones emitted with the binded message.
//
// The binding execution context can be passed if the function first argument is
// a BindContext or a context.Context.
//
// It panics if f is not a function.
func (b *Binding) Do(f interface{}) *Binding {
	return b.do(f, true)
}

// DoAsync adds the given function to the actions of the binding. The function
// will be executed on the goroutine used to perform the binding.
//
// If the function is the first one to be added to the binding, its arguments
// are mapped with the ones emitted with the binded message.
//
// The binding execution context can be passed if the function first argument is
// a BindContext or a context.Context.
//
// It panics if f is not a function.
func (b *Binding) DoAsync(f interface{}) *Binding {
	return b.do(f, false)
}

// WhenCancel setup the given function to be called when the binding is
// cancelled with a BindContext. The function will be executed on the UI
// goroutine.
func (b *Binding) WhenCancel(f func(BindContext)) {
	b.whenCancel = f
}

func (b *Binding) do(f interface{}, callOnUI bool) *Binding {
	v := reflect.ValueOf(f)
	if v.Kind() != reflect.Func {
		log.Error("adding function to binding failed").
			T("reason", "argument passed in not a function").
			T("argument type", v.Type()).
			T("message", b.msg).
			T("index", len(b.actions)).
			Panic()
	}

	b.actions = append(b.actions, do{
		function: v,
		callOnUI: callOnUI,
	})
	return b
}

// Wait delays the execution of the next function registered with Do or DoOnUI
// for the given duration.
func (b *Binding) Wait(d time.Duration) *Binding {
	b.actions = append(b.actions, d)
	return b
}

// Defer postpones the execution of the binding after the given duration. If a
// message occurs while the binding is deferred, arguments and duration are
// updated with the latest message emitted.
func (b *Binding) Defer(d time.Duration) *Binding {
	b.deferDuration = d
	return b
}

func (b *Binding) exec(args ...interface{}) {
	if b.deferDuration == 0 {
		b.execActions(args...)
		return
	}

	b.mutex.Lock()
	b.args = args
	if b.deferTimer != nil {
		b.deferTimer.Reset(b.deferDuration)
		b.mutex.Unlock()
		return
	}
	b.deferTimer = time.NewTimer(b.deferDuration)
	b.mutex.Unlock()

	<-b.deferTimer.C
	b.execActions(b.args...)
	b.deferTimer = nil
}

func (b *Binding) execActions(args ...interface{}) {
	ctx := newBindContext()
	defer ctx.Cancel(nil)
	ctxv := reflect.ValueOf(ctx)

	argsv := make([]reflect.Value, 0, len(args)+1)
	argsv = append(argsv, ctxv)
	for _, arg := range args {
		argsv = append(argsv, reflect.ValueOf(arg))
	}

	for idx, a := range b.actions {
		select {
		default:
		case <-ctx.Done():
			if b.whenCancel != nil {
				b.callOnUI(func() { b.whenCancel(ctx) })
			}
			return
		}

		switch action := a.(type) {
		case time.Duration:
			time.Sleep(action)

		case do:
			next := b.execDo(idx, action, argsv)
			if !next {
				return
			}
			argsv = []reflect.Value{ctxv}
		}
	}
}

func (b *Binding) execDo(idx int, do do, args []reflect.Value) bool {
	fnval := do.function
	fntype := fnval.Type()

	i := 0
	ctxHandled := false

	for i < fntype.NumIn() && i < len(args) {
		intype := fntype.In(i)

		// Ignore missing context in function args.
		if !ctxHandled && !intype.Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
			args = args[1:]
			ctxHandled = true
			continue
		} else {
			ctxHandled = true
		}

		if !args[i].Type().AssignableTo(intype) {
			log.Error("executing binding function failed").
				T("reason", "non assignable arg").
				T("message", b.msg).
				T("function index", idx).
				T("function type", fntype).
				T("arg index", i).
				T("expected type", intype).
				T("arg type", args[i].Type())
			return false
		}
		i++
	}

	for i < fntype.NumIn() {
		log.Warn("binding function argument missing").
			T("message", b.msg).
			T("function index", idx).
			T("function type", fntype).
			T("arg index", i).
			T("arg type", fntype.In(i)).
			T("pkg fix", "assigned to zero value")
		args = append(args, reflect.Zero(fntype.In(i)))
		i++
	}

	args = args[:i]
	if do.callOnUI {
		retchan := make(chan []reflect.Value, 1)
		defer close(retchan)

		b.callOnUI(func() {
			retchan <- fnval.Call(args)
		})
		args = <-retchan

	} else {
		args = fnval.Call(args)
	}
	return true
}

type do struct {
	function reflect.Value
	callOnUI bool
}

type messenger struct {
	mutex    sync.RWMutex
	bindings map[string][]*Binding
	callExec func(func(...interface{}), ...interface{})
	callOnUI func(func())
}

func (m *messenger) emit(msg string, args ...interface{}) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, b := range m.bindings[msg] {
		m.callExec(b.exec, args...)
	}
}

func (m *messenger) bind(msg string, c Compo) (*Binding, func()) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.bindings == nil {
		m.bindings = make(map[string][]*Binding)
	}

	b := &Binding{
		msg:      msg,
		compo:    c,
		callOnUI: m.callOnUI,
	}
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

// BindContext is the interface that describes a context passed when a binding
// is executed.
type BindContext interface {
	context.Context

	// Get returns a value from the current binding chain.
	Get(k string) interface{}

	// Lookup looks for a value from the current binding chain.
	Lookup(k string) (interface{}, bool)

	// Set sets a value in the current binding chain.
	Set(k string, v interface{})

	// Cancel stop the binding chain execution.
	Cancel(error)
}

type bindContext struct {
	context.Context

	mu     sync.RWMutex
	values map[string]interface{}
	cancel func()
	err    error
}

func newBindContext() *bindContext {
	ctx, cancel := context.WithCancel(context.Background())
	return &bindContext{
		Context: ctx,
		cancel:  cancel,
	}
}

func (ctx *bindContext) Get(k string) interface{} {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()

	v := ctx.values[k]
	return v
}

func (ctx *bindContext) Lookup(k string) (interface{}, bool) {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()

	v, ok := ctx.values[k]
	return v, ok
}

func (ctx *bindContext) Set(k string, v interface{}) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	if ctx.values == nil {
		ctx.values = make(map[string]interface{})
	}
	ctx.values[k] = v
}

func (ctx *bindContext) Cancel(err error) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	if err != nil {
		ctx.err = err
	}
	ctx.cancel()
}

func (ctx *bindContext) Err() error {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()

	if ctx.err != nil {
		return ctx.err
	}
	return ctx.Context.Err()
}
