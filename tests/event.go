package tests

import (
	"testing"

	"github.com/murlokswarm/app"
)

func TestEventRegistry(t *testing.T, newRegistry func() app.EventRegistry) {
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
				err := recover()

				if err != nil && !test.panic {
					t.Error(err)
				}
			}()

			called := false

			r := newRegistry()
			unsub := r.Subscribe(test.subName, test.handler(&called))
			defer unsub()

			r.Dispatch(test.dispName, test.dispArg)

			if called != test.called {
				t.Error("called expected:", test.called)
				t.Error("called:         ", called)
			}

			if test.panic {
				t.Error("no panic")
			}
		})
	}
}
