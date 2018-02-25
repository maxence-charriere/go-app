package objutil

import "fmt"

// ParseInt parses a decimanl representation of an int64 from b.
//
// The function is equivalent to calling strconv.ParseInt(string(b), 10, 64) but
// it prevents Go from making a memory allocation for converting a byte slice to
// a string (escape analysis fails due to the error returned by strconv.ParseInt).
//
// Because it only works with base 10 the function is also significantly faster
// than strconv.ParseInt.
func ParseInt(b []byte) (int64, error) {
	var val int64

	if len(b) == 0 {
		return 0, errorInvalidUint64(b)
	}

	if b[0] == '-' {
		const max = Int64Min
		const lim = max / 10

		if b = b[1:]; len(b) == 0 {
			return 0, errorInvalidUint64(b)
		}

		for _, d := range b {
			if !(d >= '0' && d <= '9') {
				return 0, errorInvalidInt64(b)
			}

			if val < lim {
				return 0, errorOverflowInt64(b)
			}

			val *= 10
			x := int64(d - '0')

			if val < (max + x) {
				return 0, errorOverflowInt64(b)
			}

			val -= x
		}
	} else {
		const max = Int64Max
		const lim = max / 10

		for _, d := range b {
			if !(d >= '0' && d <= '9') {
				return 0, errorInvalidInt64(b)
			}
			x := int64(d - '0')

			if val > lim {
				return 0, errorOverflowInt64(b)
			}

			if val *= 10; val > (max - x) {
				return 0, errorOverflowInt64(b)
			}

			val += x
		}
	}

	return val, nil
}

// ParseUintHex parses a hexadecimanl representation of a uint64 from b.
//
// The function is equivalent to calling strconv.ParseUint(string(b), 16, 64) but
// it prevents Go from making a memory allocation for converting a byte slice to
// a string (escape analysis fails due to the error returned by strconv.ParseUint).
//
// Because it only works with base 16 the function is also significantly faster
// than strconv.ParseUint.
func ParseUintHex(b []byte) (uint64, error) {
	const max = Uint64Max
	const lim = max / 0x10
	var val uint64

	if len(b) == 0 {
		return 0, errorInvalidUint64(b)
	}

	for _, d := range b {
		var x uint64

		switch {
		case d >= '0' && d <= '9':
			x = uint64(d - '0')

		case d >= 'A' && d <= 'F':
			x = uint64(d-'A') + 0xA

		case d >= 'a' && d <= 'f':
			x = uint64(d-'a') + 0xA

		default:
			return 0, errorInvalidUint64(b)
		}

		if val > lim {
			return 0, errorOverflowUint64(b)
		}

		if val *= 0x10; val > (max - x) {
			return 0, errorOverflowUint64(b)
		}

		val += x
	}

	return val, nil
}

func errorInvalidInt64(b []byte) error {
	return fmt.Errorf("objconv: %#v is not a valid decimal representation of a signed 64 bits integer", string(b))
}

func errorOverflowInt64(b []byte) error {
	return fmt.Errorf("objconv: %#v overflows the maximum values of a signed 64 bits integer", string(b))
}

func errorInvalidUint64(b []byte) error {
	return fmt.Errorf("objconv: %#v is not a valid decimal representation of an unsigned 64 bits integer", string(b))
}

func errorOverflowUint64(b []byte) error {
	return fmt.Errorf("objconv: %#v overflows the maximum values of an unsigned 64 bits integer", string(b))
}
