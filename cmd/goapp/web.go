package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/murlokswarm/app"

	"github.com/pkg/errors"
	"github.com/segmentio/conf"
)

type webInitConfig struct {
	Verbose bool `conf:"v" help:"Verbose mode."`
}

type webBuildConfig struct {
	Verbose bool `conf:"v" help:"Verbose mode."`
	Minify  bool `conf:"m" help:"Minify gopherjs file."`
}

func web(ctx context.Context, args []string) {
	ld := conf.Loader{
		Name: "goapp web",
		Args: args,
		Commands: []conf.Command{
			{Name: "help", Help: "Show the goapp web help"},
			{Name: "init", Help: "Create the required files and directories to build a web app."},
			{Name: "build", Help: "Build the web server and generate Gopher.js file."},
		},
	}

	switch cmd, args := conf.LoadWith(nil, ld); cmd {
	case "help":
		ld.PrintHelp(nil)

	case "init":
		initWeb(ctx, args)

	case "build":
		buildWeb(ctx, args)

	default:
		panic("unreachable")
	}

}

func initWeb(ctx context.Context, args []string) {
	config := webInitConfig{}

	ld := conf.Loader{
		Name:    "web init",
		Args:    args,
		Usage:   "[options...] [packages...]",
		Sources: []conf.Source{conf.NewEnvSource("GOAPP", os.Environ()...)},
	}

	defer func() {
		err := recover()
		if err != nil {
			ld.PrintHelp(config)
			ld.PrintError(errors.Errorf("%s", err))
			os.Exit(-1)
		}
	}()

	_, unusedArgs := conf.LoadWith(&config, ld)
	roots, err := packageRoots(unusedArgs)
	if err != nil {
		panic(err)
	}

	for _, root := range roots {
		if err = initPackage(root); err != nil {
			panic(err)
		}

		if err = goGetGopherJS(config); err != nil {
			panic(err)
		}
	}
}

func goGetGopherJS(config webInitConfig) error {
	args := []string{
		"get",
		"-u",
	}
	if config.Verbose {
		args = append(args, "-v")
	}
	args = append(args, "github.com/gopherjs/gopherjs")

	return execute("go", args...)
}

func buildWeb(ctx context.Context, args []string) {
	config := webBuildConfig{
		Minify: true,
	}

	ld := conf.Loader{
		Name:    "web build",
		Args:    args,
		Usage:   "[options...] [package]",
		Sources: []conf.Source{conf.NewEnvSource("GOAPP", os.Environ()...)},
	}

	_, roots := conf.LoadWith(&config, ld)
	if len(roots) == 0 {
		roots = []string{"."}
	}

	root := roots[0]

	if err := goBuild(root, config.Verbose); err != nil {
		app.Error("go build:", err)
	}

	if err := gopherJSBuild(root, config.Verbose, config.Minify); err != nil {
		app.Error("gopherjs build:", err)
	}
}

func gopherJSBuild(target string, verbose, minify bool) error {
	args := []string{"build"}
	if verbose {
		args = append(args, "-v")
	}
	if minify {
		args = append(args, "-m")
	}
	args = append(args, "-o", filepath.Join(target, "resources", "goapp.js"))
	args = append(args, target)
	return execute("gopherjs", args...)
}
