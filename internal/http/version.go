package http

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"
)

// GenerateEtag generates an etag.
func GenerateEtag() string {
	t := time.Now().UTC().String()
	return fmt.Sprintf(`%x`, sha1.Sum([]byte(t)))
}

// GetEtag returns the etag for the given web directory.
func GetEtag(webDir string) string {
	filename := filepath.Join(webDir, ".etag")

	etag, err := ioutil.ReadFile(filename)
	if err != nil {
		return ""
	}
	return string(etag)
}
