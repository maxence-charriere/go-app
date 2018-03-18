package bridge

import "testing"

func TestGoRPC(t *testing.T) {
	tests := []struct {
		scenario string
		call     string
		method   string
		handler  GoRPCHandler
		expected string
		err      bool
	}{
		{
			scenario: "call without return",
			call:     `{"Method":"tests.WithoutReturn"}`,
			method:   "tests.WithoutReturn",
			handler: func(in map[string]interface{}) interface{} {
				return nil
			},
		},
		{
			scenario: "call with return",
			call:     `{"Method":"tests.WithReturn"}`,
			method:   "tests.WithReturn",
			handler: func(in map[string]interface{}) interface{} {
				return struct {
					Hello string
				}{
					Hello: "world",
				}
			},
			expected: `{"Hello":"world"}`,
		},
		{
			scenario: "call with bad input",
			call:     `}{`,
			err:      true,
		},
		{
			scenario: "call not handled",
			call:     `{"Method":"tests.NotHandled"}`,
			err:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			var rpc GoRPC

			rpc.Handle(test.method, test.handler)

			ret, err := rpc.Call(test.call)
			if test.err && err == nil {
				t.Fatal("error is nil")
			} else if test.err && err != nil {
				return
			}

			if ret != test.expected {
				t.Error("expected:", test.expected)
				t.Error("returned:", ret)
			}
		})
	}
}
