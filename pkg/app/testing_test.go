package app

import (
	"os"
	"runtime"
	"testing"

	"github.com/maxence-charriere/go-app/v10/pkg/logs"
	"github.com/stretchr/testify/require"
)

func TestTestMatch(t *testing.T) {
	t.Run("match different types returns an error", func(t *testing.T) {
		err := Match(Div(), Span())
		require.Error(t, err)
		t.Log(err)
	})

	t.Run("match text succeeds", func(t *testing.T) {
		require.NoError(t, Match(Text("hello"), Text("hello")))
	})

	t.Run("match text returns an error", func(t *testing.T) {
		err := Match(Text("hello"), Text("bye"))
		require.Error(t, err)
		t.Log(err)
	})

	t.Run("match html succeeds", func(t *testing.T) {
		require.NoError(t, Match(Div(), Div()))
	})

	t.Run("match html with different tag returns an error", func(t *testing.T) {
		err := Match(Elem("div"), Elem("span"))
		require.Error(t, err)
		t.Log(err)
	})

	t.Run("match html with attributes succeeds", func(t *testing.T) {
		require.NoError(t, Match(Div().Class("hi"), Div().Class("hi")))
	})

	t.Run("match html with missing attributes returns an error", func(t *testing.T) {
		err := Match(Div().Class("hi"), Div())
		require.Error(t, err)
		t.Log(err)
	})

	t.Run("match html with unexpected attributes returns an error", func(t *testing.T) {
		err := Match(Div(), Div().Class("bye"))
		require.Error(t, err)
		t.Log(err)
	})

	t.Run("match html with different attribute values returns an error", func(t *testing.T) {
		err := Match(Div().Class("hi"), Div().Class("bye"))
		require.Error(t, err)
		t.Log(err)
	})

	t.Run("match html with event handlers succeeds", func(t *testing.T) {
		eventHandler := func(Context, Event) {}
		require.NoError(t, Match(Div().OnClick(eventHandler), Div().OnClick(eventHandler)))
	})

	t.Run("match html with missing event handlers returns an error", func(t *testing.T) {
		eventHandler := func(Context, Event) {}
		err := Match(Div().OnClick(eventHandler), Div())
		require.Error(t, err)
		t.Log(err)
	})

	t.Run("match html with unexpected event handlers returns an error", func(t *testing.T) {
		eventHandler := func(Context, Event) {}
		err := Match(Div(), Div().OnClick(eventHandler))
		require.Error(t, err)
		t.Log(err)
	})

	t.Run("match component succeeds", func(t *testing.T) {
		err := Match(&hello{Greeting: "hello"}, &hello{Greeting: "hello"})
		require.NoError(t, err)
	})

	t.Run("match component with different field values returns an error", func(t *testing.T) {
		err := Match(&hello{Greeting: "hello"}, &hello{Greeting: "bye"})
		require.Error(t, err)
		t.Log(err)
	})

	t.Run("match raw html succeeds", func(t *testing.T) {
		err := Match(Raw("<img>"), Raw("<img>"))
		require.NoError(t, err)
	})

	t.Run("match raw html with different values returns an error", func(t *testing.T) {
		err := Match(Raw("<img>"), Raw("<br>"))
		require.Error(t, err)
		t.Log(err)
	})

	t.Run("match path by html succeeds", func(t *testing.T) {
		err := Match(Text("hello"), Div().Body(
			Span(),
			Text("hello"),
		), 1)
		require.NoError(t, err)
	})

	t.Run("match path by html with bad path returns an error", func(t *testing.T) {
		err := Match(Text("hello"), Div().Body(
			Span(),
			Text("hello"),
		), 2)
		require.Error(t, err)
		t.Log(t, err)
	})

	e := newTestEngine()

	t.Run("match path by mounted component returns an error", func(t *testing.T) {
		compo := &compoWithCustomRoot{Root: Text("hello")}
		e.Load(compo)

		err := Match(Text("hello"), compo, 0)
		require.NoError(t, err)
	})

	t.Run("match path by mounted component and bad path returns an error", func(t *testing.T) {
		compo := &compoWithCustomRoot{Root: Text("hello")}
		e.Load(compo)

		err := Match(Text("hello"), compo, 1)
		require.Error(t, err)
		t.Log(err)
	})

	t.Run("match path by dismounted component returns an error", func(t *testing.T) {
		err := Match(Text("hello"), &compoWithCustomRoot{Root: Text("hello")}, 0)
		require.Error(t, err)
		t.Log(err)
	})

	t.Run("match unsupported elements returns an error", func(t *testing.T) {
		err := Match(nil, nil)
		require.Error(t, err)
		t.Log(err)
	})
}

func testSkipNonWasm(t *testing.T) {
	if goarch := runtime.GOARCH; goarch != "wasm" {
		t.Skip(logs.New("skipping test").
			WithTag("reason", "unsupported architecture").
			WithTag("required-architecture", "wasm").
			WithTag("current-architecture", goarch),
		)
	}
}

func testSkipWasm(t *testing.T) {
	if goarch := runtime.GOARCH; goarch == "wasm" {
		t.Skip(logs.New("skipping test").
			WithTag("reason", "unsupported architecture").
			WithTag("required-architecture", "!= than wasm").
			WithTag("current-architecture", goarch),
		)
	}
}

func testCreateDir(t *testing.T, path string) func() {
	err := os.MkdirAll(path, 0755)
	require.NoError(t, err)

	return func() {
		os.RemoveAll(path)
	}
}

func testCreateFile(t *testing.T, path, content string) {
	err := os.WriteFile(path, []byte(content), 0666)
	require.NoError(t, err)
}
