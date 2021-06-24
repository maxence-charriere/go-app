// Package ui provides a set of components to organize an application layout.
package ui

import (
	"strconv"
)

func pxToString(px int) string {
	return strconv.Itoa(px) + "px"
}
