package main

import (
	"bufio"
	"context"
	"io"
	"os"
	"os/exec"
)

func execute(ctx context.Context, cmd string, args ...string) error {
	command := exec.CommandContext(ctx, cmd, args...)

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
			printErr("%s", err)
			continue
		}

		if verbose {
			output.Write([]byte("    "))
		}

		output.Write(b[:n])
	}
}
