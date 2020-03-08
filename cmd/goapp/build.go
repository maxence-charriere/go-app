package main

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/maxence-charriere/go-app/v5/internal/http"
	"github.com/pkg/errors"
	"github.com/segmentio/conf"
)

type buildConfig struct {
	Name    string `conf:"name"  help:"The name of the app."`
	Force   bool   `conf:"force" help:"Force rebuilding of package that are already up-to-date."`
	Race    bool   `conf:"race"  help:"Enable data race detection."`
	Verbose bool   `conf:"v"     help:"Enable verbose mode."`

	rootDir   string
	serverDir string
	webDir    string
	wasmDir   string
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

	if c.Name == "" {
		c.Name = filepath.Base(rootDir)
	}

	if err := build(ctx, c); err != nil {
		fail("%s", err)
	}

	success("build succeeded")
}

func build(ctx context.Context, c buildConfig) error {
	c.serverDir = filepath.Join(c.rootDir, "cmd", c.Name+"-server")
	c.webDir = filepath.Join(c.serverDir, "web")
	c.wasmDir = filepath.Join(c.rootDir, "cmd", c.Name+"-wasm")

	log("building wasm app")
	if err := buildWasm(ctx, c); err != nil {
		return err
	}

	log("building server")
	if err := buildServer(ctx, c); err != nil {
		return err
	}

	log("installing go wasm support file")
	if err := installWasmExec(c); err != nil {
		return err
	}

	log("generating etag")
	etag := http.GenerateEtag()
	if err := generateEtag(etag, c.webDir); err != nil {
		return err
	}

	log("generating service worker")
	if err := generateServiceWorker(etag, c.webDir); err != nil {
		return err
	}

	log("generating icons")
	return generateProgressiveAppIcons(c)
}

func buildWasm(ctx context.Context, c buildConfig) error {
	out := filepath.Join(c.webDir, "goapp.wasm")

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

	cmd = append(cmd, c.wasmDir)
	return execute(ctx, cmd[0], cmd[1:]...)
}

func buildServer(ctx context.Context, c buildConfig) error {
	out := filepath.Join(c.serverDir, c.Name+"-server")
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

	cmd = append(cmd, c.serverDir)
	return execute(ctx, cmd[0], cmd[1:]...)
}

func installWasmExec(c buildConfig) error {
	webWasmExec := filepath.Join(c.webDir, "wasm_exec.js")
	return generateTemplate(webWasmExec, wasmExecJS, nil)
}

func generateServiceWorker(etag, webDir string) error {
	filename := filepath.Join(webDir, "goapp.js")
	cachePaths := []string{}

	walk := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		cachePath := strings.Replace(path, webDir, "", 1)
		if strings.HasPrefix(cachePath, "/.") {
			return nil
		}
		if cachePath == "/goapp.js" {
			return nil
		}

		cachePaths = append(cachePaths, cachePath)
		return nil
	}

	if err := filepath.Walk(webDir, walk); err != nil {
		return errors.Wrap(err, "getting caching routes failed")
	}

	return generateTemplate(filename, goappJS, struct {
		ETag  string
		Paths []string
	}{
		ETag:  etag,
		Paths: cachePaths,
	})
}

func generateEtag(etag, webDir string) error {
	etagname := filepath.Join(webDir, ".etag")
	if err := ioutil.WriteFile(etagname, []byte(etag), 0666); err != nil {
		return errors.Wrap(err, "generating etag failed")
	}
	return nil
}

func generateProgressiveAppIcons(c buildConfig) error {
	iconname := filepath.Join(c.webDir, "icon.png")
	if _, err := os.Stat(iconname); err != nil {
		iconname = filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "maxence-charriere", "go-app", "logo.png")
	}

	return generateIcons(iconname,
		iconInfo{
			Name:   filepath.Join(c.webDir, "icon-192.png"),
			Width:  192,
			Height: 192,
			Scale:  1,
		},
		iconInfo{
			Name:   filepath.Join(c.webDir, "icon-512.png"),
			Width:  512,
			Height: 512,
			Scale:  1,
		},
	)
}
