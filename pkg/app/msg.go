package app

import (
	"reflect"
	"sync"
	"time"

	"github.com/maxence-charriere/app/pkg/log"
)

// Binding represents a series of functions that are executed successively when
// a message is emitted.
type Binding struct {
	msg      string
	actions  []interface{}
	callOnUI func(func())
}

// Do adds a function to be executed when a message is emitted.
//
// Functions arguments are mapped to the emitted arguments or the return values
// of the previous function added with Do.
//
// Functions added with Do are executed on a different goroutine.
//
// It panics if f is not a function.
func (b *Binding) Do(f interface{}) *Binding {
	return b.do(f, false)
}

// DoOnUI is like Do but the added function is executed on the UI goroutine.
func (b *Binding) DoOnUI(f interface{}) *Binding {
	return b.do(f, true)
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
		callOnUI: callOnUI},
	)
	return b
}

// Wait adds the given duration before the execution of the next function added
// with Do.
func (b *Binding) Wait(d time.Duration) *Binding {
	b.actions = append(b.actions, d)
	return b
}

func (b *Binding) exec(args ...interface{}) {
	argsv := make([]reflect.Value, 0, len(args))
	for _, arg := range args {
		argsv = append(argsv, reflect.ValueOf(arg))
	}

	for idx, a := range b.actions {
		switch action := a.(type) {
		case time.Duration:
			time.Sleep(action)

		case do:
			ok := false
			if argsv, ok = b.execDo(idx, action, argsv); !ok {
				return
			}
		}
	}
}

func (b *Binding) execDo(idx int, do do, args []reflect.Value) ([]reflect.Value, bool) {
	fnval := do.function
	fntype := fnval.Type()

	i := 0
	for i < fntype.NumIn() && i < len(args) {
		if !args[i].Type().AssignableTo(fntype.In(i)) {
			log.Error("executing binding function failed").
				T("reason", "non assignable arg").
				T("message", b.msg).
				T("function index", idx).
				T("function type", fntype).
				T("arg index", i).
				T("expected type", fntype.In(i)).
				T("arg type", args[i].Type())
			return nil, false
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
	return args, true
}

type do struct {
	function reflect.Value
	callOnUI bool
}

type messenger struct {
	mutex    sync.RWMutex
	bindings map[string][]*Binding
}

func (m *messenger) emit(msg string, args ...interface{}) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, b := range m.bindings[msg] {
		b.exec(args...)
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
