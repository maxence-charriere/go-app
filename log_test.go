package app

import (
	"bytes"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	buffer := &bytes.Buffer{}
	logger := NewLogger(buffer, buffer, true, false)

	logger.Log("a message")
	logger.Log("a message with args: %v", 42)
	logger.Log("a message with line return\nhere is")
	logger.Log("an error: %s", errors.New("error"))
	logger.Debug("a debug message")
	logger.Debug("a debug message with args: %v", 42)

	logger.WhenDebug(func() {
		logger.Debug("yoda is strong")
	})
	assert.Contains(t, buffer.String(), "yoda is strong")

	logger = NewLogger(buffer, buffer, false, true)
	logger.WhenDebug(func() {
		logger.Debug("vader is strong")
	})
	assert.NotContains(t, buffer.String(), "vader is strong")
	logger.Log("a message")
	logger.Log("an error: %s", errors.New("error"))
	logger.Debug("a debug message")

	t.Log(buffer.String())
}
