package objconv

import (
	"errors"
	"fmt"
)

func typeConversionError(from Type, to Type) error {
	return fmt.Errorf("objconv: cannot convert from %s to %s", from, to)
}

var (
	// End is expected to be returned to indicate that a function has completed
	// its work, this is usually employed in generic algorithms.
	End = errors.New("end")

	// This error value is used as a building block for reflection and is never
	// returned by the package.
	errBase = errors.New("")
)
