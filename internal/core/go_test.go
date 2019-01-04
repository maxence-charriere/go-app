package core

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGo(t *testing.T) {
	tests := []struct {
		scenario string
		method   string
		inStr    string
		in       map[string]interface{}
		err      bool
	}{
		{
			scenario: "call succeed",
			method:   "test.Greet",
		},
		{
			scenario: "call succeed",
			method:   "test.GreetWithName",
			in:       map[string]interface{}{"Name": "Maxence"},
		},
		{
			scenario: "call method with bad json input returns an error",
			method:   "test.GreetWithName",
			inStr:    "}{",
			err:      true,
		},
		{
			scenario: "call unhandled method returns an error",
			method:   "test.Unknown",
			err:      true,
		},
	}

	golang := Go{}

	golang.Handle("test.Greet", func(in map[string]interface{}) {
		t.Log("hello")
	})

	golang.Handle("test.GreetWithName", func(in map[string]interface{}) {
		t.Log("hello", in["Name"].(string))
	})

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			if len(test.inStr) == 0 {
				call := goCall{
					Method: test.method,
					In:     test.in,
				}

				b, _ := json.Marshal(call)
				test.inStr = string(b)
			}

			err := golang.Call(test.inStr)
			if test.err {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
