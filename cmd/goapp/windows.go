// +build windows

package main

import "github.com/segmentio/conf"

func commands() []conf.Command {
	return []conf.Command{
		{Name: "web", Help: "Build app for web."},
		{Name: "help", Help: "Show the help."},
	}
}

func openCommand() string {
	return "explorer"
}
