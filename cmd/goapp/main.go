//go:generate go run gen.go
//go:generate go fmt

package main

import (
	"context"
	"os"
	"os/signal"

	_ "github.com/maxence-charriere/app"
	"github.com/segmentio/conf"
)

func main() {
	ld := conf.Loader{
		Name: "goapp",
		Args: os.Args[1:],
		Commands: []conf.Command{
			{Name: "init", Help: "Init the project layout and install the wasm dependencies."},
			{Name: "build", Help: "Build the wasm app and its server."},
			{Name: "run", Help: "Build and run the wasm app and its server."},
			{Name: "clean", Help: "Delete the wasm app and its server."},
			{Name: "update", Help: "Update to the latest version."},
			{Name: "help", Help: "Show the help."},
		},
	}

	ctx, cancel := ctxWithSignals(context.Background(), os.Interrupt)
	defer cancel()

	switch cmd, args := conf.LoadWith(nil, ld); cmd {
	case "init":
		initProject(ctx, args)

	case "build":
		buildProject(ctx, args)

	case "run":
		runProject(ctx, args)

	case "clean":
		cleanProject(ctx, args)

	case "update":
		update(ctx, args)

	case "help":
		ld.PrintHelp(nil)

	default:
		panic("unreachable")
	}
}

func ctxWithSignals(ctx context.Context, sigs ...os.Signal) (context.Context, func()) {
	ctx, cancel := context.WithCancel(ctx)
	sigc := make(chan os.Signal)
	signal.Notify(sigc, sigs...)

	go func() {
		defer close(sigc)
		<-sigc
		cancel()
	}()

	return ctx, cancel
}
