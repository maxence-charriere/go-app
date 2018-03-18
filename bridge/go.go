package bridge

import (
	"encoding/json"
	"net/url"
	"strings"

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

// ----------------------- TMP ------------------------------

// GoBridge is the interface that describes a bridge to communicate from the
// underlying platform to Go.
type GoBridge interface {
	// Handle registers the handler for the given pattern.
	// It panics if pattern doesn't start with '/'
	Handle(pattern string, handler GoHandler)

	// Request issues a request to the specified URL with the payload.
	Request(url string, p Payload)

	// Request issues a request to the specified URL with the payload and
	// returns to have a response.
	RequestWithResponse(url string, p Payload) (res Payload)
}

// GoHandler describes the func that will handle requests to Go.
type GoHandler func(u *url.URL, p Payload) (res Payload)

// NewGoBridge creates a Go bridge.
func NewGoBridge(uichan chan func()) GoBridge {
	return newGoBridge(uichan)
}

type goBridge struct {
	handlers map[string]GoHandler
	uichan   chan func()
}

func newGoBridge(uichan chan func()) *goBridge {
	return &goBridge{
		handlers: make(map[string]GoHandler),
		uichan:   uichan,
	}
}

func (b *goBridge) Handle(pattern string, handler GoHandler) {
	if len(pattern) == 0 || pattern[0] != '/' {
		panic("go handle pattern doesn't begin by '/'")
	}

	if handler == nil {
		panic("go handler can't be nil")
	}

	b.handlers[pattern] = handler
}

func (b *goBridge) Request(rawurl string, p Payload) {
	u, err := url.Parse(rawurl)
	if err != nil {
		panic(errors.Wrap(err, "parsing URL failed"))
	}

	b.uichan <- func() {
		b.handle(u, p)
	}
}

func (b *goBridge) RequestWithResponse(rawurl string, p Payload) (res Payload) {
	u, err := url.Parse(rawurl)
	if err != nil {
		panic(errors.Wrap(err, "parsing URL failed"))
	}

	reschan := make(chan Payload, 1)

	// Here we donc execute the handler in the ui goroutine because it can
	// cause a deadlock if RequestWithResponse is called while some platform
	// requests are waiting for an async result.
	go func() {
		reschan <- b.handle(u, p)
	}()

	res = <-reschan
	return
}

func (b *goBridge) handle(u *url.URL, p Payload) (res Payload) {
	pattern := u.Path

	for {
		handler, ok := b.handlers[pattern]
		if ok {
			res = handler(u, p)
			return
		}

		sep := strings.LastIndexByte(pattern, '/')
		if sep == -1 {
			panic(errors.Errorf("go request %v is not handled", u))
		}
		pattern = pattern[:sep]
	}
}
