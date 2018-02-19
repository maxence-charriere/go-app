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
		notFound     bool
	}{
		{
			scenario:     "error is a not supported error",
			err:          NewErrNotSupported("test"),
			notSupported: true,
		},
		{
			scenario: "error is a not found error",
			err:      NewErrNotFound("test"),
			notFound: true,
		},
		{
			scenario: "error is a normal error",
			err:      errors.New("a simple error"),
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			t.Log("error:", test.err)

			if notSupported := NotSupported(test.err); notSupported != test.notSupported {
				t.Errorf("not supported is not %v: %v", test.notSupported, notSupported)
			}

			if notFound := NotFound(test.err); notFound != test.notFound {
				t.Errorf("not found is not %v: %v", test.notFound, notFound)
			}
		})
	}
}
