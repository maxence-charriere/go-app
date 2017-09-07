package bridge

import (
	"net/url"

	"github.com/pkg/errors"
)

// GoBridge is the interface that describes a bridge to communicate from the
// underlying platform to Go.
type GoBridge interface {
	Request(url string, p Payload)

	RequestWithResponse(url string, p Payload) (res Payload)
}

// GoHandler decribes the func that will handle requests to Go.
type GoHandler func(u *url.URL, p Payload) (res Payload)

// NewGoBridge creates a Go bridge.
func NewGoBridge(handler GoHandler, uichan chan func()) GoBridge {
	return newGoBridge(handler, uichan)
}

type goBridge struct {
	handler GoHandler
	uichan  chan func()
}

func newGoBridge(h GoHandler, uichan chan func()) *goBridge {
	return &goBridge{
		handler: h,
		uichan:  uichan,
	}
}

func (b *goBridge) Request(rawurl string, p Payload) {
	u, err := url.Parse(rawurl)
	if err != nil {
		panic(errors.Wrap(err, "calling callback failed"))
	}

	b.uichan <- func() {
		b.handler(u, p)
	}
}

func (b *goBridge) RequestWithResponse(rawurl string, p Payload) (res Payload) {
	u, err := url.Parse(rawurl)
	if err != nil {
		panic(errors.Wrap(err, "calling callback with response failed"))
	}

	reschan := make(chan Payload, 1)

	b.uichan <- func() {
		reschan <- b.handler(u, p)
	}

	res = <-reschan
	return
}
