// +build darwin,amd64

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/murlokswarm/app/internal/file"
	"github.com/pkg/errors"
)

type macPackage struct {
	workingDir     string
	buildDir       string
	buildResources string
	goExec         string
	name           string
	contents       string
	macOS          string
	resources      string
	config         macBuildConfig
	bundle         bundle
}

func newMacPackage(buildDir, name string) (*macPackage, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	if buildDir, err = filepath.Abs(buildDir); err != nil {
		return nil, err
	}

	if len(name) == 0 {
		name = filepath.Base(buildDir) + ".app"
		name = filepath.Join(wd, name)
	}

	if !strings.HasSuffix(name, ".app") {
		name += ".app"
	}

	return &macPackage{
		workingDir:     wd,
		buildDir:       buildDir,
		buildResources: filepath.Join(buildDir, "resources"),
		goExec:         filepath.Join(wd, filepath.Base(buildDir)),
		name:           name,
		contents:       filepath.Join(name, "Contents"),
		macOS:          filepath.Join(name, "Contents", "MacOS"),
		resources:      filepath.Join(name, "Contents", "Resources"),
	}, nil
}

func (pkg *macPackage) Build(ctx context.Context, c macBuildConfig) error {
	pkg.config = c

	printVerbose("building go executable")
	if err := pkg.buildGoExecutable(ctx); err != nil {
		return err
	}

	printVerbose("reading bundle info")
	if err := pkg.readBundleInfo(ctx); err != nil {
		return err
	}

	name := filepath.Base(pkg.name)

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

	if len(c.Sign) == 0 {
		if c.AppStore {
			return errors.New("app store apps require to be signed")
		}

		return nil
	}

	printVerbose("signing package")
	if err := pkg.signPackage(ctx); err != nil {
		return err
	}

	if !c.Sandbox {
		return errors.New("app store apps require sandbox mode")
	}

	if c.AppStore {
		printVerbose("packaging for app store")
		return pkg.packForAppStore(ctx)
	}

	return nil
}

func (pkg *macPackage) buildGoExecutable(ctx context.Context) error {
	os.Setenv("MACOSX_DEPLOYMENT_TARGET", pkg.config.DeploymentTarget)

	cmd := []string{"go", "build",
		"-ldflags", "-s",
		"-o", pkg.goExec,
	}

	if verbose {
		cmd = append(cmd, "-v")
	}

	if pkg.config.Force {
		cmd = append(cmd, "-a")
	}

	if pkg.config.Race {
		cmd = append(cmd, "-race")
	}

	return execute(ctx, cmd[0], cmd[1:]...)
}

func (pkg *macPackage) readBundleInfo(ctx context.Context) error {
	bundleJSON := filepath.Join(pkg.workingDir, ".bundle.json")
	os.Setenv("GOAPP_BUNDLE", bundleJSON)
	defer os.Remove(bundleJSON)
	defer os.Unsetenv("GOAPP_BUNDLE")

	if err := execute(ctx, pkg.goExec); err != nil {
		return err
	}

	data, err := ioutil.ReadFile(bundleJSON)
	if err != nil {
		return err
	}

	var b bundle
	if err = json.Unmarshal(data, &b); err != nil {
		return err
	}

	b.ExecName = filepath.Base(pkg.goExec)
	b.Sandbox = pkg.config.Sandbox
	b.AppName = stringWithDefault(b.AppName, filepath.Base(pkg.goExec))
	b.ID = stringWithDefault(b.ID, fmt.Sprintf("%v.%v", os.Getenv("USER"), b.AppName))
	b.URLScheme = stringWithDefault(b.URLScheme, strings.ToLower(b.AppName))
	b.Version = stringWithDefault(b.Version, "1.0")
	b.BuildNumber = intWithDefault(b.BuildNumber, 1)
	b.DevRegion = stringWithDefault(b.DevRegion, "en")
	b.DeploymentTarget = pkg.config.DeploymentTarget
	b.Category = stringWithDefault(b.Category, "public.app-category.developer-tools")
	b.Copyright = stringWithDefault(b.Copyright, fmt.Sprintf("Copyright © %v %s. All rights reserved.",
		time.Now().Year(),
		os.Getenv("USER"),
	))
	b.Role = stringWithDefault(b.Role, "None")

	if b.Sandbox && len(pkg.config.Sign) == 0 {
		printWarn("desactivating sandbox: sanboxed apps require to be signed")
		b.Sandbox = false
	}

	d, _ := json.MarshalIndent(b, "", "    ")
	printVerbose("bundle: %s", d)

	pkg.bundle = b
	return nil
}

func (pkg *macPackage) createPackage() error {
	os.RemoveAll(filepath.Join(pkg.contents, "_CodeSignature"))

	dirs := []string{
		pkg.name,
		pkg.contents,
		pkg.macOS,
		pkg.resources,
	}

	for _, d := range dirs {
		if err := os.MkdirAll(d, os.ModeDir|0755); err != nil {
			return err
		}
	}

	if err := file.Copy(
		filepath.Join(pkg.macOS, pkg.bundle.ExecName),
		pkg.goExec,
	); err != nil {
		return err
	}

	infoPlist := filepath.Join(pkg.contents, "Info.plist")
	return generatePlist(infoPlist, plist, pkg.bundle)
}

func (pkg *macPackage) syncResources() error {
	return file.Sync(pkg.resources, pkg.buildResources)
}

func (pkg *macPackage) generateIcons(ctx context.Context) error {
	if len(pkg.bundle.Icon) == 0 {
		return nil
	}

	icon := filepath.Join(pkg.buildResources, pkg.bundle.Icon)

	iconset := filepath.Base(icon)
	iconset = strings.TrimSuffix(iconset, filepath.Ext(iconset))
	iconset = filepath.Join(pkg.resources, iconset) + ".iconset"

	if err := os.Mkdir(iconset, os.ModeDir|0755); err != nil {
		return err
	}
	defer os.RemoveAll(iconset)

	retinaIcon := func(w, h, s int) string {
		return filepath.Join(iconset, fmt.Sprintf("icon_%vx%v@%vx.png", w, h, s))
	}

	standardIcon := func(w, h int) string {
		return filepath.Join(iconset, fmt.Sprintf("icon_%vx%v.png", w, h))
	}

	if err := generateIcons(icon, []iconInfo{
		{Name: retinaIcon(512, 512, 2), Width: 512, Height: 512, Scale: 2},
		{Name: standardIcon(512, 512), Width: 512, Height: 512, Scale: 1},

		{Name: retinaIcon(256, 256, 2), Width: 256, Height: 256, Scale: 2},
		{Name: standardIcon(256, 256), Width: 256, Height: 256, Scale: 1},

		{Name: retinaIcon(128, 128, 2), Width: 128, Height: 128, Scale: 2},
		{Name: standardIcon(128, 128), Width: 128, Height: 128, Scale: 1},

		{Name: retinaIcon(32, 32, 2), Width: 32, Height: 32, Scale: 2},
		{Name: standardIcon(32, 32), Width: 32, Height: 32, Scale: 1},

		{Name: retinaIcon(16, 16, 2), Width: 16, Height: 16, Scale: 2},
		{Name: standardIcon(16, 16), Width: 16, Height: 16, Scale: 1},
	}); err != nil {
		return err
	}

	return execute(ctx, "iconutil", "-c", "icns", iconset)
}

func (pkg *macPackage) signPackage(ctx context.Context) error {
	ents := filepath.Join(pkg.workingDir, ".entitlements")

	if err := generatePlist(ents, entitlements, pkg.bundle); err != nil {
		return err
	}
	defer os.Remove(ents)

	signEntsCmd := []string{
		"codesign",
		"--force",
		"--sign",
		pkg.config.Sign,
		"--entitlements",
		ents,
		pkg.name,
	}

	if verbose {
		signEntsCmd = append(signEntsCmd, "-v")
	}

	if err := execute(ctx, signEntsCmd[0], signEntsCmd[1:]...); err != nil {
		return err
	}

	cmd := []string{
		"codesign",
		"--verify",
		"--deep",
		"--strict",
		pkg.name,
	}

	if verbose {
		cmd = append(cmd, "--verbose=2")
	}

	return execute(ctx, cmd[0], cmd[1:]...)
}

func (pkg *macPackage) packForAppStore(ctx context.Context) error {
	name := filepath.Base(pkg.name)
	name = strings.TrimSuffix(name, ".app")

	return execute(ctx,
		"productbuild",
		"--component",
		pkg.name,
		"/Applications",
		"--sign",
		pkg.config.Sign,
		name+".pkg",
	)
}
