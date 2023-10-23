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

	// Executes the given function on the UI goroutine when the observed value
	// changes. Multiple functions can be executed by successively calling
	// OnChange().
	OnChange(fn func()) Observer

	// Stores the value associated with the observed state into the given
	// receiver. Panics when the receiver is not a pointer or nil.
	//
	// The receiver is updated each time the associated state changes. It is
	// unchanged when its pointed value has a different type than the associated
	// state value.
	Value(recv any)
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

	value     any
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
	onChanges  []func()
	receiver   any
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

func (o *observer) OnChange(fn func()) Observer {
	o.onChanges = append(o.onChanges, fn)
	return o
}

func (o *observer) Value(recv any) {
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

func (s *store) Set(key string, v any, opts ...StateOption) {
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
				WithTag("state", key).
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
				WithTag("state", key).
				Wrap(err))
			return
		}
	}

	for obs := range state.observers {
		o := obs

		if !o.element.Mounted() {
			delete(state.observers, o)
			continue
		}

		s.disp.Dispatch(Dispatch{
			Mode:   Update,
			Source: o.element,
			Function: func(ctx Context) {
				if !o.isObserving() {
					s.mutex.Lock()
					delete(state.observers, o)
					s.mutex.Unlock()
					return
				}

				if err := storeValue(o.receiver, v); err != nil {
					Log(errors.New("notifying observer failed").
						WithTag("state", key).
						WithTag("element", reflect.TypeOf(o.element)).
						Wrap(err))
					return
				}

				for _, fn := range o.onChanges {
					fn()
				}
			},
		})
	}
}

func (s *store) Get(key string, recv any) {
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
			WithTag("state", key).
			Wrap(err))
	}
}

func (s *store) Del(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.states, key)
	s.disp.getLocalStorage().Del(key)
}

func (s *store) Observe(key string, elem UI) Observer {
	return newObserver(elem, func(o *observer) {
		s.mutex.Lock()
		defer s.mutex.Unlock()

		if err := s.subscribe(key, o); err != nil {
			Log(errors.New("notifying observer failed").
				WithTag("state", key).
				WithTag("element", reflect.TypeOf(elem)).
				Wrap(err))
		}
	})
}

func (s *store) Cleanup() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.removeUnusedObservers()
	s.expireExpiredValues()
	s.expireExpriredPersistentValues()
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

func (s *store) getPersistent(key string, recv any) error {
	var state persistentState
	s.disp.getLocalStorage().Get(key, &state)

	if state.EncryptedValue == nil && state.Value == nil && state.ExpiresAt == (time.Time{}) {
		return nil
	}

	if state.isExpired(time.Now()) {
		s.disp.getLocalStorage().Del(key)
		return nil
	}

	if len(state.EncryptedValue) == 0 {
		return json.Unmarshal(state.Value, recv)
	}
	return s.disp.Context().Decrypt(state.EncryptedValue, recv)
}

func (s *store) setPersistent(key string, encrypt bool, expiresAt time.Time, v any) error {
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

	return s.disp.getLocalStorage().Set(key, state)
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

func (s *store) expireExpriredPersistentValues() {
	object := Window().Get("Object")
	if !object.Truthy() {
		return
	}

	keys := object.Call("keys", Window().Get("localStorage"))
	for i, l := 0, keys.Get("length").Int(); i < l; i++ {
		key := keys.Index(i).String()

		var state persistentState
		s.disp.getLocalStorage().Get(key, &state)

		if state.EncryptedValue == nil && state.Value == nil && state.ExpiresAt == (time.Time{}) {
			continue
		}

		if state.isExpired(time.Now()) {
			s.disp.getLocalStorage().Del(key)
		}
	}
}

func (s *store) expire(key string, state State) State {
	s.disp.getLocalStorage().Del(key)
	state.value = nil
	return state
}

func (s *store) initBroadcast() {
	broadcastChannel := Window().Get("BroadcastChannel")
	if !broadcastChannel.Truthy() {
		s.onBroadcastClose = func() {}
		return
	}
	broadcastChannel = broadcastChannel.New("go-app-broadcast-states")
	s.broadcastChannel = broadcastChannel

	onBroadcast := FuncOf(func(this Value, args []Value) any {
		s.onBroadcast(args[0].Get("data"))
		return nil
	})
	s.onBroadcastClose = onBroadcast.Release

	broadcastChannel.Set("onmessage", onBroadcast)
}

func (s *store) broadcast(key string, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	if s.broadcastChannel != nil {
		s.broadcastChannel.Call("postMessage", map[string]any{
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

	s.mutex.Lock()
	defer s.mutex.Unlock()

	key := event.Get("State").String()
	v := []byte(event.Get("Value").String())
	state := s.states[key]

	for obs := range state.observers {
		o := obs

		if !o.element.Mounted() {
			delete(state.observers, o)
			continue
		}

		s.disp.Dispatch(Dispatch{
			Mode:   Update,
			Source: o.element,
			Function: func(ctx Context) {
				if !o.isObserving() {
					s.mutex.Lock()
					delete(state.observers, o)
					s.mutex.Unlock()
					return
				}

				if err := json.Unmarshal(v, o.receiver); err != nil {
					Log(errors.New("notifying observer failed").
						WithTag("state", key).
						WithTag("element", reflect.TypeOf(o.element)).
						Wrap(err))
					return
				}

				for _, fn := range o.onChanges {
					fn()
				}
			},
		})
	}
}

type persistentState struct {
	Value          json.RawMessage `json:",omitempty"`
	EncryptedValue []byte          `json:",omitempty"`
	ExpiresAt      time.Time       `json:",omitempty"`
}

func (s *persistentState) isExpired(now time.Time) bool {
	return s.ExpiresAt != time.Time{} && now.After(s.ExpiresAt)
}

// -----------------------------------------------------------------------------

type StateX struct {
	value     any
	expiresAt time.Time

	ctx       Context
	name      string
	expire    func(StateX, time.Time) StateX
	persist   func(StateX, bool) StateX
	broadcast func(StateX) StateX
}

func (s StateX) ExpiresIn(v time.Duration) StateX {
	return s.expire(s, time.Now().Add(v))
}

func (s StateX) ExpiresAt(v time.Time) StateX {
	return s.expire(s, v)
}

func (s StateX) Persist() StateX {
	return s.persist(s, false)
}

func (s StateX) PersistWithEncryption() StateX {
	return s.persist(s, true)
}

func (s StateX) Broadcast() StateX {
	return s.broadcast(s)
}

type storableState struct {
	Value          json.RawMessage `json:",omitempty"`
	EncryptedValue []byte          `json:",omitempty"`
	ExpiresAt      time.Time       `json:",omitempty"`
}

type stateManager struct {
	mutex            sync.Mutex
	states           map[string]StateX
	observers        map[string]map[UI]observer
	broadcastStoreID string
	broadcastChannel Value
}

func (m *stateManager) Get(ctx nodeContext, state string, receiver any) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	value, exists := m.states[state]
	if !exists {
		if err := m.getStoredState(ctx, state, receiver); err != nil {
			Log(errors.New("getting state from local storage failed").
				WithTag("state", state).
				Wrap(err))
		}
		return
	}

	if !value.expiresAt.IsZero() && value.expiresAt.Before(time.Now()) {
		delete(m.states, state)
	}

	if err := storeValue(receiver, value.value); err != nil {
		Log(errors.New("getting state failed").
			WithTag("state", state).
			Wrap(err))
	}
}

func (m *stateManager) getStoredState(ctx Context, state string, receiver any) error {
	var value storableState
	if err := ctx.LocalStorage().Get(state, &value); err != nil {
		return err
	}

	if !value.ExpiresAt.IsZero() && value.ExpiresAt.Before(time.Now()) {
		ctx.LocalStorage().Del(state)
		return nil
	}

	if len(value.EncryptedValue) != 0 {
		return ctx.Decrypt(value.EncryptedValue, receiver)
	}
	return json.Unmarshal(value.Value, receiver)
}

func (m *stateManager) Set(ctx nodeContext, state string, v any) StateX {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.states == nil {
		m.states = make(map[string]StateX)
	}

	value := StateX{value: v}
	m.states[state] = value

	for _, observer := range m.observers[state] {
		o := observer
		ctx.Dispatch(func(ctx Context) {
			if !o.isObserving() {
				m.mutex.Lock()
				delete(m.observers[state], o.element)
				m.mutex.Unlock()
				return
			}

			if err := m.notifyChange(o, v); err != nil {
				Log(errors.New("notifying state change failed").
					WithTag("state", state).
					WithTag("observer-type", reflect.TypeOf(o.element)).
					Wrap(err))
			}
		})
	}

	return StateX{
		value:     v,
		ctx:       ctx,
		name:      state,
		expire:    m.setExpiration,
		persist:   m.persist,
		broadcast: m.broadcast,
	}
}

func (m *stateManager) notifyChange(o observer, v any) error {
	if err := storeValue(o.receiver, v); err != nil {
		return errors.New("storing value into receiver failed").Wrap(err)
	}

	for _, handleChange := range o.onChanges {
		handleChange()
	}
	return nil
}

func (m *stateManager) setExpiration(s StateX, v time.Time) StateX {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	s.expiresAt = v
	if !s.expiresAt.IsZero() && v.Before(time.Now()) {
		delete(m.states, s.name)
		return s
	}

	value := m.states[s.name]
	value.expiresAt = v
	m.states[s.name] = value

	return s
}

func (m *stateManager) persist(s StateX, encrypt bool) StateX {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if !s.expiresAt.IsZero() && s.expiresAt.Before(time.Now()) {
		s.ctx.LocalStorage().Del(s.name)
		return s
	}

	value := storableState{ExpiresAt: s.expiresAt}
	if encrypt {
		b, err := s.ctx.Encrypt(s.value)
		if err != nil {
			Log(errors.New("persisting encrypted state failed").
				WithTag("state", s.name).
				Wrap(err))
			return s
		}
		value.EncryptedValue = b
	} else {
		b, err := json.Marshal(s.value)
		if err != nil {
			Log(errors.New("persisting state failed").
				WithTag("state", s.name).
				Wrap(err))
			return s
		}
		value.Value = b
	}

	if err := s.ctx.LocalStorage().Set(s.name, value); err != nil {
		Log(errors.New("persisting state failed").
			WithTag("state", s.name).
			Wrap(err))
	}
	return s
}

func (m *stateManager) broadcast(s StateX) StateX {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.broadcastChannel == nil {
		Log(errors.New("persisting state failed").
			WithTag("state", s.name).
			Wrap(errors.New("broadcast is not supported")))
		return s
	}

	b, err := json.Marshal(s.value)
	if err != nil {
		Log(errors.New("persisting state failed").
			WithTag("state", s.name).
			Wrap(err))
		return s
	}

	m.broadcastChannel.Call("postMessage", map[string]any{
		"StoreID": m.broadcastStoreID,
		"State":   s.name,
		"Value":   string(b),
	})
	return s
}

func (m *stateManager) InitBroadcast(ctx Context) {
	broadcastChannel := Window().Get("BroadcastChannel")
	if !broadcastChannel.Truthy() {
		return
	}
	broadcastChannel = broadcastChannel.New("go-app-broadcast-states")
	m.broadcastChannel = broadcastChannel
	m.broadcastStoreID = uuid.NewString()

	handleBroadcast := FuncOf(func(this Value, args []Value) any {
		m.handleBroadcast(ctx, args[0].Get("data"))
		return nil
	})
	broadcastChannel.Set("onmessage", handleBroadcast)
}

func (m *stateManager) handleBroadcast(ctx Context, data Value) {
	if storeID := data.Get("StoreID").String(); storeID != m.broadcastStoreID {
		return
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	state := data.Get("State").String()
	value := []byte(data.Get("Value").String())

	for _, observer := range m.observers[state] {
		o := observer
		ctx.Dispatch(func(ctx Context) {
			if !o.isObserving() {
				m.mutex.Lock()
				delete(m.observers[state], o.element)
				m.mutex.Unlock()
				return
			}

			if err := m.notifyBroadcastChange(o, value); err != nil {
				Log(errors.New("notifying broadcast state change failed").
					WithTag("state", state).
					WithTag("observer-type", reflect.TypeOf(o.element)).
					Wrap(err))
			}
		})
	}
}

func (m *stateManager) notifyBroadcastChange(o observer, v []byte) error {
	if err := json.Unmarshal(v, o.receiver); err != nil {
		return errors.New("storing value into receiver failed").Wrap(err)
	}

	for _, handleChange := range o.onChanges {
		handleChange()
	}
	return nil
}

func storeValue(recv, v any) error {
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
			WithTag("value", src.Type()).
			WithTag("receiver", dst.Type())
	}
	dst.Set(src)
	return nil
}
