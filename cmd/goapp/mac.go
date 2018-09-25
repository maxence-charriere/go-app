// +build darwin,amd64

package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/murlokswarm/app/internal/logs"

	"github.com/segmentio/conf"
)

func mac(ctx context.Context, args []string) {
	ld := conf.Loader{
		Name: "goapp mac",
		Args: args,
		Commands: []conf.Command{
			{Name: "init", Help: "Download MacOS SDK and create required files and directories."},
			{Name: "build", Help: "Build the MacOS app."},
			{Name: "run", Help: "Run a MacOS app and capture its logs."},
			{Name: "help", Help: "Show the MacOS help"},
		},
	}

	switch cmd, args := conf.LoadWith(nil, ld); cmd {
	case "init":
		initMac(ctx, args)

	case "build":
		buildMac(ctx, args)

	case "run":
		runMac(ctx, args)

	case "help":
		ld.PrintHelp(nil)

	default:
		panic("unreachable")
	}
}

type macInitConfig struct {
	Verbose bool `conf:"v" help:"Enable verbose mode."`
}

func initMac(ctx context.Context, args []string) {
	c := macInitConfig{}

	ld := conf.Loader{
		Name:    "mac init",
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

	printVerbose("checking for xcode-select...")
	execute(ctx, "xcode-select", "--install")

	for _, root := range roots {
		if err = initPackage(root); err != nil {
			fail("init %s failed: %s", root, err)
		}
	}

	printSuccess("init succeeded")
}

type macBuildConfig struct {
	Output           string `conf:"o"                 help:"The output."`
	DeploymentTarget string `conf:"deployment-target" help:"The MacOS version."`
	Sign             string `conf:"sign"              help:"The signing identifier to sign the app.\n\t\033[95msecurity find-identity -v -p codesigning\033[00m to see signing identifiers.\n\thttps://developer.apple.com/library/content/documentation/Security/Conceptual/CodeSigningGuide/Procedures/Procedures.html to create one."`
	Sandbox          bool   `conf:"sandbox"           help:"Setup the app to run in sandbox mode."`
	AppStore         bool   `conf:"appstore"          help:"Report whether the app will be packaged to be uploaded on the app store."`
	Force            bool   `conf:"a"                 help:"Force rebuilding of packages that are already up-to-date."`
	Race             bool   `conf:"race"              help:"Enable data race detection."`
	Verbose          bool   `conf:"v"                 help:"Enable verbose mode."`
}

func buildMac(ctx context.Context, args []string) {
	c := macBuildConfig{
		DeploymentTarget: macOSVersion(),
	}

	ld := conf.Loader{
		Name:    "mac build",
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
	pkg, err := newMacPackage(roots[0], c.Output)
	if err != nil {
		fail("%s", err)
	}

	if err = pkg.Build(ctx, c); err != nil {
		fail("%s", err)
	}

	printSuccess("build succeeded")
}

type macRunConfig struct {
	LogsAddr         string `conf:"logs-addr"         help:"The address used to listen app logs." validate:"nonzero"`
	DeploymentTarget string `conf:"deployment-target" help:"The MacOS version."`
	Sign             string `conf:"sign"              help:"The signing identifier to sign the app.\n\t\033[95msecurity find-identity -v -p codesigning\033[00m to see signing identifiers.\n\thttps://developer.apple.com/library/content/documentation/Security/Conceptual/CodeSigningGuide/Procedures/Procedures.html to create one."`
	Sandbox          bool   `conf:"sandbox"           help:"Setup the app to run in sandbox mode."`
	Debug            bool   `conf:"d"                 help:"Enable debug mode is enabled."`
	Force            bool   `conf:"a"                 help:"Force rebuilding of packages that are already up-to-date."`
	Race             bool   `conf:"race"              help:"Enable data race detection."`
	Verbose          bool   `conf:"v"                 help:"Enable verbose mode."`
}

func runMac(ctx context.Context, args []string) {
	c := macRunConfig{
		LogsAddr:         ":9000",
		DeploymentTarget: macOSVersion(),
	}

	ld := conf.Loader{
		Name:    "mac run",
		Args:    args,
		Usage:   "[options...] [.app path]",
		Sources: []conf.Source{conf.NewEnvSource("GOAPP", os.Environ()...)},
	}

	_, roots := conf.LoadWith(&c, ld)
	verbose = c.Verbose

	if len(roots) == 0 {
		roots = []string{"."}
	}

	appname := roots[0]

	if !strings.HasSuffix(appname, ".app") {
		printVerbose("building package")
		pkg, err := newMacPackage(roots[0], "")
		if err != nil {
			fail("%s", err)
		}

		if err = pkg.Build(ctx, macBuildConfig{
			DeploymentTarget: c.DeploymentTarget,
			Sign:             c.Sign,
			Sandbox:          c.Sandbox,
			Force:            c.Force,
			Race:             c.Race,
			Verbose:          c.Verbose,
		}); err != nil {
			fail("%s", err)
		}

		appname = pkg.name
	}

	go listenLogs(ctx, c.LogsAddr)
	time.Sleep(time.Millisecond * 200)

	os.Unsetenv("GOAPP_BUNDLE")
	os.Setenv("GOAPP_LOGS_ADDR", c.LogsAddr)
	os.Setenv("GOAPP_DEBUG", fmt.Sprintf("%v", c.Debug))

	printVerbose("running %s", appname)
	if err := execute(ctx, "open", "--wait-apps", appname); err != nil {
		printErr("%s", err)
	}
}

func listenLogs(ctx context.Context, addr string) {
	logs := logs.GoappServer{
		Addr:   addr,
		Writer: os.Stderr,
	}

	err := logs.ListenAndLog(ctx)
	if err != nil {
		printErr("listening logs failed: %s", err)
	}
}

func macOSVersion() string {
	return executeString("sw_vers", "-productVersion")
}

func win(ctx context.Context, args []string) {
	printErr("you are not on Windows!")
	os.Exit(-1)
}
