package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/murlokswarm/app/internal/file"
	"github.com/pkg/errors"
)

// WinPackage represents a package for a Windows app.
// It implements the Package interface.
type WinPackage struct {
	// The path where the sources are.
	// It must refer to a Go main package.
	// Default is ".".
	Sources string

	// The path where the package is saved.
	// If not set, the ".app" extension is added.
	Output string

	// The targetted architecture:
	// - x64
	// - x86
	Architecture string

	// The path of the Windows 10 SDK directory.
	SDK string

	// Creates a .appx to be uploaed on the Windows Store.
	AppStore bool

	// Force rebuilding of package that are already up-to-date.
	Force bool

	// Enable data race detection.
	Race bool

	// Enable verbose mode.
	Verbose bool

	// The function to log events.
	Log func(string, ...interface{})

	name                string
	workingDir          string
	sourcesResourcesDir string
	tmpDir              string
	tmpExecutable       string
	assetsDir           string
	resourcesDir        string
	executable          string
	settings            winSettings
}

// Init satisfies the Package interface.
func (pkg *WinPackage) Init(ctx context.Context) error {
	if err := pkg.init(); err != nil {
		return err
	}

	pkg.Log("creating resources directory")
	if err := os.MkdirAll(filepath.Join(pkg.Sources, "resources", "css"), 0755); err != nil {
		return err
	}

	pkg.Log("Windows 10 SDK manual install required: https://developer.microsoft.com/en-US/windows/downloads/windows-10-sdk")
	pkg.Log("mingw64 manual install required: http://mingw-w64.org/doku.php/download/mingw-builds")
	return nil
}

func (pkg *WinPackage) init() (err error) {
	if runtime.GOOS != "windows" {
		return errors.New("operating system is not Windows")
	}

	if len(pkg.Architecture) == 0 {
		pkg.Architecture = runtime.GOARCH
	}

	if len(pkg.SDK) == 0 {
		pkg.SDK = winSDKDirectory(defaultWinSSDKRoot)
	}

	if len(pkg.Sources) == 0 || pkg.Sources == "." || pkg.Sources == "./" {
		pkg.Sources = "."
	}
	if pkg.Sources, err = filepath.Abs(pkg.Sources); err != nil {
		return err
	}

	name := filepath.Base(pkg.Sources)
	execName := name + ".exe"

	if len(pkg.Output) == 0 {
		pkg.Output = name
	}
	if !strings.HasSuffix(pkg.Output, ".app") {
		pkg.Output += ".app"
	}

	pkg.name = filepath.Base(pkg.Output)

	if pkg.workingDir, err = os.Getwd(); err != nil {
		return err
	}

	pkg.sourcesResourcesDir = filepath.Join(pkg.Sources, "resources")

	if pkg.tmpDir = os.Getenv("TEMP"); len(pkg.tmpDir) == 0 {
		return errors.New("tmp dir not set")
	}
	pkg.tmpDir = filepath.Join(pkg.tmpDir, "goapp", name)
	pkg.tmpExecutable = filepath.Join(pkg.tmpDir, execName)

	pkg.assetsDir = filepath.Join(pkg.Output, "Assets")
	pkg.resourcesDir = filepath.Join(pkg.Output, "Resources")
	pkg.executable = filepath.Join(pkg.Output, execName)
	return nil
}

// Build satisfies the Package interface.
func (pkg *WinPackage) Build(ctx context.Context) error {
	if err := pkg.init(); err != nil {
		return err
	}

	pkg.Log("creating %s", pkg.name)
	if err := pkg.create(); err != nil {
		return err
	}

	pkg.Log("building executable")
	if err := pkg.buildExecutable(ctx); err != nil {
		return err
	}

	pkg.Log("reading settings")
	if err := pkg.readSettings(ctx); err != nil {
		return err
	}

	pkg.Log("generating AppxManifest.xml")
	if err := pkg.generateAppxManifest(); err != nil {
		return err
	}

	pkg.Log("syncing resources")
	if err := pkg.syncResources(); err != nil {
		return err
	}

	pkg.Log("generating icons")
	if err := pkg.generateIcons(ctx); err != nil {
		return err
	}

	return nil
}

func (pkg *WinPackage) create() error {
	dirs := []string{
		pkg.Output,
		pkg.assetsDir,
		pkg.resourcesDir,
	}

	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return err
		}
	}

	return nil
}

func (pkg *WinPackage) buildExecutable(ctx context.Context) error {
	args := []string{
		"go", "build",
		"-o", pkg.tmpExecutable,
	}

	ldflags := []string{"-X github.com/murlokswarm/app.Kind=desktop"}

	if pkg.Verbose {
		args = append(args, "-v")
	} else {
		ldflags = append(ldflags, "-H=windowsgui")
	}

	args = append(args, "-ldflags", strings.Join(ldflags, " "))

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

func (pkg *WinPackage) readSettings(ctx context.Context) error {
	settingsPath := filepath.Join(pkg.tmpDir, "settings.json")
	defer os.RemoveAll(settingsPath)

	os.Setenv("GOAPP_BUILD", settingsPath)
	defer os.Unsetenv("GOAPP_BUILD")

	if err := execute(ctx, pkg.tmpExecutable); err != nil {
		return err
	}

	b, err := ioutil.ReadFile(settingsPath)
	if err != nil {
		return err
	}

	s := winSettings{}
	if err = json.Unmarshal(b, &s); err != nil {
		return err
	}

	s.Executable = filepath.Base(pkg.executable)
	name := strings.TrimSuffix(s.Executable, ".exe")
	user := filepath.Base(os.Getenv("USERPROFILE"))
	user = strings.Replace(user, " ", "", -1)
	user = strings.Replace(user, "\t", "", -1)

	s.Name = stringWithDefault(s.Name, name)
	s.ID = stringWithDefault(s.ID, fmt.Sprintf("%v.%v", user, s.Name))
	s.EntryPoint = strings.Replace(s.Executable, ".exe", ".app", 1)
	s.Description = stringWithDefault(s.Description, s.Name)
	s.Publisher = stringWithDefault(s.Publisher, user)
	s.URLScheme = stringWithDefault(s.URLScheme, name)
	s.Icon = stringWithDefault(s.Icon, "logo.png")

	pkg.settings = s

	if b, err = json.MarshalIndent(s, "", "    "); err != nil {
		return err
	}

	pkg.Log("settings: %s", b)
	return nil
}

func (pkg *WinPackage) generateAppxManifest() error {
	manifest := filepath.Join(pkg.Output, "AppxManifest.xml")
	return generateTemplatedFile(manifest, appxManifestTmpl, pkg.settings)
}

func (pkg *WinPackage) syncResources() error {
	return file.Sync(pkg.resourcesDir, pkg.sourcesResourcesDir)
}

func (pkg *WinPackage) generateIcons(ctx context.Context) error {
	icon := filepath.Join(pkg.resourcesDir, pkg.settings.Icon)
	if _, err := os.Stat(icon); os.IsNotExist(err) {
		file.Copy(icon, filepath.Join(murlokswarm(), "logo.png"))
	}

	scaled := func(n string, s int) string {
		if s <= 1 {
			return filepath.Join(pkg.assetsDir, fmt.Sprintf("%s.png", n))
		}

		return filepath.Join(pkg.assetsDir, fmt.Sprintf("%s.scale-%v.png", n, s))
	}

	return generateIcons(icon, []iconInfo{
		{Name: scaled("Square44x44Logo", 1), Width: 44, Height: 44, Scale: 1, Padding: true},
		{Name: scaled("Square44x44Logo", 100), Width: 44, Height: 44, Scale: 1},
		{Name: scaled("Square44x44Logo", 125), Width: 44, Height: 44, Scale: 1.25},
		{Name: scaled("Square44x44Logo", 150), Width: 44, Height: 44, Scale: 1.5},
		{Name: scaled("Square44x44Logo", 200), Width: 44, Height: 44, Scale: 2},
		{Name: scaled("Square44x44Logo", 400), Width: 44, Height: 44, Scale: 4},

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

// Run satisfies the Package interface.
func (pkg *WinPackage) Run(ctx context.Context) error {
	panic("not implemented")
}

// Clean satisfies the Package interface.
func (pkg *WinPackage) Clean(ctx context.Context) error {
	if err := pkg.init(); err != nil {
		return err
	}

	pkg.Log("removing %s", pkg.Output)
	if err := os.RemoveAll(pkg.Output); err != nil {
		return err
	}

	pkg.Log("removing %s", pkg.tmpDir)
	return os.RemoveAll(pkg.tmpDir)
}

type winSettings struct {
	Executable  string
	Name        string
	ID          string
	EntryPoint  string
	Description string
	Publisher   string
	URLScheme   string
	Icon        string
}

var (
	defaultWinSSDKRoot = filepath.Join(`C:\`, "Program Files (x86)", "Windows Kits", "10", "bin")
)

func winSDKDirectory(sdkRoot string) string {
	dirs, err := ioutil.ReadDir(sdkRoot)
	if err != nil {
		return ""
	}

	builds := make([]string, 0, len(dirs))

	for _, fi := range dirs {
		if name := fi.Name(); strings.HasPrefix(name, "10.") {
			builds = append(builds, fi.Name())
		}
	}

	if len(builds) == 0 {
		return ""
	}

	sort.SliceStable(builds, func(i, j int) bool {
		return strings.Compare(builds[i], builds[j]) > 0
	})

	return filepath.Join(sdkRoot, builds[0])
}

func defaultWinArchitecture() string {
	switch arch := runtime.GOARCH; arch {
	case "386":
		return "x86"

	case "amd64":
		return "x64"

	default:
		return arch
	}
}
