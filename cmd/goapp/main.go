package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/segmentio/conf"
)

const (
	dirPerm os.FileMode = 0777
)

var (
	verbose = false
)

func main() {
	ld := conf.Loader{
		Name: "goapp",
		Args: os.Args[1:],
		Commands: []conf.Command{
			{Name: "mac", Help: "Build app for MacOS."},
			{Name: "web", Help: "Build app for web."},
			{Name: "win", Help: "Build app for Windows."},
			{Name: "update", Help: "Update goapp to the latest version."},
			{Name: "help", Help: "Show the help."},
		},
	}

	ctx, cancel := ctxWithSignals(context.Background(), os.Interrupt)
	defer cancel()

	switch cmd, args := conf.LoadWith(nil, ld); cmd {
	case "mac":
		mac(ctx, args)

	case "web":
		web(ctx, args)

	case "win":
		win(ctx, args)

	case "update":
		update(ctx, args)

	case "help":
		ld.PrintHelp(nil)

	default:
		panic("unreachable")
	}
}

func ctxWithSignals(parent context.Context, s ...os.Signal) (ctx context.Context, cancel func()) {
	ctx, cancel = context.WithCancel(parent)
	sigc := make(chan os.Signal)
	signal.Notify(sigc, s...)

	go func() {
		defer close(sigc)

		<-sigc
		cancel()
	}()

	return ctx, cancel
}

type updateConfig struct {
	Verbose bool `conf:"v" help:"Enable verbose mode."`
}

func update(ctx context.Context, args []string) {
	c := updateConfig{}

	ld := conf.Loader{
		Name:    "goapp update",
		Args:    args,
		Usage:   "[options...]",
		Sources: []conf.Source{conf.NewEnvSource("GOAPP", os.Environ()...)},
	}

	conf.LoadWith(&c, ld)
	verbose = c.Verbose

	cmd := []string{"go", "get", "-u"}
	if verbose {
		cmd = append(cmd, "-v")
	}
	cmd = append(cmd, "github.com/murlokswarm/app/...")

	printVerbose("get https://github.com/murlokswarm/app latest version")
	if err := execute(ctx, cmd[0], cmd[1:]...); err != nil {
		fail("%s", err)
	}

	printSuccess("goapp successfully updated")
}

func packageRoots(packages []string) ([]string, error) {
	if len(packages) == 0 {
		packages = []string{"."}
	}

	roots := make([]string, len(packages))

	for i, p := range packages {
		dir, err := filepath.Abs(p)
		if err != nil {
			return nil, err
		}

		var info os.FileInfo
		if info, err = os.Stat(dir); err != nil {
			return nil, err
		}

		if !info.IsDir() {
			return nil, errors.Errorf("%s is not a directory", dir)
		}

		roots[i] = dir
	}

	return roots, nil
}

func initPackage(root string) error {
	printVerbose("set up resources")
	if err := os.MkdirAll(
		filepath.Join(root, "resources"),
		dirPerm,
	); err != nil && !os.IsExist(err) {
		return err
	}

	printVerbose("set up resources/css")
	if err := os.Mkdir(
		filepath.Join(root, "resources", "css"),
		dirPerm,
	); err != nil && !os.IsExist(err) {
		return err
	}

	return nil
}

func goBuild(ctx context.Context, buildDir string, args ...string) error {
	args = append([]string{"build"}, args...)

	if verbose {
		args = append(args, "-v")
	}

	args = append(args, buildDir)
	return execute(ctx, "go", args...)
}

var (
	greenColor   = "\033[92m"
	redColor     = "\033[91m"
	orangeColor  = "\033[93m"
	defaultColor = "\033[00m"
)

func printVerbose(format string, v ...interface{}) {
	if verbose {
		format = "‣ " + format
		fmt.Printf(format, v...)
		fmt.Println()
	}
}

func printSuccess(format string, v ...interface{}) {
	fmt.Print(greenColor)
	format = "✔ " + format
	fmt.Printf(format, v...)
	fmt.Println(defaultColor)
}

func printErr(format string, v ...interface{}) {
	fmt.Print(redColor)
	format = "x " + format
	fmt.Printf(format, v...)
	fmt.Println(defaultColor)
}

func printWarn(format string, v ...interface{}) {
	fmt.Print(orangeColor)
	format = "! " + format
	fmt.Printf(format, v...)
	fmt.Println(defaultColor)
}

func fail(format string, v ...interface{}) {
	printErr(format, v...)
	os.Exit(-1)
}

func failWithHelp(ld *conf.Loader, format string, v ...interface{}) {
	ld.PrintHelp(nil)
	ld.PrintError(errors.Errorf(format, v...))
	os.Exit(-1)
}

func verboseFlag(v bool) string {
	if v {
		return "-v"
	}

	return ""
}

func stringWithDefault(value, defaultValue string) string {
	if len(value) == 0 {
		return defaultValue
	}

	return value
}

func intWithDefault(value, defaultValue int) int {
	if value == 0 {
		return defaultValue
	}

	return value
}
