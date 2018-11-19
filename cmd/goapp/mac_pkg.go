package main

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/murlokswarm/app/internal/file"

	"github.com/pkg/errors"
)

// MacPackage represents a package for a MacOS app.
// It implements the Package interface.
type MacPackage struct {
	// The path where the sources are.
	// It must refer to a Go main package.
	// Default is ".".
	Sources string

	// The path where the package is saved.
	// If not set, the ".app" extension is added.
	Output string

	// Enable verbose mode.
	Verbose bool

	// Force rebuilding of packages that are already up-to-date.
	Force bool

	// Enable data race detection.
	Race bool

	// The version on MacOS the build is for.
	DeploymentTarget string

	// The function to log events.
	Log func(string, ...interface{}) (int, error)

	name          string
	workingDir    string
	tmpDir        string
	tmpExecutable string
	contentsDir   string
	macOSDir      string
	resourcesDir  string
	executable    string
	settings      macSettings
}

// Build satisfies the Package interface.
func (pkg *MacPackage) Build(ctx context.Context) error {
	pkg.init()
	panic("not implemented")
}

func (pkg *MacPackage) init() (err error) {
	if len(pkg.Sources) == 0 || pkg.Sources == "." || pkg.Sources == "./" {
		pkg.Sources = "."
	}
	if pkg.Sources, err = filepath.Abs(pkg.Sources); err != nil {
		return err
	}

	execName := filepath.Base(pkg.Sources)

	if len(pkg.Output) == 0 {
		pkg.Output = execName
	}
	if !strings.HasSuffix(pkg.Output, ".app") {
		pkg.Output += ".app"
	}

	pkg.name = filepath.Base(pkg.Output)

	if pkg.workingDir, err = os.Getwd(); err != nil {
		return err
	}

	if pkg.tmpDir = os.Getenv("TMPDIR"); len(pkg.tmpDir) == 0 {
		return errors.New("tmp dir not set")
	}
	pkg.tmpDir = filepath.Join(pkg.tmpDir, "goapp")
	pkg.tmpExecutable = filepath.Join(pkg.tmpDir, execName)

	pkg.contentsDir = filepath.Join(pkg.Output, "Contents")
	pkg.macOSDir = filepath.Join(pkg.Output, "Contents", "MacOS")
	pkg.resourcesDir = filepath.Join(pkg.Output, "Contents", "Resources")
	pkg.executable = filepath.Join(pkg.Output, "Contents", "MacOS", execName)
	return nil
}

// Run satisfies the Package interface.
func (pkg *MacPackage) Run(ctx context.Context) error {
	if runtime.GOOS != "darwin" {
		return errors.New("operating system is not MacOS")
	}

	if err := pkg.init(); err != nil {
		return err
	}

	pkg.Log("creating %s", pkg.name)
	if err := pkg.create(); err != nil {
		return err
	}

	pkg.Log("build executable")
	if err := pkg.buildExecutable(ctx); err != nil {
		return err
	}

	return nil
}

func (pkg *MacPackage) create() error {
	if err := os.RemoveAll(filepath.Join(pkg.contentsDir, "_CodeSignature")); err != nil {
		return err
	}

	dirs := []string{
		pkg.Output,
		pkg.contentsDir,
		pkg.macOSDir,
		pkg.resourcesDir,
	}

	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return err
		}
	}

	return nil
}

func (pkg *MacPackage) buildExecutable(ctx context.Context) error {
	os.Setenv("MACOSX_DEPLOYMENT_TARGET", pkg.DeploymentTarget)

	args := []string{
		"go", "build",
		"-ldflags", "-s",
		"-o", pkg.tmpExecutable,
	}

	if pkg.Verbose {
		args = append(args, "-v")
	}

	if pkg.Force {
		args = append(args, "-a")
	}

	if pkg.Race {
		args = append(args, "-race")
	}

	if err := execute(ctx, args[0], args[1:]...); err != nil {
		return err
	}

	return file.Copy(pkg.executable, pkg.tmpExecutable)
}

// func (pkg *MacPackage)

type macSettings struct{}
