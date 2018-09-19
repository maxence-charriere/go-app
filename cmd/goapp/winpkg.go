// +build windows

package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/murlokswarm/app/internal/file"
)

type winPackage struct {
	workingDir     string
	buildDir       string
	buildResources string
	goPackageName  string
	goExec         string
	name           string
	resources      string
	assets         string
	config         winBuilConfig
	manifest       manifest
}

type manifest struct {
	Name        string
	Executable  string
	Description string
	Publisher   string
	Icon        string
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
		goPackageName:  goPackageName,
		goExec:         filepath.Join(tmpDir, goPackageName+".exe"),
		resources:      filepath.Join(name, "Resources"),
		assets:         filepath.Join(name, "Assets"),
		name:           name,
	}, nil
}

func (pkg *winPackage) Build(ctx context.Context, c winBuilConfig) error {
	pkg.config = c
	name := filepath.Base(pkg.name)

	printVerbose("building go executable")
	if err := pkg.buildGoExecutable(ctx); err != nil {
		return err
	}

	printVerbose("reading settings")
	if err := pkg.readSettings(ctx); err != nil {
		return err
	}

	printVerbose("creating %s", name)
	if err := pkg.createPackage(); err != nil {
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

func (pkg *winPackage) readSettings(ctx context.Context) error {
	settings := filepath.Join(pkg.workingDir, ".settings.win.json")
	defer os.Remove(settings)

	os.Setenv("GOAPP_BUILD", settings)
	defer os.Unsetenv("GOAPP_BUILD")

	if err := execute(ctx, pkg.goExec); err != nil {
		return err
	}

	data, err := ioutil.ReadFile(settings)
	if err != nil {
		return err
	}

	var m manifest
	if err = json.Unmarshal(data, &m); err != nil {
		return err
	}

	m.Name = stringWithDefault(m.Name, pkg.goPackageName)
	m.Executable = filepath.Base(pkg.goExec)
	m.Description = stringWithDefault(m.Description, m.Name)
	m.Publisher = stringWithDefault(m.Publisher, "goapp")

	d, _ := json.MarshalIndent(m, "", "    ")
	printVerbose("settings: %s", d)
	pkg.manifest = m
	return nil
}

func (pkg *winPackage) createPackage() error {
	dirs := []string{
		pkg.name,
		pkg.resources,
		pkg.assets,
	}

	for _, d := range dirs {
		if err := os.MkdirAll(d, os.ModeDir|0755); err != nil {
			return err
		}
	}

	if err := file.Copy(
		filepath.Join(pkg.name, filepath.Base(pkg.goExec)),
		pkg.goExec,
	); err != nil {
		return err
	}

	appxManifest := filepath.Join(pkg.name, "AppxManifest.xml")
	return generateTemplate(appxManifest, appxManifestTmpl, pkg.manifest)
}

func (pkg *winPackage) syncResources() error {
	return file.Sync(pkg.resources, pkg.buildResources)
}

func (pkg *winPackage) convertToAppx(ctx context.Context) error {
	// if err := os.RemoveAll(pkg.name); err != nil {
	// 	return err
	// }

	// if err := os.MkdirAll(pkg.name, 0755); err != nil {
	// 	return err
	// }

	// cmd := []string{"powershell", "DesktopAppConverter.exe",
	// 	"-Installer", pkg.tmpDir,
	// 	"-AppExecutable", filepath.Base(pkg.goExec),
	// 	"-Destination", pkg.workingDir,
	// 	"-PackageName", filepath.Base(pkg.name),
	// 	"-Publisher", "CN=goapp",
	// 	"-Version", "1.0.0.0",
	// 	"-MakeAppx",
	// 	"-Sign",
	// }

	// if verbose {
	// 	cmd = append(cmd, "-Verbose")
	// }

	// return execute(ctx, cmd[0], cmd[1:]...)
	return nil
}
