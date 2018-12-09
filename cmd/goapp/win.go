package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/murlokswarm/app/internal/file"
	"github.com/pkg/errors"
	"github.com/segmentio/conf"
)

type winInitConfig struct {
	Verbose bool `conf:"v" help:"Enable verbose mode."`
}

type winBuildConfig struct {
	Output       string `conf:"o"             help:"The path where the package is saved."`
	Architecture string `conf:"arch"          help:"The targetted architecture."`
	SDK          string `conf:"sdk"           help:"The path of the Windows 10 SDK directory."`
	Appx         bool   `conf:"appx"          help:"Generates an .appx package."`
	AppxPfx      string `conf:"appx-pfx"      help:"The path of the certificate to sign the app."`
	AppxPassword string `conf:"appx-password" help:"The passord to use with the provided certificate."`
	Force        bool   `conf:"force"         help:"Force rebuilding of package that are already up-to-date."`
	Race         bool   `conf:"race"          help:"Enable data race detection."`
	Verbose      bool   `conf:"v"             help:"Enable verbose mode."`
	Dev          bool   `conf:"dev"           help:"Enable goapp dev mode."`
}

type winRunConfig struct {
	Output       string `conf:"o"     help:"The path where the package is saved."`
	Architecture string `conf:"arch"  help:"The targetted architecture."`
	SDK          string `conf:"sdk"   help:"The path of the Windows 10 SDK directory."`
	Force        bool   `conf:"force" help:"Force rebuilding of package that are already up-to-date."`
	Race         bool   `conf:"race"  help:"Enable data race detection."`
	Verbose      bool   `conf:"v"     help:"Enable verbose mode."`
	Dev          bool   `conf:"dev"   help:"Enable goapp dev mode."`
}

type winCleanConfig struct {
	Output  string `conf:"o" help:"The path where the package is saved."`
	Verbose bool   `conf:"v" help:"Enable verbose mode."`
}

func win(ctx context.Context, args []string) {
	ld := conf.Loader{
		Name: "goapp win",
		Args: args,
		Commands: []conf.Command{
			{Name: "init", Help: "Download the Windows 10 dev tools and create required directories."},
			{Name: "build", Help: "Build the Windows app."},
			{Name: "run", Help: "Run a Windows app."},
			{Name: "clean", Help: "Delete a Windows app and its temporary build files."},
			{Name: "help", Help: "Show the Windows help"},
		},
	}

	switch cmd, args := conf.LoadWith(nil, ld); cmd {
	case "init":
		initWin(ctx, args)

	case "build":
		buildWin(ctx, args)

	case "run":
		runWin(ctx, args)

	case "clean":
		cleanWin(ctx, args)

	case "help":
		ld.PrintHelp(nil)

	default:
		panic("unreachable")
	}
}

func initWin(ctx context.Context, args []string) {
	c := winInitConfig{}

	ld := conf.Loader{
		Name:    "win init",
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

	pkg := WinPackage{
		Sources: sources,
		Verbose: c.Verbose,
		Log:     printVerbose,
	}

	if err := pkg.Init(ctx); err != nil {
		fail("%s", err)
	}

	printSuccess("init succeeded")
}

func buildWin(ctx context.Context, args []string) {
	c := winBuildConfig{
		SDK: winSDKDirectory(defaultWinSSDKRoot),
	}

	ld := conf.Loader{
		Name:    "win build",
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

	pkg := WinPackage{
		Sources:      sources,
		Output:       c.Output,
		SDK:          c.SDK,
		Appx:         c.Appx,
		AppxPfx:      c.AppxPfx,
		AppxPassword: c.AppxPassword,
		Verbose:      c.Verbose,
		Force:        c.Force,
		Race:         c.Race,
		dev:          c.Dev,
		Log:          printVerbose,
	}

	if err := pkg.Build(ctx); err != nil {
		fail("%s", err)
	}

	printSuccess("build succeeded")
}

func runWin(ctx context.Context, args []string) {
	c := winRunConfig{
		SDK: winSDKDirectory(defaultWinSSDKRoot),
	}

	ld := conf.Loader{
		Name:    "win run",
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

	pkg := WinPackage{
		Sources: sources,
		Output:  c.Output,
		SDK:     c.SDK,
		Verbose: c.Verbose,
		Force:   c.Force,
		Race:    c.Race,
		dev:     c.Dev,
		Log:     printVerbose,
	}

	if err := pkg.Run(ctx); err != nil {
		fail("%s", err)
	}
}

func cleanWin(ctx context.Context, args []string) {
	c := winCleanConfig{}

	ld := conf.Loader{
		Name:    "win clean",
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

	pkg := WinPackage{
		Sources: sources,
		Output:  c.Output,
		Verbose: c.Verbose,
		Log:     printVerbose,
	}

	if err := pkg.Clean(ctx); err != nil {
		fail("%s", err)
	}

	printSuccess("clean succeeded")
}

func init() {
	greenColor = ""
	redColor = ""
	orangeColor = ""
	defaultColor = ""
}

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

	// The path of the Windows 10 SDK directory.
	SDK string

	// Reports whether to creates a .appx package.
	Appx bool

	// The path of the certificate (.pfx) to sign the app.
	// Sings the package when set.
	AppxPfx string

	// The password to use with the .pfx in order to sign the app.
	AppxPassword string

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
	manifest            string
	resourcesPri        string
	appx                string
	settings            winSettings
	dev                 bool
	logsURL             string
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

	if len(pkg.SDK) == 0 {
		pkg.SDK = winSDKDirectory(defaultWinSSDKRoot)
	}

	if len(pkg.Sources) == 0 || pkg.Sources == "." || pkg.Sources == "./" {
		pkg.Sources = "."
	}
	if pkg.Sources, err = filepath.Abs(pkg.Sources); err != nil {
		return err
	}

	if len(pkg.AppxPfx) != 0 {
		if pkg.AppxPfx, err = filepath.Abs(pkg.AppxPfx); err != nil {
			return err
		}
	}

	name := filepath.Base(pkg.Sources)
	if name == "uwp" {
		return errors.New("sources dir can't be named uwp")
	}

	execName := name + ".exe"

	if len(pkg.Output) == 0 {
		pkg.Output = name
	}
	if !strings.HasSuffix(pkg.Output, ".app") {
		pkg.Output += ".app"
	}

	pkg.Output, err = filepath.Abs(pkg.Output)
	if err != nil {
		return err
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
	pkg.manifest = filepath.Join(pkg.Output, "AppxManifest.xml")
	pkg.resourcesPri = filepath.Join(pkg.Output, "resources.pri")
	pkg.appx = pkg.Output + "x"
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

	pkg.Log("installing uwp dependencies")
	if err := pkg.installUWPDeps(); err != nil {
		return err
	}

	pkg.Log("reading settings")
	if err := pkg.readSettings(ctx); err != nil {
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

	if err := pkg.generateSupportedFilesIcons(ctx); err != nil {
		return err
	}

	pkg.Log("generating AppxManifest.xml")
	if err := pkg.generateAppxManifest(); err != nil {
		return err
	}

	printVerbose("generating resources.pri")
	if err := pkg.generatePri(ctx); err != nil {
		return err
	}

	if !pkg.Appx {
		printVerbose("deploying dev version")
		return pkg.deployDev(ctx)
	}

	printVerbose("creating %sx", pkg.name)
	if err := pkg.createAppx(ctx); err != nil {
		return err
	}

	if len(pkg.AppxPfx) != 0 {
		printVerbose("signing %sx", pkg.name)
		return pkg.signAppx(ctx)
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
		ldflags = append(ldflags, "-X github.com/murlokswarm/app/drivers/win.debug=true")
	}

	if pkg.dev {
		ldflags = append(ldflags, "-X github.com/murlokswarm/app/drivers/win.dev=true")
	} else {
		if len(pkg.logsURL) != 0 {
			ldflags = append(ldflags, "-X github.com/murlokswarm/app/drivers/win.logsURL="+pkg.logsURL)
		}

		ldflags = append(ldflags, "-H=windowsgui")

	}

	args = append(args, "-ldflags", strings.Join(ldflags, " "))

	if pkg.Force {
		args = append(args, "-a")
	}

	if pkg.Race {
		args = append(args, "-race")
	}

	args = append(args, pkg.Sources)

	if err := execute(ctx, args[0], args[1:]...); err != nil {
		return err
	}

	return file.Copy(pkg.executable, pkg.tmpExecutable)
}

func (pkg *WinPackage) installUWPDeps() error {
	uwpDir := file.RepoPath("cmd", "goapp", "uwp", winArch())

	files, err := ioutil.ReadDir(uwpDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		src := filepath.Join(uwpDir, f.Name())
		dst := filepath.Join(pkg.Output, f.Name())

		if err = file.Copy(dst, src); err != nil {
			return err
		}
	}

	return nil
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
	s.URLScheme = stringWithDefault(s.URLScheme, "goapp-"+name)
	s.Icon = stringWithDefault(s.Icon, "logo.png")

	if err = validateWinFileTypes(s.SupportedFiles); err != nil {
		return err
	}

	pkg.settings = s

	if b, err = json.MarshalIndent(s, "", "    "); err != nil {
		return err
	}

	pkg.Log("settings: %s", b)
	return nil
}

func (pkg *WinPackage) syncResources() error {
	return file.Sync(pkg.resourcesDir, pkg.sourcesResourcesDir)
}

func (pkg *WinPackage) generateIcons(ctx context.Context) error {
	scaled := func(n string, s int) string {
		if s <= 1 {
			return filepath.Join(pkg.assetsDir, fmt.Sprintf("%s.png", n))
		}

		return filepath.Join(pkg.assetsDir, fmt.Sprintf("%s.scale-%v.png", n, s))
	}

	icon := filepath.Join(pkg.resourcesDir, pkg.settings.Icon)
	if _, err := os.Stat(icon); os.IsNotExist(err) {
		file.Copy(icon, file.RepoPath("logo.png"))
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

func (pkg *WinPackage) generateSupportedFilesIcons(ctx context.Context) error {
	scaled := func(n string, s int) string {
		if s <= 1 {
			return filepath.Join(pkg.assetsDir, fmt.Sprintf("%s.png", n))
		}

		return filepath.Join(pkg.assetsDir, fmt.Sprintf("%s.scale-%v.png", n, s))
	}

	for _, f := range pkg.settings.SupportedFiles {
		if len(f.Icon) == 0 {
			continue
		}

		icon := filepath.Join(pkg.resourcesDir, f.Icon)
		if _, err := os.Stat(icon); os.IsNotExist(err) {
			return err
		}

		ext := filepath.Ext(f.Icon)
		iconName := strings.TrimSuffix(f.Icon, ext)

		if err := generateIcons(icon, []iconInfo{
			{Name: scaled(iconName, 1), Width: 44, Height: 44, Scale: 1},
			{Name: scaled(iconName, 100), Width: 44, Height: 44, Scale: 1},
			{Name: scaled(iconName, 125), Width: 44, Height: 44, Scale: 1.25},
			{Name: scaled(iconName, 150), Width: 44, Height: 44, Scale: 1.5},
			{Name: scaled(iconName, 200), Width: 44, Height: 44, Scale: 2},
			{Name: scaled(iconName, 400), Width: 44, Height: 44, Scale: 4},
		}); err != nil {
			return err
		}
	}

	return nil
}

func (pkg *WinPackage) generateAppxManifest() error {
	for i, f := range pkg.settings.SupportedFiles {
		if len(f.Icon) != 0 {
			f.Icon = filepath.Join("Assets", f.Icon)
			pkg.settings.SupportedFiles[i] = f
		}
	}

	return generateTemplatedFile(pkg.manifest, appxManifestTmpl, pkg.settings)
}

func (pkg *WinPackage) generatePri(ctx context.Context) error {
	os.Chdir(pkg.sdkDir())
	defer os.Chdir(pkg.workingDir)

	configName := filepath.Join(pkg.tmpDir, "priconfig.xml")
	config := []string{
		"makepri.exe", "createconfig",
		"/cf", configName,
		"/dq", "lang-en-US",
		"/pv", "10.0.0",
		"/o",
	}

	if _, err := executeQuiet(ctx, config[0], config[1:]...); err != nil {
		return errors.Wrap(err, "generating pri configuration failed")
	}
	defer os.RemoveAll(configName)

	new := []string{
		"makepri.exe", "new",
		"/cf", configName,
		"/pr", pkg.Output,
		"/mn", pkg.manifest,
		"/of", pkg.resourcesPri,
		"/o",
	}

	_, err := executeQuiet(ctx, new[0], new[1:]...)
	return err
}

func (pkg *WinPackage) sdkDir() string {
	return filepath.Join(pkg.SDK, winArch())
}

func (pkg *WinPackage) deployDev(ctx context.Context) error {
	cmd := []string{"powershell",
		"Add-AppxPackage",
		"-Path", pkg.manifest,
		"-Register",
	}

	return execute(ctx, cmd[0], cmd[1:]...)
}

func (pkg *WinPackage) createAppx(ctx context.Context) error {
	os.Chdir(pkg.sdkDir())
	defer os.Chdir(pkg.workingDir)

	cmd := []string{
		"makeappx.exe", "pack",
		"-d", pkg.Output,
		"-p", pkg.appx,
		"-o",
	}

	if verbose {
		cmd = append(cmd, "-v")
	}

	return execute(ctx, cmd[0], cmd[1:]...)
}

func (pkg *WinPackage) signAppx(ctx context.Context) error {
	os.Chdir(pkg.sdkDir())
	defer os.Chdir(pkg.workingDir)

	cmd := []string{
		"signtool.exe", "sign",
		"/fd", "SHA256", "/a",
		"/f", pkg.AppxPfx,
		"/p", pkg.AppxPassword,
	}

	if verbose {
		cmd = append(cmd, "/v")
	}

	cmd = append(cmd, pkg.appx)
	return execute(ctx, cmd[0], cmd[1:]...)
}

// Run satisfies the Package interface.
func (pkg *WinPackage) Run(ctx context.Context) error {
	wg := sync.WaitGroup{}
	wg.Add(1)

	handleLogs := func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/close" {
			wg.Done()
			return
		}

		line, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return
		}

		if pkg.Verbose {
			fmt.Fprint(os.Stderr, "    ")
		}

		fmt.Fprintln(os.Stderr, string(line))
	}

	server := httptest.NewServer(http.HandlerFunc(handleLogs))
	defer server.Close()
	pkg.logsURL = server.URL

	if err := pkg.Build(ctx); err != nil {
		return err
	}

	if err := execute(ctx, "powershell", "start", fmt.Sprintf("%s://", pkg.settings.URLScheme)); err != nil {
		return err
	}

	if !pkg.dev {
		wg.Wait()
	}

	return nil
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

	pkg.Log("removing %s", pkg.appx)
	if err := os.RemoveAll(pkg.appx); err != nil {
		return err
	}

	pkg.Log("removing %s", pkg.tmpDir)
	return os.RemoveAll(pkg.tmpDir)
}

type winSettings struct {
	Executable     string        `json:",omitempty"`
	Name           string        `json:",omitempty"`
	ID             string        `json:",omitempty"`
	EntryPoint     string        `json:",omitempty"`
	Description    string        `json:",omitempty"`
	Publisher      string        `json:",omitempty"`
	URLScheme      string        `json:",omitempty"`
	Icon           string        `json:",omitempty"`
	SupportedFiles []winFileType `json:",omitempty"`
}

type winFileType struct {
	Name       string             `json:",omitempty"`
	Help       string             `json:",omitempty"`
	Icon       string             `json:",omitempty"`
	Extensions []winFileExtension `json:",omitempty"`
}

type winFileExtension struct {
	Ext  string `json:",omitempty"`
	Mime string `json:",omitempty"`
}

func validateWinFileTypes(fileTypes []winFileType) error {
	for i, f := range fileTypes {
		if len(f.Name) == 0 {
			return errors.Errorf("file type at index %v: name is not set", i)
		}

		if len(f.Extensions) == 0 {
			return errors.Errorf("file type at index %v: no extensions", i)
		}

		for j, e := range f.Extensions {
			if len(e.Ext) == 0 {
				return errors.Errorf("file type at index %v: extension at index %v: ext not set", i, j)
			}

			if !strings.HasPrefix(e.Ext, ".") {
				return errors.Errorf(`file type at index %v: extension at index %v: ext not valid, change it to ".%s"`, i, j, e.Ext)
			}
		}
	}

	return nil
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

func winArch() string {
	switch arch := runtime.GOARCH; arch {
	case "386":
		return "x86"

	case "amd64":
		return "x64"

	case "arm":
		return "ARM"

	default:
		panic(errors.Errorf("architecture not supported: %s", arch))
	}
}
