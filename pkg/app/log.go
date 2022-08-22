package app

import (
	"fmt"
	"runtime"
	"strings"
)

var (
	// DefaultLogger is the logger used to log info and errors.
	DefaultLogger func(format string, v ...any)

	defaultColor string
	errorColor   string
	infoColor    string
)

func init() {
	goarch := runtime.GOARCH
	if goarch == "wasm" {
		DefaultLogger = clientLog
		return
	}

	if goarch != "window" {
		defaultColor = "\033[00m"
		errorColor = "\033[91m"
		infoColor = "\033[94m"
	}
	DefaultLogger = serverLog
}

// Log logs using the default formats for its operands. Spaces are always added
// between operands.
func Log(v ...any) {
	var b strings.Builder
	for i := 0; i < len(v); i++ {
		if i != 0 {
			b.WriteByte(' ')
		}
		b.WriteString("%v")
	}
	Logf(b.String(), v...)
}

// Logf logs according to a format specifier.
func Logf(format string, v ...any) {
	DefaultLogger(format, v...)
}

func serverLog(format string, v ...any) {
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

func clientLog(format string, v ...any) {
	isErrorLevel := false
	for _, a := range v {
		if _, isErr := a.(error); isErr {
			isErrorLevel = true
			break
		}
	}

	if isErrorLevel {
		Window().Get("console").Call("error", fmt.Sprintf(format, v...))
		return
	}
	Window().Get("console").Call("log", fmt.Sprintf(format, v...))
}
