// +build windows

package uwp

/*
#include "uwp.hpp"
*/
import "C"

import (
	"strconv"
	"syscall"
	"unsafe"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/bridge"
	"github.com/pkg/errors"
)

var (
	callOnUI func(func())
	goRPC    bridge.GoRPC
	winRPC   = bridge.PlatformRPC{
		Handler: winCall,
	}

	dll              *syscall.DLL
	winCallReturnPtr unsafe.Pointer
	goCallPtr        unsafe.Pointer
)

// Connect connects the package to the app connection service.
func Connect() func() {
	var err error
	if dll, err = syscall.LoadDLL("goapp.dll"); err != nil {
		panic(err)
	}

	winCallReturnPtr = C.winCallReturn
	if _, err = callDllFunc("Bridge_SetWinCallReturn", winCallReturnPtr); err != nil {
		panic(errors.Wrap(err, "init winReturn func failed"))
	}

	goCallPtr = C.goCall
	if _, err = callDllFunc("Bridge_SetGoCall", goCallPtr); err != nil {
		panic(errors.Wrap(err, "init goCall func failed"))
	}

	if _, err := callDllFunc("Bridge_Init"); err != nil {
		panic(errors.Wrap(err, "init bridge connection failed"))
	}

	return func() {
		dll.Release()
	}
}

// RPC returns rpc objects to achieve two way communication between uwp and Go.
func RPC(ui func(func())) (*bridge.PlatformRPC, *bridge.GoRPC) {
	callOnUI = ui
	return &winRPC, &goRPC
}

func winCall(call string) error {
	ccall := C.CString(call)
	ptr := unsafe.Pointer(ccall)
	_, err := callDllFunc("Bridge_Call", ptr)
	C.free(ptr)
	return err
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

//export onWinCallReturn
func onWinCallReturn(retID, ret, err *C.char) {
	winRPC.Return(
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
		callOnUI(func() {
			if _, err := goRPC.Call(call); err != nil {
				app.Panic(errors.Wrap(err, "go call failed"))
			}
		})

		return nil
	}

	ret, err := goRPC.Call(call)
	if err != nil {
		app.Panic(errors.Wrap(err, "go call failed"))
	}

	// Returned string must be free in c++ code.
	return C.CString(ret)
}
