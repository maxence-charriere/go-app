package file

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCSS(t *testing.T) {
	dir := "css-test"
	os.MkdirAll(dir, 0777)
	defer os.RemoveAll(dir)

	assert.Len(t, CSS(dir), 0)

	os.MkdirAll(filepath.Join(dir, "sub"), 0777)
	os.Create(filepath.Join(dir, "test.css"))
	os.Create(filepath.Join(dir, "test.scss"))
	os.Create(filepath.Join(dir, "sub", "sub.css"))

	assert.Contains(t, CSS(dir), filepath.Join(dir, "test.css"))
	assert.NotContains(t, CSS(dir), filepath.Join(dir, "test.scss"))
	assert.Contains(t, CSS(dir), filepath.Join(dir, "sub", "sub.css"))
}
