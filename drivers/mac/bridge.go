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

	"github.com/murlokswarm/app"
	"github.com/pkg/errors"
)

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
				app.Panic(errors.Wrap(err, "go call failed"))
			}
		})

		return nil
	}

	ret, err := driver.goRPC.Call(call)
	if err != nil {
		app.Panic(errors.Wrap(err, "go call failed"))
	}

	// Returned string must be free in objc code.
	return C.CString(ret)
}
