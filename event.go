package app

import (
	"github.com/murlokswarm/app/key"
)

// MouseEvent represents an onmouse event arg.
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

// WheelEvent represents an onwheel event arg.
type WheelEvent struct {
	DeltaX    float64
	DeltaY    float64
	DeltaZ    float64
	DeltaMode DeltaMode
}

// DeltaMode is an indication of the units of measurement for a delta value.
type DeltaMode uint64

// KeyboardEvent represents an onkey event arg.
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

// DragAndDropEvent represents an ondrop event arg.
type DragAndDropEvent struct {
	Files         []string
	Data          string
	DropEffect    string
	EffectAllowed string
}
