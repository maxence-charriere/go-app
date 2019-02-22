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

type buildConfig struct {
	Force   bool `conf:"force" help:"Force rebuilding of package that are already up-to-date."`
	Race    bool `conf:"race"  help:"Enable data race detection."`
	Verbose bool `conf:"v"     help:"Enable verbose mode."`

	rootDir string
}

func buildProject(ctx context.Context, args []string) {
	c := buildConfig{}

	ld := conf.Loader{
		Name:    "goapp build",
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
	c.rootDir = rootDir

	if err := build(ctx, c); err != nil {
		fail("%s", err)
	}

	success("build succeeded")
}

func build(ctx context.Context, c buildConfig) error {
	log("building wasm app")
	if err := buildWasm(ctx, c); err != nil {
		return err
	}

	log("building server")
	if err := buildServer(ctx, c); err != nil {
		fail("%s", err)
	}

	log("installing wasm_exec.js")
	return installWasmExec(c.rootDir)
}

func buildWasm(ctx context.Context, c buildConfig) error {
	pkgName := filepath.Base(c.rootDir) + "-wasm"
	pkg := filepath.Join(c.rootDir, "cmd", pkgName)
	out := filepath.Join(c.rootDir, "web", "goapp.wasm")

	os.Setenv("GOOS", "js")
	os.Setenv("GOARCH", "wasm")
	defer os.Unsetenv("GOOS")
	defer os.Unsetenv("GOARCH")

	cmd := []string{
		"go", "build",
		"-o", out,
	}

	if c.Force {
		cmd = append(cmd, "-a")
	}

	if c.Verbose {
		cmd = append(cmd, "-v")
	}

	cmd = append(cmd, pkg)
	return execute(ctx, cmd[0], cmd[1:]...)
}

func buildServer(ctx context.Context, c buildConfig) error {
	pkgName := filepath.Base(c.rootDir) + "-server"
	pkg := filepath.Join(c.rootDir, "cmd", pkgName)

	out := filepath.Join(c.rootDir, pkgName)
	if runtime.GOOS == "windows" {
		out += ".exe"
	}

	cmd := []string{
		"go", "build",
		"-o", out,
	}

	if c.Force {
		cmd = append(cmd, "-a")
	}

	if c.Race {
		cmd = append(cmd, "-race")
	}

	if c.Verbose {
		cmd = append(cmd, "-v")
	}

	cmd = append(cmd, pkg)
	return execute(ctx, cmd[0], cmd[1:]...)
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
