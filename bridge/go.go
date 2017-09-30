package bridge

import (
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

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

// GoHandler decribes the func that will handle requests to Go.
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
		panic("go handle pattern should begin with '/'")
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
