package app

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemoryStorage(t *testing.T) {
	testBrowserStorage(t, newMemoryStorage())
}

func TestJSLocalStorage(t *testing.T) {
	testSkipNonWasm(t)
	testBrowserStorage(t, newJSStorage("localStorage"))
}

func TestJSSessionStorage(t *testing.T) {
	testSkipNonWasm(t)
	testBrowserStorage(t, newJSStorage("sessionStorage"))
}

type obj struct {
	Foo int
	Bar string
}

func testBrowserStorage(t *testing.T, s BrowserStorage) {
	tests := []struct {
		scenario string
		function func(*testing.T, BrowserStorage)
	}{
		{
			scenario: "key does not exists",
			function: testBrowserStorageGetNotExists,
		},
		{
			scenario: "key is set and get",
			function: testBrowserStorageSetGet,
		},
		{
			scenario: "key is deleted",
			function: testBrowserStorageDel,
		},
		{
			scenario: "storage is cleared",
			function: testBrowserStorageClear,
		},
		{
			scenario: "set a non json value returns an error",
			function: testBrowserStorageSetError,
		},
		{
			scenario: "get with non json value receiver returns an error",
			function: testBrowserStorageGetError,
		},
		{
			scenario: "len returns the storage length",
			function: testBrowserStorageLen,
		},
		{
			scenario: "foreach iterates over each storage keys",
			function: testBrowserStorageForEach,
		},

		{
			scenario: "contains reports an existing key",
			function: testBrowserStorageContains,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			test.function(t, s)
		})
	}
}

func testBrowserStorageGetNotExists(t *testing.T, s BrowserStorage) {
	var o obj
	err := s.Get("/notexists", &o)
	require.NoError(t, err)
	require.Zero(t, o)
}

func testBrowserStorageSetGet(t *testing.T, s BrowserStorage) {
	var o obj
	err := s.Set("/exists", obj{
		Foo: 42,
		Bar: "hello",
	})
	require.NoError(t, err)

	err = s.Get("/exists", &o)
	require.NoError(t, err)
	require.Equal(t, 42, o.Foo)
	require.Equal(t, "hello", o.Bar)
}

func testBrowserStorageDel(t *testing.T, s BrowserStorage) {
	var o obj
	err := s.Set("/deleted", obj{
		Foo: 42,
		Bar: "bye",
	})
	require.NoError(t, err)

	s.Del("/deleted")
	err = s.Get("/deleted", &o)
	require.NoError(t, err)
	require.Zero(t, o)
}

func testBrowserStorageClear(t *testing.T, s BrowserStorage) {
	var o obj
	err := s.Set("/cleared", obj{
		Foo: 42,
		Bar: "sayonara",
	})
	require.NoError(t, err)

	s.Clear()
	err = s.Get("/cleared", &o)
	require.NoError(t, err)
	require.Zero(t, o)
}

func testBrowserStorageSetError(t *testing.T, s BrowserStorage) {
	err := s.Set("/func", func() {})
	require.Error(t, err)
}

func testBrowserStorageGetError(t *testing.T, s BrowserStorage) {
	err := s.Set("/value", obj{
		Foo: 42,
		Bar: "omae",
	})
	require.NoError(t, err)

	var f func()
	err = s.Get("/value", &f)
	require.Error(t, err)
}

func testBrowserStorageFull(t *testing.T, s BrowserStorage) {
	testSkipNonWasm(t)

	var err error
	data := make([]byte, 4096)
	i := 0

	for {
		key := fmt.Sprintf("/key_%d", i)

		if err = s.Set(key, data); err != nil {
			break
		}

		i++
	}

	require.Error(t, err)
	t.Log(err)
}

func testBrowserStorageLen(t *testing.T, s BrowserStorage) {
	s.Clear()

	s.Set("hello", 42)
	s.Set("world", 42)
	s.Set("bye", 42)

	require.Equal(t, 3, s.Len())
}

func testBrowserStorageForEach(t *testing.T, s BrowserStorage) {
	s.Clear()

	keys := []string{
		"starwars",
		"startrek",
		"alien",
		"marvel",
		"dune",
		"lords of the rings",
	}
	for _, k := range keys {
		s.Set(k, 3000)
	}
	require.Equal(t, s.Len(), len(keys))

	keyMap := make(map[string]struct{})
	s.ForEach(func(key string) {
		keyMap[key] = struct{}{}
	})
	require.Len(t, keyMap, len(keys))

	for _, k := range keys {
		require.Contains(t, keyMap, k)
	}
}

func testBrowserStorageContains(t *testing.T, s BrowserStorage) {
	s.Clear()

	require.False(t, s.Contains("lightsaber"))
	s.Set("lightsaber", true)
	require.True(t, s.Contains("lightsaber"))
}
