package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpdateManagerAdd(t *testing.T) {
	t.Run("component is queued", func(t *testing.T) {
		var m updateManager

		compo := &hello{}
		m.Add(compo, 1)
		require.Len(t, m.pending, 100)
		_, ok := m.pending[0][compo]
		require.True(t, ok)
	})

	t.Run("pending components size is increased", func(t *testing.T) {
		var m updateManager

		compo := &hello{}
		m.Add(compo, 1)
		require.Len(t, m.pending, 100)

		compo2 := &bar{Compo: Compo{treeDepth: 200}}
		m.Add(compo2, 1)
		require.Len(t, m.pending, 201)

		_, added := m.pending[0][compo]
		require.True(t, added)

		_, added2 := m.pending[200][compo2]
		require.True(t, added2)
	})
}

func TestUpdateManagerDone(t *testing.T) {
	t.Run("component is removed from pending", func(t *testing.T) {
		var m updateManager

		compo := &hello{}
		m.Add(compo, 1)
		_, ok := m.pending[0][compo]
		require.True(t, ok)

		m.Done(compo)
		_, ok = m.pending[0][compo]
		require.False(t, ok)
	})

	t.Run("non added component is skipped", func(t *testing.T) {
		var m updateManager
		m.Done(&hello{})
	})
}

func TestUpdateManagerForEach(t *testing.T) {
	var m updateManager

	compo := &hello{}
	m.Add(compo, 1)

	m.ForEach(func(c Composer) {
		m.Done(c)
	})
	require.Empty(t, m.pending[0])
}
