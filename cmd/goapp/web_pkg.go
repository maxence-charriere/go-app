package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/murlokswarm/app/internal/file"
	"github.com/pkg/errors"
)

// WebPackage represents a directory that contains a website.
// It implements the Package interface.
type WebPackage struct {
	// The path where the sources are.
	// It must refer to a Go main package.
	// Default is ".".
	Sources string

	// The path where the package is saved.
	// If not set, the ".app" extension is added.
	Output string

	// Minify the gopher js client.
	Minify bool

	// The bind address used when run.
	Addr string

	// Rune the client with the default browser.
	Browser bool

	// Run the client with chrome.
	Chrome bool

	// Enable verbose mode.
	Verbose bool

	// Force rebuilding of package that are already up-to-date.
	Force bool

	// Enable data race detection.
	Race bool

	// The function to log events.
	Log func(string, ...interface{})

	name                string
	workingDir          string
	sourcesResourcesDir string
	resourcesDir        string
	executable          string
	goappjs             string
}

// Init satisfies the Package interface.
func (pkg *WebPackage) Init(ctx context.Context) error {
	if err := pkg.init(); err != nil {
		return err
	}

	pkg.Log("creating resources directory")
	if err := os.MkdirAll(filepath.Join(pkg.Sources, "resources", "css"), 0755); err != nil {
		return err
	}

	pkg.Log("installing gopherjs")
	return pkg.installGopherJS(ctx)
}

func (pkg *WebPackage) installGopherJS(ctx context.Context) error {
	cmd := []string{"go", "get", "-u"}

	if pkg.Verbose {
		cmd = append(cmd, "-v")
	}

	cmd = append(cmd, "github.com/gopherjs/gopherjs")
	return execute(ctx, cmd[0], cmd[1:]...)
}

func (pkg *WebPackage) init() (err error) {
	if len(pkg.Sources) == 0 || pkg.Sources == "." || pkg.Sources == "./" {
		pkg.Sources = "."
	}
	if pkg.Sources, err = filepath.Abs(pkg.Sources); err != nil {
		return err
	}

	execName := filepath.Base(pkg.Sources)

	if len(pkg.Output) == 0 {
		pkg.Output = execName
	}
	if !strings.HasSuffix(pkg.Output, ".wapp") {
		pkg.Output += ".wapp"
	}

	pkg.name = filepath.Base(pkg.Output)

	if pkg.workingDir, err = os.Getwd(); err != nil {
		return err
	}

	pkg.sourcesResourcesDir = filepath.Join(pkg.Sources, "resources")
	pkg.resourcesDir = filepath.Join(pkg.Output, "resources")

	pkg.executable = filepath.Join(pkg.Output, execName)
	if runtime.GOOS == "windows" {
		pkg.executable += ".exe"
	}

	pkg.goappjs = filepath.Join(pkg.resourcesDir, "goapp.js")
	return nil
}

// Build satisfies the Package interface.
func (pkg *WebPackage) Build(ctx context.Context) error {
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

	pkg.Log("building javascript client")
	if err := pkg.buildJavascriptClient(ctx); err != nil {
		return err
	}

	pkg.Log("syncing resources")
	if err := pkg.syncResources(); err != nil {
		return err
	}

	return nil
}

func (pkg *WebPackage) create() error {
	dirs := []string{
		pkg.Output,
		pkg.resourcesDir,
	}

	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return err
		}
	}

	return nil
}

func (pkg *WebPackage) buildExecutable(ctx context.Context) error {
	args := []string{"go", "build", "-o", pkg.executable}

	if pkg.Verbose {
		args = append(args, "-v")
	}

	if pkg.Force {
		args = append(args, "-a")
	}

	if pkg.Race {
		args = append(args, "-race")
	}

	return execute(ctx, args[0], args[1:]...)
}

func (pkg *WebPackage) buildJavascriptClient(ctx context.Context) error {
	if runtime.GOOS == "windows" {
		os.Setenv("GOOS", "linux")
		defer os.Unsetenv("GOOS")
	}

	args := []string{"gopherjs", "build", "-o", pkg.goappjs}

	if pkg.Minify {
		args = append(args, "-m")
	}

	if pkg.Verbose {
		args = append(args, "-v")
	}

	return execute(ctx, args[0], args[1:]...)
}

func (pkg *WebPackage) syncResources() error {
	return file.Sync(pkg.resourcesDir, pkg.sourcesResourcesDir)
}

// Run satisfies the Package interface.
func (pkg *WebPackage) Run(ctx context.Context) error {
	if err := pkg.Build(ctx); err != nil {
		return err
	}

	executable, err := filepath.Abs(pkg.executable)
	if err != nil {
		return err
	}

	if err = os.Chdir(pkg.Output); err != nil {
		return err
	}

	if pkg.Browser || pkg.Chrome {
		go func() {
			time.Sleep(time.Millisecond * 250)

			if err := pkg.launchWithBrowser(ctx); err != nil {
				panic(err)
			}
		}()
	}

	pkg.Log("running server")
	return execute(ctx, executable)
}

func (pkg *WebPackage) launchWithBrowser(ctx context.Context) error {
	addr := pkg.Addr
	if !strings.HasPrefix(addr, "http://") {
		addr = "http://" + addr
	}

	u, err := url.Parse(addr)
	if err != nil {
		return err
	}
	if strings.HasPrefix(u.Host, ":") {
		u.Host = "127.0.0.1" + u.Host
	}
	addr = u.String()

	fmt.Println("ADDR:", addr)

	if pkg.Chrome {
		pkg.Log("running client in google chrome")
		return pkg.launchWithChrome(ctx, addr)
	}

	pkg.Log("running client in default browser")
	return pkg.launchWithDefaultBrowser(ctx, addr)
}

func (pkg *WebPackage) launchWithChrome(ctx context.Context, url string) error {
	var cmd []string

	switch runtime.GOOS {
	case "darwin":
		cmd = []string{"open", "-a", "Google Chrome", url}

	case "windows":
		cmd = []string{"powershell", "start", "chrome", url}

	case "linux":
		cmd = []string{"google-chrome", url}

	default:
		return errors.New("unsuported operation system")
	}

	return execute(ctx, cmd[0], cmd[1:]...)
}

func (pkg *WebPackage) launchWithDefaultBrowser(ctx context.Context, url string) error {
	var cmd []string

	switch runtime.GOOS {
	case "darwin":
		cmd = []string{"open", url}

	case "windows":
		cmd = []string{"powershell", "start", url}

	case "linux":
		cmd = []string{"xdg-open", url}

	default:
		return errors.New("unsuported operation system")
	}

	return execute(ctx, cmd[0], cmd[1:]...)
}
