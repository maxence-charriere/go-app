// +build windows

package main

import (
	"context"
	"os"
)

func openCommand() string {
	return "explorer"
}

func win(ctx context.Context, args []string) {
	printErr("work in progress")
}

func mac(ctx context.Context, args []string) {
	printErr("you are not on MacOS!")
	os.Exit(-1)
}

func init() {
	greenColor = ""
	redColor = ""
	defaultColor = ""
}
