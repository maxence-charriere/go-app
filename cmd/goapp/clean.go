package main

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/maxence-charriere/go-app/internal/http"
	"github.com/segmentio/conf"
)

type cleanConfig struct {
	Name    string `conf:"name" help:"The name of the app."`
	Verbose bool   `conf:"v"    help:"Enable verbose mode."`
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

	if c.Name == "" {
		c.Name = filepath.Base(rootDir)
	}

	pkgDir := filepath.Join(rootDir, "cmd", c.Name+"-server")
	webDir := filepath.Join(pkgDir, "web")

	serverExec := filepath.Join(pkgDir, c.Name+"-server")
	if runtime.GOOS == "windows" {
		serverExec += ".exe"
	}

	filenames := []string{
		serverExec,
		filepath.Join(webDir, "goapp.wasm"),
		filepath.Join(webDir, "goapp.js"),
		filepath.Join(webDir, "wasm_exec.js"),
		filepath.Join(webDir, ".etag"),
		filepath.Join(webDir, "icon-192.png"),
		filepath.Join(webDir, "icon-512.png"),
	}

	for _, f := range filenames {
		log("removing %s", f)
		os.Remove(f)
	}

	if err := cleanCompressedStaticResources(webDir); err != nil {
		fail("%s", err)
	}

	success("cleaning succeeded")
}

func cleanCompressedStaticResources(webDir string) error {
	walk := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		staticExt := ".gz"
		if etag := http.GetEtag(webDir); etag != "" {
			staticExt = "." + etag + staticExt
		}

		if !strings.HasSuffix(path, staticExt) {
			return nil
		}

		log("removing %s", path)
		return os.Remove(path)
	}

	return filepath.Walk(webDir, walk)
}
