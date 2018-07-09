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
		li.ModTime() == ri.ModTime() &&
		li.IsDir() == ri.IsDir()

}

// Copy copy src into dst.
func Copy(dst, src string) error {
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()

	var df *os.File
	if df, err = os.Create(dst); err != nil {
		return err
	}
	defer df.Close()

	if _, err = io.Copy(df, sf); err != nil {
		return err
	}

	if err = df.Sync(); err != nil {
		return err
	}

	var sinfo os.FileInfo
	if sinfo, err = sf.Stat(); err != nil {
		return err
	}

	if err = os.Chmod(dst, sinfo.Mode()); err != nil {
		return err
	}

	return os.Chtimes(dst, time.Now(), sinfo.ModTime())
}

// Sync synchronize src and dst.
func Sync(dst, src string) error {
	walkCopy := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		dstPath := strings.Replace(path, src, dst, 1)

		if info.IsDir() {
			var fi os.FileInfo
			if fi, err = os.Stat(dstPath); err != nil && !os.IsNotExist(err) {
				return err
			}

			if os.IsNotExist(err) {
				return os.Mkdir(dstPath, info.Mode())
			}

			if !fi.IsDir() {
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
