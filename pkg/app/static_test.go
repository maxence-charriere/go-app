//go:build !wasm
// +build !wasm

package app

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"

	"github.com/maxence-charriere/go-app/v11/pkg/errors"
	"github.com/stretchr/testify/require"
)

type staticTestCompo struct {
	Compo
	path string
}

func (c *staticTestCompo) Render() UI {
	return Div().Text("Static page: " + c.path)
}

func registerStaticTestRoutes(t testing.TB, paths ...string) {
	for _, path := range paths {
		routes.mu.Lock()
		previous, existed := routes.routes[path]
		routes.routes[path] = func() Composer { return &staticTestCompo{path: path} }
		routes.mu.Unlock()
		t.Cleanup(func() {
			routes.mu.Lock()
			defer routes.mu.Unlock()
			if existed {
				routes.routes[path] = previous
			} else {
				delete(routes.routes, path)
			}
		})
	}
}

func TestGenerateStaticWebsite(t *testing.T) {
	dir := t.TempDir()
	registerStaticTestRoutes(t, "/", "/hello", "/world", "/nested/foo", "/custom.html")

	err := GenerateStaticWebsite(dir,
		&Handler{
			Name:      "Static Go-app",
			Title:     "Static test",
			Resources: GitHubPages("go-app"),
		},
		"/hello",
		"world",
		"/nested/foo",
	)
	require.NoError(t, err)

	files := []string{
		filepath.Join(dir),
		filepath.Join(dir, "web"),
		filepath.Join(dir, "index.html"),
		filepath.Join(dir, "wasm_exec.js"),
		filepath.Join(dir, "app.js"),
		filepath.Join(dir, "app-worker.js"),
		filepath.Join(dir, "manifest.webmanifest"),
		filepath.Join(dir, "app.css"),
		filepath.Join(dir, "hello.html"),
		filepath.Join(dir, "world.html"),
		filepath.Join(dir, "nested", "foo.html"),
		filepath.Join(dir, "custom.html"),
	}

	for _, f := range files {
		t.Run(f, func(t *testing.T) {
			_, err := os.Stat(f)
			require.NoError(t, err)
		})
	}
	for _, path := range []string{"index", "hello", "world", "nested/foo", "custom"} {
		body, err := os.ReadFile(filepath.Join(dir, path+".html"))
		require.NoError(t, err)
		require.Contains(t, string(body), "Static page:")
		require.Contains(t, string(body), "<title>Static test</title>")
	}
}

func TestCreateStaticPageStatus(t *testing.T) {
	for _, status := range []int{http.StatusOK, http.StatusNoContent, http.StatusPartialContent, http.StatusNotFound, http.StatusInternalServerError} {
		t.Run(strconv.Itoa(status), func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(status)
				io.WriteString(w, "response body")
			}))
			defer s.Close()
			body, err := createStaticPage(s.URL + "/page")
			if status == http.StatusOK {
				require.NoError(t, err)
				require.Equal(t, "response body", string(body))
			} else {
				require.Error(t, err)
				require.Empty(t, body)
				require.Equal(t, status, errors.Tag(err, "status-code"))
				require.Equal(t, s.URL+"/page", errors.Tag(err, "path"))
			}
		})
	}
}

func TestGenerateStaticWebsitePreservesFailedPage(t *testing.T) {
	for _, exists := range []bool{false, true} {
		t.Run(fmt.Sprint(exists), func(t *testing.T) {
			dir := t.TempDir()
			filename := filepath.Join(dir, "missing-static-page.html")
			if exists {
				require.NoError(t, os.WriteFile(filename, []byte("previous output"), 0600))
			}
			err := GenerateStaticWebsite(dir, &Handler{}, "/missing-static-page")
			require.Error(t, err)
			require.Equal(t, http.StatusNotFound, errors.Tag(err, "status-code"))
			require.Equal(t, "/missing-static-page", errors.Tag(err, "path"))
			body, readErr := os.ReadFile(filename)
			if exists {
				require.NoError(t, readErr)
				require.Equal(t, "previous output", string(body))
			} else {
				require.True(t, os.IsNotExist(readErr), "failed pages must not create output files")
			}
		})
	}
}

func TestCreateStaticDirFailures(t *testing.T) {
	t.Run("existing file", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "file")
		require.NoError(t, os.WriteFile(path, []byte("existing file"), 0600))
		require.Error(t, createStaticDir(path, ""))
		body, err := os.ReadFile(path)
		require.NoError(t, err)
		require.Equal(t, "existing file", string(body))
	})
	t.Run("stat failure", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "loop")
		if err := os.Symlink("loop", path); err != nil {
			t.Skipf("symlinks unavailable: %v", err)
		}
		require.Error(t, createStaticDir(path, ""))
	})
}

func TestGenerateStaticWebsitePreservesFailedResource(t *testing.T) {
	for _, status := range []int{http.StatusInternalServerError, http.StatusOK} {
		t.Run(strconv.Itoa(status), func(t *testing.T) {
			dir := t.TempDir()
			require.NoError(t, os.Mkdir(filepath.Join(dir, "web"), 0755))
			filename := filepath.Join(dir, "web", "failure.txt")
			require.NoError(t, os.WriteFile(filename, []byte("previous output"), 0600))
			h := &Handler{Resources: forwardingTestResolver{http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if status == http.StatusOK {
					w.Header().Set("Content-Length", "100") // Force an incomplete response body.
				}
				w.WriteHeader(status)
				io.WriteString(w, "failed response")
			})}}
			err := GenerateStaticWebsite(dir, h, "/web/failure.txt")
			require.Error(t, err)
			require.Equal(t, "/web/failure.txt", errors.Tag(err, "path"))
			if status == http.StatusOK {
				require.True(t, errors.Is(err, io.ErrUnexpectedEOF), "%v", err)
			} else {
				require.Equal(t, status, errors.Tag(err, "status-code"))
			}
			body, err := os.ReadFile(filename)
			require.NoError(t, err)
			require.Equal(t, "previous output", string(body))
		})
	}
}

func TestGenerateStaticWebsiteFileCreationFailure(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(dir, "index.html"), 0755))
	err := GenerateStaticWebsite(dir, &Handler{})
	require.Error(t, err)
	require.Equal(t, "/", errors.Tag(err, "path"))
	require.Equal(t, "/index.html", errors.Tag(err, "filename"))
}

type staticPageTestWriter struct {
	n          int
	writeErr   error
	closeErr   error
	closeCalls int
}

func (w *staticPageTestWriter) Write(page []byte) (int, error) {
	return min(w.n, len(page)), w.writeErr
}

func (w *staticPageTestWriter) Close() error {
	w.closeCalls++
	return w.closeErr
}

func TestWriteStaticPage(t *testing.T) {
	page := []byte("page")
	for _, test := range []struct {
		name     string
		n        int
		writeErr error
		closeErr error
	}{
		{"success", len(page), nil, nil},
		{"write failure", 2, io.ErrClosedPipe, nil},
		{"short write", 2, nil, nil},
		{"close failure", len(page), nil, os.ErrInvalid},
		{"write and close failure", 2, io.ErrClosedPipe, os.ErrInvalid},
	} {
		t.Run(test.name, func(t *testing.T) {
			w := &staticPageTestWriter{n: test.n, writeErr: test.writeErr, closeErr: test.closeErr}
			err := writeStaticPage(w, page)
			require.Equal(t, 1, w.closeCalls, "files must be closed even after a write failure")
			if test.writeErr == nil && test.closeErr == nil && test.n == len(page) {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			if test.writeErr != nil {
				require.True(t, errors.Is(err, test.writeErr), "%v", err)
			} else if test.n < len(page) {
				require.True(t, errors.Is(err, io.ErrShortWrite), "%v", err)
			}
			if test.closeErr != nil {
				require.True(t, errors.Is(err, test.closeErr), "%v", err)
			}
		})
	}
}

func staticTestFDDirectory(t testing.TB) string {
	t.Helper()
	for _, path := range []string{"/proc/self/fd", "/dev/fd"} {
		if _, err := os.ReadDir(path); err == nil {
			return path
		}
	}
	t.Skip("file descriptor counts unavailable")
	return ""
}

func TestGenerateStaticWebsiteBoundsOpenFiles(t *testing.T) {
	fdDirectory := staticTestFDDirectory(t)
	paths := make([]string, 128)
	for i := range paths {
		paths[i] = fmt.Sprintf("/static-many/%d", i)
	}
	registerStaticTestRoutes(t, paths...)
	before, err := os.ReadDir(fdDirectory)
	require.NoError(t, err)
	peak := len(before)
	h := &Handler{HTML: func() HTMLHtml {
		fds, err := os.ReadDir(fdDirectory)
		if err != nil {
			t.Errorf("reading descriptors: %v", err)
		}
		peak = max(peak, len(fds))
		return Html()
	}}
	dir := t.TempDir()
	require.NoError(t, GenerateStaticWebsite(dir, h))
	require.LessOrEqual(t, peak-len(before), 12, "output files must be closed between pages")
	for _, path := range paths {
		body, err := os.ReadFile(filepath.Join(dir, path+".html"))
		require.NoError(t, err)
		require.Contains(t, string(body), "Static page: "+path)
	}
}

func BenchmarkGenerateStaticWebsite(b *testing.B) {
	for _, count := range []int{100, 1000} {
		b.Run(strconv.Itoa(count), func(b *testing.B) {
			paths := make([]string, count)
			for i := range paths {
				paths[i] = fmt.Sprintf("/static-benchmark/%d", i)
			}
			registerStaticTestRoutes(b, paths...)
			dir := b.TempDir()
			fdDirectory := staticTestFDDirectory(b)
			fds, err := os.ReadDir(fdDirectory)
			require.NoError(b, err)
			baseFDs, peakFDs := len(fds), len(fds)
			var memory runtime.MemStats
			runtime.ReadMemStats(&memory)
			baseHeap, peakHeap := memory.HeapAlloc, memory.HeapAlloc
			h := &Handler{HTML: func() HTMLHtml {
				fds, err := os.ReadDir(fdDirectory)
				if err == nil {
					peakFDs = max(peakFDs, len(fds))
				}
				runtime.ReadMemStats(&memory)
				peakHeap = max(peakHeap, memory.HeapAlloc)
				return Html()
			}}
			require.NoError(b, GenerateStaticWebsite(dir, h))
			body, err := os.ReadFile(filepath.Join(dir, "static-benchmark", "0.html"))
			require.NoError(b, err)
			require.Contains(b, string(body), "Static page: /static-benchmark/0")
			// Sample peaks separately so descriptor/memory inspection is not timed.
			h.HTML = func() HTMLHtml { return Html() }
			b.ReportAllocs()
			for b.Loop() {
				if err := GenerateStaticWebsite(dir, h); err != nil {
					b.Fatal(err)
				}
			}
			b.ReportMetric(float64(peakFDs-baseFDs), "peak-extra-fds")
			b.ReportMetric(float64(peakHeap-baseHeap), "sampled-heap-growth-B")
		})
	}
}
