// +build linux

package main

import (
	"context"
	"os"
)

func openCommand() string {
	return "xdg-open"
}

func mac(ctx context.Context, args []string) {
	printErr("you are not on MacOS!")
	os.Exit(-1)
}

func win(ctx context.Context, args []string) {
	printErr("you are not on Windows!")
	os.Exit(-1)
}
