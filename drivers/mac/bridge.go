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

	"github.com/murlokswarm/app/bridge"
	"github.com/pkg/errors"
)

func handleMacOSRequest(rawurl string, p bridge.Payload, returnID string) (res bridge.Payload, err error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		panic(errors.Wrap(err, "handling MacOS request failed"))
	}

	if len(returnID) != 0 {
		u.Query().Set("return-id", returnID)
	}

	var bres C.bridge_result
	if p == nil {
		bres = C.macosRequest(C.CString(u.String()), nil)
	} else {
		bres = C.macosRequest(C.CString(u.String()), C.CString(p.String()))
	}

	if bres.payload != nil {
		res = bridge.PayloadFromString(C.GoString(bres.payload))
		C.free(unsafe.Pointer(bres.payload))
	}

	if bres.err != nil {
		err = errors.Errorf("handling MacOS request failed: %s", C.GoString(bres.err))
		C.free(unsafe.Pointer(bres.err))
	}
	return
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
