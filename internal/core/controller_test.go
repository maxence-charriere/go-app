package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestController(t *testing.T) {
	c := &Controller{}
	c.Close()
	assert.Error(t, c.Err())
}
