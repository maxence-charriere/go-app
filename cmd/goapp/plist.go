// +build darwin,amd64

package main

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type bundle struct {
	ExecName         string
	AppName          string
	ID               string
	URLScheme        string
	Version          string
	BuildNumber      int
	Icon             string
	DevRegion        string
	DeploymentTarget string
	Copyright        string
	Role             string
	Category         string
	Sandbox          bool
	Background       bool
	Server           bool
	Camera           bool
	Microphone       bool
	USB              bool
	Printers         bool
	Bluetooth        bool
	Contacts         bool
	Location         bool
	Calendar         bool
	FilePickers      string
	Downloads        string
	Pictures         string
	Music            string
	Movies           string
	SupportedFiles   []string
}

const plist = `
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleDevelopmentRegion</key>
	<string>{{.DevRegion}}</string>

	<key>CFBundleExecutable</key>
	<string>{{.ExecName}}</string>

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
	<key>com.apple.security.files.user-selected.{{.FilePickers}}</key>
	<true/>
	{{end}}
	{{if len .Downloads}}
	<key>com.apple.security.files.downloads.{{.Downloads}}</key>
	<true/>
	{{end}}
	{{if len .Pictures}}
	<key>com.apple.security.assets.pictures.{{.Pictures}}/key>
	<true/>
	{{end}}
	{{if len .Music}}
	<key>com.apple.security.assets.music.{{.Music}}</key>
	<true/>
	{{end}}
	{{if len .Movies}}
	<key>com.apple.security.assets.movies.{{.Movies}}/key>
	<true/>
	{{end}}
</dict>
</plist>
`

func generatePlist(filename string, plistTmpl string, b bundle) error {
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

	return tmpl.Execute(f, b)
}
