// +build darwin,amd64

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

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

	fmt.Println("bundle configuration:", string(data))

	var bundle driver.Bundle
	if err := json.Unmarshal(data, &bundle); err != nil {
		return err
	}

	return nil
}

func openCommand() string {
	return "open"
}
