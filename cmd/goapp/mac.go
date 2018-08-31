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

	"github.com/murlokswarm/app/internal/logs"

	driver "github.com/murlokswarm/app/drivers/mac"
	"github.com/murlokswarm/app/internal/file"
	"github.com/segmentio/conf"
)

func mac(ctx context.Context, args []string) {
	ld := conf.Loader{
		Name: "goapp mac",
		Args: args,
		Commands: []conf.Command{
			{Name: "init", Help: "Download MacOS SDK and create required file and directories."},
			{Name: "build", Help: "Build the MacOS app."},
			{Name: "run", Help: "Run a MacOS app and capture its logs."},
			{Name: "help", Help: "Show the MacOS help"},
		},
	}

	switch cmd, args := conf.LoadWith(nil, ld); cmd {
	case "init":
		initMac(ctx, args)

	case "build":
		buildMac(ctx, args)

	case "run":
		runMac(ctx, args)

	case "help":
		ld.PrintHelp(nil)

	default:
		panic("unreachable")
	}
}

type macInitConfig struct {
	Verbose bool `conf:"v" help:"Enable verbose mode."`
}

func initMac(ctx context.Context, args []string) {
	c := macInitConfig{}

	ld := conf.Loader{
		Name:    "mac init",
		Args:    args,
		Usage:   "[options...] [packages...]",
		Sources: []conf.Source{conf.NewEnvSource("GOAPP", os.Environ()...)},
	}

	_, unusedArgs := conf.LoadWith(&c, ld)
	verbose = c.Verbose

	roots, err := packageRoots(unusedArgs)
	if err != nil {
		failWithHelp(&ld, "%s", err)
	}

	printVerbose("checking for xcode-select...")
	execute(ctx, "xcode-select", "--install")

	for _, root := range roots {
		if err = initPackage(root); err != nil {
			fail("init %s failed: %s", root, err)
		}
	}

	printSuccess("init succeeded")
}

type macBuildConfig struct {
	Sign     string `conf:"sign"     help:"The signing identifier to sign the app.\n\t\033[95msecurity find-identity -v -p codesigning\033[00m to see signing identifiers.\n\thttps://developer.apple.com/library/content/documentation/Security/Conceptual/CodeSigningGuide/Procedures/Procedures.html to create one."`
	AppStore bool   `conf:"appstore" help:"Report whether the app will be packaged to be uploaded on the app store."`
	Verbose  bool   `conf:"v"        help:"Enable verbose mode."`
}

func buildMac(ctx context.Context, args []string) {
	c := macBuildConfig{}

	ld := conf.Loader{
		Name:    "mac build",
		Args:    args,
		Usage:   "[options...] [package]",
		Sources: []conf.Source{conf.NewEnvSource("GOAPP", os.Environ()...)},
	}

	_, roots := conf.LoadWith(&c, ld)
	verbose = c.Verbose

	if len(roots) == 0 {
		roots = []string{"."}
	}

	root := roots[0]

	if _, err := buildMacApp(ctx, root, c); err != nil {
		fail("%s", err)
	}

	printSuccess("build succeeded")
}

func buildMacApp(ctx context.Context, root string, c macBuildConfig) (string, error) {
	printVerbose("building go exec")
	if err := goBuild(ctx, root, "-ldflags", "-s"); err != nil {
		return "", err
	}

	printVerbose("building .app")
	return bundleMacApp(ctx, root, c)
}

func bundleMacApp(ctx context.Context, root string, c macBuildConfig) (string, error) {
	err := os.Setenv("GOAPP_BUNDLE", "true")
	if err != nil {
		return "", err
	}

	if root, err = filepath.Abs(root); err != nil {
		return "", err
	}

	var wd string
	if wd, err = os.Getwd(); err != nil {
		return "", err
	}

	if err = execute(ctx, filepath.Join(
		wd,
		filepath.Base(root),
	)); err != nil {
		return "", err
	}

	var data []byte
	if data, err = ioutil.ReadFile("goapp-mac.json"); err != nil {
		return "", err
	}
	defer os.Remove("goapp-mac.json")

	var bundle driver.Bundle
	if err = json.Unmarshal(data, &bundle); err != nil {
		return "", err
	}

	if bundle.Sandbox && len(c.Sign) == 0 {
		printWarn("sanboxed app require to be signed")
		printWarn("sandbox set to false")
		bundle.Sandbox = false
	}

	if c.AppStore && !bundle.Sandbox {
		fail("app store require app to run in sandbox mode")
	}

	bundle = fillBundle(bundle, root)
	data, _ = json.MarshalIndent(bundle, "", "    ")
	printVerbose("bundle configuration %s", data)

	appName := bundle.AppName + ".app"
	if err = createAppBundle(ctx, bundle, root, appName); err != nil {
		os.RemoveAll(appName)
		return "", err
	}

	if len(c.Sign) != 0 {
		printVerbose("signing app")
		if err = signAppBundle(ctx, bundle, c.Sign, root, appName); err != nil {
			os.RemoveAll(appName)
			return "", err
		}

		if c.AppStore {
			printVerbose("packaging for the app store")
			return appName, createAppPkg(ctx, bundle, c.Sign, appName)
		}
	}

	return appName, nil
}

func fillBundle(b driver.Bundle, root string) driver.Bundle {
	if len(b.AppName) == 0 {
		b.AppName = filepath.Base(root)
	}

	if len(b.ID) == 0 {
		b.ID = fmt.Sprintf("%v.%v", os.Getenv("USER"), b.AppName)
	}

	if len(b.URLScheme) == 0 {
		b.URLScheme = strings.ToLower(b.AppName)
	}

	if len(b.Version) == 0 {
		b.Version = "1.0"
	}

	if b.BuildNumber == 0 {
		b.BuildNumber = 1
	}

	if len(b.DevRegion) == 0 {
		b.DevRegion = "en"
	}

	if len(b.DeploymentTarget) == 0 {
		b.DeploymentTarget = "10.13"
	}

	if len(b.Category) == 0 {
		b.Category = driver.DeveloperToolsApp
	}

	if len(b.Copyright) == 0 {
		b.Copyright = fmt.Sprintf("Copyright © %v %s. All rights reserved.",
			time.Now().Year(),
			os.Getenv("USER"),
		)
	}

	if len(b.Role) == 0 {
		b.Role = driver.NoRole
	}

	return b
}

func createAppBundle(ctx context.Context, bundle driver.Bundle, root, appName string) error {
	os.RemoveAll(appName)

	appContents := filepath.Join(appName, "Contents")
	appMacOS := filepath.Join(appName, "Contents", "MacOS")
	appResources := filepath.Join(appName, "Contents", "Resources")
	resources := filepath.Join(root, "resources")

	dirs := []string{
		appContents,
		appMacOS,
		appResources,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, os.ModeDir|0755); err != nil {
			return err
		}
	}

	if err := generatePlist(filepath.Join(appContents, "Info.plist"), plist, bundle); err != nil {
		return err
	}

	execName := filepath.Base(root)
	if err := file.Copy(filepath.Join(appName, "Contents", "MacOS", bundle.AppName), execName); err != nil {
		return err
	}

	printVerbose("syncing resources")
	if err := file.Sync(appResources, resources); err != nil {
		return err
	}

	printVerbose("generating icons")
	if len(bundle.Icon) != 0 {
		return generateAppIcons(ctx, bundle.Icon, appResources)
	}
	return nil
}

func generateAppIcons(ctx context.Context, icon, appResources string) error {
	appIcon := filepath.Join(appResources, icon)

	iconset := filepath.Base(icon)
	iconset = strings.TrimSuffix(iconset, filepath.Ext(iconset))
	iconset = filepath.Join(appResources, iconset)
	iconset += ".iconset"

	if err := os.Mkdir(iconset, os.ModeDir|0755); err != nil {
		return err
	}
	defer os.RemoveAll(iconset)

	standardName := func(w, h int) string {
		return filepath.Join(iconset, fmt.Sprintf("icon_%vx%v.png", w, h))
	}

	retinaName := func(w, h, s int) string {
		return filepath.Join(iconset, fmt.Sprintf("icon_%vx%v@%vx.png", w, h, s))
	}

	if err := generateIcons(appIcon, []iconInfo{
		{Name: retinaName(512, 512, 2), Width: 512, Height: 512, Scale: 2},
		{Name: standardName(512, 512), Width: 512, Height: 512, Scale: 1},

		{Name: retinaName(256, 256, 2), Width: 256, Height: 256, Scale: 2},
		{Name: standardName(256, 256), Width: 256, Height: 256, Scale: 1},

		{Name: retinaName(128, 128, 2), Width: 128, Height: 128, Scale: 2},
		{Name: standardName(128, 128), Width: 128, Height: 128, Scale: 1},

		{Name: retinaName(32, 32, 2), Width: 32, Height: 32, Scale: 2},
		{Name: standardName(32, 32), Width: 32, Height: 32, Scale: 1},

		{Name: retinaName(16, 512, 2), Width: 16, Height: 16, Scale: 2},
		{Name: standardName(16, 512), Width: 16, Height: 16, Scale: 1},
	}); err != nil {
		return err
	}

	return execute(ctx, "iconutil", "-c", "icns", iconset)
}

func signAppBundle(ctx context.Context, bundle driver.Bundle, sign, root, appName string) error {
	entitlementsName := filepath.Join(root, ".entitlements")
	if err := generatePlist(entitlementsName, entitlements, bundle); err != nil {
		return err
	}
	defer os.Remove(entitlementsName)

	if err := execute(ctx, "codesign",
		"--force",
		"--sign",
		sign,
		"--entitlements",
		entitlementsName,
		appName,
	); err != nil {
		return err
	}

	return execute(ctx, "codesign",
		"--verify",
		"--deep",
		"--strict",
		"--verbose=2",
		appName,
	)
}

func createAppPkg(ctx context.Context, bundle driver.Bundle, sign, appName string) error {
	return execute(ctx,
		"productbuild",
		"--component",
		appName,
		"/Applications",
		"--sign",
		sign,
		bundle.AppName+".pkg",
	)
}

type macRunConfig struct {
	LogsAddr string `conf:"logs-addr" help:"The address used to listen app logs." validate:"nonzero"`
	Sign     string `conf:"sign"      help:"The signing identifier to sign the app.\n\t\033[95msecurity find-identity -v -p codesigning\033[00m to see signing identifiers.\n\thttps://developer.apple.com/library/content/documentation/Security/Conceptual/CodeSigningGuide/Procedures/Procedures.html to create one."`
	Debug    bool   `conf:"d"         help:"Enable debug mode is enabled."`
	Verbose  bool   `conf:"v"         help:"Enable verbose mode."`
}

func runMac(ctx context.Context, args []string) {
	c := macRunConfig{
		LogsAddr: ":9000",
	}

	ld := conf.Loader{
		Name:    "mac run",
		Args:    args,
		Usage:   "[options...] [.app path]",
		Sources: []conf.Source{conf.NewEnvSource("GOAPP", os.Environ()...)},
	}

	_, roots := conf.LoadWith(&c, ld)
	verbose = c.Verbose

	if len(roots) == 0 {
		roots = []string{"."}
	}

	appname := roots[0]
	if !strings.HasSuffix(appname, ".app") {
		var err error
		if appname, err = buildMacApp(ctx, appname, macBuildConfig{
			Sign: c.Sign,
		}); err != nil {
			fail("building app failed: %s", err)
		}
	}

	go listenLogs(ctx, c.LogsAddr)
	time.Sleep(time.Millisecond * 500)

	os.Unsetenv("GOAPP_BUNDLE")
	os.Setenv("GOAPP_LOGS_ADDR", c.LogsAddr)
	os.Setenv("GOAPP_DEBUG", fmt.Sprintf("%v", c.Debug))

	printVerbose("running %s", appname)
	if err := execute(ctx, "open", "--wait-apps", appname); err != nil {
		printErr("%s", err)
	}
}

func listenLogs(ctx context.Context, addr string) {
	logs := logs.GoappServer{
		Addr:   addr,
		Writer: os.Stderr,
	}

	err := logs.ListenAndLog(ctx)
	if err != nil {
		printErr("listening logs failed: %s", err)
	}
}

func win(ctx context.Context, args []string) {
	printErr("you are not on Windows!")
	os.Exit(-1)
}
