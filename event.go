package app

// ChangeArg represents the data passed in a onchange event.
type ChangeArg struct {
	Value  string
	Target DOMElement
}

// MouseArg represents data fired when interacting
// with a pointing device (such as a mouse).
type MouseArg struct {
	ClientX  float64
	ClientY  float64
	PageX    float64
	PageY    float64
	ScreenX  float64
	ScreenY  float64
	Button   int
	Detail   int
	AltKey   bool
	CtrlKey  bool
	MetaKey  bool
	ShiftKey bool
	Target   DOMElement
}

// WheelArg represents data fired when a wheel button of a
// pointing device (usually a mouse) is rotated.
type WheelArg struct {
	DeltaX    float64
	DeltaY    float64
	DeltaZ    float64
	DeltaMode DeltaMode
	Target    DOMElement
}

// DeltaMode is an indication of the units of measurement for a delta value.
type DeltaMode uint64

// KeyboardArg represents data fired when the keyboard is used.
type KeyboardArg struct {
	CharCode rune
	KeyCode  KeyCode
	Location KeyLocation
	AltKey   bool
	CtrlKey  bool
	MetaKey  bool
	ShiftKey bool
	Target   DOMElement
}

// EventArg represents the data passed in events.
type EventArg struct {
	Target DOMElement
}
