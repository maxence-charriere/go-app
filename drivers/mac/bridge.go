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
	"unsafe"

	"github.com/murlokswarm/app/bridge"
	"github.com/pkg/errors"
)

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
	return res
}

func macCall(call string) error {
	ccall := C.CString(call)
	defer C.free(unsafe.Pointer(ccall))
	C.macCall(ccall)
	return nil
}

//export macCallReturn
func macCallReturn(retID, ret, err *C.char) {
	driver.macRPC.Return(
		C.GoString(retID),
		C.GoString(ret),
		C.GoString(err),
	)
}

//export goCall
func goCall(ccall *C.char, ui C.BOOL) (cout *C.char) {
	call := C.GoString(ccall)

	if ui == 1 {
		driver.CallOnUIGoroutine(func() {
			if _, err := driver.goRPC.Call(call); err != nil {
				panic(errors.Wrap(err, "go call"))
			}
		})
		return nil
	}

	ret, err := driver.goRPC.Call(call)
	if err != nil {
		panic(errors.Wrap(err, "go call"))
	}

	// Returned string must be free in objc code.
	return C.CString(ret)
}

func stringSlice(v interface{}) []string {
	src := v.([]interface{})
	s := make([]string, 0, len(src))

	for _, item := range src {
		s = append(s, item.(string))
	}
	return s
}
