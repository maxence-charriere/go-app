package core

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// GoHandler describes a function that handle requests from underlying platform.
type GoHandler func(in map[string]interface{})

// Go is a struct that provides remote procedure calls from the underlying
// platform to Go.
type Go struct {
	handlers map[string]GoHandler
}

// Handle set up the given handler to handle the named method.
func (g *Go) Handle(method string, h GoHandler) {
	if (g.handlers) == nil {
		g.handlers = make(map[string]GoHandler)
	}

	g.handlers[method] = h
}

// Call perform the described call.
func (g *Go) Call(call string) error {
	gocall := goCall{}
	if err := json.Unmarshal([]byte(call), &gocall); err != nil {
		return err
	}

	h, ok := g.handlers[gocall.Method]
	if !ok {
		return errors.Errorf("method %q does not have a handler", gocall.Method)
	}

	h(gocall.In)
	return nil
}

type goCall struct {
	Method string
	In     map[string]interface{}
}
