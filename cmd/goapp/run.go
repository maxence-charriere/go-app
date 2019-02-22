package main

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/segmentio/conf"
)

type runConfig struct {
	Browser string `conf:"b"     help:"The browser to open to display the app. (default, chrome, firefox or safari)"`
	Force   bool   `conf:"force" help:"Force rebuilding of package that are already up-to-date."`
	URL     string `conf:"url"   help:"The URL to load in the browser."`
	Race    bool   `conf:"race"  help:"Enable data race detection."`
	Verbose bool   `conf:"v"     help:"Enable verbose mode."`

	rootDir string
}

func runProject(ctx context.Context, args []string) {
	c := runConfig{
		URL: "http://localhost:3000",
	}

	ld := conf.Loader{
		Name:    "goapp run",
		Args:    args,
		Usage:   "[options...] [package]",
		Sources: []conf.Source{conf.NewEnvSource("GOAPP", os.Environ()...)},
	}

	_, args = conf.LoadWith(&c, ld)
	verbose = c.Verbose

	pkg := "."
	if len(args) != 0 {
		pkg = args[0]
	}

	rootDir, err := filepath.Abs(pkg)
	if err != nil {
		fail("%s", err)
	}
	c.rootDir = rootDir

	if err := build(ctx, buildConfig{
		Force:   c.Force,
		Race:    c.Race,
		Verbose: c.Verbose,
		rootDir: rootDir,
	}); err != nil {
		fail("%s", err)
	}

	if c.Browser != "" {
		go func() {
			log("opening %s browser on %s", c.Browser, c.URL)
			launchBrowser(ctx, c)
		}()
	}

	log("running server")
	if err := runServer(ctx, rootDir); err != nil {
		fail("%s", err)
	}
}

func runServer(ctx context.Context, rootDir string) error {
	serverName := filepath.Base(rootDir) + "-server"
	serverPath := filepath.Join(rootDir, serverName)
	os.Chdir(rootDir)
	return execute(ctx, serverPath)
}

func launchBrowser(ctx context.Context, c runConfig) {
	key := strings.ToLower(c.Browser)

	switch runtime.GOOS {
	case "darwin":
		runMacOSBrowser(ctx, key, c.URL)

	case "linux":
		runLinuxBrowser(ctx, key, c.URL)

	case "windows":
		runWindowsBrowser(ctx, key, c.URL)

	default:
		warn("launching browser on %s is not supported", runtime.GOOS)
	}
}

func runMacOSBrowser(ctx context.Context, key string, url string) {
	var cmd []string

	switch key {
	case "chrome":
		cmd = append(cmd, "open", "-a", "Google Chrome", url)

	case "firefox":
		cmd = append(cmd, "open", "-a", "Firefox", url)

	case "safari":
		cmd = append(cmd, "open", "-a", "Safari", url)

	default:
		cmd = append(cmd, "open", url)
	}

	execute(ctx, cmd[0], cmd[1:]...)
}

func runLinuxBrowser(ctx context.Context, key string, url string) {
	var cmd []string

	switch key {
	case "chrome":
		cmd = append(cmd, "google-chrome", url)

	case "firefox":
		cmd = append(cmd, "firefox", url)

	default:
		cmd = append(cmd, "xdg-open", url)
	}

	execute(ctx, cmd[0], cmd[1:]...)
}

func runWindowsBrowser(ctx context.Context, key string, url string) {
	var cmd []string

	switch key {
	case "chrome":
		cmd = append(cmd, "powershell", "start", "chrome", url)

	case "firefox":
		cmd = append(cmd, "powershell", "start", "firefox", url)

	default:
		cmd = append(cmd, "powershell", "start", url)
	}

	execute(ctx, cmd[0], cmd[1:]...)
}
