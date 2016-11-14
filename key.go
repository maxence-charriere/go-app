package app

// Keyboard keys.
const (
	KeyBackspace          KeyCode = 8
	KeyTab                KeyCode = 9
	KeyEnter              KeyCode = 13
	KeyShift              KeyCode = 16
	KeyCtrl               KeyCode = 17
	KeyAlt                KeyCode = 18
	KeyPauseBreak         KeyCode = 19
	KeyCapsLock           KeyCode = 20
	KeyEsc                KeyCode = 27
	KeySpace              KeyCode = 32
	KeyPageUp             KeyCode = 33
	KeyPageDown           KeyCode = 34
	KeyEnd                KeyCode = 35
	KeyHome               KeyCode = 36
	KeyLeft               KeyCode = 37
	KeyUp                 KeyCode = 38
	KeyRight              KeyCode = 39
	KeyDown               KeyCode = 40
	KeyPrintScreen        KeyCode = 44
	KeyInsert             KeyCode = 45
	KeyDelete             KeyCode = 46
	Key0                  KeyCode = 48
	Key1                  KeyCode = 49
	Key2                  KeyCode = 50
	Key3                  KeyCode = 51
	Key4                  KeyCode = 52
	Key5                  KeyCode = 53
	Key6                  KeyCode = 54
	Key7                  KeyCode = 55
	Key8                  KeyCode = 56
	Key9                  KeyCode = 57
	KeyA                  KeyCode = 65
	KeyB                  KeyCode = 66
	KeyC                  KeyCode = 67
	KeyD                  KeyCode = 68
	KeyE                  KeyCode = 69
	KeyF                  KeyCode = 70
	KeyG                  KeyCode = 71
	KeyH                  KeyCode = 72
	KeyI                  KeyCode = 73
	KeyJ                  KeyCode = 74
	KeyK                  KeyCode = 75
	KeyL                  KeyCode = 76
	KeyM                  KeyCode = 77
	KeyN                  KeyCode = 78
	KeyO                  KeyCode = 79
	KeyP                  KeyCode = 80
	KeyQ                  KeyCode = 81
	KeyR                  KeyCode = 82
	KeyS                  KeyCode = 83
	KeyT                  KeyCode = 84
	KeyU                  KeyCode = 85
	KeyV                  KeyCode = 86
	KeyW                  KeyCode = 87
	KeyX                  KeyCode = 88
	KeyY                  KeyCode = 88
	KeyZ                  KeyCode = 90
	KeyMeta               KeyCode = 91
	KeyMenu               KeyCode = 93
	KeyNumPad0            KeyCode = 96
	KeyNumPad1            KeyCode = 97
	KeyNumPad2            KeyCode = 98
	KeyNumPad3            KeyCode = 99
	KeyNumPad4            KeyCode = 100
	KeyNumPad5            KeyCode = 101
	KeyNumPad6            KeyCode = 102
	KeyNumPad7            KeyCode = 103
	KeyNumPad8            KeyCode = 104
	KeyNumPad9            KeyCode = 105
	KeyNumPadMult         KeyCode = 106
	KeyNumPadPlus         KeyCode = 107
	KeyNumPadMin          KeyCode = 109
	KeyNumPadDot          KeyCode = 110
	KeyNumPadDecimal      KeyCode = 111
	KeyF1                 KeyCode = 112
	KeyF2                 KeyCode = 113
	KeyF3                 KeyCode = 114
	KeyF4                 KeyCode = 115
	KeyF5                 KeyCode = 116
	KeyF6                 KeyCode = 117
	KeyF7                 KeyCode = 118
	KeyF8                 KeyCode = 119
	KeyF9                 KeyCode = 120
	KeyF10                KeyCode = 121
	KeyF11                KeyCode = 122
	KeyF12                KeyCode = 123
	KeyNumLock            KeyCode = 144
	KeyMute               KeyCode = 173
	KeyVolumeDown         KeyCode = 174
	KeyVolumeUp           KeyCode = 175
	KeySemicolon          KeyCode = 186
	KeyEqual              KeyCode = 187
	KeyComa               KeyCode = 188
	KeyDash               KeyCode = 189
	KeyDot                KeyCode = 190
	KeySlash              KeyCode = 191
	KeyBackquote          KeyCode = 192
	KeySquareBracketLeft  KeyCode = 219
	KeyBackslash          KeyCode = 220
	KeySquareBracketRight KeyCode = 221
	KeyQuote              KeyCode = 222
)

// Keyboard locations.
const (
	KeyLocationStandard KeyLocation = iota
	KeyLocationLeft
	KeyLocationRight
	KeyLocationNumpad
)

// KeyCode represents a system and implementation dependent numerical
// code identifying the unmodified value of the pressed key.
type KeyCode uint8

// KeyLocation represents the location of the key on the keyboard or
// other input device.
type KeyLocation uint8
