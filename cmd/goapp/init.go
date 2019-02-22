package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/segmentio/conf"
)

type initConfig struct {
	Verbose bool `conf:"v" help:"Enable verbose mode."`
}

func initProject(ctx context.Context, args []string) {
	c := initConfig{}

	ld := conf.Loader{
		Name:    "goapp init",
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

	log("initializing project layout")
	if err = initProjectDirectories(rootDir); err != nil {
		fail("%s", err)
	}

	success("initialization succeeded")
}

func initProjectDirectories(rootDir string) error {
	pkgName := filepath.Base(rootDir)

	dirs := []string{
		filepath.Join(rootDir, "cmd", pkgName+"-server"),
		filepath.Join(rootDir, "cmd", pkgName+"-wasm"),
		filepath.Join(rootDir, "web"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fail("%s", err)
		}
	}

	return nil
}
