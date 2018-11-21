package main

import (
	"context"
)

func openCommand() string {
	return "explorer"
}

func win(ctx context.Context, args []string) {
	printErr("work in progress")
}

func init() {
	greenColor = ""
	redColor = ""
	defaultColor = ""
}
