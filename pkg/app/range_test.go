package app

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRangeConditionSlice(t *testing.T) {
	s := []string{
		"foo",
		"bar",
		"boo",
	}

	rs := Range(s).
		Slice(func(i int) Node {
			return Text(s[i])
		})
	require.Equal(t, reflect.TypeOf(rs), rs.nodeType())
	require.Len(t, rs.nodes(), 3)
	for i := range s {
		require.Equal(t, s[i], rs.nodes()[i].(textNode).text())
	}

	require.Panics(t, func() {
		Range(42).
			Slice(func(int) Node {
				return nil
			})
	})
}

func TestRangeConditionMap(t *testing.T) {
	m := map[string]string{
		"foo": "maxxy",
		"bar": "maxoo",
		"boo": "max",
	}

	rm := Range(m).
		Map(func(k string) Node {
			return Text(m[k])
		})
	require.Len(t, rm.nodes(), 3)

	require.Panics(t, func() {
		Range(42).
			Map(func(string) Node {
				return nil
			})
	})

	require.Panics(t, func() {
		Range(map[int]string{42: ""}).
			Map(func(string) Node {
				return nil
			})
	})
}
