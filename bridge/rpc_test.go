package bridge

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

type input struct {
	Name string
}

type asyncInput struct {
	Name string
}

func (i asyncInput) Async() bool {
	return true
}

type invalidInput struct {
	Name string
	Func func()
}

type output struct {
	Greeting string
}

func TestRPC(t *testing.T) {
	tests := []struct {
		scenario       string
		method         string
		input          interface{}
		expectedOutput output
		returnErr      bool
	}{
		{
			scenario: "method",
			method:   "test.Greet",
			input: input{
				Name: "Maxence",
			},
			expectedOutput: output{
				Greeting: "Hello, Maxence",
			},
		},
		{
			scenario: "async method",
			method:   "test.AsyncGreet",
			input: asyncInput{
				Name: "Maxence",
			},
			expectedOutput: output{
				Greeting: "Hello, Maxence",
			},
		},
		{
			scenario:  "async method error",
			method:    "test.AsyncGreetErr",
			input:     asyncInput{},
			returnErr: true,
		},
		{
			scenario:  "unknown method",
			method:    "test.Unkown",
			input:     input{},
			returnErr: true,
		},
		{
			scenario:  "invalid input",
			method:    "test.Greet",
			input:     invalidInput{},
			returnErr: true,
		},
	}

	var rpc RPC

	handler := func(rawCall string) (ret string, err error) {
		var call call
		if err = json.Unmarshal([]byte(rawCall), &call); err != nil {
			return "", err
		}

		name := call.Input.(map[string]interface{})["Name"].(string)

		var out []byte
		if out, err = json.Marshal(output{
			Greeting: "Hello, " + name,
		}); err != nil {
			return "", err
		}

		switch call.Method {
		case "test.Greet":
			return string(out), nil

		case "test.AsyncGreet":
			go func() {
				rpc.Return(call.ReturnID, string(out), "")
			}()
			return "", nil

		case "test.AsyncGreetErr":
			go func() {
				rpc.Return(call.ReturnID, "", "simulated err")
			}()
			return "", nil

		default:
			return "", errors.Errorf("%s: unknown rpc method", call.Method)
		}
	}

	rpc.Handler = handler

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			var out output
			err := rpc.Call(test.method, test.input, &out)
			if test.returnErr && err == nil {
				t.Fatal("error is nil")
			} else if test.returnErr && err != nil {
				return
			}
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(test.expectedOutput, out) {
				t.Errorf("expected: %+v", test.expectedOutput)
				t.Errorf("output  : %+v", out)
			}
		})
	}
}

func TestRPCReturnPanic(t *testing.T) {
	defer func() {
		recover()
	}()

	rpc := RPC{}
	rpc.Return("test", "", "")
	t.Error("test did not panic")
}
