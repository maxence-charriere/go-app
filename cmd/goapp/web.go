package main

import (
	"context"
	"os"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
	"github.com/segmentio/conf"
)

type webInitConfig struct {
}

type webBuildConfig struct {
	Minify bool `conf:"m" help:"Minify gopherjs file."`
}

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

	if err = goGetGopherJS(); err != nil {
		panic(err)
	}

	for _, root := range roots {
		if err = initPackage(root); err != nil {
			panic(err)
		}
	}

	printSuccess("init succeeded")
}

func goGetGopherJS() error {
	return execute("go",
		"get",
		"-u",
		"-v",
		"github.com/gopherjs/gopherjs",
	)
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

	if err := goBuild(root); err != nil {
		printErr("go build failed: %s", err)
		return
	}

	if err := gopherJSBuild(root, config.Minify); err != nil {
		printErr("gopherjs build failed: %s", err)
		return
	}

	printSuccess("build succeeded")
}

func gopherJSBuild(target string, minify bool) error {
	cmd := []string{}

	if runtime.GOOS == "windows" {
		os.Setenv("GOOS", "darwin")
	}

	cmd = append(cmd, "gopherjs", "build", "-v")

	if minify {
		cmd = append(cmd, "-m")
	}

	cmd = append(cmd, "-o", filepath.Join(target, "resources", "goapp.js"))
	cmd = append(cmd, target)
	return execute(cmd[0], cmd[1:]...)
}
