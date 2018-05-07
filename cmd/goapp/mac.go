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
	"github.com/pkg/errors"
	"github.com/segmentio/conf"
)

type macBuildConfig struct {
	Bundle bool   `conf:"bundle" help:"Bundles the application into a .app."`
	SignID string `conf:"signID" help:"The signing identifier to sign the app.\n\t\033[95msecurity find-identity -v -p codesigning\033[00m to see signing identifiers.\n\thttps://developer.apple.com/library/content/documentation/Security/Conceptual/CodeSigningGuide/Procedures/Procedures.html to create one."`
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
		}
	}
}

func bundleMacApp(root string, c macBuildConfig) error {
	err := os.Setenv("GOAPP_BUNDLE", "true")
	if err != nil {
		return err
	}

	if root, err = filepath.Abs(root); err != nil {
		return err
	}

	if err := execute(filepath.Join(".", filepath.Base(root))); err != nil {
		return err
	}

	data, err := ioutil.ReadFile("goapp-mac.json")
	if err != nil {
		return err
	}
	defer os.Remove("goapp-mac.json")

	var bundle driver.Bundle
	if err := json.Unmarshal(data, &bundle); err != nil {
		return err
	}

	bundle = fillBundle(bundle, root)
	data, _ = json.MarshalIndent(bundle, "", "  ")
	fmt.Println("bundle configuration:", string(data))

	appName := bundle.AppName + ".app"

	// TODO
	// create app bundle

	if bundle.Sandbox {
		if len(c.SignID) == 0 {
			return errors.New("sandboxed app require to be bundled with a singID")
		}

		entitlementsName := root + ".entitlements"
		if err = generatePlist(entitlementsName, entitlements, bundle); err != nil {
			os.Remove(appName)
			return err
		}
		defer os.Remove(entitlementsName)

		if err = signApp(c.SignID, entitlementsName, appName); err != nil {
			os.Remove(appName)
			return err
		}

		if err = signCheck(appName); err != nil {
			os.Remove(appName)
			return err
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

func generatePlist(filename string, plistTmpl string, bundle driver.Bundle) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	tmpl := template.Must(template.New("plist").Parse(plistTmpl))
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
	<string>{{.Icon}}</string>

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
	<key>com.apple.security.app-sandbox</key>
	<true/>

    <!-- Network -->
    {{if .Server}}
    <key>com.apple.security.network.server</key>
	<true/>
    {{end}}
    {{if .Client}}
	<key>com.apple.security.network.client</key>
	<true/>
    {{end}}

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

func signApp(signID, entitlementsName, appName string) error {
	return execute("codesign",
		"-s",
		signID,
		"--entitlements",
		entitlementsName,
		appName,
	)
}

func signCheck(appName string) error {
	return execute("codesign",
		"--verify",
		"--deep",
		"--strict",
		"--verbose=2",
		appName,
	)
}

func openCommand() string {
	return "open"
}
