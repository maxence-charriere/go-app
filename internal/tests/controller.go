package tests

import (
	"testing"

	"github.com/murlokswarm/app"
	"github.com/stretchr/testify/assert"
)

func testController(t *testing.T, c app.Controller) {
	c.Close()
	assert.Error(t, c.Err())
}
