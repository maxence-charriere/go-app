package app

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/maxence-charriere/go-app/v6/pkg/errors"
	"github.com/maxence-charriere/go-app/v6/pkg/logs"
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

func TestMountAndDismount(t *testing.T) {
	utests := []struct {
		scenario string
		node     UI
	}{
		{
			scenario: "text",
			node:     Text("hello"),
		},
		{
			scenario: "html element",
			node: Div().
				Class("hello").
				OnClick(func(Context, Event) {}),
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			testSkipNoWasm(t)

			n := u.node
			err := mount(n)
			require.NoError(t, err)
			testMounted(t, n)

			dismount(u.node)
			testDismounted(t, n)
		})
	}
}

func testMounted(t *testing.T, n UI) {
	require.NotNil(t, n.JSValue())
	require.True(t, n.Mounted())

	switch n.Kind() {
	case HTML, Component:
		require.NoError(t, n.context().Err())
		require.NotNil(t, n.self())
	}

	for _, c := range n.children() {
		require.Equal(t, n, c.parent())
		testMounted(t, c)
	}
}

func testDismounted(t *testing.T, n UI) {
	require.Nil(t, n.JSValue())
	require.False(t, n.Mounted())

	switch n.Kind() {
	case HTML, Component:
		require.Error(t, n.context().Err())
		require.Nil(t, n.self())
	}

	for _, c := range n.children() {
		testDismounted(t, c)
	}
}

func testSkipNoWasm(t *testing.T) {
	if goarch := runtime.GOARCH; goarch != "wasm" {
		t.Skip(logs.New("skipping test").
			Tag("reason", "unsupported architecture").
			Tag("required-architecture", "wasm").
			Tag("current-architecture", goarch),
		)
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

func TestUpdate(t *testing.T) {
	utests := []struct {
		scenario   string
		a          UI
		b          UI
		matches    []TestUIDescriptor
		replaceErr bool
	}{
		{
			scenario:   "text element returns replace error when updated with an html element",
			a:          Text("hello"),
			b:          Div(),
			replaceErr: true,
		},
		{
			scenario:   "html element returns replace error when updated with a text element",
			a:          Div(),
			b:          Text("hello"),
			replaceErr: true,
		},
		{
			scenario: "text element is updated",
			a:        Text("hello"),
			b:        Text("world"),
			matches: []TestUIDescriptor{
				{
					Expected: Text("world"),
				},
			},
		},
		{
			scenario: "html element attributes are updated",
			a: Div().
				ID("max").
				Class("foo").
				AccessKey("test"),
			b: Div().
				ID("max").
				Class("bar").
				Lang("fr"),
			matches: []TestUIDescriptor{
				{
					Expected: Div().
						ID("max").
						Class("bar").
						Lang("fr"),
				},
			},
		},
		{
			scenario: "html element event handlers are updated",
			a: Div().
				OnClick(func(Context, Event) {}).
				OnBlur(func(Context, Event) {}),
			b: Div().
				OnClick(func(Context, Event) {}).
				OnChange(func(Context, Event) {}),
			matches: []TestUIDescriptor{
				{
					Expected: Div().
						OnClick(nil).
						OnChange(nil),
				},
			},
		},
		{
			scenario: "html element is replaced by a text",
			a: Div().Body(
				H2().Text("hello"),
			),
			b: Div().Body(
				Text("hello"),
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},
				{
					Path:     TestPath(0),
					Expected: Text("hello"),
				},
			},
		},
		{
			scenario: "text is replaced by a html elem",
			a: Div().Body(
				Text("hello"),
			),
			b: Div().Body(
				H2().Text("hello"),
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},
				{
					Path:     TestPath(0),
					Expected: H2(),
				},
				{
					Path:     TestPath(0, 0),
					Expected: Text("hello"),
				},
			},
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			testSkipNoWasm(t)

			err := mount(u.a)
			require.NoError(t, err)
			defer dismount(u.a)

			err = u.a.update(u.b)
			if u.replaceErr {
				require.Error(t, err)
				require.True(t, isErrReplace(err))
				return
			}

			require.NoError(t, err)

			for _, d := range u.matches {
				require.NoError(t, TestMatch(u.a, d))
			}
		})
	}
}
