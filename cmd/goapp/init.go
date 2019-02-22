package main

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
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

	log("installing wasm_exec.js")
	if err = installWasmExec(rootDir); err != nil {
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

func installWasmExec(rootDir string) error {
	wasmExec := filepath.Join(runtime.GOROOT(), "misc", "wasm", "wasm_exec.js")
	webWasmExec := filepath.Join(rootDir, "web", filepath.Base(wasmExec))

	src, err := os.Open(wasmExec)
	if err != nil {
		return errors.Wrapf(err, "opening %q failed", wasmExec)
	}
	defer src.Close()

	dst, err := os.Create(webWasmExec)
	if err != nil {
		return errors.Wrapf(err, "creating %q failed", webWasmExec)
	}
	defer src.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return errors.Wrapf(err, "copying %q failed", wasmExec)
	}
	return nil
}
