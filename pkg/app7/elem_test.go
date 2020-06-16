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
			key:            "class",
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
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			e := &elem{}
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

	registeredHandler := e.eventHandlers["click"]
	require.Equal(t, h, registeredHandler)
}
