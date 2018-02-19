package app

import "github.com/pkg/errors"

const (
	errNotSupported int = iota
	errNotFound
)

// ErrNotSupported is the interface that reports if an error is not supported
// error.
type ErrNotSupported interface {
	NotSupported() bool
}

// ErrNotFound is the interface that reports if an error is not found error.
type ErrNotFound interface {
	NotFound() bool
}

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

func (e appErr) NotFound() bool {
	return e.kind == errNotFound
}

// NewErrNotSupported creates a not supported error from the given feature.
func NewErrNotSupported(feature string) error {
	return appErr{
		kind: errNotSupported,
		base: errors.Errorf("%s is not supported", feature),
	}
}

// NewErrNotFound creates a not found error from the given object.
func NewErrNotFound(object string) error {
	return appErr{
		kind: errNotFound,
		base: errors.Errorf("%s is not found", object),
	}
}

// NotSupported is a helper fuction that reports if the given error is an
// not supported error.
func NotSupported(err error) bool {
	if ns, ok := err.(ErrNotSupported); ok {
		return ns.NotSupported()
	}
	return false
}

// NotFound is a helper fuction that reports if the given error is an
// not found error.
func NotFound(err error) bool {
	if nf, ok := err.(ErrNotFound); ok {
		return nf.NotFound()
	}
	return false
}
