// +build darwin,amd64

package main

import (
	"context"
	"os"

	"github.com/murlokswarm/app"
	"github.com/pkg/errors"
	"github.com/segmentio/conf"
)

func mac(ctx context.Context, args []string) {
	ld := conf.Loader{
		Name: "goapp mac",
		Args: args,
		Commands: []conf.Command{
			{Name: "help", Help: "Show the macOS help"},
			{Name: "init", Help: "Download macOS SDK and create required file and directories."},
		},
	}

	switch cmd, args := conf.LoadWith(nil, ld); cmd {
	case "help":
		ld.PrintHelp(nil)

	case "init":
		initMac(ctx, args)

	default:
		panic("unreachable")
	}
}

func initMac(ctx context.Context, args []string) {
	config := struct{}{}

	ld := conf.Loader{
		Name:    "mac init",
		Args:    args,
		Usage:   "[options...] [packages...]",
		Sources: []conf.Source{conf.NewEnvSource("GOAPP", os.Environ()...)},
	}

	defer func() {
		err := recover()
		if err != nil {
			ld.PrintHelp(nil)
			ld.PrintError(errors.Errorf("%s", err))
			os.Exit(-1)
		}
	}()

	_, unusedArgs := conf.LoadWith(&config, ld)
	roots, err := packageRoots(unusedArgs)
	if err != nil {
		panic(err)
	}

	if err := execute("xcode-select", "--install"); err != nil {
		app.Error(errors.Wrap(err, "xcode-select"))
		return
	}

	for _, root := range roots {
		if err = initPackage(root); err != nil {
			app.Error(errors.Wrap(err, "init package"))
			return
		}
	}
}

func openCommand() string {
	return "open"
}
