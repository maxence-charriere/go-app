package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/murlokswarm/app/internal/file"
	"github.com/segmentio/conf"
)

type winInitConfig struct {
	Verbose bool `conf:"v" help:"Enable verbose mode."`
}

type winBuildConfig struct {
	Output       string `conf:"o"        help:"The path where the package is saved."`
	Architecture string `conf:"arch"     help:"The targetted architecture."`
	SDK          string `conf:"sdk"      help:"The path of the Windows 10 SDK directory."`
	AppStore     bool   `conf:"appstore" help:"Creates a .pkg to be uploaded on the app store."`
	Force        bool   `conf:"force"    help:"Force rebuilding of package that are already up-to-date."`
	Race         bool   `conf:"race"     help:"Enable data race detection."`
	Verbose      bool   `conf:"v"        help:"Enable verbose mode."`
}

type winRunConfig struct {
	Output       string `conf:"o"     help:"The path where the package is saved."`
	Architecture string `conf:"arch"  help:"The targetted architecture."`
	SDK          string `conf:"sdk"   help:"The path of the Windows 10 SDK directory."`
	Force        bool   `conf:"force" help:"Force rebuilding of package that are already up-to-date."`
	Race         bool   `conf:"race"  help:"Enable data race detection."`
	Verbose      bool   `conf:"v"     help:"Enable verbose mode."`
}

type winCleanConfig struct {
	Output  string `conf:"o" help:"The path where the package is saved."`
	Verbose bool   `conf:"v" help:"Enable verbose mode."`
}

func win(ctx context.Context, args []string) {
	ld := conf.Loader{
		Name: "goapp win",
		Args: args,
		Commands: []conf.Command{
			{Name: "init", Help: "Download the Windows 10 dev tools and create required directories."},
			{Name: "build", Help: "Build the Windows app."},
			{Name: "run", Help: "Run a Windows app."},
			{Name: "clean", Help: "Delete a Windows app and its temporary build files."},
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

	case "clean":
		cleanWin(ctx, args)

	case "help":
		ld.PrintHelp(nil)

	default:
		panic("unreachable")
	}
}

func initWin(ctx context.Context, args []string) {
	c := winInitConfig{}

	ld := conf.Loader{
		Name:    "win init",
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

	pkg := WinPackage{
		Sources: sources,
		Verbose: c.Verbose,
		Log:     printVerbose,
	}

	if err := pkg.Init(ctx); err != nil {
		fail("%s", err)
	}

	printSuccess("init succeeded")
}

func buildWin(ctx context.Context, args []string) {
	c := winBuildConfig{
		Architecture: defaultWinArchitecture(),
		SDK:          winSDKDirectory(defaultWinSSDKRoot),
	}

	ld := conf.Loader{
		Name:    "win build",
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

	pkg := WinPackage{
		Sources:      sources,
		Output:       c.Output,
		Architecture: c.Architecture,
		SDK:          c.SDK,
		AppStore:     c.AppStore,
		Verbose:      c.Verbose,
		Force:        c.Force,
		Race:         c.Race,
		Log:          printVerbose,
	}

	if err := pkg.Build(ctx); err != nil {
		fail("%s", err)
	}

	printSuccess("build succeeded")
}

func runWin(ctx context.Context, args []string) {
	c := winRunConfig{
		Architecture: defaultWinArchitecture(),
		SDK:          winSDKDirectory(defaultWinSSDKRoot),
	}

	ld := conf.Loader{
		Name:    "win run",
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

	pkg := WinPackage{
		Sources:      sources,
		Output:       c.Output,
		Architecture: c.Architecture,
		SDK:          c.SDK,
		Verbose:      c.Verbose,
		Force:        c.Force,
		Race:         c.Race,
		Log:          printVerbose,
	}

	if err := pkg.Run(ctx); err != nil {
		fail("%s", err)
	}
}

func cleanWin(ctx context.Context, args []string) {
	c := winCleanConfig{}

	ld := conf.Loader{
		Name:    "win clean",
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

	pkg := WinPackage{
		Sources: sources,
		Output:  c.Output,
		Verbose: c.Verbose,
		Log:     printVerbose,
	}

	if err := pkg.Clean(ctx); err != nil {
		fail("%s", err)
	}

	printSuccess("clean succeeded")
}

// func runWin(ctx context.Context, args []string) {
// 	c := winRunConfig{
// 		LogsAddr: ":9000",
// 	}

// 	ld := conf.Loader{
// 		Name:    "win run",
// 		Args:    args,
// 		Usage:   "[options...] [app name]",
// 		Sources: []conf.Source{conf.NewEnvSource("GOAPP", os.Environ()...)},
// 	}

// 	_, roots := conf.LoadWith(&c, ld)
// 	verbose = c.Verbose

// 	if len(roots) == 0 {
// 		roots = []string{"."}
// 	}

// 	appname := roots[0]

// 	if !strings.HasSuffix(appname, ".app") {
// 		printVerbose("building package")
// 		pkg, err := newWinPackage(roots[0], "")
// 		if err != nil {
// 			fail("%s", err)
// 		}

// 		if err = pkg.Build(ctx, winBuilConfig{}); err != nil {
// 			fail("%s", err)
// 		}

// 		appname = pkg.manifest.Scheme
// 	}

// 	_, appname = filepath.Split(appname)
// 	appname = strings.TrimSuffix(appname, ".app")

// 	go listenLogs(ctx, c.LogsAddr)
// 	time.Sleep(time.Millisecond * 500)

// 	os.Setenv("GOAPP_LOGS_ADDR", c.LogsAddr)
// 	os.Setenv("GOAPP_DEBUG", fmt.Sprintf("%v", c.Debug))

// 	printVerbose("running %s", appname)
// 	if err := execute(ctx, "powershell", "start", fmt.Sprintf("%s://goapp", appname)); err != nil {
// 		fail("%s", err)
// 	}

// 	<-ctx.Done()
// 	if err := ctx.Err(); err != nil {
// 		printErr("%s", ctx.Err())
// 	}
// }

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
		file.RepoPath(),
		"cmd",
		"goapp",
		"certificates",
		"win.cer",
	)
}
