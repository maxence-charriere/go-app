package app

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			scenario: "register and emit without arg",
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
			scenario: "register without arg and emit with arg",
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
			scenario: "register and emit with arg",
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
			scenario: "register and emit without enough args",
			subName:  "test",
			handler: func(called *bool) interface{} {
				return func(string, bool) {
					*called = true
				}
			},
			called:   false,
			dispName: "test",
			dispArg:  "hello",
		},
		{
			scenario: "register and emit with bad arg",
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

			ui := make(chan func(), 32)
			defer close(ui)

			r := newEventRegistry(ui)

			unsub := r.subscribe(test.subName, test.handler(&called))
			defer unsub()

			r.Emit(test.dispName, test.dispArg)

			select {
			case f := <-ui:
				f()

			default:
			}

			assert.Equal(t, test.called, called)
		})
	}
}

func TestSubscriber(t *testing.T) {
	s := NewSubscriber()
	defer s.Close()

	s.Subscribe("test-event-subscriber", func() {})
}
