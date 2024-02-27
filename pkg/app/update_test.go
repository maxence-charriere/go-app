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

	t.Run("removed updates are skipped", func(t *testing.T) {
		var m updateManager

		compo := &hello{}
		m.Add(compo, 1)
		_, ok := m.pending[0][compo]
		require.True(t, ok)

		m.pending[0] = nil
		_, ok = m.pending[0][compo]
		require.False(t, ok)

		m.Done(compo)
		_, ok = m.pending[0][compo]
		require.False(t, ok)
	})

	t.Run("non added component is skipped", func(t *testing.T) {
		var m updateManager
		m.Done(&hello{})
	})
}

func TestUpdateManagerUpdateForEach(t *testing.T) {
	t.Run("component with positive counter is updated", func(t *testing.T) {
		var m updateManager

		compo := &hello{}
		m.Add(compo, 1)
		require.NotEmpty(t, m.pending[0])

		var updates int
		m.UpdateForEach(func(c Composer) {
			updates++
		})
		require.Equal(t, 1, updates)
		require.Empty(t, m.pending[0])
	})

	t.Run("component with negative counter is skipped", func(t *testing.T) {
		var m updateManager

		compo := &hello{}
		m.Add(compo, -1)
		require.NotEmpty(t, m.pending[0])

		var updates int
		m.UpdateForEach(func(c Composer) {
			updates++
		})
		require.Zero(t, updates)
		require.Empty(t, m.pending[0])
	})
}
