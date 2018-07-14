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
	"text/template"
	"time"

	driver "github.com/murlokswarm/app/drivers/mac"
	"github.com/murlokswarm/app/internal/file"
	"github.com/pkg/errors"
	"github.com/segmentio/conf"
)

type macBuildConfig struct {
	Bundle   bool   `conf:"bundle"   help:"Bundles the application into a .app."`
	Sign     string `conf:"sign"     help:"The signing identifier to sign the app.\n\t\033[95msecurity find-identity -v -p codesigning\033[00m to see signing identifiers.\n\thttps://developer.apple.com/library/content/documentation/Security/Conceptual/CodeSigningGuide/Procedures/Procedures.html to create one."`
	AppStore bool   `conf:"appstore" help:"Report whether the app will be packaged to be uploaded on the app store."`
}

func commands() []conf.Command {
	return []conf.Command{
		{Name: "web", Help: "Build app for web."},
		{Name: "mac", Help: "Build app for macOS."},
		{Name: "help", Help: "Show the help."},
	}
}

func mac(ctx context.Context, args []string) {
	ld := conf.Loader{
		Name: "goapp mac",
		Args: args,
		Commands: []conf.Command{
			{Name: "help", Help: "Show the macOS help"},
			{Name: "init", Help: "Download macOS SDK and create required file and directories."},
			{Name: "build", Help: "Build the macOS app."},
		},
	}

	switch cmd, args := conf.LoadWith(nil, ld); cmd {
	case "help":
		ld.PrintHelp(nil)

	case "init":
		initMac(ctx, args)

	case "build":
		buildMac(ctx, args)

	default:
		panic("unreachable")
	}
}

func initMac(ctx context.Context, args []string) {
	config := struct{}{}

	ld := conf.Loader{
		Name:    "mac init",
		Args:    args,
		Usage:   "[options...] [packages...]",
		Sources: []conf.Source{conf.NewEnvSource("GOAPP", os.Environ()...)},
	}

	defer func() {
		err := recover()
		if err != nil {
			ld.PrintHelp(nil)
			ld.PrintError(errors.Errorf("%s", err))
			os.Exit(-1)
		}
	}()

	_, unusedArgs := conf.LoadWith(&config, ld)
	roots, err := packageRoots(unusedArgs)
	if err != nil {
		panic(err)
	}

	fmt.Println("checking for xcode-select...")
	execute("xcode-select", "--install")

	for _, root := range roots {
		if err = initPackage(root); err != nil {
			printErr("%s", errors.Wrap(err, "init package"))
			return
		}
	}

	printSuccess("init succeeded")
}

func buildMac(ctx context.Context, args []string) {
	config := macBuildConfig{}

	ld := conf.Loader{
		Name:    "mac build",
		Args:    args,
		Usage:   "[options...] [package]",
		Sources: []conf.Source{conf.NewEnvSource("GOAPP", os.Environ()...)},
	}

	_, roots := conf.LoadWith(&config, ld)
	if len(roots) == 0 {
		roots = []string{"."}
	}
	root := roots[0]

	if err := goBuild(root, "-ldflags", "-s"); err != nil {
		printErr("%s", err)
		return
	}

	if config.Bundle {
		if err := bundleMacApp(root, config); err != nil {
			printErr("%s", err)
			return
		}
	}
	printSuccess("build succeeded")
}

func bundleMacApp(root string, c macBuildConfig) error {
	err := os.Setenv("GOAPP_BUNDLE", "true")
	if err != nil {
		return err
	}

	if root, err = filepath.Abs(root); err != nil {
		return err
	}

	var wd string
	if wd, err = os.Getwd(); err != nil {
		return err
	}

	if err = execute(filepath.Join(
		wd,
		filepath.Base(root),
	)); err != nil {
		return err
	}

	var data []byte
	if data, err = ioutil.ReadFile("goapp-mac.json"); err != nil {
		return err
	}
	defer os.Remove("goapp-mac.json")

	var bundle driver.Bundle
	if err = json.Unmarshal(data, &bundle); err != nil {
		return err
	}

	if bundle.Sandbox && len(c.Sign) == 0 {
		return errors.New("sanboxed app require to be signed")
	}

	bundle = fillBundle(bundle, root)
	data, _ = json.MarshalIndent(bundle, "", "  ")
	fmt.Println("bundle configuration:", string(data))

	appName := bundle.AppName + ".app"
	if err = createAppBundle(bundle, root, appName); err != nil {
		os.RemoveAll(appName)
		return err
	}

	if len(c.Sign) != 0 {
		if err = signAppBundle(bundle, c.Sign, root, appName); err != nil {
			os.RemoveAll(appName)
			return err
		}

		if c.AppStore {
			return createAppPkg(bundle, c.Sign, appName)
		}
	}
	return nil
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

func createAppBundle(bundle driver.Bundle, root, appName string) error {
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
	if err := os.Rename(execName, filepath.Join(appName, "Contents", "MacOS", execName)); err != nil {
		return err
	}

	if err := file.Sync(appResources, resources); err != nil {
		return err
	}

	if len(bundle.Icon) != 0 {
		return generateAppIcons(bundle.Icon, appResources)
	}
	return nil
}

func generateAppIcons(icon, appResources string) error {
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

	return execute("iconutil", "-c", "icns", iconset)
}

func signAppBundle(bundle driver.Bundle, sign, root, appName string) error {
	entitlementsName := filepath.Join(root, ".entitlements")
	if err := generatePlist(entitlementsName, entitlements, bundle); err != nil {
		return err
	}
	defer os.Remove(entitlementsName)

	if err := execute("codesign",
		"--force",
		"--sign",
		sign,
		"--entitlements",
		entitlementsName,
		appName,
	); err != nil {
		return err
	}

	return execute("codesign",
		"--verify",
		"--deep",
		"--strict",
		"--verbose=2",
		appName,
	)
}

func createAppPkg(bundle driver.Bundle, sign, appName string) error {
	return execute("productbuild",
		"--component",
		appName,
		"/Applications",
		"--sign",
		sign,
		bundle.AppName+".pkg",
	)
}

func generatePlist(filename string, plistTmpl string, bundle driver.Bundle) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	tmpl := template.Must(template.
		New("plist").
		Funcs(template.FuncMap{
			"icon": func(icon string) string {
				icon = filepath.Base(icon)
				return strings.TrimSuffix(icon, filepath.Ext(icon))
			},
		}).
		Parse(plistTmpl))

	return tmpl.Execute(f, bundle)
}

const plist = `
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleDevelopmentRegion</key>
	<string>{{.DevRegion}}</string>

	<key>CFBundleExecutable</key>
	<string>{{.AppName}}</string>

	<key>CFBundleIconFile</key>
	<string>{{icon .Icon}}</string>

	<key>CFBundleIdentifier</key>
	<string>{{.ID}}</string>

	<key>CFBundleInfoDictionaryVersion</key>
	<string>6.0</string>

	<key>CFBundleName</key>
	<string>{{.AppName}}</string>

	<key>CFBundlePackageType</key>
	<string>APPL</string>

	<key>CFBundleSupportedPlatforms</key>
	<array>
		<string>MacOSX</string>
	</array>

	<key>CFBundleShortVersionString</key>
	<string>{{.Version}}</string>

	<key>CFBundleVersion</key>
	<string>{{.BuildNumber}}</string>

	<key>LSMinimumSystemVersion</key>
	<string>{{.DeploymentTarget}}</string>

	<key>LSApplicationCategoryType</key>
	<string>{{.Category}}</string>

	{{if .Background}}
	<key>LSUIElement</key>
	<true/>
	{{end}}

	<key>NSHumanReadableCopyright</key>
	<string>{{html .Copyright}}</string>

	<key>NSPrincipalClass</key>
	<string>NSApplication</string>

	<key>NSAppTransportSecurity</key>
	<dict>
		<key>NSAllowsArbitraryLoadsInWebContent</key>
		<true/>
	</dict>

	<key>CFBundleDocumentTypes</key>
	<array>
		<dict>
			<key>CFBundleTypeName</key>
			<string>Supported files</string>
			<key>CFBundleTypeRole</key>
			<string>{{.Role}}</string>
			<key>LSItemContentTypes</key>
			<array>{{range .SupportedFiles}}
				<string>{{.}}</string>{{end}}
			</array>
		</dict>
	</array>

	<key>CFBundleURLTypes</key>
	<array>
		<dict>
			<key>CFBundleURLName</key>
			<string>{{.ID}}</string>
			<key>CFBundleTypeRole</key>
			<string>{{.Role}}</string>
			<key>CFBundleURLSchemes</key>
			<array>
				<string>{{.URLScheme}}</string>
			</array>
		</dict>
	</array>
</dict>
</plist>
`

const entitlements = `
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	{{if .Sandbox}}
	<key>com.apple.security.app-sandbox</key>
	<true/>
	{{end}}

	<!-- Network -->
	{{if .Server}}
	<key>com.apple.security.network.server</key>
	<true/>
	{{end}}
	<key>com.apple.security.network.client</key>
	<true/>

	<!-- Hadrware -->
	{{if .Camera}}
	<key>com.apple.security.device.camera</key>
	<true/>
	{{end}}
	{{if .Microphone}}
	<key>com.apple.security.device.microphone</key>
	<true/>
	{{end}}
	{{if .USB}}
	<key>com.apple.security.device.usb</key>
	<true/>
	{{end}}
	{{if .Printers}}
	<key>com.apple.security.print</key>
	<true/>
	{{end}}
	{{if .Bluetooth}}
	<key>com.apple.security.device.bluetooth</key>
	<true/>
	{{end}}

	<!-- AppData -->
	{{if .Contacts}}
	<key>com.apple.security.personal-information.addressbook</key>
	<true/>
	{{end}}
	{{if .Location}}
	<key>com.apple.security.personal-information.location</key>
	<true/>
	{{end}}
	{{if .Calendar}}
	<key>com.apple.security.personal-information.calendars</key>
	<true/>
	{{end}}

	<!-- FileAccess -->
	{{if len .FilePickers}}
	<key>com.apple.security.files.user-selected.{{.FileAccess.UserSelected}}</key>
	<true/>
	{{end}}
	{{if len .Downloads}}
	<key>com.apple.security.files.downloads.{{.FileAccess.Downloads}}</key>
	<true/>
	{{end}}
	{{if len .Pictures}}
	<key>com.apple.security.assets.pictures.{{.FileAccess.Pictures}}/key>
	<true/>
	{{end}}
	{{if len .Music}}
	<key>com.apple.security.assets.music.{{.FileAccess.Music}}</key>
	<true/>
	{{end}}
	{{if len .Movies}}
	<key>com.apple.security.assets.movies.{{.FileAccess.Movies}}/key>
	<true/>
	{{end}}
</dict>
</plist>
`

func openCommand() string {
	return "open"
}
