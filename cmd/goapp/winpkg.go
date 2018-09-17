// +build windows

package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/murlokswarm/app/internal/file"
)

type winPackage struct {
	workingDir     string
	buildDir       string
	buildResources string
	tmpDir         string
	tmpResources   string
	goPackageName  string
	goExec         string
	name           string
	config         winBuilConfig
}

func newWinPackage(buildDir, name string) (*winPackage, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	if buildDir, err = filepath.Abs(buildDir); err != nil {
		return nil, err
	}

	goPackageName := filepath.Base(buildDir)
	tmpDir := filepath.Join(os.Getenv("TEMP"), "goapp", goPackageName)

	if len(name) == 0 {
		name = goPackageName + ".app"
		name = filepath.Join(wd, name)
	}

	if !strings.HasSuffix(name, ".app") {
		name += ".app"
	}

	return &winPackage{
		workingDir:     wd,
		buildDir:       buildDir,
		buildResources: filepath.Join(buildDir, "resources"),
		tmpDir:         tmpDir,
		tmpResources:   filepath.Join(tmpDir, "Assets"),
		goPackageName:  goPackageName,
		goExec:         filepath.Join(tmpDir, goPackageName+".exe"),
		name:           name,
	}, nil
}

func (pkg *winPackage) Build(ctx context.Context, c winBuilConfig) error {
	pkg.config = c

	printVerbose("building go executable")
	if err := pkg.buildGoExecutable(ctx); err != nil {
		return err
	}

	printVerbose("syncing resources")
	if err := pkg.syncResources(); err != nil {
		return err
	}

	printVerbose("converting to appx")
	return pkg.convertToAppx(ctx)
}

func (pkg *winPackage) buildGoExecutable(ctx context.Context) error {
	return goBuild(ctx, pkg.buildDir, "-ldflags", "-s", "-o", pkg.goExec)
}

func (pkg *winPackage) syncResources() error {
	return file.Sync(pkg.tmpResources, pkg.buildResources)
}

func (pkg *winPackage) convertToAppx(ctx context.Context) error {
	if err := os.RemoveAll(pkg.name); err != nil {
		return err
	}

	if err := os.MkdirAll(pkg.name, 0755); err != nil {
		return err
	}

	cmd := []string{"powershell", "DesktopAppConverter.exe",
		"-Installer", pkg.tmpDir,
		"-AppExecutable", filepath.Base(pkg.goExec),
		"-Destination", pkg.workingDir,
		"-PackageName", filepath.Base(pkg.name),
		"-Publisher", "CN=goapp",
		"-Version", "1.0.0.0",
		"-MakeAppx",
		"-Sign",
	}

	if verbose {
		cmd = append(cmd, "-Verbose")
	}

	return execute(ctx, cmd[0], cmd[1:]...)
}
