package dom

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type M struct {
	String              string
	Int                 int
	IntWithMethod       MappingInt
	IntPtr              *int
	Struct              MappingStruct
	Map                 map[string]string
	MapWithMethod       MappingMap
	Slice               []int
	SliceWithMethod     MappingSlice
	Array               [5]int
	Func                func()
	FuncWithArg         func(i int)
	FuncWithMultipleArg func(x, y int)
	method              func()
}

func (m *M) Method() {
	m.method()
}

func (m *M) Render() string {
	return `<div>Some mappings</div>`
}

type MappingStruct struct {
	Exported   int
	unexported int
	method     func()
}

func (s MappingStruct) Method() {
	s.method()
}

type MappingMap map[string]func()

func (m MappingMap) Method() {
	m["method"]()
}

type MappingSlice []func()

func (s MappingSlice) Method() {
	s[0]()
}

type MappingInt int

func (i MappingInt) Method(nb int) {
	mappedInt = nb
}

func TestMapping(t *testing.T) {
	intv := 42

	tests := []struct {
		scenario string
		mapping  Mapping
		expected M
		err      bool
		isFunc   bool
	}{
		{
			scenario: "map invalid field or method",
			mapping:  Mapping{FieldOrMethod: "String..Hello"},
			err:      true,
		},
		{
			scenario: "map field",
			mapping: Mapping{
				FieldOrMethod: "String",
				JSONValue:     `"hello"`,
			},
			expected: M{String: "hello"},
		},
		{
			scenario: "map method",
			mapping: Mapping{
				FieldOrMethod: "Method",
			},
		},
		{
			scenario: "map unexported method",
			mapping: Mapping{
				FieldOrMethod: "method",
			},
			err: true,
		},
		{
			scenario: "map pointer",
			mapping: Mapping{
				FieldOrMethod: "IntPtr",
				JSONValue:     "42",
			},
			expected: M{IntPtr: &intv},
		},
		{
			scenario: "map struct",
			mapping: Mapping{
				FieldOrMethod: "Struct",
				JSONValue:     `{"Exported": 42}`,
			},
			expected: M{Struct: MappingStruct{Exported: 42}},
		},
		{
			scenario: "map struct field",
			mapping: Mapping{
				FieldOrMethod: "Struct.Exported",
				JSONValue:     "42",
			},
			expected: M{Struct: MappingStruct{Exported: 42}},
		},
		{
			scenario: "map struct unexported field",
			mapping: Mapping{
				FieldOrMethod: "Struct.unexported",
				JSONValue:     "42",
			},
			err: true,
		},
		{
			scenario: "map nonexistent struct field",
			mapping: Mapping{
				FieldOrMethod: "Struct.Nonexistent",
			},
			err: true,
		},
		{
			scenario: "map struct method",
			mapping: Mapping{
				FieldOrMethod: "Struct.Method",
			},
			isFunc: true,
		},
		{
			scenario: "map map",
			mapping: Mapping{
				FieldOrMethod: "Map",
				JSONValue:     `{"foo": "bar"}`,
			},
			expected: M{Map: map[string]string{
				"foo": "bar",
			}},
		},
		{
			scenario: "map map method",
			mapping: Mapping{
				FieldOrMethod: "MapWithMethod.Method",
				JSONValue:     `{"foo": "bar"}`,
			},
			isFunc: true,
		},
		{
			scenario: "map nonexistent map key",
			mapping: Mapping{
				FieldOrMethod: "Map.hello",
			},
			err: true,
		},
		{
			scenario: "map slice",
			mapping: Mapping{
				FieldOrMethod: "Slice",
				JSONValue:     `[1, 2, 3, 4, 5]`,
			},
			expected: M{
				Slice: []int{1, 2, 3, 4, 5},
			},
		},
		{
			scenario: "map slice method",
			mapping: Mapping{
				FieldOrMethod: "SliceWithMethod.Method",
			},
			isFunc: true,
		},
		{
			scenario: "map slice out of range",
			mapping: Mapping{
				FieldOrMethod: "Slice.0",
			},
			err: true,
		},
		{
			scenario: "map array",
			mapping: Mapping{
				FieldOrMethod: "Array",
				JSONValue:     `[1, 2, 3, 4, 5]`,
			},
			expected: M{
				Array: [5]int{1, 2, 3, 4, 5},
			},
		},
		{
			scenario: "map func with arg",
			mapping: Mapping{
				FieldOrMethod: "FuncWithArg",
				JSONValue:     `42`,
			},
			isFunc: true,
		},
		{
			scenario: "map func with arg impossible field",
			mapping: Mapping{
				FieldOrMethod: "FuncWithArg.Unkown",
			},
			err: true,
		},
		{
			scenario: "map func with multiple arg",
			mapping: Mapping{
				FieldOrMethod: "FuncWithMultipleArg",
			},
			err: true,
		},
		{
			scenario: "map func with arg with bad json",
			mapping: Mapping{
				FieldOrMethod: "FuncWithArg",
				JSONValue:     `}{`,
			},
			err: true,
		},
		{
			scenario: "map value with bad json",
			mapping: Mapping{
				FieldOrMethod: "Int",
				JSONValue:     `}{`,
			},
			err: true,
		},
		{
			scenario: "map value method",
			mapping: Mapping{
				FieldOrMethod: "IntWithMethod.Method",
				JSONValue:     `42`,
			},
			isFunc: true,
		},
		{
			scenario: "map value nonexported method",
			mapping: Mapping{
				FieldOrMethod: "IntWithMethod.method",
			},
			err: true,
		},
		{
			scenario: "map value undefined method",
			mapping: Mapping{
				FieldOrMethod: "IntWithMethod.UndefinedMethod",
			},
			err: true,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			mappedInt = 0
			fn := func() {
				mappedInt = 42
			}

			c := &M{
				method: fn,
				Struct: MappingStruct{
					method: fn,
				},
				MapWithMethod: MappingMap{
					"method": fn,
				},
				SliceWithMethod: MappingSlice{fn},
				FuncWithArg: func(n int) {
					mappedInt = n
				},
			}

			f, err := test.mapping.Map(c)

			if test.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assertM(t, test.expected, *c)

			if test.isFunc {
				require.NotNil(t, f)
				f()
				require.Equal(t, 42, mappedInt)
			}
		})
	}
}

var mappedInt int

func assertM(t *testing.T, expected, actual M) {
	assert.Equal(t, expected.String, actual.String)
	assert.Equal(t, expected.Int, actual.Int)
	assert.Equal(t, expected.IntWithMethod, actual.IntWithMethod)

	if expected.IntPtr != nil {
		assert.Equal(t, expected.IntPtr, actual.IntPtr)
	}

	assert.Equal(t, expected.Struct.Exported, actual.Struct.Exported)
	assert.Equal(t, expected.Map, actual.Map)
	assert.Equal(t, expected.Slice, actual.Slice)
	assert.Equal(t, expected.Array, actual.Array)
}

func TestPipeline(t *testing.T) {
	tests := []struct {
		scenario         string
		target           string
		expectedPipeline []string
		err              bool
	}{
		{
			scenario:         "parses target",
			target:           "Hello",
			expectedPipeline: []string{"Hello"},
		},
		{
			scenario:         "parses target with multiple elements",
			target:           "Hello.World",
			expectedPipeline: []string{"Hello", "World"},
		},
		{
			scenario: "parses empyt target returns an error",
			err:      true,
		},
		{
			scenario: "parses target with empty element returns an error",
			target:   ".Hello.World",
			err:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			p, err := pipeline(test.target)

			if test.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.expectedPipeline, p)
		})
	}
}
