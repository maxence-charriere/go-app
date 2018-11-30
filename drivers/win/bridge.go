// +build windows

package win

/*
#include "bridge.hpp"
*/
import "C"
import (
	"strconv"
	"syscall"
	"unsafe"

	"github.com/murlokswarm/app"
	"github.com/pkg/errors"
)

var (
	dll              *syscall.DLL
	winCallReturnPtr unsafe.Pointer
	goCallPtr        unsafe.Pointer
)

func loadDLL() error {
	var err error
	dll, err = syscall.LoadDLL("goapp.dll")
	return err
}

func closeDLL() {
	dll.Release()
}

//go:uintptrescapes
func callDllFunc(name string, a ...unsafe.Pointer) (uintptr, error) {
	args := make([]uintptr, len(a))
	for i, arg := range a {
		args[i] = uintptr(arg)
	}

	proc, err := dll.FindProc(name)
	if err != nil {
		return 0, err
	}

	r, _, _ := proc.Call(args...)
	return r, nil
}

func initBridge() error {
	if _, err := callDllFunc("Bridge_Init"); err != nil {
		return errors.Wrap(err, "init bridge connection failed")
	}

	winCallReturnPtr = C.winCallReturn
	if _, err := callDllFunc("Bridge_SetWinCallReturn", winCallReturnPtr); err != nil {
		return errors.Wrap(err, "init winReturn func failed")
	}

	goCallPtr = C.goCall
	if _, err := callDllFunc("Bridge_SetGoCall", goCallPtr); err != nil {
		return errors.Wrap(err, "init goCall func failed")
	}

	return nil
}

func winCall(call string) error {
	ccall := C.CString(call)
	ptr := unsafe.Pointer(ccall)
	_, err := callDllFunc("Bridge_Call", ptr)
	C.free(ptr)
	return err
}

//export onWinCallReturn
func onWinCallReturn(retID, ret, err *C.char) {
	driver.winRPC.Return(
		C.GoString(retID),
		C.GoString(ret),
		C.GoString(err),
	)
}

//export onGoCall
func onGoCall(ccall *C.char, cui *C.char) (cout *C.char) {
	call := C.GoString(ccall)
	ui, _ := strconv.ParseBool(C.GoString(cui))

	if ui {
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

	// Returned string must be free in c++ code.
	return C.CString(ret)
}
