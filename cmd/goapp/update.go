package main

import (
	"context"
	"os"

	"github.com/segmentio/conf"
)

type updateConfig struct {
	Verbose bool `conf:"v" help:"Enable verbose mode."`
}

func update(ctx context.Context, args []string) {
	c := initConfig{}

	ld := conf.Loader{
		Name:    "goapp update",
		Args:    args,
		Usage:   "[options...]",
		Sources: []conf.Source{conf.NewEnvSource("GOAPP", os.Environ()...)},
	}

	_, args = conf.LoadWith(&c, ld)
	verbose = c.Verbose

	cmd := []string{
		"go", "get", "-u",
	}

	if verbose {
		cmd = append(cmd, "-v")
	}

	cmd = append(cmd, "github.com/maxence-charriere/app/cmd/goapp")

	log("updating to the latest version")
	if err := execute(ctx, cmd[0], cmd[1:]...); err != nil {
		fail("%s", err)
	}

	success("initialization succeeded")
}
