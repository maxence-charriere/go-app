package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/segmentio/conf"
)

type initConfig struct {
	Name    string `conf:"name" help:"The name of the app."`
	Verbose bool   `conf:"v" help:"Enable verbose mode."`
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
	if err = initProjectLayout(rootDir, c.Name); err != nil {
		fail("%s", err)
	}

	success("initialization succeeded")
}

func initProjectLayout(rootDir, name string) error {
	if name == "" {
		name = filepath.Base(rootDir)
	}

	serverdir := filepath.Join(rootDir, "cmd", name+"-server")
	wasmdir := filepath.Join(rootDir, "cmd", name+"-wasm")

	dirs := []string{
		serverdir,
		wasmdir,
		filepath.Join(serverdir, "web"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fail("%s", err)
		}
	}

	serverMainName := filepath.Join(serverdir, "main.go")
	wasmMainName := filepath.Join(wasmdir, "main.go")

	if err := initMain(serverMainName, mainServer); err != nil {
		return err
	}
	return initMain(wasmMainName, mainWasm)
}

func initMain(filename, tmpl string) error {
	if _, err := os.Stat(filename); err == nil {
		return nil
	}

	log("generating %s", filename)
	return generateTemplate(filename, tmpl, nil)
}
