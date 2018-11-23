package main

import (
	"context"
	"os"

	"github.com/segmentio/conf"
)

type webInitConfig struct {
	Verbose bool `conf:"v" help:"Enable verbose mode."`
}

type webBuildConfig struct {
	Output  string `conf:"o"      help:"The path where the package is saved."`
	Minify  bool   `conf:"minify" help:"Minify gopherjs file."`
	Force   bool   `conf:"force"  help:"Force rebuilding of package that are already up-to-date."`
	Race    bool   `conf:"race"   help:"Enable data race detection."`
	Verbose bool   `conf:"v"      help:"Enable verbose mode."`
}

type webRunConfig struct {
	Output  string `conf:"o"      help:"The path where the package is saved."`
	Addr    string `conf:"addr"   help:"The server bind address."`
	Browser bool   `conf:"b"      help:"Run the client."`
	Chrome  bool   `conf:"chrome" help:"Run the client with Google Chrome."`
	Minify  bool   `conf:"minify" help:"Minify gopherjs file."`
	Force   bool   `conf:"force"  help:"Force rebuilding of package that are already up-to-date."`
	Race    bool   `conf:"race"   help:"Enable data race detection."`
	Verbose bool   `conf:"v"      help:"Enable verbose mode."`
}

type webCleanConfig struct {
	Output  string `conf:"o" help:"The path where the package is saved."`
	Verbose bool   `conf:"v" help:"Enable verbose mode."`
}

func web(ctx context.Context, args []string) {
	ld := conf.Loader{
		Name: "goapp web",
		Args: args,
		Commands: []conf.Command{
			{Name: "init", Help: "Download gopherjs and create the required files and directories."},
			{Name: "build", Help: "Build a web app."},
			{Name: "run", Help: "Run the server and launch the client in the default browser."},
			{Name: "clean", Help: "Delete a web app."},
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

	case "clean":
		cleanWeb(ctx, args)

	default:
		panic("unreachable")
	}
}

func initWeb(ctx context.Context, args []string) {
	c := webInitConfig{}

	ld := conf.Loader{
		Name:    "web init",
		Args:    args,
		Usage:   "[options...] [package]",
		Sources: []conf.Source{conf.NewEnvSource("GOAPP", os.Environ()...)},
	}

	_, args = conf.LoadWith(&c, ld)
	verbose = c.Verbose

	sources := "."
	if len(args) != 0 {
		sources = args[0]
	}

	pkg := WebPackage{
		Sources: sources,
		Verbose: c.Verbose,
		Log:     printVerbose,
	}

	if err := pkg.Init(ctx); err != nil {
		fail("%s", err)
	}

	printSuccess("init succeeded")
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

	_, args = conf.LoadWith(&c, ld)
	verbose = c.Verbose

	sources := "."
	if len(args) != 0 {
		sources = args[0]
	}

	pkg := WebPackage{
		Sources: sources,
		Output:  c.Output,
		Minify:  c.Minify,
		Force:   c.Force,
		Race:    c.Race,
		Verbose: c.Verbose,
		Log:     printVerbose,
	}

	if err := pkg.Build(ctx); err != nil {
		fail("%s", err)
	}

	printSuccess("build succeeded")
}

func runWeb(ctx context.Context, args []string) {
	c := webRunConfig{
		Addr:   ":9001",
		Minify: true,
	}

	ld := conf.Loader{
		Name:    "web run",
		Args:    args,
		Usage:   "[options...] [package]",
		Sources: []conf.Source{conf.NewEnvSource("GOAPP", os.Environ()...)},
	}

	_, args = conf.LoadWith(&c, ld)
	verbose = c.Verbose

	sources := "."
	if len(args) != 0 {
		sources = args[0]
	}

	pkg := WebPackage{
		Sources: sources,
		Output:  c.Output,
		Addr:    c.Addr,
		Chrome:  c.Chrome,
		Browser: c.Browser,
		Minify:  c.Minify,
		Force:   c.Force,
		Race:    c.Race,
		Verbose: c.Verbose,
		Log:     printVerbose,
	}

	os.Setenv("GOAPP_SERVER_ADDR", c.Addr)
	defer os.Unsetenv("GOAPP_SERVER_ADDR")

	if err := pkg.Run(ctx); err != nil {
		fail("%s", err)
	}
}

func cleanWeb(ctx context.Context, args []string) {
	c := webCleanConfig{}

	ld := conf.Loader{
		Name:    "web run",
		Args:    args,
		Usage:   "[options...] [package]",
		Sources: []conf.Source{conf.NewEnvSource("GOAPP", os.Environ()...)},
	}

	_, args = conf.LoadWith(&c, ld)
	verbose = c.Verbose

	sources := "."
	if len(args) != 0 {
		sources = args[0]
	}

	pkg := WebPackage{
		Sources: sources,
		Output:  c.Output,
		Verbose: c.Verbose,
		Log:     printVerbose,
	}

	if err := pkg.Clean(ctx); err != nil {
		fail("%s", err)
	}
}
