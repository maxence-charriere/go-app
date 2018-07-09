package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEventRegistry(t *testing.T) {
	tests := []struct {
		scenario string
		subName  string
		handler  func(*bool) interface{}
		called   bool
		dispName string
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
				if err := recover(); err != nil {
					require.True(t, test.panic)
				}
			}()

			r := newEventRegistry(func(f func()) {
				f()
			})

			called := false
			unsub := r.Subscribe(test.subName, test.handler(&called))
			defer unsub()

			r.Dispatch(test.dispName, test.dispArg)
			require.Equal(t, test.called, called)

			if test.panic {
				t.Error("no panic")
			}
		})
	}
}

func TestEventSubscriber(t *testing.T) {
	s := NewEventSubscriber()
	defer s.Close()

	s.Subscribe("test-event-subscriber", func() {})
}
