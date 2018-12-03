package file

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCaptureOutput(t *testing.T) {
	result := bytes.Buffer{}
	expected := bytes.Buffer{}

	cancel, err := CaptureOutput(&result)
	require.NoError(t, err)

	fmt.Fprintln(os.Stdout, "hello")
	time.Sleep(time.Millisecond * 5)
	fmt.Fprintln(os.Stderr, "world")
	fmt.Fprint(&expected, "helloworld")

	time.Sleep(time.Millisecond * 100)

	cancel()
	require.Equal(t, expected.String(), result.String())
}

func TestHTTPWriter(t *testing.T) {
	var wg sync.WaitGroup
	var result []byte
	var err error

	wg.Add(1)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/close" {
			wg.Done()
			return
		}

		result, err = ioutil.ReadAll(r.Body)
	}))
	defer ts.Close()

	w := &HTTPWriter{
		URL:    ts.URL,
		Client: http.DefaultClient,
	}

	expected := []byte(fmt.Sprint("hello"))
	n, werr := w.Write(expected)
	require.NoError(t, werr)
	require.Equal(t, len(expected), n)
	require.NoError(t, w.Close())

	wg.Wait()

	require.NoError(t, err)
	require.Equal(t, expected, result)
}
