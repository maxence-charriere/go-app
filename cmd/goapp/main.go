package main

import (
	"context"
	"fmt"
	"os"
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
			{Name: "help", Help: "Show the help."},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	switch cmd, args := conf.LoadWith(nil, ld); cmd {
	case "mac":
		mac(ctx, args)

	case "web":
		web(ctx, args)

	case "win":
		win(ctx, args)

	case "help":
		ld.PrintHelp(nil)

	default:
		panic("unreachable")
	}
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

func goBuild(ctx context.Context, target string, args ...string) error {
	args = append([]string{"build"}, args...)

	if verbose {
		args = append(args, "-v", target)
	}

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
