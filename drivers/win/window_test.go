// +build windows

package win

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourcesDir(t *testing.T) {
	assert.Equal(t, "ms-appx-web:///Resources/hello/world", resourcesDir("hello", "world"))
	assert.Equal(t, "ms-appx-web:///Resources/hello", resourcesDir("/hello"))
}
