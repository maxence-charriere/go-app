package main

import (
	"context"
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
		Name: "goapp",
		Args: os.Args[1:],
		Commands: []conf.Command{
			{Name: "web", Help: "Build app on web."},
			{Name: "help", Help: "Show the help."},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	switch cmd, args := conf.LoadWith(nil, ld); cmd {
	case "web":
		web(ctx, args)

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

func goBuild(target string, verbose bool) error {
	args := []string{"build"}
	if verbose {
		args = append(args, "-v")
	}
	args = append(args, target)
	return execute("go", args...)
}
