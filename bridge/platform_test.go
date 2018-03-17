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
		skipOutput     bool
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
			scenario: "method on goroutine",
			method:   "test.GreetOnGoroutine",
			input: input{
				Name: "Maxence",
			},
			expectedOutput: output{
				Greeting: "Hello, Maxence",
			},
		},
		{
			scenario: "method without output",
			method:   "test.NoGreet",
			input: input{
				Name: "Maxence",
			},
			expectedOutput: output{},
		},
		{
			scenario:  "async method error",
			method:    "test.GreetErr",
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

	var rpc PlatformRPC

	handler := func(rawCall string) error {
		var call PlatformCall
		err := json.Unmarshal([]byte(rawCall), &call)
		if err != nil {
			return err
		}

		name := call.Input.(map[string]interface{})["Name"].(string)

		var out []byte
		if out, err = json.Marshal(output{
			Greeting: "Hello, " + name,
		}); err != nil {
			return err
		}

		switch call.Method {
		case "test.Greet":
			rpc.Return(call.ReturnID, string(out), "")
			return nil

		case "test.GreetOnGoroutine":
			go rpc.Return(call.ReturnID, string(out), "")
			return nil

		case "test.NoGreet":
			rpc.Return(call.ReturnID, "", "")
			return nil

		case "test.GreetErr":
			rpc.Return(call.ReturnID, "", "simulated err")
			return nil

		default:
			return errors.Errorf("%s: unknown rpc method", call.Method)
		}
	}

	rpc.Handler = handler

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			var out output

			err := rpc.Call(test.method, &out, test.input)
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

	rpc := PlatformRPC{}
	rpc.Return("test", "", "")
	t.Error("test did not panic")
}
