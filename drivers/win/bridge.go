package win

/*
#include "bridge.hpp"
*/
import "C"
import (
	"syscall"
	"unsafe"

	"github.com/pkg/errors"
)

var (
	dll          *syscall.DLL
	winReturnPtr unsafe.Pointer
)

func loadDLL() error {
	var err error
	dll, err = syscall.LoadDLL("goapp.dll")
	return err
}

func closeDLL() {
	dll.Release()
}

func callDllFunc(name string, a ...uintptr) (uintptr, error) {
	proc, err := dll.FindProc(name)
	if err != nil {
		return 0, err
	}

	r, _, _ := proc.Call(a...)
	return r, nil
}

func initBridge() error {
	if _, err := callDllFunc("Bridge_Init"); err != nil {
		return errors.Wrap(err, "init bridge connection failed")
	}

	winReturnPtr = C.winReturn
	if _, err := callDllFunc("Bridge_SetReturn", uintptr(winReturnPtr)); err != nil {
		return errors.Wrap(err, "init win return func failed")
	}

	return nil
}

func winCall(call string) error {
	c := []byte(call)
	ptr := unsafe.Pointer(&c[0])
	_, err := callDllFunc("Bridge_Call", uintptr(ptr))
	return err
}

//export winCallReturn
func winCallReturn(retID, ret, err *C.char) {
	driver.winRPC.Return(
		C.GoString(retID),
		C.GoString(ret),
		C.GoString(err),
	)
}
