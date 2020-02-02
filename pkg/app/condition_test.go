package app

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIfConditionIf(t *testing.T) {
	nodes := []Node{Div()}

	ifTrue := If(true, nodes...)
	require.False(t, ifTrue.eval)
	require.Equal(t, indirect(nodes...), ifTrue.nodes())
	require.Equal(t, reflect.TypeOf(ifTrue), ifTrue.nodeType())

	ifFalse := If(false, nodes...)
	require.True(t, ifFalse.eval)
	require.Empty(t, ifFalse.nodes())
}

func TestIfConditionElseIf(t *testing.T) {
	ifElems := []Node{Div()}
	elseIfElems := []Node{P()}

	noElseIf := If(true, ifElems...).
		ElseIf(true, elseIfElems...)
	require.False(t, noElseIf.eval)
	require.Equal(t, indirect(ifElems...), noElseIf.nodes())

	elseIfTrue := If(false, ifElems...).
		ElseIf(true, elseIfElems...)
	require.False(t, elseIfTrue.eval)
	require.Equal(t, indirect(elseIfElems...), elseIfTrue.nodes())

	elseIfFalse := If(false, ifElems...).
		ElseIf(false, elseIfElems...)
	require.True(t, elseIfFalse.eval)
	require.Empty(t, elseIfFalse.nodes())
}

func TestIfConditionElse(t *testing.T) {
	ifElems := []Node{Div()}
	elseElems := []Node{Text("hello")}

	noElse := If(true, ifElems...).
		Else(elseElems...)
	require.False(t, noElse.eval)
	require.Equal(t, indirect(ifElems...), noElse.nodes())

	elseTrue := If(false, ifElems...).
		Else(elseElems...)
	require.False(t, elseTrue.eval)
	require.Equal(t, indirect(elseElems...), elseTrue.nodes())
}

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
