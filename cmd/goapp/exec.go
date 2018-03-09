package main

import (
	"bufio"
	"io"
	"os"
	"os/exec"

	"github.com/murlokswarm/app"
)

func execute(cmd string, args ...string) error {
	command := exec.Command(cmd, args...)

	cmdout, err := command.StdoutPipe()
	if err != nil {
		return err
	}

	cmderr, err := command.StderrPipe()
	if err != nil {
		return err
	}

	go printOutput(cmdout, os.Stdout)
	go printOutput(cmderr, os.Stderr)

	if err = command.Start(); err != nil {
		return err
	}

	err = command.Wait()
	return err
}

func printOutput(r io.Reader, output io.Writer) {
	reader := bufio.NewReader(r)
	b := make([]byte, 1024)

	for {
		n, err := reader.Read(b)
		if err == io.EOF {
			return
		}
		if err != nil {
			app.Error(err)
			continue
		}
		output.Write(b[:n])
	}
}
