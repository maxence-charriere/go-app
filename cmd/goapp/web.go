package main

import (
	"context"
	"os"
	"path/filepath"
	"runtime"

	"github.com/segmentio/conf"
)

func web(ctx context.Context, args []string) {
	ld := conf.Loader{
		Name: "goapp web",
		Args: args,
		Commands: []conf.Command{
			{Name: "help", Help: "Show the web help"},
			{Name: "init", Help: "Download gopherjs and create the required files and directories."},
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

type webInitConfig struct {
	Verbose bool `conf:"v" help:"Enable verbose mode."`
}

func initWeb(ctx context.Context, args []string) {
	c := webInitConfig{}

	ld := conf.Loader{
		Name:    "web init",
		Args:    args,
		Usage:   "[options...] [packages...]",
		Sources: []conf.Source{conf.NewEnvSource("GOAPP", os.Environ()...)},
	}

	_, unusedArgs := conf.LoadWith(&c, ld)
	verbose = c.Verbose

	roots, err := packageRoots(unusedArgs)
	if err != nil {
		failWithHelp(&ld, "%s", err)
	}

	printVerbose("get gopherjs")
	if err = goGetGopherJS(ctx); err != nil {
		failWithHelp(&ld, "%s", err)
	}

	for _, root := range roots {
		if err = initPackage(root); err != nil {
			failWithHelp(&ld, "%s", err)
		}
	}

	printSuccess("init succeeded")
}

func goGetGopherJS(ctx context.Context) error {
	args := []string{"get", "-u"}

	if verbose {
		args = append(args, "-v")
	}

	args = append(args, "github.com/gopherjs/gopherjs")
	return execute(ctx, "go", args...)
}

type webBuildConfig struct {
	Minify  bool `conf:"m" help:"Minify gopherjs file."`
	Verbose bool `conf:"v" help:"Enable verbose mode."`
}

func buildWeb(ctx context.Context, args []string) {
	c := webBuildConfig{
		Minify: true,
	}

	ld := conf.Loader{
		Name:    "web build",
		Args:    args,
		Usage:   "[options...] [package]",
		Sources: []conf.Source{conf.NewEnvSource("GOAPP", os.Environ()...)},
	}

	_, roots := conf.LoadWith(&c, ld)
	verbose = c.Verbose

	if len(roots) == 0 {
		roots = []string{"."}
	}

	root := roots[0]

	printVerbose("building go server")
	if err := goBuild(ctx, root); err != nil {
		printErr("go build failed: %s", err)
		return
	}

	printVerbose("building gopherjs client")
	if err := gopherJSBuild(ctx, root, c.Minify); err != nil {
		printErr("gopherjs build failed: %s", err)
		return
	}

	printSuccess("build succeeded")
}

func gopherJSBuild(ctx context.Context, target string, minify bool) error {
	cmd := []string{}

	if runtime.GOOS == "windows" {
		os.Setenv("GOOS", "darwin")
	}

	cmd = append(cmd, "gopherjs", "build", "-v")

	if minify {
		cmd = append(cmd, "-m")
	}

	if verbose {
		cmd = append(cmd, "-v")
	}

	cmd = append(cmd, "-o", filepath.Join(target, "resources", "goapp.js"))
	cmd = append(cmd, target)
	return execute(ctx, cmd[0], cmd[1:]...)
}
