package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAttributesSet(t *testing.T) {
	t.Run("set style", func(t *testing.T) {
		attributes := make(attributes)
		attributes.Set("style", "width:42px")
		require.Equal(t, "width:42px;", attributes["style"])
	})

	t.Run("set class", func(t *testing.T) {
		attributes := make(attributes)
		attributes.Set("class", "foo")
		require.Equal(t, "foo", attributes["class"])
	})

	t.Run("set multiple classes", func(t *testing.T) {
		attributes := make(attributes)
		attributes.Set("class", "foo")
		attributes.Set("class", "bar")
		require.Equal(t, "foo bar", attributes["class"])
	})

	t.Run("set srcset", func(t *testing.T) {
		attributes := make(attributes)
		attributes.Set("srcset", "/hi")
		require.Equal(t, "/hi", attributes["srcset"])
	})

	t.Run("set multiple srcset", func(t *testing.T) {
		attributes := make(attributes)
		attributes.Set("srcset", "/hi")
		attributes.Set("srcset", "/bye")
		require.Equal(t, "/hi, /bye", attributes["srcset"])
	})

	t.Run("set common attribute", func(t *testing.T) {
		attributes := make(attributes)
		attributes.Set("value", "foo")
		require.Equal(t, "foo", attributes["value"])
	})
}

func TestAttributesMount(t *testing.T) {
	attributes := make(attributes)
	attributes.Set("value", "foo")

	div := Div()
	client := NewClientTester(div)
	defer client.Close()

	attributes.Mount(div.JSValue(), func(s string) string {
		return s
	})
}

func TestAttributesUpdate(t *testing.T) {
	div := Div()
	client := NewClientTester(div)
	defer client.Close()

	resolveURL := func(s string) string {
		return s
	}

	t.Run("attribute is deleted", func(t *testing.T) {
		attributes := make(attributes)
		attributes.Set("value", "foo")
		attributes.Update(div.JSValue(), nil, resolveURL)
		require.Empty(t, attributes)
	})

	t.Run("same attribute value is skipped", func(t *testing.T) {
		a := make(attributes)
		a.Set("value", "foo")

		b := make(attributes)
		b.Set("value", "foo")

		a.Update(div.JSValue(), b, resolveURL)
	})

	t.Run("attribute is updated", func(t *testing.T) {
		a := make(attributes)
		a.Set("value", "foo")

		b := make(attributes)
		b.Set("value", "bar")

		a.Update(div.JSValue(), b, resolveURL)
		require.Equal(t, "bar", a["value"])
	})
}

func TestToAttributeValue(t *testing.T) {
	utests := []struct {
		scenario string
		in       string
		out      string
	}{
		{
			scenario: "spaces around",
			in:       "   \n  foo       \n",
			out:      "foo",
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			require.Equal(t, u.out, toAttributeValue(u.in))
		})
	}
}

func TestResolveAttributeURLValue(t *testing.T) {
	utests := []struct {
		name          string
		value         string
		resolvedValue string
	}{
		{
			name:          "value",
			value:         "bar",
			resolvedValue: "bar",
		},
		{
			name:          "cite",
			value:         "bar",
			resolvedValue: "/foo/bar",
		},
		{
			name:          "data",
			value:         "bar",
			resolvedValue: "/foo/bar",
		},
		{
			name:          "href",
			value:         "bar",
			resolvedValue: "/foo/bar",
		},
		{
			name:          "src",
			value:         "bar",
			resolvedValue: "/foo/bar",
		},
		{
			name:          "srcset",
			value:         "bar",
			resolvedValue: "/foo/bar",
		},
		{
			name:          "srcset",
			value:         "hi, bye",
			resolvedValue: "/foo/hi, /foo/bye",
		},
	}

	for _, u := range utests {
		t.Run(u.name, func(t *testing.T) {
			require.Equal(t, u.resolvedValue, resolveAttributeURLValue(
				u.name,
				u.value,
				func(s string) string {
					return "/foo/" + s
				}))
		})
	}
}

func TestSetDeleteJSAttribute(t *testing.T) {
	utests := []struct {
		name  string
		value string
	}{
		{
			name:  "value",
			value: "foo",
		},
		{
			name:  "class",
			value: "foo",
		},
		{
			name:  "contenteditable",
			value: "true",
		},
		{
			name:  "ismap",
			value: "true",
		},
		{
			name:  "readonly",
			value: "true",
		},
		{
			name:  "async",
			value: "true",
		},
		{
			name:  "autofocus",
			value: "true",
		},
		{
			name:  "autoplay",
			value: "true",
		},
		{
			name:  "checked",
			value: "true",
		},
		{
			name:  "default",
			value: "true",
		},
		{
			name:  "defer",
			value: "true",
		},
		{
			name:  "disabled",
			value: "true",
		},
		{
			name:  "hidden",
			value: "true",
		},
		{
			name:  "loop",
			value: "true",
		},
		{
			name:  "multiple",
			value: "true",
		},
		{
			name:  "muted",
			value: "true",
		},
		{
			name:  "open",
			value: "true",
		},
		{
			name:  "required",
			value: "true",
		},
		{
			name:  "reversed",
			value: "true",
		},
		{
			name:  "selected",
			value: "true",
		},
		{
			name:  "id",
			value: "foo",
		},
	}

	div := Div()
	client := NewClientTester(div)
	defer client.Close()

	for _, u := range utests {
		t.Run(u.name, func(t *testing.T) {
			t.Run("delete undefined", func(t *testing.T) {
				deleteJSAttribute(div.JSValue(), u.name)
			})

			t.Run("set and delete", func(t *testing.T) {
				setJSAttribute(div.JSValue(), u.name, u.value)
				deleteJSAttribute(div.JSValue(), u.name)
			})
		})
	}
}
