// Package ui provides a set of components to organize an application layout.
package ui

import (
	"strconv"
)

var (
	// The padding of block-like components in px.
	BlockPadding = 30

	// The padding of block-like components in px when app width is <= 480px.
	BlockMobilePadding = 18

	// The content width of block-like components in px.
	BlockContentWidth = 580

	// The horizontal padding of base-like components in px.
	BaseHPadding = 36

	// The horizontal padding of base-like components in px when app width is <= 480px.
	BaseMobileHPadding = 12

	// The horizontal padding of base-like ad components in px.
	BaseAdHPadding = BaseHPadding / 2

	// The vertical padding of base-like components in px.
	BaseVPadding = 12

	// The default icon size in px.
	DefaultIconSize = 24

	// The default icon space.
	DefaultIconSpace = 6

	// The default width for flow items in px.
	DefaultFlowItemWidth = 372
)

const (
	defaultHeaderHeight = 90
)

func pxToString(px int) string {
	return strconv.Itoa(px) + "px"
}

type alignment int

const (
	stretch alignment = iota
	top
	right
	bottom
	left
	middle
)

type style struct {
	key   string
	value string
}
