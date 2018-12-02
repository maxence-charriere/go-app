package logs

import (
	"context"
	"io"
	"net"
	"sync"

	"github.com/murlokswarm/app"
)

// GoappClient represents a client that send logs to goapp.
type GoappClient struct {
	conn   net.Conn
	logger Logger
}

// NewGoappClient creates a goapp log client that connects to the goapp log
// server with the given address.
func NewGoappClient(addr string, prompt func(Logger) Logger) *GoappClient {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		app.Panic("connection to goapp failed:", err)
	}

	return &GoappClient{
		conn:   conn,
		logger: prompt(ToWriter(conn)),
	}
}

// Logger returns the logger that produces logs.
func (c *GoappClient) Logger() Logger {
	return c.logger
}

// WaitForStop is waiting for goapp to send a stop signal in order to call
// the given stop function.
func (c *GoappClient) WaitForStop(stop func()) {
	sig := make([]byte, 1)
	_, err := c.conn.Read(sig)

	if err != nil {
		app.Panic("reading goapp stop signal failed:", err)
	}

	stop()
}

// Close close the established connection.
func (c *GoappClient) Close() error {
	return c.conn.Close()
}

// GoappServer represents a server that receives and prints logs from an
// app that is run by goapp.
type GoappServer struct {
	// The address used for listen connections.
	Addr string

	// The writer where the logs are written.
	Writer io.Writer
}

// ListenAndLog listen app log and write them to the server writer.
func (s *GoappServer) ListenAndLog(ctx context.Context) error {
	conn, err := net.ListenPacket("udp", s.Addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	var once sync.Once

	addrc := make(chan net.Addr)
	errc := make(chan error)

	go func() {
		defer close(addrc)
		defer close(errc)

		for {
			log := make([]byte, 1024)

			n, addr, err := conn.ReadFrom(log)
			if err != nil {
				errc <- err
				return
			}

			once.Do(func() {
				addrc <- addr
			})

			if _, err = s.Writer.Write(log[:n]); err != nil {
				errc <- err
				return
			}
		}
	}()

	var addr net.Addr

	select {
	case addr = <-addrc:
	case err = <-errc:
	}

	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		_, err = conn.WriteTo([]byte("q"), addr)

	case err = <-errc:
	}

	return err
}
