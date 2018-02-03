// +build darwin,amd64

package mac

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Cocoa
#cgo LDFLAGS: -framework WebKit
#cgo LDFLAGS: -framework CoreImage
#cgo LDFLAGS: -framework Security
#include "bridge.h"
*/
import "C"
import (
	"net/url"
	"unsafe"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/bridge"
	"github.com/pkg/errors"
)

func handleMacOSRequest(rawurl string, p bridge.Payload, returnID string) (res bridge.Payload, err error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		panic(errors.Wrap(err, "handling MacOS request failed"))
	}

	if len(returnID) != 0 {
		q := u.Query()
		q.Set("return-id", returnID)
		u.RawQuery = q.Encode()
	}

	var bres C.bridge_result
	if p == nil {
		bres = C.macosRequest(C.CString(u.String()), nil)
	} else {
		bres = C.macosRequest(C.CString(u.String()), C.CString(p.String()))
	}
	return parseBridgeResult(bres)
}

func parseBridgeResult(res C.bridge_result) (p bridge.Payload, err error) {
	if res.payload != nil {
		p = bridge.PayloadFromString(C.GoString(res.payload))
		C.free(unsafe.Pointer(res.payload))
	}

	if res.err != nil {
		err = errors.Errorf("handling MacOS request failed: %s", C.GoString(res.err))
		C.free(unsafe.Pointer(res.err))
	}
	return
}

//export macosRequestResult
func macosRequestResult(rawretID *C.char, res C.bridge_result) {
	retID := C.GoString(rawretID)
	C.free(unsafe.Pointer(rawretID))

	payload, err := parseBridgeResult(res)
	driver.macos.Return(retID, payload, err)
}

//export goRequest
func goRequest(url *C.char, payload *C.char) {
	driver.golang.Request(
		C.GoString(url),
		bridge.PayloadFromString(C.GoString(payload)),
	)
}

//export goRequestWithResult
// res should be free after each call of goRequestWithResult.
func goRequestWithResult(url *C.char, payload *C.char) (res *C.char) {
	pret := driver.golang.RequestWithResponse(
		C.GoString(url),
		bridge.PayloadFromString(C.GoString(payload)),
	)

	if pret != nil {
		res = C.CString(pret.String())
	}
	return
}

func windowHandler(h func(w *Window, u *url.URL, p bridge.Payload) (res bridge.Payload)) bridge.GoHandler {
	return func(u *url.URL, p bridge.Payload) (res bridge.Payload) {
		id, err := uuid.Parse(u.Query().Get("id"))
		if err != nil {
			panic(errors.Wrap(err, "creating window handler failed"))
		}

		elem, ok := driver.elements.Element(id)
		if !ok {
			app.DefaultLogger.Logf("%v: window with id %v doesn't exists", u.Path, id)
			return nil
		}

		win, ok := elem.(*Window)
		if !ok {
			panic(errors.Errorf("creating window handler failed: element with id %v is not a window", id))
		}

		return h(win, u, p)
	}
}

func menuHandler(h func(m *Menu, u *url.URL, p bridge.Payload) (res bridge.Payload)) bridge.GoHandler {
	return func(u *url.URL, p bridge.Payload) (res bridge.Payload) {
		id, err := uuid.Parse(u.Query().Get("id"))
		if err != nil {
			panic(errors.Wrap(err, "creating menu handler failed"))
		}

		elem, ok := driver.elements.Element(id)
		if !ok {
			app.DefaultLogger.Logf("%v: menu with id %v doesn't exists", u.Path, id)
			return nil
		}

		menu, ok := elem.(*Menu)
		if !ok {
			panic(errors.Errorf("creating menu handler failed: element with id %v is not a menu", id))
		}

		return h(menu, u, p)
	}
}
