package maestro

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type foo struct {
	Int       int
	Nested    bar
	nonexport int
}

func (f foo) Method()     {}
func (f *foo) PtrMethod() {}

type bar struct {
	Int int
}

type foomap map[string]interface{}

func (f foomap) Method() {}

type fooslice []interface{}

func (f fooslice) Method() {}

type fooint int

func (f fooint) Method()     {}
func (f *fooint) PtrMethod() {}

func TestGetReceiver(t *testing.T) {
	tests := []struct {
		scenario string
		value    interface{}
		target   string
		recvKind reflect.Kind
		err      bool
	}{
		{
			scenario: "field from ptr",
			value:    &foo{},
			target:   "Int",
			recvKind: reflect.Int,
		},
		{
			scenario: "method from ptr",
			value:    &foo{},
			target:   "PtrMethod",
			recvKind: reflect.Func,
		},

		{
			scenario: "field from struct",
			value:    foo{},
			target:   "Int",
			recvKind: reflect.Int,
		},
		{
			scenario: "func from struct",
			value:    foo{},
			target:   "Method",
			recvKind: reflect.Func,
		},
		{
			scenario: "field from nested struct",
			value:    foo{},
			target:   "Nested.Int",
			recvKind: reflect.Int,
		},
		{
			scenario: "non exported field from struct",
			value:    foo{},
			target:   "nonexport",
			err:      true,
		},
		{
			scenario: "unknown field or method from struct",
			value:    foo{},
			target:   "Unknown",
			err:      true,
		},

		{
			scenario: "value from map",
			value:    foomap{"hello": 42},
			target:   "hello",
			recvKind: reflect.Int,
		},
		{
			scenario: "func from map",
			value:    foomap{"hello": 42},
			target:   "Method",
			recvKind: reflect.Func,
		},
		{
			scenario: "nested ptr value from map",
			value: foomap{
				"foo": &foo{},
			},
			target:   "foo.Int",
			recvKind: reflect.Int,
		},
		{
			scenario: "nested value from map",
			value: foomap{
				"foo": foo{},
			},
			target:   "foo.Int",
			recvKind: reflect.Int,
		},
		{
			scenario: "unknown value or method from map",
			value:    foomap{},
			target:   "unknown",
			err:      true,
		},

		{
			scenario: "value from slice",
			value:    fooslice{"hello"},
			target:   "0",
			recvKind: reflect.String,
		},
		{
			scenario: "func from slice",
			value:    fooslice{},
			target:   "Method",
			recvKind: reflect.Func,
		},
		{
			scenario: "nested value from slice",
			value: fooslice{
				"",
				foo{},
			},
			target:   "1.Int",
			recvKind: reflect.Int,
		},
		{
			scenario: "out of range value from slice",
			value:    fooslice{"hello"},
			target:   "1",
			err:      true,
		},
		{
			scenario: "non number index value from slice",
			value:    fooslice{"hello"},
			target:   "[]+_",
			err:      true,
		},

		{
			scenario: "common value with empty target",
			value:    fooint(42),
			target:   "",
			recvKind: reflect.Int,
		},
		{
			scenario: "func from common value",
			value:    fooint(42),
			target:   "Method",
			recvKind: reflect.Func,
		},
		{
			scenario: "func from common value",
			value:    new(fooint),
			target:   "PtrMethod",
			recvKind: reflect.Func,
		},
		{
			scenario: "common value within a target",
			value:    fooint(42),
			target:   "Method.Hello",
			err:      true,
		},
		{
			scenario: "unknown func form common value",
			value:    fooint(42),
			target:   "Unknown",
			err:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			recv, err := getReceiver(test.value, test.target)
			if test.err {
				assert.Error(t, err)
				t.Log(err)
				return
			}
			require.NoError(t, err)

			if recv.Kind() == reflect.Interface {
				recv = recv.Elem()
			}

			require.Equal(t, test.recvKind, recv.Kind())
		})
	}
}
