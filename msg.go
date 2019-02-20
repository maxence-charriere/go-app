package app

import (
	"sync"
)

// Msg is the interface that describes message.
type Msg interface {
	// The message key.
	Key() string

	// The message value.
	Value() interface{}

	// Sets the message value.
	WithValue(interface{}) Msg

	// Posts the message.
	// It will be handled in another goroutine.
	Post()
}

// MsgHandler is the interface that describes a message handler.
// It is used to respond to a Msg.
type MsgHandler func(Msg)

type msg struct {
	key   string
	value interface{}
}

func (m *msg) Key() string {
	return m.key
}

func (m *msg) Value() interface{} {
	return m.value
}

func (m *msg) WithValue(v interface{}) Msg {
	m.value = v
	return m
}

func (m *msg) Post() {
	messages.post(m)
}

type msgRegistry struct {
	mutex sync.RWMutex
	msgs  map[string]MsgHandler
}

func newMsgRegistry() *msgRegistry {
	return &msgRegistry{
		msgs: make(map[string]MsgHandler),
	}
}

func (r *msgRegistry) handle(key string, h MsgHandler) {
	r.mutex.Lock()
	r.msgs[key] = h
	r.mutex.Unlock()
}

func (r *msgRegistry) post(msgs ...Msg) {
	go func() {
		for _, m := range msgs {
			r.exec(m)
		}
	}()
}

func (r *msgRegistry) exec(m Msg) {
	r.mutex.RLock()
	h, ok := r.msgs[m.Key()]
	r.mutex.RUnlock()

	if ok {
		h(m)
	}
}
