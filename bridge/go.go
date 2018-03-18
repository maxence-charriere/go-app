package bridge

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// GoCall represents a Go call.
type GoCall struct {
	Method string
	Input  map[string]interface{}
}

// GoRPCHandler represents the handler that will perform the Go call.
type GoRPCHandler func(in map[string]interface{}) interface{}

// GoRPC is a struct that implements the remote procedure call from
// underlying platform to Go.
type GoRPC struct {
	handlers map[string]GoRPCHandler
}

// Handle registers the handler for the given method.
func (r *GoRPC) Handle(method string, handler GoRPCHandler) {
	if r.handlers == nil {
		r.handlers = make(map[string]GoRPCHandler)
	}
	r.handlers[method] = handler
}

// Call calls the described method.
func (r *GoRPC) Call(call string) (string, error) {
	var c GoCall
	if err := json.Unmarshal([]byte(call), &c); err != nil {
		return "", err
	}

	h, ok := r.handlers[c.Method]
	if !ok {
		return "", errors.Errorf("%s is not handled", c.Method)
	}

	ret := h(c.Input)
	if ret == nil {
		return "", nil
	}

	data, err := json.Marshal(ret)
	return string(data), err
}
