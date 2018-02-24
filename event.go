package app

import (
	"github.com/murlokswarm/app/key"
)

// MouseEvent represents data fired when interacting with a pointing device
// (such as a mouse).
type MouseEvent struct {
	ClientX   float64
	ClientY   float64
	PageX     float64
	PageY     float64
	ScreenX   float64
	ScreenY   float64
	Button    int
	Detail    int
	AltKey    bool
	CtrlKey   bool
	MetaKey   bool
	ShiftKey  bool
	InnerText string
}

// WheelEvent represents data fired when a wheel button of a pointing device
// (usually a mouse) is rotated.
type WheelEvent struct {
	DeltaX    float64
	DeltaY    float64
	DeltaZ    float64
	DeltaMode DeltaMode
}

// DeltaMode is an indication of the units of measurement for a delta value.
type DeltaMode uint64

// KeyboardEvent represents data fired when the keyboard is used.
type KeyboardEvent struct {
	CharCode  rune
	KeyCode   key.Code
	Location  key.Location
	AltKey    bool
	CtrlKey   bool
	MetaKey   bool
	ShiftKey  bool
	InnerText string
}
