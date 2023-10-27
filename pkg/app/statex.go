package app

import (
	"encoding/json"
	"reflect"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

// StateX represents a state with additional features such as expiration,
// persistence, and broadcasting capabilities.
type StateX struct {
	value     any
	expiresAt time.Time

	ctx       Context
	name      string
	expire    func(StateX, time.Time) StateX
	persist   func(StateX, bool) StateX
	broadcast func(StateX) StateX
}

// ExpiresIn sets the expiration time for the state by specifying a duration
// from the current time.
func (s StateX) ExpiresIn(v time.Duration) StateX {
	return s.expire(s, time.Now().Add(v))
}

// ExpiresAt sets the exact expiration time for the state.
func (s StateX) ExpiresAt(v time.Time) StateX {
	return s.expire(s, v)
}

// Persist ensures the state is persisted into the local storage.
func (s StateX) Persist() StateX {
	return s.persist(s, false)
}

// PersistWithEncryption ensures the state is persisted into the local storage
// with encryption.
func (s StateX) PersistWithEncryption() StateX {
	return s.persist(s, true)
}

// Broadcast signals that changes to the state will be broadcasted to other
// browser tabs and windows sharing the same origin when it is supported.
func (s StateX) Broadcast() StateX {
	return s.broadcast(s)
}

type storableState struct {
	Value          json.RawMessage `json:",omitempty"`
	EncryptedValue []byte          `json:",omitempty"`
	ExpiresAt      time.Time       `json:",omitempty"`
}

// ObserverX represents a mechanism to monitor and react to changes in a state.
type ObserverX struct {
	source        UI
	receiver      any
	condition     func() bool
	changeHandler func()

	state       string
	setObserver func(ObserverX) ObserverX
}

// While sets a condition for the observer, determining whether it observes
// a state. The condition is periodically checked. Observation stops when the
// condition returns false.
func (o ObserverX) While(condition func() bool) ObserverX {
	o.condition = condition
	return o.setObserver(o)
}

// OnChange sets a callback function to be executed each time the observer
// detects a change in the associated state value.
func (o ObserverX) OnChange(h func()) ObserverX {
	o.changeHandler = h
	return o.setObserver(o)
}

func (o ObserverX) observing() bool {
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
	mutex            sync.RWMutex
	states           map[string]StateX
	observers        map[string]map[UI]ObserverX
	broadcastStoreID string
	broadcastChannel Value
}

// Observe initiates observation for a specified state, ensuring the state
// is fetched and set into the given receiver. The returned observer object
// offers methods for advanced observation configurations.
func (m *stateManager) Observe(ctx Context, state string, receiver any) ObserverX {
	m.Get(ctx, state, receiver)
	return m.setObserver(ObserverX{
		source:      ctx.Src(),
		receiver:    receiver,
		state:       state,
		setObserver: m.setObserver,
	})
}

func (m *stateManager) setObserver(v ObserverX) ObserverX {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.observers == nil {
		m.observers = make(map[string]map[UI]ObserverX)
	}

	observers := m.observers[v.state]
	if observers == nil {
		observers = map[UI]ObserverX{}
		m.observers[v.state] = observers
	}
	observers[v.source] = ObserverX{
		source:        v.source,
		receiver:      v.receiver,
		condition:     v.condition,
		changeHandler: v.changeHandler,
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

	if !value.expiresAt.IsZero() && value.expiresAt.Before(time.Now()) {
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

	if !value.ExpiresAt.IsZero() && value.ExpiresAt.Before(time.Now()) {
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

func (m *stateManager) Set(ctx Context, state string, v any) StateX {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.states == nil {
		m.states = make(map[string]StateX)
	}

	value := StateX{value: v}
	m.states[state] = value

	for _, observer := range m.observers[state] {
		o := observer
		ctx.sourceElement = o.source

		ctx.Dispatch(func(ctx Context) {
			m.mutex.RLock()
			value := m.states[state]
			m.mutex.RUnlock()
			if !value.expiresAt.IsZero() && value.expiresAt.Before(time.Now()) {
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

	return StateX{
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
func (m *stateManager) setExpiration(s StateX, v time.Time) StateX {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	s.expiresAt = v

	value := m.states[s.name]
	value.expiresAt = v
	m.states[s.name] = value

	return s
}

func (m *stateManager) persist(s StateX, encrypt bool) StateX {
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

// InitBroadcast initializes a broadcast channel to share state changes
// across browser tabs or windows.
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
		ctx.dispatch(func() {
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

// TODO:
// - expire all expired values from local storage
// - test broadcast

// func TestExpireExpriredPersistentValues(t *testing.T) {
// 	if IsServer {
// 		t.Skip()
// 	}

// 	d := NewClientTester(&foo{})
// 	defer d.Close()
// 	localStorage := d.getLocalStorage()

// 	s := newStore(d)
// 	defer s.Close()

// 	t.Run("non expired state is not removed", func(t *testing.T) {
// 		localStorage.Clear()
// 		s.setPersistent("/hello", false, time.Now().Add(time.Minute), "hello")
// 		require.Equal(t, 1, localStorage.Len())

// 		s.expireExpriredPersistentValues()
// 		require.Equal(t, 1, localStorage.Len())
// 	})

// 	t.Run("expired state is removed", func(t *testing.T) {
// 		localStorage.Clear()
// 		s.setPersistent("/bye", false, time.Now().Add(-time.Minute), "bye")
// 		require.Equal(t, 1, localStorage.Len())

// 		s.expireExpriredPersistentValues()
// 		require.Equal(t, 0, localStorage.Len())
// 	})

// 	t.Run("non state value is not removed", func(t *testing.T) {
// 		localStorage.Clear()
// 		localStorage.Set("/hi", "hi")
// 		require.Equal(t, 1, localStorage.Len())

// 		s.expireExpriredPersistentValues()
// 		require.Equal(t, 1, localStorage.Len())
// 	})
// }

// func TestStoreBroadcast(t *testing.T) {
// 	d1 := NewClientTester(&foo{})
// 	s1 := newStore(d1)
// 	defer d1.Close()
// 	defer s1.Close()

// 	bar := &bar{}
// 	d2 := NewClientTester(bar)
// 	s2 := newStore(d2)
// 	defer d2.Close()
// 	defer s2.Close()

// 	require.NotEqual(t, s1.id, s2.id)

// 	key := "/test/store/broadcast"
// 	t.Run("state is not broadcasted", func(t *testing.T) {
// 		var v int
// 		s2.Observe(key, bar).Value(&v)

// 		s1.Set(key, func() {}, Broadcast)
// 		d2.Consume()
// 		require.Zero(t, v)
// 	})

// 	t.Run("state is broadcasted", func(t *testing.T) {
// 		if IsServer {
// 			t.Skip()
// 		}

// 		var v int

// 		s2.Observe(key, bar).Value(&v)
// 		s1.Set(key, 42, Broadcast)

// 		time.Sleep(time.Millisecond * 100)
// 		d2.Consume()
// 		require.Equal(t, 42, v)
// 	})

// 	t.Run("broadcasted state is not observed", func(t *testing.T) {
// 		if IsServer {
// 			t.Skip()
// 		}

// 		var v func()

// 		s2.Observe(key, bar).Value(&v)
// 		s1.Set(key, 42, Broadcast)

// 		time.Sleep(time.Millisecond * 50)
// 		d2.Consume()
// 		require.Zero(t, v)
// 	})
// }
