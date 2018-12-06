package file

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Matches reports whether given files match.
func Matches(l, r string) bool {
	li, err := os.Stat(l)
	if err != nil {
		return false
	}

	var ri os.FileInfo
	if ri, err = os.Stat(r); err != nil {
		return false
	}

	return li.Name() == ri.Name() &&
		li.Size() == ri.Size() &&
		li.Mode() == ri.Mode() &&
		li.ModTime().Equal(ri.ModTime()) &&
		li.IsDir() == ri.IsDir()

}

// Copy copy src into dst.
func Copy(dst, src string) error {
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()

	var sinfo os.FileInfo
	if sinfo, err = sf.Stat(); err != nil {
		return err
	}

	var df *os.File
	if df, err = os.OpenFile(dst, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, sinfo.Mode()); err != nil {
		return err
	}
	defer df.Close()

	if _, err = io.Copy(df, sf); err != nil {
		return err
	}

	if err = df.Close(); err != nil {
		return err
	}

	return os.Chtimes(dst, time.Now(), sinfo.ModTime())
}

// Sync synchronize src and dst.
func Sync(dst, src string) error {
	dst = filepath.Clean(dst)
	src = filepath.Clean(src)

	walkCopy := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, path[len(src):])

		if info.IsDir() {
			var dstInfo os.FileInfo
			dstInfo, err = os.Stat(dstPath)

			if os.IsNotExist(err) {
				return os.Mkdir(dstPath, info.Mode())
			}

			if err != nil {
				return err
			}

			if !dstInfo.IsDir() {
				os.Remove(dstPath)
				return os.Mkdir(dstPath, info.Mode())
			}

			return nil
		}

		if !Matches(path, dstPath) {
			if err = os.RemoveAll(dstPath); err != nil {
				return err
			}
			return Copy(dstPath, path)
		}
		return nil
	}

	if err := filepath.Walk(src, walkCopy); err != nil {
		return err
	}

	walkRemove := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		srcPath := strings.Replace(path, dst, src, 1)

		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			return os.RemoveAll(path)
		}
		return nil
	}

	return filepath.Walk(dst, walkRemove)
}

// Filenames returns the filenames within the directory and its subdirectories
// that have the given extensions.
func Filenames(dirname string, extensions ...string) []string {
	var filenames []string

	exts := make(map[string]struct{})
	for _, ext := range extensions {
		exts[ext] = struct{}{}
	}

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if _, ok := exts[filepath.Ext(path)]; !ok {
			return nil
		}

		filenames = append(filenames, path)
		return nil
	}

	filepath.Walk(dirname, walker)
	return filenames
}

// RepoPath returns the app package repository path.
func RepoPath(p ...string) string {
	path := []string{
		os.Getenv("GOPATH"),
		"src",
		"github.com",
		"murlokswarm",
		"app",
	}

	path = append(path, p...)
	return filepath.Join(path...)
}
