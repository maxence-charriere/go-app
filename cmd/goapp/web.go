package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/murlokswarm/app/internal/file"
	"github.com/segmentio/conf"
)

func web(ctx context.Context, args []string) {
	ld := conf.Loader{
		Name: "goapp web",
		Args: args,
		Commands: []conf.Command{
			{Name: "init", Help: "Download gopherjs and create the required files and directories."},
			{Name: "build", Help: "Build the web server and generate Gopher.js file."},
			{Name: "run", Help: "Run the server and launch the client in the default browser."},
			{Name: "help", Help: "Show the web help"},
		},
	}

	switch cmd, args := conf.LoadWith(nil, ld); cmd {
	case "help":
		ld.PrintHelp(nil)

	case "init":
		initWeb(ctx, args)

	case "build":
		buildWeb(ctx, args)

	case "run":
		runWeb(ctx, args)

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

	pkg, err := newWebPackage(roots[0])
	if err != nil {
		fail("%s", err)
	}

	if err = pkg.Build(ctx, c); err != nil {
		fail("%s", err)
	}

	printSuccess("build succeeded")
}

type webRunConfig struct {
	Addr    string   `conf:"addr"    help:"The server bind address."`
	Args    []string `conf:"args"    help:"The arguments to launch the server."`
	Browser string   `conf:"browser" help:"The browser to use."`
	Minify  bool     `conf:"m"       help:"Minify gopherjs file."`
	Verbose bool     `conf:"v"       help:"Enable verbose mode."`
}

func runWeb(ctx context.Context, args []string) {
	c := webRunConfig{
		Addr:   "http://127.0.0.1:7042",
		Minify: true,
	}

	ld := conf.Loader{
		Name:    "web run",
		Args:    args,
		Usage:   "[options...] [*.wapp]",
		Sources: []conf.Source{conf.NewEnvSource("GOAPP", os.Environ()...)},
	}

	_, roots := conf.LoadWith(&c, ld)
	verbose = c.Verbose

	wappname := "."
	if len(roots) != 0 {
		wappname = roots[0]
	}

	if !strings.HasSuffix(wappname, ".wapp") {
		printVerbose("building package")
		pkg, err := newWebPackage(wappname)
		if err != nil {
			fail("%s", err)
		}

		if err = pkg.Build(ctx, webBuildConfig{
			Minify: c.Minify,
		}); err != nil {
			fail("%s", err)
		}

		wappname = pkg.name
	}

	server := filepath.Base(wappname)
	server = strings.TrimSuffix(server, ".wapp")
	server = filepath.Join(wappname, server)

	go launchNavigator(ctx, c.Addr)

	printVerbose("starting server")
	if err := execute(ctx, server, args...); err != nil {
		fail("%s", err)
	}
}

func launchNavigator(ctx context.Context, url string) {
	time.Sleep(time.Millisecond * 200)
	printVerbose("starting client")
	
	// switch runtime.GOOS{
	// case "darwin":
	// 	open = "open"

	// case "windows":

	// case "linux":

	// default:
	// 	fail("you are not on Linux, MacOS or Windows")
	// }

	execute(ctx, "open", url)
}

type webPackage struct {
	workingDir       string
	buildDir         string
	gopherJSBuildDir string
	buildResources   string
	name             string
	resources        string
	goExec           string
	gopherJS         string
	minify           bool
}

func newWebPackage(buildDir string) (*webPackage, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	gopherJSBuildDIr := buildDir

	if buildDir, err = filepath.Abs(buildDir); err != nil {
		return nil, err
	}

	name := filepath.Base(buildDir) + ".wapp"

	return &webPackage{
		workingDir:       wd,
		buildDir:         buildDir,
		gopherJSBuildDir: gopherJSBuildDIr,
		buildResources:   filepath.Join(buildDir, "resources"),
		name:             filepath.Join(wd, name),
		resources:        filepath.Join(wd, name, "resources"),
		goExec:           filepath.Join(wd, name, filepath.Base(buildDir)),
		gopherJS:         filepath.Join(wd, name, "resources", "goapp.js"),
	}, nil
}

func (pkg *webPackage) Build(ctx context.Context, c webBuildConfig) error {
	pkg.minify = c.Minify
	name := filepath.Base(pkg.name)

	printVerbose("creating %s", name)
	if err := pkg.createPackage(); err != nil {
		return err
	}

	printVerbose("building go server")
	if err := pkg.buildGoExec(ctx); err != nil {
		return err
	}

	printVerbose("building gopherjs client")
	if err := pkg.buildGopherJS(ctx); err != nil {
		return err
	}

	printVerbose("syncing resources")
	return pkg.syncResources()
}

func (pkg *webPackage) createPackage() error {
	dirs := []string{
		pkg.name,
		pkg.resources,
	}

	for _, d := range dirs {
		if err := os.MkdirAll(d, os.ModeDir|0755); err != nil {
			return err
		}
	}

	return nil
}

func (pkg *webPackage) buildGoExec(ctx context.Context) error {
	if err := goBuild(ctx, pkg.buildDir, "-o", pkg.goExec); err != nil {
		return err
	}

	return nil
}

func (pkg *webPackage) buildGopherJS(ctx context.Context) error {
	if runtime.GOOS == "windows" {
		os.Setenv("GOOS", "linux")
		defer os.Unsetenv("GOOS")
	}

	cmd := []string{"gopherjs", "build", "-o", pkg.gopherJS}

	if pkg.minify {
		cmd = append(cmd, "-m")
	}

	if verbose {
		cmd = append(cmd, "-v")
	}

	cmd = append(cmd, pkg.gopherJSBuildDir)

	fmt.Println(cmd)
	return execute(ctx, cmd[0], cmd[1:]...)
}

func (pkg *webPackage) syncResources() error {
	return file.Sync(pkg.resources, pkg.buildResources)
}
