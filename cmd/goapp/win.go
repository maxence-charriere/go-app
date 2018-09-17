// +build windows

package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/segmentio/conf"
)

func win(ctx context.Context, args []string) {
	ld := conf.Loader{
		Name: "goapp win",
		Args: args,
		Commands: []conf.Command{
			{Name: "init", Help: "Download the Windows dev tools and create required files and directories."},
			{Name: "build", Help: "Build the Windows app."},
			{Name: "run", Help: "Run a Windows app and capture its logs."},
			{Name: "help", Help: "Show the Windows help"},
		},
	}

	switch cmd, args := conf.LoadWith(nil, ld); cmd {
	case "init":
		initWin(ctx, args)

	case "build":
		buildWin(ctx, args)

	case "run":
		runWin(ctx, args)

	case "help":
		ld.PrintHelp(nil)

	default:
		panic("unreachable")
	}
}

type winInitConfig struct {
	Verbose bool `conf:"v" help:"Enable verbose mode."`
}

func initWin(ctx context.Context, args []string) {
	c := winInitConfig{}

	ld := conf.Loader{
		Name:    "win init",
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

	for _, root := range roots {
		if err = initPackage(root); err != nil {
			fail("init %s failed: %s", root, err)
		}
	}

	printWarn("install Windows 10 SDK: https://developer.microsoft.com/en-US/windows/downloads/windows-10-sdk")
	printWarn("install Desktop app converter: https://aka.ms/converter")

	printVerbose("installing dev certificate")
	os.Chdir(certMgr())

	if err = execute(ctx, "powershell",
		`.\Certmgr.exe`,
		"/add", certificate(),
		"/s", "/r",
		"localMachine",
		"root",
	); err != nil {
		fail("installing dev certificate failed: %s", err)
	}

	printSuccess("init succeeded")
}

type winBuilConfig struct {
	Output  string `conf:"o" help:"The output."`
	Verbose bool   `conf:"v" help:"Enable verbose mode."`
}

func buildWin(ctx context.Context, args []string) {
	c := winBuilConfig{}

	ld := conf.Loader{
		Name:    "win build",
		Args:    args,
		Usage:   "[options...] [package]",
		Sources: []conf.Source{conf.NewEnvSource("GOAPP", os.Environ()...)},
	}

	_, roots := conf.LoadWith(&c, ld)
	verbose = c.Verbose

	if len(roots) == 0 {
		roots = []string{"."}
	}

	printVerbose("building package")
	pkg, err := newWinPackage(roots[0], c.Output)
	if err != nil {
		fail("%s", err)
	}

	if err = pkg.Build(ctx, c); err != nil {
		fail("%s", err)
	}

	printSuccess("build succeeded")
}

func runWin(ctx context.Context, args []string) {
	panic("not implemented")
}

func mac(ctx context.Context, args []string) {
	printErr("you are not on MacOS!")
	os.Exit(-1)
}

func init() {
	greenColor = ""
	redColor = ""
	orangeColor = ""
	defaultColor = ""
}

func certMgr() string {
	return filepath.Join(
		os.Getenv("ProgramFiles(x86)"),
		"Windows Kits", "10", "bin", "10.0.17134.0", "x64",
	)
}

func certificate() string {
	return filepath.Join(
		os.Getenv("GOPATH"),
		"src",
		"github.com",
		"murlokswarm",
		"app",
		"cmd",
		"goapp",
		"certificates",
		"win.cer",
	)
}
