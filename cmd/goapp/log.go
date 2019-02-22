package main

import (
	"fmt"
	"os"
	"runtime"
)

var (
	successColor = "\033[92m"
	errorColor   = "\033[91m"
	warnColor    = "\033[93m"
	defaultColor = "\033[00m"
	verbose      = false
)

func init() {
	if runtime.GOOS == "windows" {
		successColor = ""
		errorColor = ""
		warnColor = ""
		defaultColor = ""
	}
}

func log(format string, v ...interface{}) {
	if verbose {
		format = "‣ " + format
		fmt.Printf(format, v...)
		fmt.Println()
	}
}

func success(format string, v ...interface{}) {
	fmt.Print(successColor)
	format = "✔ " + format
	fmt.Printf(format, v...)
	fmt.Println(defaultColor)
}

func warn(format string, v ...interface{}) {
	fmt.Print(warnColor)
	format = "! " + format
	fmt.Printf(format, v...)
	fmt.Println(defaultColor)
}

func fail(format string, v ...interface{}) {
	fmt.Print(errorColor)
	format = "x " + format
	fmt.Printf(format, v...)
	fmt.Println(defaultColor)
	os.Exit(-1)
}
