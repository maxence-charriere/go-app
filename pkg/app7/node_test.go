package app

import (
	"runtime"
	"testing"

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
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			testSkipNoWasm(t)

			n := u.node

			err := n.mount()
			require.NoError(t, err)
			testMounted(t, n)

			u.node.dismount()
			testDismounted(t, n)
		})
	}
}

func testMounted(t *testing.T, n UI) {
	require.NotNil(t, n.JSValue())
	require.True(t, n.Mounted())

	for _, c := range n.children() {
		require.Equal(t, n, c.parent())
		testMounted(t, c)
	}
}

func testDismounted(t *testing.T, n UI) {
	require.Nil(t, n.JSValue())
	require.False(t, n.Mounted())

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
