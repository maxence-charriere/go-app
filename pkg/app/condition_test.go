package app

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIfConditionIf(t *testing.T) {
	nodes := []Node{Div()}

	ifTrue := If(true, nodes...)
	require.False(t, ifTrue.isSatisfied())
	require.Equal(t, indirect(nodes...), ifTrue.nodes())
	require.Equal(t, reflect.TypeOf(ifTrue), ifTrue.nodeType())

	ifFalse := If(false, nodes...)
	require.True(t, ifFalse.isSatisfied())
	require.Empty(t, ifFalse.nodes())
}

func TestIfConditionElseIf(t *testing.T) {
	ifElems := []Node{Div()}
	elseIfElems := []Node{P()}

	noElseIf := If(true, ifElems...).
		ElseIf(true, elseIfElems...)
	require.False(t, noElseIf.isSatisfied())
	require.Equal(t, indirect(ifElems...), noElseIf.nodes())

	elseIfTrue := If(false, ifElems...).
		ElseIf(true, elseIfElems...)
	require.False(t, elseIfTrue.isSatisfied())
	require.Equal(t, indirect(elseIfElems...), elseIfTrue.nodes())

	elseIfFalse := If(false, ifElems...).
		ElseIf(false, elseIfElems...)
	require.True(t, elseIfFalse.isSatisfied())
	require.Empty(t, elseIfFalse.nodes())
}

func TestIfConditionElse(t *testing.T) {
	ifElems := []Node{Div()}
	elseElems := []Node{Text("hello")}

	noElse := If(true, ifElems...).
		Else(elseElems...)
	require.False(t, noElse.isSatisfied())
	require.Equal(t, indirect(ifElems...), noElse.nodes())

	elseTrue := If(false, ifElems...).
		Else(elseElems...)
	require.False(t, elseTrue.isSatisfied())
	require.Equal(t, indirect(elseElems...), elseTrue.nodes())
}
