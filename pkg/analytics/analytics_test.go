package analytics

import "testing"

func TestAnalytics(t *testing.T) {
	testingProps := func() map[string]interface{} {
		return map[string]interface{}{
			"string": 42,
			"uint":   uint(23),
			"int":    42,
			"float":  42.2,
			"slice":  []interface{}{"hello", 42},
			"map":    map[string]interface{}{"foo": "bar"},
			"struct": struct{ Foo string }{Foo: "bar"},
		}
	}

	providers := []struct {
		name    string
		backend Backend
	}{
		{
			name:    "google analytics",
			backend: NewGoogleAnalytics(),
		},
	}

	for _, p := range providers {
		t.Run(p.name, func(t *testing.T) {
			Add(p.backend)
			defer func() {
				backends = nil
			}()

			t.Run("identify", func(t *testing.T) {
				Identify("Maxoo", nil)
			})
			t.Run("identify with traits", func(t *testing.T) {
				Identify("Maxoo", testingProps())
			})

			t.Run("event", func(t *testing.T) {
				Track("test", nil)
			})
			t.Run("event with properties", func(t *testing.T) {
				Track("test", testingProps())
			})

			t.Run("page", func(t *testing.T) {
				Page("Test", nil)
			})
			t.Run("page with properties", func(t *testing.T) {
				Page("Test", testingProps())
			})
		})
	}
}
