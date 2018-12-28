// +build darwin

package objc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRPC(t *testing.T) {
	macr, gor := RPC(func(func()) {})
	assert.NotNil(t, macr)
	assert.NotNil(t, gor)
}
