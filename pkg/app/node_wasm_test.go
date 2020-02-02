package app

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIndirect(t *testing.T) {
	tests := []struct {
		scenario string
		node     Node
		expected []reflect.Type
	}{
		{
			scenario: "indirect standard node returns a standard node",
			node:     Div(),
			expected: []reflect.Type{
				Div().nodeType(),
			},
		},
		{
			scenario: "indirect condition node returns a standard node",
			node:     If(true, Div()),
			expected: []reflect.Type{
				Div().nodeType(),
			},
		},
		{
			scenario: "indirect if condition node returns a standard node",
			node:     If(true, Div()),
			expected: []reflect.Type{
				Div().nodeType(),
			},
		},
		{
			scenario: "indirect range condition node returns a standard nodes",
			node: Range([]int{1, 2, 3}).
				Slice(func(i int) Node {
					return Div()
				}),
			expected: []reflect.Type{
				Div().nodeType(),
				Div().nodeType(),
				Div().nodeType(),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			nodes := indirect(test.node)
			for i, n := range nodes {
				require.Equal(t, test.expected[i], n.nodeType())
			}
		})
	}
}

func TestMount(t *testing.T) {
	tests := []struct {
		scenario string
		node     Node
		err      bool
	}{
		{
			scenario: "text node is mounted",
			node:     Text("foo"),
		},
		{
			scenario: "standard node is mounted",
			node:     Div(),
		},
		{
			scenario: "compo node is mounted",
			node:     &foo{},
		},
		{
			scenario: "mounting if condition node returns an error",
			node:     If(true, Div()),
			err:      true,
		},
		{
			scenario: "mounting range condition node returns an error",
			node: Range([]int{42}).
				Slice(func(int) Node {
					return Div()
				}),
			err: true,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			err := mount(test.node)
			if test.err {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

type foo struct {
	Compo
	Text        string
	Bar         bar
	notExported bool
}

func (f *foo) Render() ValueNode {
	return Text(f.Text)
}

type bar struct {
	Int int
}

func TestUpdateStandardNode(t *testing.T) {
	a := Div()
	err := mount(a)
	require.NoError(t, err)
	require.Empty(t, a.attrs)
	require.Empty(t, a.events)

	b := Div().
		Class("foo").
		OnClick(func(src Value, e Event) {})
	err = update(a, b)
	require.NoError(t, err)
	require.Equal(t, b.attrs, a.attrs)
	require.Len(t, a.events, 1)
}

func TestUpdateStandardNodeAddChild(t *testing.T) {
	a := Div()
	err := mount(a)
	require.NoError(t, err)
	require.Empty(t, a.children())

	b := Div().Body(
		Text("hello"),
		Br(),
	)
	err = update(a, b)
	require.NoError(t, err)
	require.Len(t, a.children(), 2)
	require.Equal(t, reflect.TypeOf(Text("")), a.body[0].nodeType())
	require.NotNil(t, a.body[0].JSValue())
	require.Equal(t, reflect.TypeOf(Br()), a.body[1].nodeType())
	require.NotNil(t, a.body[1].JSValue())
}

func TestUpdateStandardNodeRemoveChild(t *testing.T) {
	a := Div().Body(
		Br(),
		Text("hi"),
	)
	err := mount(a)
	require.NoError(t, err)
	require.Len(t, a.children(), 2)

	b := Div()
	err = update(a, b)
	require.NoError(t, err)
	require.Empty(t, a.children())
}

func TestUpdateStandardNodeChild(t *testing.T) {
	text := Text("foo")
	a := Div().Body(
		text,
	)
	err := mount(a)
	require.NoError(t, err)

	b := Div().Body(
		Text("bar"),
	)
	err = update(a, b)
	require.NoError(t, err)
	require.Equal(t, "bar", text.(textNode).text())
}

func TestUpdateChildComponent(t *testing.T) {
	c := &foo{Text: "boo"}
	a := Div().Body(
		c,
	)
	err := mount(a)
	require.NoError(t, err)

	b := Div().Body(
		&foo{
			Text: "foo",
			Bar:  bar{Int: 42},
		},
	)

	err = update(a, b)
	require.NoError(t, err)
	require.Equal(t, "foo", c.Text)
	require.Equal(t, 42, c.Bar.Int)
}
