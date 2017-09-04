package bridge

import (
	"net/url"

	"github.com/google/uuid"
)

type GoBridge interface {
	Request(url string, payload []byte) (res []byte, err error)

	RequestWithAsyncResponse(url string, payload []byte) (res []byte, err error)

	Return(returnID uuid.UUID, payload interface{}, err error)
}

type GoHandler func(u *url.URL, payload Payload) (response interface{}, err error)

func NewGoBridge(handler GoHandler) GoBridge {
	return newGoBridge(handler)
}

type goBridge struct {
	handler GoHandler
}

func newGoBridge(h GoHandler) *goBridge {
	return &goBridge{
		handler: h,
	}
}

func (b *goBridge) Request(url string, payload []byte) (res []byte, err error) {
	return
}

func (b *goBridge) RequestWithAsyncResponse(url string, payload []byte) (res []byte, err error) {
	return
}

func (b *goBridge) Return(returnID uuid.UUID, payload interface{}, err error) {
	return
}

// Handler:
// parse
// execute sur uigoroutine

// Handler avec result
// parse
// cree un return
// appel fonction sur uigoroutine
// wait for return
