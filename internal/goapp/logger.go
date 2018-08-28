package goapp

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/murlokswarm/app"
)

const (
	// Addr is the address to use with a logger or a server.
	Addr = ":7042"
)

// NewLogger creates a logger that send logs to a log server.
func NewLogger(debug, colors bool, stop func()) app.Logger {
	conn, err := net.Dial("udp", Addr)
	if err != nil {
		panic(err)
	}

	r, w := io.Pipe()
	logc := make(chan []byte)
	stopc := make(chan struct{})

	readLogs := func() {
		for {
			r := bufio.NewReader(r)

			log, err := r.ReadBytes('\n')
			if err != nil {
				return
			}

			logc <- log
		}
	}

	listenForStop := func() {
		r := bufio.NewReader(conn)
		r.ReadByte()
		stopc <- struct{}{}
	}

	sendLogs := func() {
		defer conn.Close()

		for {
			select {
			case log := <-logc:
				if _, err := conn.Write(log); err != nil {
					fmt.Println("client ->", err)
					return
				}

			case <-stopc:
				stop()
			}
		}
	}

	go listenForStop()
	go readLogs()
	go sendLogs()

	return app.NewLogger(w, w, debug, colors)
}

// ListenAndWriteLogs listens for logs and write them on standard output.
func ListenAndWriteLogs(ctx context.Context) error {
	conn, err := net.ListenPacket("udp", ":7042")
	if err != nil {
		return err
	}
	defer conn.Close()

	var appAddr net.Addr
	var once sync.Once

	readLogs := func() {
		for {
			log := make([]byte, 1024)

			n, addr, err := conn.ReadFrom(log)
			if err != nil {
				return
			}

			once.Do(func() {
				appAddr = addr
			})

			fmt.Print(string(log[:n]))
		}
	}

	go readLogs()

	<-ctx.Done()
	_, err = conn.WriteTo([]byte("bye"), appAddr)
	return err
}
