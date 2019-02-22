package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/segmentio/conf"
)

type cleanConfig struct {
	Verbose bool `conf:"v" help:"Enable verbose mode."`
}

func cleanProject(ctx context.Context, args []string) {
	c := cleanConfig{}

	ld := conf.Loader{
		Name:    "goapp clean",
		Args:    args,
		Usage:   "[options...] [package]",
		Sources: []conf.Source{conf.NewEnvSource("GOAPP", os.Environ()...)},
	}

	_, args = conf.LoadWith(&c, ld)
	verbose = c.Verbose

	pkg := "."
	if len(args) != 0 {
		pkg = args[0]
	}

	rootDir, err := filepath.Abs(pkg)
	if err != nil {
		fail("%s", err)
	}

	pkgName := filepath.Base(rootDir)

	filenames := []string{
		filepath.Join(rootDir, pkgName+"-server"),
		filepath.Join(rootDir, "web", "goapp.wasm"),
		filepath.Join(rootDir, "web", "wasm_exec.js"),
	}

	for _, f := range filenames {
		log("removing %s", f)
		if err := os.Remove(f); err != nil {
			warn("%s", err)
		}
	}
}
