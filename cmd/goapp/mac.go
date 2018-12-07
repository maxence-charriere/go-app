package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/murlokswarm/app/internal/file"
	"github.com/pkg/errors"
	"github.com/segmentio/conf"
)

type macInitConfig struct {
	Verbose bool `conf:"v" help:"Enable verbose mode."`
}

type macBuildConfig struct {
	Output           string `conf:"o"                 help:"The path where the package is saved."`
	DeploymentTarget string `conf:"deployment-target" help:"The version on MacOS the build is for."`
	SignID           string `conf:"sign-id"           help:"The signing identifier to sign the app.\n\t\033[95msecurity find-identity -v -p codesigning\033[00m to see signing identifiers.\n\thttps://developer.apple.com/library/content/documentation/Security/Conceptual/CodeSigningGuide/Procedures/Procedures.html to create one."`
	Sandbox          bool   `conf:"sandbox"           help:"Configure the app to run in sandbox mode."`
	AppStore         bool   `conf:"appstore"          help:"Creates a .pkg to be uploaded on the app store."`
	Force            bool   `conf:"force"             help:"Force rebuilding of package that are already up-to-date."`
	Race             bool   `conf:"race"              help:"Enable data race detection."`
	Verbose          bool   `conf:"v"                 help:"Enable verbose mode."`
}

type macRunConfig struct {
	Output           string `conf:"o"                 help:"The path where the package is saved."`
	DeploymentTarget string `conf:"deployment-target" help:"The version on MacOS the build is for."`
	SignID           string `conf:"sign-id"           help:"The signing identifier to sign the app.\n\t\033[95msecurity find-identity -v -p codesigning\033[00m to see signing identifiers.\n\thttps://developer.apple.com/library/content/documentation/Security/Conceptual/CodeSigningGuide/Procedures/Procedures.html to create one."`
	Sandbox          bool   `conf:"sandbox"           help:"Configure the app to run in sandbox mode."`
	Force            bool   `conf:"force"             help:"Force rebuilding of package that are already up-to-date."`
	Race             bool   `conf:"race"              help:"Enable data race detection."`
	Verbose          bool   `conf:"v"                 help:"Enable verbose mode."`
}

type macCleanConfig struct {
	Output  string `conf:"o" help:"The path where the package is saved."`
	Verbose bool   `conf:"v" help:"Enable verbose mode."`
}

func mac(ctx context.Context, args []string) {
	ld := conf.Loader{
		Name: "goapp mac",
		Args: args,
		Commands: []conf.Command{
			{Name: "init", Help: "Download MacOS SDK and create required directories."},
			{Name: "build", Help: "Build a MacOS app."},
			{Name: "run", Help: "Run a MacOS app."},
			{Name: "clean", Help: "Delete a MacOS app and its temporary build files."},
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

	case "clean":
		cleanMac(ctx, args)

	case "help":
		ld.PrintHelp(nil)

	default:
		panic("unreachable")
	}
}

func initMac(ctx context.Context, args []string) {
	c := macInitConfig{}

	ld := conf.Loader{
		Name:    "mac init",
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

	pkg := MacPackage{
		Sources: sources,
		Verbose: c.Verbose,
		Log:     printVerbose,
	}

	if err := pkg.Init(ctx); err != nil {
		fail("%s", err)
	}

	printSuccess("init succeeded")
}

func buildMac(ctx context.Context, args []string) {
	c := macBuildConfig{
		DeploymentTarget: macOSVersion(),
	}

	ld := conf.Loader{
		Name:    "mac build",
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

	pkg := MacPackage{
		Sources:          sources,
		Output:           c.Output,
		DeploymentTarget: c.DeploymentTarget,
		SignID:           c.SignID,
		Sandbox:          c.Sandbox,
		AppStore:         c.AppStore,
		Verbose:          c.Verbose,
		Force:            c.Force,
		Race:             c.Race,
		Log:              printVerbose,
	}

	if err := pkg.Build(ctx); err != nil {
		fail("%s", err)
	}

	printSuccess("build succeeded")
}

func runMac(ctx context.Context, args []string) {
	c := macRunConfig{
		DeploymentTarget: macOSVersion(),
	}

	ld := conf.Loader{
		Name:    "mac run",
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

	pkg := MacPackage{
		Sources:          sources,
		Output:           c.Output,
		DeploymentTarget: c.DeploymentTarget,
		SignID:           c.SignID,
		Sandbox:          c.Sandbox,
		Verbose:          c.Verbose,
		Force:            c.Force,
		Race:             c.Race,
		Log:              printVerbose,
	}

	if err := pkg.Run(ctx); err != nil {
		fail("%s", err)
	}
}

func cleanMac(ctx context.Context, args []string) {
	c := macCleanConfig{}

	ld := conf.Loader{
		Name:    "mac clean",
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

	pkg := MacPackage{
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

// MacPackage represents a package for a MacOS app.
// It implements the Package interface.
type MacPackage struct {
	// The path where the sources are.
	// It must refer to a Go main package.
	// Default is ".".
	Sources string

	// The path where the package is saved.
	// If not set, the ".app" extension is added.
	Output string

	// The version on MacOS the build is for.
	DeploymentTarget string

	// The signing identifier to sign the package.
	SignID string

	// Configure the app to run in sandbox mode.
	Sandbox bool

	// Creates a .pkg to be uploaed on the App Store.
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
	contentsDir         string
	macOSDir            string
	resourcesDir        string
	executable          string
	settings            macSettings
}

// Init satisfies the Package interface.
func (pkg *MacPackage) Init(ctx context.Context) error {
	if err := pkg.init(); err != nil {
		return err
	}

	pkg.Log("creating resources directory")
	if err := os.MkdirAll(filepath.Join(pkg.Sources, "resources", "css"), 0755); err != nil {
		return err
	}

	pkg.Log("installing Xcode command line tools")
	execute(ctx, "xcode-select", "--install")
	return nil
}

func (pkg *MacPackage) init() (err error) {
	if runtime.GOOS != "darwin" {
		return errors.New("operating system is not MacOS")
	}

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
	if !strings.HasSuffix(pkg.Output, ".app") {
		pkg.Output += ".app"
	}

	pkg.name = filepath.Base(pkg.Output)

	if pkg.workingDir, err = os.Getwd(); err != nil {
		return err
	}

	pkg.sourcesResourcesDir = filepath.Join(pkg.Sources, "resources")

	if pkg.tmpDir = os.Getenv("TMPDIR"); len(pkg.tmpDir) == 0 {
		return errors.New("tmp dir not set")
	}
	pkg.tmpDir = filepath.Join(pkg.tmpDir, "goapp", execName)
	pkg.tmpExecutable = filepath.Join(pkg.tmpDir, execName)

	pkg.contentsDir = filepath.Join(pkg.Output, "Contents")
	pkg.macOSDir = filepath.Join(pkg.Output, "Contents", "MacOS")
	pkg.resourcesDir = filepath.Join(pkg.Output, "Contents", "Resources")
	pkg.executable = filepath.Join(pkg.Output, "Contents", "MacOS", execName)
	return nil
}

// Build satisfies the Package interface.
func (pkg *MacPackage) Build(ctx context.Context) error {
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

	pkg.Log("syncing resources")
	if err := pkg.syncResources(); err != nil {
		return err
	}

	pkg.Log("generating icons")
	if err := pkg.generateIcons(ctx); err != nil {
		return err
	}

	pkg.Log("generating Info.plist")
	if err := pkg.generateInfoPlist(); err != nil {
		return err
	}

	if len(pkg.SignID) == 0 {
		if pkg.AppStore {
			return errors.New("app store requires a sign id")
		}

		return nil
	}

	pkg.Log("signing %s", pkg.name)
	if err := pkg.signing(ctx); err != nil {
		return err
	}

	if !pkg.Verbose {
		return nil
	}

	pkg.Log("packing for app store", pkg.name)
	return pkg.packForAppStore(ctx)
}

func (pkg *MacPackage) create() error {
	if err := os.RemoveAll(filepath.Join(pkg.contentsDir, "_CodeSignature")); err != nil {
		return err
	}

	dirs := []string{
		pkg.Output,
		pkg.contentsDir,
		pkg.macOSDir,
		pkg.resourcesDir,
	}

	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return err
		}
	}

	return nil
}

func (pkg *MacPackage) buildExecutable(ctx context.Context) error {
	os.Setenv("MACOSX_DEPLOYMENT_TARGET", pkg.DeploymentTarget)

	args := []string{
		"go", "build",
		"-ldflags", "-s -X github.com/murlokswarm/app.Kind=desktop",
		"-o", pkg.tmpExecutable,
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

	if err := execute(ctx, args[0], args[1:]...); err != nil {
		return err
	}

	return file.Copy(pkg.executable, pkg.tmpExecutable)
}

func (pkg *MacPackage) readSettings(ctx context.Context) error {
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

	s := macSettings{}
	if err = json.Unmarshal(b, &s); err != nil {
		return err
	}

	s.Executable = filepath.Base(pkg.executable)
	s.Name = stringWithDefault(s.Name, s.Executable)
	s.ID = stringWithDefault(s.ID, fmt.Sprintf("%v.%v", os.Getenv("USER"), s.Name))
	s.URLScheme = stringWithDefault(s.URLScheme, strings.ToLower(s.Name))
	s.Version = stringWithDefault(s.Version, "1.0")
	s.BuildNumber = intWithDefault(s.BuildNumber, 1)
	s.Icon = stringWithDefault(s.Icon, "logo.png")
	s.DevRegion = stringWithDefault(s.DevRegion, "en")
	s.Category = stringWithDefault(s.Category, "public.app-category.developer-tools")
	s.Copyright = stringWithDefault(s.Copyright, fmt.Sprintf("Copyright © %v %s. All rights reserved.",
		time.Now().Year(),
		os.Getenv("USER"),
	))
	s.DeploymentTarget = pkg.DeploymentTarget
	s.Sandbox = pkg.Sandbox

	if s.Sandbox && len(pkg.SignID) == 0 {
		return errors.New("sandbox requires a sign id")
	}

	if err = validateMacFileTypes(s.SupportedFiles...); err != nil {
		return err
	}

	pkg.settings = s

	if b, err = json.MarshalIndent(s, "", "    "); err != nil {
		return err
	}

	pkg.Log("settings: %s", b)
	return nil
}

func (pkg *MacPackage) syncResources() error {
	return file.Sync(pkg.resourcesDir, pkg.sourcesResourcesDir)
}

func (pkg *MacPackage) generateIcons(ctx context.Context) error {
	appIcon := filepath.Join(pkg.resourcesDir, pkg.settings.Icon)
	if _, err := os.Stat(appIcon); os.IsNotExist(err) {
		file.Copy(appIcon, file.RepoPath("logo.png"))
	}

	icons := []string{
		appIcon,
	}

	for _, i := range icons {
		if err := pkg.generateIcon(ctx, i); err != nil {
			return errors.Wrapf(err, "generating icon for %q failed", i)
		}
	}

	return nil
}

func (pkg *MacPackage) generateIcon(ctx context.Context, path string) error {
	if _, err := os.Stat(path); err != nil {
		return err
	}

	iconset := trimExt(path) + ".iconset"

	if err := os.Mkdir(iconset, 0755); err != nil {
		return err
	}
	defer os.RemoveAll(iconset)

	retinaIcon := func(w, h, s int) string {
		return filepath.Join(iconset, fmt.Sprintf("icon_%vx%v@%vx.png", w, h, s))
	}

	standardIcon := func(w, h int) string {
		return filepath.Join(iconset, fmt.Sprintf("icon_%vx%v.png", w, h))
	}

	if err := generateIcons(path, []iconInfo{
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

func (pkg *MacPackage) generateInfoPlist() error {
	pkg.settings.Icon = trimExt(pkg.settings.Icon)

	for i, f := range pkg.settings.SupportedFiles {
		f.Icon = trimExt(f.Icon)
		pkg.settings.SupportedFiles[i] = f
	}

	plist := filepath.Join(pkg.contentsDir, "Info.plist")
	return generateTemplatedFile(plist, infoPlistTmpl, pkg.settings)
}

func (pkg *MacPackage) signing(ctx context.Context) error {
	entitlements := filepath.Join(pkg.tmpDir, ".entitlements")

	if err := generateTemplatedFile(entitlements, entitlementsPlistTmpl, pkg.settings); err != nil {
		return err
	}
	defer os.RemoveAll(entitlements)

	cmd := []string{"codesign"}

	if pkg.Verbose {
		cmd = append(cmd, "--verbose")
	}

	cmd = append(cmd,
		"--force",
		"--sign",
		pkg.SignID,
		"--entitlements",
		entitlements,
		pkg.Output,
	)

	if err := execute(ctx, cmd[0], cmd[1:]...); err != nil {
		return err
	}

	cmd = []string{"codesign"}

	if pkg.Verbose {
		cmd = append(cmd, "--verbose")
	}

	cmd = append(cmd,
		"--verify",
		"--deep",
		"--strict",
		pkg.Output,
	)

	return execute(ctx, cmd[0], cmd[1:]...)
}

func (pkg *MacPackage) packForAppStore(ctx context.Context) error {
	name := strings.Replace(pkg.Output, ".app", ".pkg", 1)

	return execute(ctx,
		"productbuild",
		"--component",
		pkg.Output,
		"/Applications",
		"--sign",
		pkg.SignID,
		name,
	)
}

// Run satisfies the Package interface.
func (pkg *MacPackage) Run(ctx context.Context) error {
	if err := pkg.Build(ctx); err != nil {
		return err
	}

	if pkg.Verbose {
		os.Setenv("GOAPP_DEBUG", "true")
	}

	pkg.Log("running %s", pkg.name)
	return execute(ctx, pkg.executable)
}

// Clean satisfies the Package interface.
func (pkg *MacPackage) Clean(ctx context.Context) error {
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

type macSettings struct {
	Executable       string `json:",omitempty"`
	Name             string `json:",omitempty"`
	ID               string `json:",omitempty"`
	URLScheme        string `json:",omitempty"`
	Version          string `json:",omitempty"`
	BuildNumber      int
	Icon             string        `json:",omitempty"`
	DevRegion        string        `json:",omitempty"`
	DeploymentTarget string        `json:",omitempty"`
	Copyright        string        `json:",omitempty"`
	Category         string        `json:",omitempty"`
	Sandbox          bool          `json:",omitempty"`
	Background       bool          `json:",omitempty"`
	Server           bool          `json:",omitempty"`
	Camera           bool          `json:",omitempty"`
	Microphone       bool          `json:",omitempty"`
	USB              bool          `json:",omitempty"`
	Printers         bool          `json:",omitempty"`
	Bluetooth        bool          `json:",omitempty"`
	Contacts         bool          `json:",omitempty"`
	Location         bool          `json:",omitempty"`
	Calendar         bool          `json:",omitempty"`
	FilePickers      string        `json:",omitempty"`
	Downloads        string        `json:",omitempty"`
	Pictures         string        `json:",omitempty"`
	Music            string        `json:",omitempty"`
	Movies           string        `json:",omitempty"`
	SupportedFiles   []macFileType `json:",omitempty"`
}

type macFileType struct {
	Name string   `json:",omitempty"`
	Role string   `json:",omitempty"`
	Icon string   `json:",omitempty"`
	UTIs []string `json:",omitempty"`
}

func validateMacFileTypes(fileTypes ...macFileType) error {
	for i, f := range fileTypes {
		if len(f.Name) == 0 {
			return errors.Errorf("file type at index %v: name is not set", i)
		}

		if len(f.Icon) > 0 && filepath.Ext(f.Icon) != ".png" {
			return errors.Errorf(`file type at index %v: icon is not a ".png"`, i)
		}

		if len(f.UTIs) == 0 {
			return errors.Errorf("file type at index %v: no uti", i)
		}

		for j, u := range f.UTIs {
			if len(u) == 0 {
				return errors.Errorf("file type at index %v: uti at index %v: uti not set", i, j)
			}
		}

		f.Role = stringWithDefault(f.Role, "None")
		fileTypes[i] = f
	}

	return nil
}

func macOSVersion() string {
	return executeString("sw_vers", "-productVersion")
}
