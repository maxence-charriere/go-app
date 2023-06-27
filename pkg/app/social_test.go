package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTwitterCardToMap(t *testing.T) {
	t.Run("empty fields", func(t *testing.T) {
		c := TwitterCard{}
		m := c.toMap()
		require.Empty(t, m["twitter:card"])
		require.Empty(t, m["twitter:site"])
		require.Empty(t, m["twitter:creator"])
		require.Empty(t, m["twitter:title"])
		require.Empty(t, m["twitter:description"])
		require.Empty(t, m["twitter:image"])
		require.Empty(t, m["twitter:image:alt"])
	})

	t.Run("fields", func(t *testing.T) {
		c := TwitterCard{
			Card:        "summary",
			Site:        "goapp",
			Creator:     "jonhymaxoo",
			Title:       "test",
			Description: "test description",
			Image:       "test image",
			ImageAlt:    "test image alt",
		}
		m := c.toMap()
		require.Equal(t, c.Card, m["twitter:card"])
		require.Equal(t, "@"+c.Site, m["twitter:site"])
		require.Equal(t, "@"+c.Creator, m["twitter:creator"])
		require.Equal(t, c.Title, m["twitter:title"])
		require.Equal(t, c.Description, m["twitter:description"])
		require.Equal(t, c.Image, m["twitter:image"])
		require.Equal(t, c.ImageAlt, m["twitter:image:alt"])
	})
}
