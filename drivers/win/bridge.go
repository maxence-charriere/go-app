package win

import (
	"syscall"
)

var (
	dll *syscall.DLL
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
	_, err := callDllFunc("Bridge_Init")
	return err
}
