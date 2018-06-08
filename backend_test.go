package app

import "github.com/pkg/errors"

var errSimulated = errors.New("simulated error")

type backend struct {
	simulateError bool
}

func (b *backend) Run(f Factory, uiChan chan func()) error {
	if b.simulateError {
		return errSimulated
	}
	return nil
}

func (b *backend) Render(v interface{}) error {
	if b.simulateError {
		return errSimulated
	}
	return nil
}
