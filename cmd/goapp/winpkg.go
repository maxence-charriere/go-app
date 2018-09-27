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
	"github.com/pkg/errors"
)

type winPackage struct {
	workingDir     string
	buildDir       string
	buildResources string
	goPackageName  string
	goExec         string
	name           string
	namex          string
	resources      string
	assets         string
	config         winBuilConfig
	manifest       manifest
}

type manifest struct {
	Name        string
	ID          string
	Executable  string
	EntryPoint  string
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
	if strings.ToLower(goPackageName) == "uwp" {
		return nil, errors.New("go package cannot be name uwp")
	}

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
		namex:          strings.Replace(name, ".app", ".appx", 1),
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

	if !c.AppX {
		printVerbose("deploying")
		return pkg.deploy(ctx)
	}

	printVerbose("creating appx")
	if err := pkg.makeToAppx(ctx); err != nil {
		return err
	}

	printVerbose("signing")
	return pkg.sign(ctx)
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

	idFmt := func(n string) string {
		n = strings.Replace(n, "-", "", -1)
		n = strings.Replace(n, "_", "", -1)
		return n
	}

	m.Name = stringWithDefault(m.Name, pkg.goPackageName)
	m.ID = idFmt(m.Name)
	m.Executable = filepath.Base(pkg.goExec)
	m.EntryPoint = strings.Replace(m.Executable, ".exe", ".app", 1)
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

	uwpDir := filepath.Join(murlokswarm(), "cmd", "goapp", "uwp")
	uwpFiles := []string{
		"clrcompression.dll",
		"resources.pri",
		"uwp.dll",
		"uwp.exe",
	}

	for _, f := range uwpFiles {
		src := filepath.Join(uwpDir, f)
		dst := filepath.Join(pkg.name, f)

		if err := file.Copy(dst, src); err != nil {
			return err
		}
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
		if s <= 1 {
			return filepath.Join(pkg.assets, fmt.Sprintf("%s.png", n))
		}

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
		{Name: scaled("Square44x44Logo", 1), Width: 44, Height: 44, Scale: 1, Padding: true},
		{Name: scaled("Square44x44Logo", 100), Width: 44, Height: 44, Scale: 1},
		{Name: scaled("Square44x44Logo", 125), Width: 44, Height: 44, Scale: 1.25},
		{Name: scaled("Square44x44Logo", 150), Width: 44, Height: 44, Scale: 1.5},
		{Name: scaled("Square44x44Logo", 200), Width: 44, Height: 44, Scale: 2},
		{Name: scaled("Square44x44Logo", 400), Width: 44, Height: 44, Scale: 4},

		{Name: targetSized("Square44x44Logo", 16, false), Width: 16, Height: 16, Scale: 1},
		{Name: targetSized("Square44x44Logo", 16, true), Width: 16, Height: 16, Scale: 1},
		{Name: targetSized("Square44x44Logo", 24, false), Width: 24, Height: 24, Scale: 1},
		{Name: targetSized("Square44x44Logo", 24, true), Width: 24, Height: 24, Scale: 1},
		{Name: targetSized("Square44x44Logo", 32, false), Width: 32, Height: 32, Scale: 1},
		{Name: targetSized("Square44x44Logo", 32, true), Width: 32, Height: 32, Scale: 1},
		{Name: targetSized("Square44x44Logo", 48, false), Width: 48, Height: 48, Scale: 1},
		{Name: targetSized("Square44x44Logo", 48, true), Width: 48, Height: 48, Scale: 1},
		{Name: targetSized("Square44x44Logo", 256, false), Width: 256, Height: 256, Scale: 1},
		{Name: targetSized("Square44x44Logo", 256, true), Width: 256, Height: 256, Scale: 1},

		{Name: scaled("Square71x71Logo", 1), Width: 71, Height: 71, Scale: 1, Padding: true},
		{Name: scaled("Square71x71Logo", 100), Width: 71, Height: 71, Scale: 1},
		{Name: scaled("Square71x71Logo", 125), Width: 71, Height: 71, Scale: 1.25},
		{Name: scaled("Square71x71Logo", 150), Width: 71, Height: 71, Scale: 1.5},
		{Name: scaled("Square71x71Logo", 200), Width: 71, Height: 71, Scale: 2},
		{Name: scaled("Square71x71Logo", 400), Width: 71, Height: 71, Scale: 4},

		{Name: scaled("Square150x150Logo", 1), Width: 150, Height: 150, Scale: 1, Padding: true},
		{Name: scaled("Square150x150Logo", 100), Width: 150, Height: 150, Scale: 1},
		{Name: scaled("Square150x150Logo", 125), Width: 150, Height: 150, Scale: 1.25},
		{Name: scaled("Square150x150Logo", 150), Width: 150, Height: 150, Scale: 1.5},
		{Name: scaled("Square150x150Logo", 200), Width: 150, Height: 150, Scale: 2},
		{Name: scaled("Square150x150Logo", 400), Width: 150, Height: 150, Scale: 4},

		{Name: scaled("Square310x310Logo", 1), Width: 310, Height: 310, Scale: 1, Padding: true},
		{Name: scaled("Square310x310Logo", 100), Width: 310, Height: 310, Scale: 1},
		{Name: scaled("Square310x310Logo", 125), Width: 310, Height: 310, Scale: 1.25},
		{Name: scaled("Square310x310Logo", 150), Width: 310, Height: 310, Scale: 1.5},
		{Name: scaled("Square310x310Logo", 200), Width: 310, Height: 310, Scale: 2},
		{Name: scaled("Square310x310Logo", 400), Width: 310, Height: 310, Scale: 4},

		{Name: scaled("StoreLogo", 1), Width: 50, Height: 50, Scale: 1, Padding: true},
		{Name: scaled("StoreLogo", 100), Width: 50, Height: 50, Scale: 1},
		{Name: scaled("StoreLogo", 125), Width: 50, Height: 50, Scale: 1.25},
		{Name: scaled("StoreLogo", 150), Width: 50, Height: 50, Scale: 1.5},
		{Name: scaled("StoreLogo", 200), Width: 50, Height: 50, Scale: 2},
		{Name: scaled("StoreLogo", 400), Width: 50, Height: 50, Scale: 4},

		{Name: scaled("Wide310x150Logo", 1), Width: 310, Height: 150, Scale: 1, Padding: true},
		{Name: scaled("Wide310x150Logo", 100), Width: 310, Height: 150, Scale: 1},
		{Name: scaled("Wide310x150Logo", 125), Width: 310, Height: 150, Scale: 1.25},
		{Name: scaled("Wide310x150Logo", 150), Width: 310, Height: 150, Scale: 1.5},
		{Name: scaled("Wide310x150Logo", 200), Width: 310, Height: 150, Scale: 2},
		{Name: scaled("Wide310x150Logo", 400), Width: 310, Height: 150, Scale: 4},
	})
}

func (pkg *winPackage) deploy(ctx context.Context) error {
	cmd := []string{"powershell",
		"Add-AppxPackage",
		"-Path", filepath.Join(pkg.name, "AppxManifest.xml"),
		"-Register",
	}

	return execute(ctx, cmd[0], cmd[1:]...)
}

func (pkg *winPackage) makeToAppx(ctx context.Context) error {
	os.Chdir(win10SDKBinX64())
	defer os.Chdir(pkg.workingDir)

	cmd := []string{
		"makeappx.exe", "pack",
		"-d", pkg.name,
		"-p", pkg.namex,
		"-o",
	}

	if verbose {
		cmd = append(cmd, "-v")
	}

	return execute(ctx, cmd[0], cmd[1:]...)
}

func (pkg *winPackage) sign(ctx context.Context) error {
	os.Chdir(win10SDKBinX64())
	defer os.Chdir(pkg.workingDir)

	cmd := []string{
		"signtool.exe", "sign", "/a",
		"/fd", "SHA256",
		"/f", filepath.Join(
			murlokswarm(),
			"cmd",
			"goapp",
			"certificates",
			"win.pfx",
		),
		"/p", "123456",
	}

	if verbose {
		cmd = append(cmd, "/v")
	}

	cmd = append(cmd, pkg.namex)
	return execute(ctx, cmd[0], cmd[1:]...)
}

func win10SDKBinX64() string {
	return filepath.Join(
		`C:\`,
		"Program Files (x86)",
		"Windows Kits",
		"10",
		"bin",
		"10.0.17134.0",
		"x64",
	)
}

func win10SDKBinX86() string {
	return filepath.Join(
		`C:\`,
		"Program Files (x86)",
		"Windows Kits",
		"10",
		"bin",
		"10.0.17134.0",
		"x86",
	)
}
