package objutil

import (
	"fmt"
	"reflect"
)

const (
	// UintMax is the maximum value of a uint.
	UintMax = ^uint(0)

	// UintMin is the minimum value of a uint.
	UintMin = 0

	// Uint8Max is the maximum value of a uint8.
	Uint8Max = 255

	// Uint8Min is the minimum value of a uint8.
	Uint8Min = 0

	// Uint16Max is the maximum value of a uint16.
	Uint16Max = 65535

	// Uint16Min is the minimum value of a uint16.
	Uint16Min = 0

	// Uint32Max is the maximum value of a uint32.
	Uint32Max = 4294967295

	// Uint32Min is the minimum value of a uint32.
	Uint32Min = 0

	// Uint64Max is the maximum value of a uint64.
	Uint64Max = 18446744073709551615

	// Uint64Min is the minimum value of a uint64.
	Uint64Min = 0

	// UintptrMax is the maximum value of a uintptr.
	UintptrMax = ^uintptr(0)

	// UintptrMin is the minimum value of a uintptr.
	UintptrMin = 0

	// IntMax is the maximum value of a int.
	IntMax = int(UintMax >> 1)

	// IntMin is the minimum value of a int.
	IntMin = -IntMax - 1

	// Int8Max is the maximum value of a int8.
	Int8Max = 127

	// Int8Min is the minimum value of a int8.
	Int8Min = -128

	// Int16Max is the maximum value of a int16.
	Int16Max = 32767

	// Int16Min is the minimum value of a int16.
	Int16Min = -32768

	// Int32Max is the maximum value of a int32.
	Int32Max = 2147483647

	// Int32Min is the minimum value of a int32.
	Int32Min = -2147483648

	// Int64Max is the maximum value of a int64.
	Int64Max = 9223372036854775807

	// Int64Min is the minimum value of a int64.
	Int64Min = -9223372036854775808

	// Float32IntMax is the maximum consecutive integer value representable by a float32.
	Float32IntMax = 16777216

	// Float32IntMin is the minimum consecutive integer value representable by a float32.
	Float32IntMin = -16777216

	// Float64IntMax is the maximum consecutive integer value representable by a float64.
	Float64IntMax = 9007199254740992

	// Float64IntMin is the minimum consecutive integer value representable by a float64.
	Float64IntMin = -9007199254740992
)

// CheckUint64Bounds verifies that v is smaller than max, t represents the
// original type of v.
func CheckUint64Bounds(v uint64, max uint64, t reflect.Type) (err error) {
	if v > max {
		err = fmt.Errorf("objconv: %d overflows the maximum value of %d for %s", v, max, t)
	}
	return
}

// CheckInt64Bounds verifies that v is within min and max, t represents the
// original type of v.
func CheckInt64Bounds(v int64, min int64, max uint64, t reflect.Type) (err error) {
	if v < min {
		err = fmt.Errorf("objconv: %d overflows the minimum value of %d for %s", v, min, t)
	}
	if v > 0 && uint64(v) > max {
		err = fmt.Errorf("objconv: %d overflows the maximum value of %d for %s", v, max, t)
	}
	return
}
