package app

import "github.com/pkg/errors"

const (
	errNotSupported int = iota
)

type appErr struct {
	kind int
	base error
}

func (e appErr) Error() string {
	return e.base.Error()
}

func (e appErr) NotSupported() bool {
	return e.kind == errNotSupported
}

// NewErrNotSupported creates a not supported error from the given feature.
func NewErrNotSupported(feature string) error {
	return appErr{
		kind: errNotSupported,
		base: errors.Errorf("%s is not supported", feature),
	}
}

// ErrNotSupported is a helper fuction that reports if the given error is an
// unsupported error.
func ErrNotSupported(err error) bool {
	appErr, ok := err.(appErr)
	if !ok {
		return false
	}
	return appErr.NotSupported()
}
