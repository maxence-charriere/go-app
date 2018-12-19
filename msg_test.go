package app

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandle(t *testing.T) {
	wg := sync.WaitGroup{}

	Handle("test", func(e Emitter, m Msg) {
		t.Logf("key: %s value: %v", m.Key(), m.Value())
		wg.Done()
	})

	wg.Add(3)

	NewMsg("test").WithValue(42).Post()

	Post(
		NewMsg("test").WithValue(21),
		NewMsg("test").WithValue(84),
	)

	wg.Wait()
}

func TestMsgRegistry(t *testing.T) {
	d := newEventRegistry(func(f func()) {
		f()
	})

	r := newMsgRegistry(d)
	wg := sync.WaitGroup{}

	r.handle("test", func(e Emitter, m Msg) {
		wg.Done()
	})

	r.post(NewMsg("unknown"))

	wg.Add(1)
	r.post(NewMsg("test").WithValue(42))

	wg.Add(3)
	r.post(
		NewMsg("test"),
		NewMsg("test").WithValue("hello"),
		NewMsg("test").WithValue(21),
	)

	wg.Wait()
}

func TestEventRegistry(t *testing.T) {
	Logger = func(format string, a ...interface{}) {
		log := fmt.Sprintf(format, a...)
		t.Log(log)
	}

	tests := []struct {
		scenario string
		subName  Event
		handler  func(*bool) interface{}
		called   bool
		dispName Event
		dispArg  interface{}
		panic    bool
	}{
		{
			scenario: "register and dispatch without arg",
			subName:  "test",
			handler: func(called *bool) interface{} {
				return func() {
					*called = true
				}
			},
			called:   true,
			dispName: "test",
			dispArg:  nil,
		},
		{
			scenario: "register without arg and dispatch with arg",
			subName:  "test",
			handler: func(called *bool) interface{} {
				return func() {
					*called = true
				}
			},
			called:   true,
			dispName: "test",
			dispArg:  "foobar",
		},
		{
			scenario: "register and dispatch with arg",
			subName:  "test",
			handler: func(called *bool) interface{} {
				return func(arg string) {
					*called = true

					if arg != "hello" {
						panic("greet is not hello")
					}
				}
			},
			called:   true,
			dispName: "test",
			dispArg:  "hello",
		},
		{
			scenario: "register and dispatch with bad arg",
			subName:  "test",
			handler: func(called *bool) interface{} {
				return func(arg int) {
					*called = true
				}
			},
			called:   false,
			dispName: "test",
			dispArg:  "hello",
		},
		{
			scenario: "register non func handler",
			subName:  "test",
			handler:  func(called *bool) interface{} { return nil },
			panic:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			defer func() {
				err := recover()
				if test.panic {
					assert.NotNil(t, err)
					return
				}
				assert.Nil(t, err)
			}()

			called := false

			r := newEventRegistry(func(f func()) {
				f()
			})

			unsub := r.subscribe(test.subName, test.handler(&called))
			defer unsub()

			r.Emit(test.dispName, test.dispArg)
			assert.Equal(t, test.called, called)
		})
	}
}

func TestSubscriber(t *testing.T) {
	s := NewSubscriber()
	defer s.Close()

	s.Subscribe("test-event-subscriber", func() {})
}
