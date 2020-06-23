package app

import (
	"fmt"
	"runtime"
)

var (
	// DefaultLogger is the logger used to log info and errors.
	DefaultLogger func(format string, v ...interface{}) = log

	defaultColor string
	errorColor   string
	infoColor    string
)

func init() {
	switch runtime.GOARCH {
	case "wasm", "window":
	default:
		defaultColor = "\033[00m"
		errorColor = "\033[91m"
		infoColor = "\033[94m"
	}
}

// Log logs according to a format specifier using the default logger.
func Log(format string, v ...interface{}) {
	DefaultLogger(format, v...)
}

func log(format string, v ...interface{}) {
	errorLevel := false

	for _, a := range v {
		if _, ok := a.(error); ok {
			errorLevel = true
			break
		}
	}

	if errorLevel {
		fmt.Printf(errorColor+"ERROR ‣ "+defaultColor+format+"\n", v...)
		return
	}

	fmt.Printf(infoColor+"INFO ‣ "+defaultColor+format+"\n", v...)
}
