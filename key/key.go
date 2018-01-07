// Package key defines the constants related to keyboard use.
package key

// Code represents a system and implementation dependent numerical
// code identifying the unmodified value of the pressed key.
type Code uint8

// board keys.
const (
	Backspace          Code = 8
	Tab                Code = 9
	Enter              Code = 13
	Shift              Code = 16
	Ctrl               Code = 17
	Alt                Code = 18
	PauseBreak         Code = 19
	CapsLock           Code = 20
	Esc                Code = 27
	Space              Code = 32
	PageUp             Code = 33
	PageDown           Code = 34
	End                Code = 35
	Home               Code = 36
	Left               Code = 37
	Up                 Code = 38
	Right              Code = 39
	Down               Code = 40
	PrintScreen        Code = 44
	Insert             Code = 45
	Delete             Code = 46
	Digit0             Code = 48
	Digit1             Code = 49
	Digit2             Code = 50
	Digit3             Code = 51
	Digit4             Code = 52
	Digit5             Code = 53
	Digit6             Code = 54
	Digit7             Code = 55
	Digit8             Code = 56
	Digit9             Code = 57
	A                  Code = 65
	B                  Code = 66
	C                  Code = 67
	D                  Code = 68
	E                  Code = 69
	F                  Code = 70
	G                  Code = 71
	H                  Code = 72
	I                  Code = 73
	J                  Code = 74
	K                  Code = 75
	L                  Code = 76
	M                  Code = 77
	N                  Code = 78
	O                  Code = 79
	P                  Code = 80
	Q                  Code = 81
	R                  Code = 82
	S                  Code = 83
	T                  Code = 84
	U                  Code = 85
	V                  Code = 86
	W                  Code = 87
	X                  Code = 88
	Y                  Code = 88
	Z                  Code = 90
	Meta               Code = 91
	Menu               Code = 93
	NumPad0            Code = 96
	NumPad1            Code = 97
	NumPad2            Code = 98
	NumPad3            Code = 99
	NumPad4            Code = 100
	NumPad5            Code = 101
	NumPad6            Code = 102
	NumPad7            Code = 103
	NumPad8            Code = 104
	NumPad9            Code = 105
	NumpadMultiply     Code = 106
	NumpadAdd          Code = 107
	NumpadComma        Code = 108
	NumpadSubtract     Code = 109
	NumpadDecimal      Code = 110
	NumpadDivide       Code = 111
	F1                 Code = 112
	F2                 Code = 113
	F3                 Code = 114
	F4                 Code = 115
	F5                 Code = 116
	F6                 Code = 117
	F7                 Code = 118
	F8                 Code = 119
	F9                 Code = 120
	F10                Code = 121
	F11                Code = 122
	F12                Code = 123
	F13                Code = 124
	F14                Code = 125
	F15                Code = 126
	F16                Code = 127
	F17                Code = 128
	F18                Code = 129
	F19                Code = 130
	F20                Code = 131
	F21                Code = 132
	F22                Code = 133
	F23                Code = 134
	F24                Code = 135
	NumLock            Code = 144
	Mute               Code = 173
	VolumeDown         Code = 174
	VolumeUp           Code = 175
	Semicolon          Code = 186
	Equal              Code = 187
	Coma               Code = 188
	Dash               Code = 189
	Dot                Code = 190
	Slash              Code = 191
	Backquote          Code = 192
	SquareBracketLeft  Code = 219
	Backslash          Code = 220
	SquareBracketRight Code = 221
	Quote              Code = 222
)

// Location represents the location of the key on the keyboard or
// other input device.
type Location uint8

// board locations.
const (
	LocationStandard Location = iota
	LocationLeft
	LocationRight
	LocationNumpad
)
