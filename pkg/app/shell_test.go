package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShellHasMenu(t *testing.T) {
	s := Shell().(*shell)
	require.False(t, s.hasMenu())

	s = Shell().
		Menu(
			Div(),
		).(*shell)
	require.True(t, s.hasMenu())
}

func TestShellHasSubmenu(t *testing.T) {
	s := Shell().(*shell)
	require.False(t, s.hasSubmenu())

	s = Shell().
		Submenu(
			Div(),
		).(*shell)
	require.True(t, s.hasSubmenu())
}

func TestShellHasOverlayMenu(t *testing.T) {
	s := Shell().(*shell)
	require.False(t, s.hasOverlayMenu())

	s = Shell().
		OverlayMenu(
			Div(),
		).(*shell)
	require.True(t, s.hasOverlayMenu())
}

func TestShellMounted(t *testing.T) {
	s := Shell().(*shell)
	require.False(t, s.mounted())

	d := NewClientTestingDispatcher(s)
	defer d.Close()

	d.Consume()
	require.True(t, s.mounted())
}
