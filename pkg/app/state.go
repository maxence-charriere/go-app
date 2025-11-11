package app

import (
	"encoding/json"
	"reflect"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/maxence-charriere/go-app/v10/pkg/errors"
)

// State represents a state with additional features such as expiration,
// persistence, and broadcasting capabilities.
type State struct {
	value     any
	expiresAt time.Time

	ctx       Context
	name      string
	expire    func(State, time.Time) State
	persist   func(State, bool) State
	broadcast func(State) State
}

// ExpiresIn sets the expiration time for the state by specifying a duration
// from the current time.
func (s State) ExpiresIn(v time.Duration) State {
	return s.expire(s, time.Now().Add(v))
}

// ExpiresAt sets the exact expiration time for the state.
func (s State) ExpiresAt(v time.Time) State {
	return s.expire(s, v)
}

// Persist ensures the state is persisted into the local storage.
func (s State) Persist() State {
	return s.persist(s, false)
}

// PersistWithEncryption ensures the state is persisted into the local storage
// with encryption.
func (s State) PersistWithEncryption() State {
	return s.persist(s, true)
}

// Broadcast signals that changes to the state will be broadcasted to other
// browser tabs and windows sharing the same origin when it is supported.
//
// Using Broadcast creates a BroadcastChannel, which prevents the page from
// being cached. This may impact the Chrome Lighthouse performance score due to
// the additional resources required to manage the broadcast channel.
func (s State) Broadcast() State {
	return s.broadcast(s)
}

type storableState struct {
	Value          json.RawMessage `json:",omitempty"`
	EncryptedValue []byte          `json:",omitempty"`
	ExpiresAt      time.Time       `json:",omitempty"`
}

// Observer represents a mechanism to monitor and react to changes in a state.
type Observer struct {
	source        UI
	receiver      any
	condition     func() bool
	changeHandler func()
	broadcast     bool

	state           string
	setObserver     func(Observer) Observer
	enableBroadcast func()
}

// While sets a condition for the observer, determining whether it observes
// a state. The condition is periodically checked. Observation stops when the
// condition returns false.
func (o Observer) While(condition func() bool) Observer {
	o.condition = condition
	return o.setObserver(o)
}

// OnChange sets a callback function to be executed each time the observer
// detects a change in the associated state value.
func (o Observer) OnChange(h func()) Observer {
	o.changeHandler = h
	return o.setObserver(o)
}

// WithBroadcast enables the observer to listen to state changes that are
// broadcasted by other browser tabs or windows. This is useful for s
// ynchronizing state across multiple open instances of a web application within
// the same browser.
//
// Calling WithBroadcast creates a BroadcastChannel, which prevents the page
// from being cached. This may impact the Chrome Lighthouse performance
// score due to the additional resources required to manage the broadcast
// channel.
func (o Observer) WithBroadcast() Observer {
	o.enableBroadcast()
	o.broadcast = true
	return o.setObserver(o)
}

func (o Observer) observing() bool {
	if o.source == nil || !o.source.Mounted() {
		return false
	}
	if o.condition != nil {
		return o.condition()
	}
	return true
}

// stateManager is responsible for managing, tracking, and notifying changes
// to state values. It supports concurrency-safe operations and provides
// functionality to observe state changes.
type stateManager struct {
	mutex             sync.RWMutex
	states            map[string]State
	observers         map[string]map[UI]Observer
	initBroadcastOnce sync.Once
	broadcastStoreID  string
	broadcastChannel  Value
}

// Observe initiates observation for a specified state, ensuring the state
// is fetched and set into the given receiver. The returned observer object
// offers methods for advanced observation configurations.
func (m *stateManager) Observe(ctx Context, state string, receiver any) Observer {
	m.Get(ctx, state, receiver)

	return m.setObserver(Observer{
		source:          ctx.Src(),
		receiver:        receiver,
		state:           state,
		setObserver:     m.setObserver,
		enableBroadcast: func() { m.initBroadcast(ctx) },
	})
}

func (m *stateManager) setObserver(v Observer) Observer {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.observers == nil {
		m.observers = make(map[string]map[UI]Observer)
	}

	observers := m.observers[v.state]
	if observers == nil {
		observers = map[UI]Observer{}
		m.observers[v.state] = observers
	}
	observers[v.source] = Observer{
		source:        v.source,
		receiver:      v.receiver,
		condition:     v.condition,
		changeHandler: v.changeHandler,
		broadcast:     v.broadcast,
	}

	return v
}

// Get retrieves the value of a specific state, setting it to the provided
// receiver.
func (m *stateManager) Get(ctx Context, state string, receiver any) {
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

	if expiredTime(value.expiresAt) {
		delete(m.states, state)
		ctx.LocalStorage().Del(state)
		return
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

	if expiredTime(value.ExpiresAt) {
		ctx.LocalStorage().Del(state)
		return nil
	}

	if len(value.EncryptedValue) != 0 {
		return ctx.Decrypt(value.EncryptedValue, receiver)
	} else if len(value.Value) != 0 {
		return json.Unmarshal(value.Value, receiver)
	}
	return nil
}

func (m *stateManager) Set(ctx Context, state string, v any) State {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.states == nil {
		m.states = make(map[string]State)
	}

	value := State{value: v}
	m.states[state] = value

	for _, observer := range m.observers[state] {
		o := observer
		ctx.sourceElement = o.source

		ctx.Dispatch(func(ctx Context) {
			m.mutex.RLock()
			value := m.states[state]
			m.mutex.RUnlock()

			if expiredTime(value.expiresAt) {
				return
			}

			if !o.observing() {
				m.mutex.Lock()
				delete(m.observers[state], o.source)
				m.mutex.Unlock()
				return
			}

			if err := storeValue(o.receiver, value.value); err != nil {
				Log(errors.New("storing state value into receiver failed").
					WithTag("state", state).
					WithTag("observer-type", reflect.TypeOf(o.source)).
					WithTag("receiver-type", reflect.TypeOf(o.receiver)).
					Wrap(err))
				return
			}

			if o.changeHandler != nil {
				o.changeHandler()
			}
		})
	}

	return State{
		value:     v,
		ctx:       ctx,
		name:      state,
		expire:    m.setExpiration,
		persist:   m.persist,
		broadcast: m.broadcast,
	}
}

// Set updates a specified state with a new value and notifies its observers.
// It returns a state object, offering methods for advanced state manipulations.
func (m *stateManager) setExpiration(s State, v time.Time) State {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	s.expiresAt = v

	value := m.states[s.name]
	value.expiresAt = v
	m.states[s.name] = value

	return s
}

func (m *stateManager) persist(s State, encrypt bool) State {
	m.mutex.Lock()
	defer m.mutex.Unlock()

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

func (m *stateManager) broadcast(s State) State {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.initBroadcast(s.ctx)

	if m.broadcastChannel == nil {
		Log(errors.New("broadcast not supported").
			WithTag("state", s.name))
		return s
	}

	b, err := json.Marshal(s.value)
	if err != nil {
		Log(errors.New("encoding broadcast state failed").
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

func (m *stateManager) initBroadcast(ctx Context) {
	m.initBroadcastOnce.Do(func() {
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
	})
}

func (m *stateManager) handleBroadcast(ctx Context, data Value) {
	if storeID := data.Get("StoreID").String(); storeID == "" || storeID == m.broadcastStoreID {
		return
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	state := data.Get("State").String()
	value := []byte(data.Get("Value").String())

	for _, observer := range m.observers[state] {
		o := observer
		if !o.broadcast {
			continue
		}

		ctx.sourceElement = o.source
		ctx.Dispatch(func(ctx Context) {
			if !o.observing() {
				m.mutex.Lock()
				delete(m.observers[state], o.source)
				m.mutex.Unlock()
				return
			}

			if err := json.Unmarshal(value, o.receiver); err != nil {
				Log(errors.New("storing broadcast state value into receiver failed").
					WithTag("state", state).
					WithTag("observer-type", reflect.TypeOf(o.source)).
					WithTag("receiver-type", reflect.TypeOf(o.receiver)).
					Wrap(err))
			}

			if o.changeHandler != nil {
				o.changeHandler()
			}
		})
	}
}

// Delete removes the specified state from the managed states and also deletes
// it from the local storage if it was previously persisted.
func (m *stateManager) Delete(ctx Context, state string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.states, state)
	ctx.LocalStorage().Del(state)
}

func (m *stateManager) UnObserve(ctx Context, state string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for src := range m.observers[state] {
		if src == ctx.Src() {
			delete(m.observers[state], src)
		}
	}
}

// Cleanup removes observers that are no longer active and cleans up any states
// without observers.
func (m *stateManager) Cleanup() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for state, observers := range m.observers {
		for _, observer := range observers {
			if !observer.observing() {
				delete(observers, observer.source)
			}
		}

		if len(observers) == 0 {
			delete(m.observers, state)
		}
	}
}

// CleanupExpiredPersistedStates traverses the local storage to identify and
// remove any persisted states that have expired. This method ensures that the
// local storage is kept clean by eliminating outdated or irrelevant state data.
func (m *stateManager) CleanupExpiredPersistedStates(ctx Context) {
	ctx.LocalStorage().ForEach(func(key string) {
		var state storableState
		ctx.LocalStorage().Get(key, &state)
		if (len(state.Value) != 0 || len(state.EncryptedValue) != 0) &&
			expiredTime(state.ExpiresAt) {
			ctx.LocalStorage().Del(key)
		}
	})
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
			WithTag("value-type", src.Type()).
			WithTag("receiver-type", dst.Type())
	}
	dst.Set(src)
	return nil
}

func expiredTime(v time.Time) bool {
	return !v.IsZero() && v.Before(time.Now())
}
