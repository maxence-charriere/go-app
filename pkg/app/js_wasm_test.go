package app

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"syscall/js"
	"testing"
)

var sampleFuncOf = func(this Value, args []Value) any { return nil }

func TestCleanArg(t *testing.T) {

	t.Run(fmt.Sprintf("cleanArg slice"), func(t *testing.T) {
		arg := cleanArg([]any{"string", FuncOf(sampleFuncOf)})
		switch ta := arg.(type) {
		case []any:
			require.Equal(t, 2, len(ta))
			require.Equal(t, "string", ta[0])
			require.IsType(t, js.Func{}, ta[1])
			return
		}
		t.Fail()
	})

	t.Run(fmt.Sprintf("cleanArg func"), func(t *testing.T) {
		require.IsType(t, js.Func{}, cleanArg(FuncOf(sampleFuncOf)))
	})

	t.Run(fmt.Sprintf("cleanArg string"), func(t *testing.T) {
		require.Equal(t, "string", cleanArg("string"))
	})

}

func TestJSValue(t *testing.T) {

	t.Run("test JSValue(app.Value) any", func(t *testing.T) {
		global := js.Global()
		require.True(t, global.Equal(js.ValueOf(JSValue(ValueOf(global)))))
	})

}
