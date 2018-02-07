package app

import (
	"testing"

	"github.com/pkg/errors"
)

func TestErrors(t *testing.T) {
	tests := []struct {
		scenario     string
		err          error
		notSupported bool
	}{
		{
			scenario:     "error is a not supported error",
			err:          NewErrNotSupported("test"),
			notSupported: true,
		},
		{
			scenario:     "error is not a not supported error",
			err:          errors.New("error test"),
			notSupported: false,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			t.Log("error:", test.err)

			if notSupported := NotSupported(test.err); notSupported != test.notSupported {
				t.Errorf("not supported is not %v: %v", test.notSupported, notSupported)
			}
		})
	}
}
