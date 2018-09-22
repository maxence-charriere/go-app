// +build windows

package main

import (
	"context"
	"encoding/json"
	"fmt"
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

	printVerbose("generating icons")
	if err := pkg.generateIcons(ctx); err != nil {
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
	m.Icon = stringWithDefault(m.Icon, filepath.Join(murlokswarm(), "logo.png"))

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

func (pkg *winPackage) generateIcons(ctx context.Context) error {
	icon := pkg.manifest.Icon

	scaled := func(n string, s int) string {
		return filepath.Join(
			pkg.assets,
			fmt.Sprintf("%s.scale-%v.png", n, s),
		)
	}

	targetSized := func(n string, s int, alt bool) string {
		altstr := ""
		if alt {
			altstr = "_altform-unplated"
		}

		return filepath.Join(
			pkg.assets,
			fmt.Sprintf("%s.targetsize-%v%s.png", n, s, altstr),
		)
	}

	return generateIcons(icon, []iconInfo{
		{Name: scaled("Square44x44Logo", 100), Width: 44, Height: 44, Scale: 1, Padding: true},
		{Name: scaled("Square44x44Logo", 125), Width: 44, Height: 44, Scale: 1.25, Padding: true},
		{Name: scaled("Square44x44Logo", 150), Width: 44, Height: 44, Scale: 1.5, Padding: true},
		{Name: scaled("Square44x44Logo", 200), Width: 44, Height: 44, Scale: 2, Padding: true},
		{Name: scaled("Square44x44Logo", 400), Width: 44, Height: 44, Scale: 4, Padding: true},

		{Name: targetSized("Square44x44Logo", 16, false), Width: 16, Height: 16, Scale: 1, Padding: true},
		{Name: targetSized("Square44x44Logo", 16, true), Width: 16, Height: 16, Scale: 1, Padding: true},
		{Name: targetSized("Square44x44Logo", 24, false), Width: 24, Height: 24, Scale: 1, Padding: true},
		{Name: targetSized("Square44x44Logo", 24, true), Width: 24, Height: 24, Scale: 1, Padding: true},
		{Name: targetSized("Square44x44Logo", 32, false), Width: 32, Height: 32, Scale: 1, Padding: true},
		{Name: targetSized("Square44x44Logo", 32, true), Width: 32, Height: 32, Scale: 1, Padding: true},
		{Name: targetSized("Square44x44Logo", 48, false), Width: 48, Height: 48, Scale: 1, Padding: true},
		{Name: targetSized("Square44x44Logo", 48, true), Width: 48, Height: 48, Scale: 1, Padding: true},
		{Name: targetSized("Square44x44Logo", 256, false), Width: 256, Height: 256, Scale: 1, Padding: true},
		{Name: targetSized("Square44x44Logo", 256, true), Width: 256, Height: 256, Scale: 1, Padding: true},

		{Name: scaled("Square71x71Logo", 100), Width: 71, Height: 71, Scale: 1, Padding: true},
		{Name: scaled("Square71x71Logo", 125), Width: 71, Height: 71, Scale: 1.25, Padding: true},
		{Name: scaled("Square71x71Logo", 150), Width: 71, Height: 71, Scale: 1.5, Padding: true},
		{Name: scaled("Square71x71Logo", 200), Width: 71, Height: 71, Scale: 2, Padding: true},
		{Name: scaled("Square71x71Logo", 400), Width: 71, Height: 71, Scale: 4, Padding: true},

		{Name: scaled("Square150x150Logo", 100), Width: 150, Height: 150, Scale: 1, Padding: true},
		{Name: scaled("Square150x150Logo", 125), Width: 150, Height: 150, Scale: 1.25, Padding: true},
		{Name: scaled("Square150x150Logo", 150), Width: 150, Height: 150, Scale: 1.5, Padding: true},
		{Name: scaled("Square150x150Logo", 200), Width: 150, Height: 150, Scale: 2, Padding: true},
		{Name: scaled("Square150x150Logo", 400), Width: 150, Height: 150, Scale: 4, Padding: true},

		{Name: scaled("Square310x310Logo", 100), Width: 310, Height: 310, Scale: 1, Padding: true},
		{Name: scaled("Square310x310Logo", 125), Width: 310, Height: 310, Scale: 1.25, Padding: true},
		{Name: scaled("Square310x310Logo", 150), Width: 310, Height: 310, Scale: 1.5, Padding: true},
		{Name: scaled("Square310x310Logo", 200), Width: 310, Height: 310, Scale: 2, Padding: true},
		{Name: scaled("Square310x310Logo", 400), Width: 310, Height: 310, Scale: 4, Padding: true},

		{Name: scaled("StoreLogo", 100), Width: 50, Height: 50, Scale: 1, Padding: true},
		{Name: scaled("StoreLogo", 125), Width: 50, Height: 50, Scale: 1.25, Padding: true},
		{Name: scaled("StoreLogo", 150), Width: 50, Height: 50, Scale: 1.5, Padding: true},
		{Name: scaled("StoreLogo", 200), Width: 50, Height: 50, Scale: 2, Padding: true},
		{Name: scaled("StoreLogo", 400), Width: 50, Height: 50, Scale: 4, Padding: true},

		{Name: scaled("Wide310x150Logo", 100), Width: 310, Height: 150, Scale: 1, Padding: true},
		{Name: scaled("Wide310x150Logo", 125), Width: 310, Height: 150, Scale: 1.25, Padding: true},
		{Name: scaled("Wide310x150Logo", 150), Width: 310, Height: 150, Scale: 1.5, Padding: true},
		{Name: scaled("Wide310x150Logo", 200), Width: 310, Height: 150, Scale: 2, Padding: true},
		{Name: scaled("Wide310x150Logo", 400), Width: 310, Height: 150, Scale: 4, Padding: true},
	})
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
