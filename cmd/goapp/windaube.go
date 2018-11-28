package main

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

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

	// The path of the Windows 10 SDK dir.
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
		return errors.New("archecture not set")
	}

	if len(pkg.SDK) == 0 {
		return errors.New("windows 10 sdk not set")
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
	pkg.tmpDir = filepath.Join(pkg.tmpDir, "goapp")
	pkg.tmpExecutable = filepath.Join(pkg.tmpDir, execName)

	pkg.resourcesDir = filepath.Join(pkg.Output, "Resources")
	pkg.executable = filepath.Join(pkg.Output, execName)
	return nil
}

type winSettings struct {
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
		builds = append(builds, fi.Name())
	}

	if len(builds) == 0 {
		return ""
	}

	sort.SliceStable(builds, func(i, j int) bool {
		return strings.Compare(builds[i], builds[j]) > 0
	})

	return filepath.Join(sdkRoot, builds[0])
}
