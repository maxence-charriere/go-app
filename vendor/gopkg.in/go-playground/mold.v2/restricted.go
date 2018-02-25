package mold

const (
	diveTag            = "dive"
	restrictedTagChars = ".[],|=+()`~!@#$%^&*\\\"/?<>{}"
	tagSeparator    = ","
	ignoreTag       = "-"
	tagKeySeparator = "="
	utf8HexComma    = "0x2C"
	keysTag            = "keys"
	endKeysTag         = "endkeys"
)

var (
	restrictedTags = map[string]struct{}{
		diveTag:   {},
		ignoreTag: {},
	}
)
