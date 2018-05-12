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

func main() {
	ld := conf.Loader{
		Name:     "goapp",
		Args:     os.Args[1:],
		Commands: commands(),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	switch cmd, args := conf.LoadWith(nil, ld); cmd {
	case "web":
		web(ctx, args)

	case "mac":
		mac(ctx, args)

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
	if err := os.MkdirAll(
		filepath.Join(root, "resources"),
		dirPerm,
	); err != nil && !os.IsExist(err) {
		return err
	}

	if err := os.Mkdir(
		filepath.Join(root, "resources", "css"),
		dirPerm,
	); err != nil && !os.IsExist(err) {
		return err
	}

	return nil
}

func goBuild(target string, args ...string) error {
	args = append([]string{"build"}, args...)
	args = append(args, "-v", target)
	return execute("go", args...)
}

func printSuccess(format string, v ...interface{}) {
	fmt.Print("\033[92m")
	fmt.Printf(format, v...)
	fmt.Println("\033[00m")
}

func printErr(format string, v ...interface{}) {
	fmt.Print("\033[91m")
	fmt.Printf(format, v...)
	fmt.Println("\033[00m")
}
