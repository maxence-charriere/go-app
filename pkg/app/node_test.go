package app

import (
	"fmt"
	"testing"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestKindString(t *testing.T) {
	utests := []struct {
		kind           Kind
		expectedString string
	}{
		{
			kind:           UndefinedElem,
			expectedString: "undefined",
		},
		{
			kind:           SimpleText,
			expectedString: "text",
		},
		{
			kind:           HTML,
			expectedString: "html",
		},
		{
			kind:           Component,
			expectedString: "component",
		},
		{
			kind:           Selector,
			expectedString: "selector",
		},
	}

	for _, u := range utests {
		t.Run(u.expectedString, func(t *testing.T) {
			require.Equal(t, u.expectedString, u.kind.String())
		})
	}
}

func TestFilterUIElems(t *testing.T) {
	var nilText *text

	simpleText := Text("hello")

	expectedResult := []UI{
		simpleText,
	}

	res := FilterUIElems(nil, nilText, simpleText)
	require.Equal(t, expectedResult, res)
}

func TestIsErrReplace(t *testing.T) {
	utests := []struct {
		scenario     string
		err          error
		isErrReplace bool
	}{
		{
			scenario:     "error is a replace error",
			err:          errors.New("test").Tag("replace", true),
			isErrReplace: true,
		},
		{
			scenario:     "error is not a replace error",
			err:          errors.New("test").Tag("test", true),
			isErrReplace: false,
		},
		{
			scenario:     "standard error is not a replace error",
			err:          fmt.Errorf("test"),
			isErrReplace: false,
		},
		{
			scenario:     "nil error is not a replace error",
			err:          nil,
			isErrReplace: false,
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			res := isErrReplace(u.err)
			require.Equal(t, u.isErrReplace, res)
		})
	}
}

type mountTest struct {
	scenario string
	node     UI
}

func testMountDismount(t *testing.T, utests []mountTest) {
	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			n := u.node

			d := NewClientTester(n)

			d.Consume()
			testMounted(t, n)

			d.Close()
			testDismounted(t, n)
		})
	}
}

func testMounted(t *testing.T, n UI) {
	require.NotNil(t, n.JSValue())
	require.NotNil(t, n.getDispatcher())
	require.True(t, n.IsMounted())

	switch n.Kind() {
	case HTML, Component:
		require.NoError(t, n.getContext().Err())
		require.NotNil(t, n.self())
	}

	for _, c := range n.getChildren() {
		require.Equal(t, n, c.getParent())
		testMounted(t, c)
	}
}

func testDismounted(t *testing.T, n UI) {
	require.Nil(t, n.JSValue())
	require.NotNil(t, n.getDispatcher())
	require.False(t, n.IsMounted())

	switch n.Kind() {
	case HTML, Component:
		require.Error(t, n.getContext().Err())
		require.Nil(t, n.self())
	}

	for _, c := range n.getChildren() {
		testDismounted(t, c)
	}
}

type updateTest struct {
	scenario   string
	a          UI
	b          UI
	matches    []TestUIDescriptor
	replaceErr bool
}

func testUpdate(t *testing.T, utests []updateTest) {
	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			d := NewClientTester(u.a)
			defer d.Close()
			d.Consume()

			err := update(u.a, u.b)
			if u.replaceErr {
				require.Error(t, err)
				require.True(t, isErrReplace(err))
				return
			}

			d.Consume()

			require.NoError(t, err)
			for _, d := range u.matches {
				require.NoError(t, TestMatch(u.a, d))
			}
		})
	}
}

func TestHTMLString(t *testing.T) {
	utests := []struct {
		scenario string
		root     UI
	}{
		{
			scenario: "hmtl element",
			root:     Div().ID("test"),
		},
		{
			scenario: "text",
			root:     Text("hello"),
		},
		{
			scenario: "component",
			root:     &hello{},
		},
		{
			scenario: "nested component",
			root:     Div().Body(&hello{}),
		},
		{
			scenario: "nested nested component",
			root:     Div().Body(&foo{Bar: "bar"}),
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			t.Log(HTMLString(u.root))
			t.Log(HTMLStringWithIndent(u.root))
		})
	}
}
