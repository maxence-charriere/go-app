package app

import (
	"reflect"
	"sync"
	"time"

	"github.com/maxence-charriere/app/pkg/log"
)

// Binding represents a serie of actions that are executed successively when a
// message is emitted.
type Binding struct {
	msg           string
	actions       []interface{}
	callOnUI      func(func())
	deferDuration time.Duration
	deferEnd      time.Time

	mutex      sync.Mutex
	args       []interface{}
	deferTimer *time.Timer
}

// Do adds the given function to the actions of the binding. The function will
// be executed on the goroutine used to perform the binding.
//
// If the function is the first one to be added to the binding, its arguments
// are mapped with the ones emitted with the binded message. Otherwise with the
// return values of the previous function registered with Do or DoOnUI.
//
// It panics if f is not a function.
func (b *Binding) Do(f interface{}) *Binding {
	return b.do(f, false, false)
}

// DoOnUI adds the given function to the actions of the binding. The function
// will be executed on the UI goroutine.
//
// If the function is the first one to be added to the binding, its arguments
// are mapped with the ones emitted with the binded message. Otherwise with the
// return values of the previous function registered with Do or DoOnUI.
//
// It panics if f is not a function.
func (b *Binding) DoOnUI(f interface{}) *Binding {
	return b.do(f, true, false)
}

// State adds the given function to the actions of the binding. The function
// will be executed on the UI goroutine and a rendering of the component is
// automatically triggered.
//
// It is meant to execute functions that only modify the component appearance.
//
// Unlike DoOnUI, return values are not passed to the next action.
func (b *Binding) State(f interface{}) *Binding {
	return b.do(f, true, true)
}

func (b *Binding) do(f interface{}, callOnUI, noLink bool) *Binding {
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
		noLink:   noLink,
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
	argsv := make([]reflect.Value, 0, len(args))
	for _, arg := range args {
		argsv = append(argsv, reflect.ValueOf(arg))
	}

	for idx, a := range b.actions {
		switch action := a.(type) {
		case time.Duration:
			time.Sleep(action)

		case do:
			retsv, next := b.execDo(idx, action, argsv)
			if !next {
				return
			}
			if action.noLink {
				continue
			}
			argsv = retsv
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
	noLink   bool
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
