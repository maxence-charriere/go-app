package app

import (
	"encoding/json"
	"reflect"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

// Observer is an observer that observes changes for a given state.
type Observer interface {
	// Defines a condition that reports whether the observer keeps observing the
	// associated state. Multiple conditions can be defined by successively
	// calling While().
	While(condition func() bool) Observer

	// Stores the value associated with the observed state into the given
	// receiver. Panics when the receiver is not a pointer or nil.
	//
	// The receiver is updated each time the associated state changes. It is
	// unchanged when its pointed value has a different type than the associated
	// state value.
	Value(recv interface{})
}

// A state represents an observable value available across the app.
type State struct {
	// Reports whether the state is persisted in local storage.
	IsPersistent bool

	// Reports whether the state is encrypted before being persisted in local
	// storage.
	IsEncrypted bool

	// The time when the state expires. The state never expires when zero value.
	ExpiresAt time.Time

	// Reports whether a state is broadcasted to other browser tabs and windows.
	IsBroadcasted bool

	value     interface{}
	observers map[*observer]struct{}
}

func (s *State) isExpired(now time.Time) bool {
	return s.ExpiresAt != time.Time{} && now.After(s.ExpiresAt)
}

// StateOption represents an option applied when a state is set.
type StateOption func(*State)

// Persist is a state option that persists a state in local storage.
//
// Be mindful to not use this option as a cache since local storage is limited
// to 5MB in a lot of web browsers.
func Persist(s *State) {
	s.IsPersistent = true
}

// Encrypt is a state option that encrypts a state before persisting it in local
// storage. Encryption is performed only when the Persist option is also set.
func Encrypt(s *State) {
	s.IsEncrypted = true
}

// ExpiresIn returns a state option that sets a state value to its zero value
// after the given duration.
//
// Values persisted to local storage with the Persist option are removed from
// it.
func ExpiresIn(d time.Duration) StateOption {
	return ExpiresAt(time.Now().Add(d))
}

// ExpiresAt returns a state option that sets a state value to its zero value at
// the given time.
//
// Values persisted to local storage with the Persist option are removed from
// it.
func ExpiresAt(t time.Time) StateOption {
	return func(s *State) {
		s.ExpiresAt = t
	}
}

// Broadcast is a state option that broadcasts a state to other browser tabs and
// windows from the same origin.
func Broadcast(s *State) {
	s.IsBroadcasted = true
}

type observer struct {
	element    UI
	subscribe  func(*observer)
	conditions []func() bool
	receiver   interface{}
}

func newObserver(elem UI, subscribe func(*observer)) *observer {
	return &observer{
		element:   elem,
		subscribe: subscribe,
	}
}

func (o *observer) While(fn func() bool) Observer {
	o.conditions = append(o.conditions, fn)
	return o
}

func (o *observer) Value(recv interface{}) {
	if reflect.ValueOf(recv).Kind() != reflect.Ptr {
		panic(errors.New("observer value receiver is not a pointer"))
	}

	o.receiver = recv
	o.subscribe(o)
}

func (o *observer) isObserving() bool {
	if !o.element.Mounted() {
		return false
	}

	for _, c := range o.conditions {
		if !c() {
			return false
		}
	}

	return true
}

type store struct {
	mutex            sync.Mutex
	id               string
	states           map[string]State
	disp             Dispatcher
	broadcastChannel Value
	onBroadcastClose func()
}

func newStore(d Dispatcher) *store {
	s := &store{
		id:     uuid.NewString(),
		states: make(map[string]State),
		disp:   d,
	}

	s.initBroadcast()
	return s
}

func (s *store) Set(key string, v interface{}, opts ...StateOption) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	state := s.states[key]
	state.value = v
	for _, o := range opts {
		o(&state)
	}
	s.states[key] = state

	if state.IsPersistent {
		if err := s.setPersistent(key, state.IsEncrypted, state.ExpiresAt, v); err != nil {
			Log(errors.New("persisting state failed").
				Tag("state", key).
				Wrap(err))
			return
		}
	}

	if state.isExpired(time.Now()) {
		state = s.expire(key, state)
		s.states[key] = state
		return
	}

	if state.IsBroadcasted {
		if err := s.broadcast(key, v); err != nil {
			Log(errors.New("broadcasting state failed").
				Tag("state", key).
				Wrap(err))
			return
		}
	}

	for o := range state.observers {
		if !o.isObserving() {
			delete(state.observers, o)
			continue
		}

		elem := o.element
		recv := o.receiver
		s.disp.Dispatch(elem, func(ctx Context) {
			if err := storeValue(recv, v); err != nil {
				Log(errors.New("notifying observer failed").
					Tag("state", key).
					Tag("element", reflect.TypeOf(elem)).
					Wrap(err))
			}
		})
	}
}

func (s *store) Get(key string, recv interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var err error
	state := s.states[key]
	if state.isExpired(time.Now()) {
		state = s.expire(key, state)
		s.states[key] = state
	}

	if state.value != nil {
		err = storeValue(recv, state.value)
	} else {
		err = s.getPersistent(key, recv)
	}
	if err != nil {
		Log(errors.New("getting state value failed").
			Tag("state", key).
			Wrap(err))
	}
}

func (s *store) Del(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.states, key)
	s.disp.localStorage().Del(key)
}

func (s *store) Observe(key string, elem UI) Observer {
	return newObserver(elem, func(o *observer) {
		s.mutex.Lock()
		defer s.mutex.Unlock()

		if err := s.subscribe(key, o); err != nil {
			Log(errors.New("notifying observer failed").
				Tag("state", key).
				Tag("element", reflect.TypeOf(elem)).
				Wrap(err))
		}
	})
}

func (s *store) Cleanup() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.removeUnusedObservers()
	s.expireExpiredValues()
}

func (s *store) Close() {
	if s.broadcastChannel != nil {
		s.broadcastChannel.Call("close")
		s.broadcastChannel = nil
	}
	s.onBroadcastClose()
}

func (s *store) subscribe(key string, o *observer) error {
	state := s.states[key]
	if state.observers == nil {
		state.observers = make(map[*observer]struct{})
	}
	state.observers[o] = struct{}{}

	if state.isExpired(time.Now()) {
		state = s.expire(key, state)
	}
	s.states[key] = state

	if state.value != nil {
		return storeValue(o.receiver, state.value)
	}
	return s.getPersistent(key, o.receiver)
}

func (s *store) removeUnusedObservers() {
	for _, state := range s.states {
		for o := range state.observers {
			if !o.isObserving() {
				delete(state.observers, o)
			}
		}
	}
}

func (s *store) getPersistent(key string, recv interface{}) error {
	var state persistentState
	s.disp.localStorage().Get(key, &state)

	if state.EncryptedValue == nil && state.Value == nil && state.ExpiresAt == (time.Time{}) {
		return nil
	}

	if state.isExpired(time.Now()) {
		s.disp.localStorage().Del(key)
		return nil
	}

	if len(state.EncryptedValue) == 0 {
		return json.Unmarshal(state.Value, recv)
	}
	return s.disp.Context().Decrypt(state.EncryptedValue, recv)
}

func (s *store) setPersistent(key string, encrypt bool, expiresAt time.Time, v interface{}) error {
	var err error

	state := persistentState{
		ExpiresAt: expiresAt,
	}
	if encrypt {
		state.EncryptedValue, err = s.disp.Context().Encrypt(v)
	} else {
		state.Value, err = json.Marshal(v)
	}
	if err != nil {
		return err
	}

	return s.disp.localStorage().Set(key, state)
}

func (s *store) expireExpiredValues() {
	now := time.Now()
	for k, state := range s.states {
		if state.isExpired(now) {
			state = s.expire(k, state)
			s.states[k] = state
		}
	}
}

func (s *store) expire(key string, state State) State {
	s.disp.localStorage().Del(key)
	state.value = nil
	return state
}

func (s *store) initBroadcast() {
	broadcastChannel := Window().Get("BroadcastChannel").New("go-app-broadcast-states")
	if !broadcastChannel.Truthy() {
		s.onBroadcastClose = func() {}
		return
	}
	s.broadcastChannel = broadcastChannel

	onBroadcast := FuncOf(func(this Value, args []Value) interface{} {
		s.mutex.Lock()
		defer s.mutex.Unlock()

		s.onBroadcast(args[0].Get("data"))
		return nil
	})
	s.onBroadcastClose = onBroadcast.Release

	broadcastChannel.Set("onmessage", onBroadcast)
}

func (s *store) broadcast(key string, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	if s.broadcastChannel != nil {
		s.broadcastChannel.Call("postMessage", map[string]interface{}{
			"StoreID": s.id,
			"State":   key,
			"Value":   string(b),
		})
	}
	return nil
}

func (s *store) onBroadcast(event Value) {
	if storeID := event.Get("StoreID").String(); storeID == "" || storeID == s.id {
		return
	}

	key := event.Get("State").String()
	v := []byte(event.Get("Value").String())
	state := s.states[key]

	for o := range state.observers {
		if !o.isObserving() {
			delete(state.observers, o)
			continue
		}

		elem := o.element
		recv := o.receiver
		s.disp.Dispatch(elem, func(ctx Context) {
			if err := json.Unmarshal(v, recv); err != nil {
				Log(errors.New("notifying observer failed").
					Tag("state", key).
					Tag("element", reflect.TypeOf(elem)).
					Wrap(err))
			}
		})
	}
}

func storeValue(recv, v interface{}) error {
	dst := reflect.ValueOf(recv)
	if dst.Kind() != reflect.Ptr {
		return errors.New("receiver is not a pointer")
	}
	dst = dst.Elem()

	src := reflect.ValueOf(v)
	switch {
	case src == reflect.Value{}:
		dst.Set(reflect.Zero(dst.Type()))
		return nil

	case src.Kind() == reflect.Ptr:
		src = src.Elem()
	}

	if src.Type() != dst.Type() {
		return errors.New("value and receiver are not of the same type").
			Tag("value", src.Type()).
			Tag("receiver", dst.Type())
	}

	dst.Set(src)
	return nil
}

type persistentState struct {
	Value          json.RawMessage `json:",omitempty"`
	EncryptedValue []byte          `json:",omitempty"`
	ExpiresAt      time.Time       `json:",omitempty"`
}

func (s *persistentState) isExpired(now time.Time) bool {
	return s.ExpiresAt != time.Time{} && now.After(s.ExpiresAt)
}

type broadcastState struct {
	StoreID string          `json:""`
	Value   json.RawMessage `json:",omitempty"`
}