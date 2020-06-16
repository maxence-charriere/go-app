package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestElemSetAttr(t *testing.T) {
	utests := []struct {
		scenario       string
		key            string
		value          interface{}
		explectedValue string
		valueNotSet    bool
	}{
		{
			scenario:       "string",
			key:            "title",
			value:          "test",
			explectedValue: "test",
		},
		{
			scenario:       "int",
			key:            "max",
			value:          42,
			explectedValue: "42",
		},
		{
			scenario:       "bool true",
			key:            "hidden",
			value:          true,
			explectedValue: "",
		},
		{
			scenario:    "bool false",
			key:         "hidden",
			value:       false,
			valueNotSet: true,
		},
		{
			scenario:       "style",
			key:            "style",
			value:          "margin:42",
			explectedValue: "margin:42;",
		},
		{
			scenario:       "set successive styles",
			key:            "style",
			value:          "padding:42",
			explectedValue: "margin:42;padding:42;",
		},
		{
			scenario:       "class",
			key:            "class",
			value:          "hello",
			explectedValue: "hello",
		},
		{
			scenario:       "set successive classes",
			key:            "class",
			value:          "world",
			explectedValue: "hello world",
		},
	}

	e := &elem{}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			e.setAttr(u.key, u.value)
			v, exists := e.attrs[u.key]
			require.Equal(t, u.explectedValue, v)
			require.Equal(t, u.valueNotSet, !exists)
		})
	}
}

func TestElemSetEventHandler(t *testing.T) {
	e := &elem{}
	h := func(Context, Event) {}
	e.setEventHandler("click", h)

	expectedHandler := elemEventHandler{
		event: "click",
		value: h,
	}

	registeredHandler := e.eventHandlers["click"]
	require.True(t, expectedHandler.equal(registeredHandler))
}
