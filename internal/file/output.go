package file

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"os"
	"sync"
)

// CaptureOutput read from os.Stdout and os.Stderr and writes on the given
// writer.
func CaptureOutput(w io.Writer) (func(), error) {
	outr, outw, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	errr, errw, err := os.Pipe()
	if err != nil {
		outw.Close()
		outr.Close()
		return nil, err
	}

	mutex := sync.Mutex{}
	stdout := os.Stdout
	stderr := os.Stderr

	cancel := func() {
		outw.Close()
		outr.Close()
		errw.Close()
		errr.Close()

		os.Stdout = stdout
		os.Stderr = stderr
	}

	os.Stdout = outw
	os.Stderr = errw

	scanOutput := func(f io.Reader) {
		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			line := scanner.Bytes()

			if scanner.Err() != nil {
				return
			}

			mutex.Lock()
			if _, err := w.Write(line); err != nil {
				mutex.Unlock()
				return
			}
			mutex.Unlock()
		}
	}

	go scanOutput(outr)
	go scanOutput(errr)
	return cancel, nil
}

// HTTPWriter is a io.Writer that send bytes to a http server.
type HTTPWriter struct {
	// The url where bytes are sent.
	URL string

	// The http client used to send bytes.
	Client *http.Client
}

func (w *HTTPWriter) Write(p []byte) (n int, err error) {
	res, err := w.Client.Post(w.URL, "text/plain", bytes.NewBuffer(p))
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	return int(res.Request.ContentLength), nil
}

// Close sends a message that request a server to shutdown.
func (w *HTTPWriter) Close() error {
	res, err := w.Client.Get(w.URL + "/close")
	if err != nil {
		return err
	}

	return res.Body.Close()
}
