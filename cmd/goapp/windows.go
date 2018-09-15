// +build windows

package main

import (
	"context"
	"os"

	"github.com/segmentio/conf"
)

func win(ctx context.Context, args []string) {
	ld := conf.Loader{
		Name: "goapp win",
		Args: args,
		Commands: []conf.Command{
			{Name: "init", Help: "Download the Windows dev tools and create required files and directories."},
			{Name: "build", Help: "Build the Windows app."},
			{Name: "run", Help: "Run a Windows app and capture its logs."},
			{Name: "help", Help: "Show the MacOS help"},
		},
	}

	switch cmd, args := conf.LoadWith(nil, ld); cmd {
	case "init":
		initWin(ctx, args)

	case "build":
		buildWin(ctx, args)

	case "run":
		runWin(ctx, args)

	case "help":
		ld.PrintHelp(nil)

	default:
		panic("unreachable")
	}
}

type winInitConfig struct {
	Verbose bool `conf:"v" help:"Enable verbose mode."`
}

func initWin(ctx context.Context, args []string) {
	panic("not implemented")
}

func buildWin(ctx context.Context, args []string) {
	panic("not implemented")
}

func runWin(ctx context.Context, args []string) {
	panic("not implemented")
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
