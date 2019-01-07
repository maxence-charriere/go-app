// +build darwin

package objc

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Cocoa
#cgo LDFLAGS: -framework WebKit
#cgo LDFLAGS: -framework CoreImage
#cgo LDFLAGS: -framework Security
#cgo LDFLAGS: -framework GameController
#include "bridge.h"
*/
import "C"

import (
	"unsafe"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/core"
	"github.com/pkg/errors"
)

var (
	callOnUI func(func())
	platform = core.Platform{Handler: macCall}
	golang   core.Go
)

// RPC returns rpc objects to achieve two way communication between Objective C
// and Go.
func RPC(ui func(func())) (*core.Platform, *core.Go) {
	callOnUI = ui
	return &platform, &golang
}

func macCall(call string) error {
	ccall := C.CString(call)
	C.macCall(ccall)
	C.free(unsafe.Pointer(ccall))
	return nil
}

//export macCallReturn
func macCallReturn(retID, ret, err *C.char) {
	platform.Return(
		C.GoString(retID),
		C.GoString(ret),
		C.GoString(err),
	)
}

//export goCall
func goCall(ccall *C.char) {
	call := C.GoString(ccall)

	callOnUI(func() {
		if err := golang.Call(call); err != nil {
			app.Panic(errors.Wrap(err, "go call failed"))
		}
	})
}
