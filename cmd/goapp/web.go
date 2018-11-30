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
	"github.com/segmentio/conf"
)

type webInitConfig struct {
	Verbose bool `conf:"v" help:"Enable verbose mode."`
}

type webBuildConfig struct {
	Output  string `conf:"o"      help:"The path where the package is saved."`
	Minify  bool   `conf:"minify" help:"Minify gopherjs file."`
	Force   bool   `conf:"force"  help:"Force rebuilding of package that are already up-to-date."`
	Race    bool   `conf:"race"   help:"Enable data race detection."`
	Verbose bool   `conf:"v"      help:"Enable verbose mode."`
}

type webRunConfig struct {
	Output  string `conf:"o"      help:"The path where the package is saved."`
	Addr    string `conf:"addr"   help:"The server bind address."`
	Browser bool   `conf:"b"      help:"Run the client."`
	Chrome  bool   `conf:"chrome" help:"Run the client with Google Chrome."`
	Minify  bool   `conf:"minify" help:"Minify gopherjs file."`
	Force   bool   `conf:"force"  help:"Force rebuilding of package that are already up-to-date."`
	Race    bool   `conf:"race"   help:"Enable data race detection."`
	Verbose bool   `conf:"v"      help:"Enable verbose mode."`
}

type webCleanConfig struct {
	Output  string `conf:"o" help:"The path where the package is saved."`
	Verbose bool   `conf:"v" help:"Enable verbose mode."`
}

func web(ctx context.Context, args []string) {
	ld := conf.Loader{
		Name: "goapp web",
		Args: args,
		Commands: []conf.Command{
			{Name: "init", Help: "Download gopherjs and create the required files and directories."},
			{Name: "build", Help: "Build a web app."},
			{Name: "run", Help: "Run the server and launch the client in the default browser."},
			{Name: "clean", Help: "Delete a web app."},
			{Name: "help", Help: "Show the web help"},
		},
	}

	switch cmd, args := conf.LoadWith(nil, ld); cmd {
	case "help":
		ld.PrintHelp(nil)

	case "init":
		initWeb(ctx, args)

	case "build":
		buildWeb(ctx, args)

	case "run":
		runWeb(ctx, args)

	case "clean":
		cleanWeb(ctx, args)

	default:
		panic("unreachable")
	}
}

func initWeb(ctx context.Context, args []string) {
	c := webInitConfig{}

	ld := conf.Loader{
		Name:    "web init",
		Args:    args,
		Usage:   "[options...] [package]",
		Sources: []conf.Source{conf.NewEnvSource("GOAPP", os.Environ()...)},
	}

	_, args = conf.LoadWith(&c, ld)
	verbose = c.Verbose

	sources := "."
	if len(args) != 0 {
		sources = args[0]
	}

	pkg := WebPackage{
		Sources: sources,
		Verbose: c.Verbose,
		Log:     printVerbose,
	}

	if err := pkg.Init(ctx); err != nil {
		fail("%s", err)
	}

	printSuccess("init succeeded")
}

func buildWeb(ctx context.Context, args []string) {
	c := webBuildConfig{
		Minify: true,
	}

	ld := conf.Loader{
		Name:    "web build",
		Args:    args,
		Usage:   "[options...] [package]",
		Sources: []conf.Source{conf.NewEnvSource("GOAPP", os.Environ()...)},
	}

	_, args = conf.LoadWith(&c, ld)
	verbose = c.Verbose

	sources := "."
	if len(args) != 0 {
		sources = args[0]
	}

	pkg := WebPackage{
		Sources: sources,
		Output:  c.Output,
		Minify:  c.Minify,
		Force:   c.Force,
		Race:    c.Race,
		Verbose: c.Verbose,
		Log:     printVerbose,
	}

	if err := pkg.Build(ctx); err != nil {
		fail("%s", err)
	}

	printSuccess("build succeeded")
}

func runWeb(ctx context.Context, args []string) {
	c := webRunConfig{
		Addr:   ":9001",
		Minify: true,
	}

	ld := conf.Loader{
		Name:    "web run",
		Args:    args,
		Usage:   "[options...] [package]",
		Sources: []conf.Source{conf.NewEnvSource("GOAPP", os.Environ()...)},
	}

	_, args = conf.LoadWith(&c, ld)
	verbose = c.Verbose

	sources := "."
	if len(args) != 0 {
		sources = args[0]
	}

	pkg := WebPackage{
		Sources: sources,
		Output:  c.Output,
		Addr:    c.Addr,
		Chrome:  c.Chrome,
		Browser: c.Browser,
		Minify:  c.Minify,
		Force:   c.Force,
		Race:    c.Race,
		Verbose: c.Verbose,
		Log:     printVerbose,
	}

	os.Setenv("GOAPP_SERVER_ADDR", c.Addr)
	defer os.Unsetenv("GOAPP_SERVER_ADDR")

	if err := pkg.Run(ctx); err != nil {
		fail("%s", err)
	}
}

func cleanWeb(ctx context.Context, args []string) {
	c := webCleanConfig{}

	ld := conf.Loader{
		Name:    "web clean",
		Args:    args,
		Usage:   "[options...] [package]",
		Sources: []conf.Source{conf.NewEnvSource("GOAPP", os.Environ()...)},
	}

	_, args = conf.LoadWith(&c, ld)
	verbose = c.Verbose

	sources := "."
	if len(args) != 0 {
		sources = args[0]
	}

	pkg := WebPackage{
		Sources: sources,
		Output:  c.Output,
		Verbose: c.Verbose,
		Log:     printVerbose,
	}

	if err := pkg.Clean(ctx); err != nil {
		fail("%s", err)
	}
}

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
	tmpDir              string
	tmpGoappjs          string
	tmpGoappjsMap       string
	resourcesDir        string
	executable          string
	goappjs             string
	goappjsMap          string
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

	sources, err := filepath.Abs(pkg.Sources)
	if err != nil {
		return err
	}

	name := filepath.Base(sources)

	if len(pkg.Output) == 0 {
		pkg.Output = name
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

	pkg.executable = filepath.Join(pkg.Output, name)
	if runtime.GOOS == "windows" {
		pkg.executable += ".exe"
	}

	tmp := ""
	switch runtime.GOOS {
	case "darwin":
		tmp = "TMPDIR"

	case "windows":
		tmp = "TEMP"

	default:
		tmp = "/tmp"
	}

	if pkg.tmpDir = os.Getenv(tmp); len(pkg.tmpDir) == 0 {
		return errors.New("tmp dir not set")
	}
	pkg.tmpDir = filepath.Join(pkg.tmpDir, "goapp", name)
	pkg.tmpGoappjs = filepath.Join(pkg.tmpDir, "goapp.js")
	pkg.tmpGoappjsMap = pkg.tmpGoappjs + ".map"

	pkg.goappjs = filepath.Join(pkg.resourcesDir, "goapp.js")
	pkg.goappjsMap = pkg.goappjs + ".map"
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

	pkg.Log("syncing resources")
	if err := pkg.syncResources(); err != nil {
		return err
	}

	pkg.Log("building javascript client")
	if err := pkg.buildJavascriptClient(ctx); err != nil {
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
	args := []string{"go", "build",
		"-ldflags", "-X github.com/murlokswarm/app.Kind=web",
		"-o", pkg.executable,
	}

	if pkg.Verbose {
		args = append(args, "-v")
	}

	if pkg.Force {
		args = append(args, "-a")
	}

	if pkg.Race {
		args = append(args, "-race")
	}

	args = append(args, pkg.Sources)
	return execute(ctx, args[0], args[1:]...)
}

func (pkg *WebPackage) syncResources() error {
	return file.Sync(pkg.resourcesDir, pkg.sourcesResourcesDir)
}

func (pkg *WebPackage) buildJavascriptClient(ctx context.Context) error {
	if runtime.GOOS == "linux" {
		os.Setenv("GOOS", "linux")
		defer os.Unsetenv("GOOS")
	}

	args := []string{"gopherjs", "build", "-o", pkg.tmpGoappjs}

	if pkg.Minify {
		args = append(args, "-m")
	}

	if pkg.Verbose {
		args = append(args, "-v")
	}

	args = append(args, pkg.Sources)

	if err := execute(ctx, args[0], args[1:]...); err != nil {
		return err
	}

	if err := file.Copy(pkg.goappjs, pkg.tmpGoappjs); err != nil {
		return err
	}

	return file.Copy(pkg.goappjsMap, pkg.tmpGoappjsMap)
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

// Clean satisfies the Package interface.
func (pkg *WebPackage) Clean(ctx context.Context) error {
	if err := pkg.init(); err != nil {
		return err
	}

	pkg.Log("removing %s", pkg.Output)
	return os.RemoveAll(pkg.Output)
}
