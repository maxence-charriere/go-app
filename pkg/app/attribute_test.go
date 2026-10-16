package app

import (
	"bytes"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func TestAttributesSet(t *testing.T) {
	t.Run("set style", func(t *testing.T) {
		attributes := make(attributes)
		attributes.Set("style", "width:42px")
		require.Equal(t, "width:42px", attributes["style"])
	})

	t.Run("set multiple style", func(t *testing.T) {
		attributes := make(attributes)
		attributes.Set("style", "width:42px")
		attributes.Set("style", "height:42px")
		require.Equal(t, "width:42px;height:42px", attributes["style"])
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

func TestToAttributeValue(t *testing.T) {
	utests := []struct {
		scenario string
		in       string
		out      string
	}{
		{
			scenario: "spaces arround",
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

	var m nodeManager
	div, err := m.Mount(makeTestContext(), 1, Div())
	require.NoError(t, err)

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

func parseEncodedElement(t testing.TB, markup string) *html.Node {
	t.Helper()
	nodes, err := html.ParseFragment(strings.NewReader(markup), &html.Node{
		Type: html.ElementNode, Data: "div", DataAtom: atom.Div,
	})
	require.NoError(t, err)
	require.Len(t, nodes, 1, "unexpected nodes in %q", markup)
	require.Equal(t, html.ElementNode, nodes[0].Type)
	require.Nil(t, nodes[0].FirstChild, "unexpected child in %q", markup)
	return nodes[0]
}

func TestEncodeAttributeValues(t *testing.T) {
	for _, value := range []string{
		`hello" autofocus onfocus=alert(1) x="`,
		`"><script>alert(1)</script><input x="`,
		`quotes " and ' with <tags> &amp; &#34;`,
		"first\nsecond\tline",
		"世界 🌍",
		`back\slash`,
		"true",
		"false",
	} {
		t.Run(value, func(t *testing.T) {
			var b bytes.Buffer
			nodeManager{}.Encode(makeTestContext(), &b, Input().Title(value))
			node := parseEncodedElement(t, b.String())
			require.Equal(t, "input", node.Data)
			require.Equal(t, []html.Attribute{{Key: "title", Val: value}}, node.Attr)
		})
	}
}

func TestEncodeResolvedAttributeValues(t *testing.T) {
	for _, name := range []string{
		"cite",
		"data",
		"href",
		"src",
		"srcset",
	} {
		t.Run(name, func(t *testing.T) {
			ctx := makeTestContext()
			ctx.resolveURL = func(value string) string {
				require.Equal(t, "/web/resource", value)
				return `https://example.com/resource?a=1&b="quoted"`
			}
			var b bytes.Buffer
			nodeManager{}.Encode(ctx, &b, Div().Attr(name, "/web/resource"))
			node := parseEncodedElement(t, b.String())
			require.Equal(t, []html.Attribute{{Key: name, Val: `https://example.com/resource?a=1&b="quoted"`}}, node.Attr)
		})
	}
}

func TestEncodeBooleanAttributes(t *testing.T) {
	for _, name := range []string{
		"allowfullscreen",
		"allowpaymentrequest",
		"async",
		"autofocus",
		"autoplay",
		"checked",
		"controls",
		"default",
		"defer",
		"disabled",
		"disablepictureinpicture",
		"disableremoteplayback",
		"formnovalidate",
		"hidden",
		"inert",
		"ismap",
		"itemscope",
		"loop",
		"multiple",
		"muted",
		"nomodule",
		"novalidate",
		"open",
		"playsinline",
		"readonly",
		"required",
		"reversed",
		"selected",
		"DISABLED",
	} {
		t.Run(name, func(t *testing.T) {
			for _, value := range []string{
				"true",
				"false",
				"",
				strings.ToLower(name),
			} {
				var b bytes.Buffer
				nodeManager{}.Encode(makeTestContext(), &b, Div().Attr(name, value))
				node := parseEncodedElement(t, b.String())
				if value == "false" {
					require.Empty(t, node.Attr)
					continue
				}
				if value == "true" {
					require.Equal(t, "<div "+name+"></div>", b.String())
					value = ""
				}
				require.Equal(t, []html.Attribute{{Key: strings.ToLower(name), Val: value}}, node.Attr)
			}
		})
	}
}

func TestEncodeLiteralAttributes(t *testing.T) {
	for _, name := range []string{
		"data-enabled",
		"aria-hidden",
		"draggable",
		"spellcheck",
		"contenteditable",
		"title",
		"value",
		"custom-flag",
	} {
		for _, value := range []string{
			"true",
			"false",
		} {
			t.Run(name+"/"+value, func(t *testing.T) {
				var b bytes.Buffer
				nodeManager{}.Encode(makeTestContext(), &b, Div().Attr(name, value))
				node := parseEncodedElement(t, b.String())
				require.Equal(t, []html.Attribute{{Key: name, Val: value}}, node.Attr)
			})
		}
	}
	for _, elem := range []HTML{Div().Attr("hidden", "until-found"), A().Download(""), Input().Value("")} {
		var b bytes.Buffer
		nodeManager{}.Encode(makeTestContext(), &b, elem)
		node := parseEncodedElement(t, b.String())
		require.Len(t, node.Attr, 1)
		require.Equal(t, elem.attrs()[node.Attr[0].Key], node.Attr[0].Val)
	}
}

func FuzzEncodeAttributeValue(f *testing.F) {
	for _, value := range []string{
		"",
		"true",
		"false",
		`hello" autofocus onfocus=alert(1) x="`,
		`"><script>alert(1)</script>`,
		"a&b'c\\d",
		"a\r\nb\nc\td",
		"世界 🌍",
		"\x00",
		"\xff",
	} {
		f.Add(value)
	}
	f.Fuzz(func(t *testing.T, value string) {
		var b bytes.Buffer
		nodeManager{}.Encode(makeTestContext(), &b, Input().Value(value))
		node := parseEncodedElement(t, b.String())
		require.Equal(t, "input", node.Data)
		require.Len(t, node.Attr, 1)
		require.Equal(t, "value", node.Attr[0].Key)
		if utf8.ValidString(value) && !strings.ContainsRune(value, 0) {
			// HTML parsing normalizes raw CR/CRLF. Only compare round trips
			// for valid UTF-8 without NUL; structure is checked for all inputs.
			value = strings.ReplaceAll(value, "\r\n", "\n")
			value = strings.ReplaceAll(value, "\r", "\n")
			require.Equal(t, value, node.Attr[0].Val)
		}
	})
}

func TestAttributeServerBrowserParity(t *testing.T) {
	testSkipNonWasm(t)
	ctx := makeTestContext()
	var m nodeManager
	for _, test := range []struct {
		name     string
		makeElem func(bool) HTML
		property string
		attr     string
	}{
		{"disabled", func(v bool) HTML { return Input().Disabled(v) }, "disabled", ""},
		{"checked", func(v bool) HTML { return Input().Type("checkbox").Checked(v) }, "checked", ""},
		{"async", func(v bool) HTML { return Script().Async(v) }, "async", ""},
		{"muted", func(v bool) HTML { return Video().Muted(v) }, "muted", ""},
		{"selected", func(v bool) HTML { return Option().Selected(v) }, "selected", ""},
		{"hidden", func(v bool) HTML { return Div().Hidden(v) }, "hidden", ""},
		{"controls", func(v bool) HTML { return Video().Controls(v) }, "controls", ""},
		{"readonly", func(v bool) HTML { return Input().ReadOnly(v) }, "readOnly", ""},
		{"ismap", func(v bool) HTML { return Img().IsMap(v) }, "isMap", ""},
		{"formnovalidate", func(v bool) HTML { return Input().FormNoValidate(v) }, "formNoValidate", ""},
		{"data", func(v bool) HTML { return Div().DataSet("enabled", v) }, "", "data-enabled"},
		{"aria", func(v bool) HTML { return Div().Aria("hidden", v) }, "", "aria-hidden"},
		{"draggable", func(v bool) HTML { return Div().Draggable(v) }, "draggable", "draggable"},
		{"spellcheck", func(v bool) HTML { return Div().Spellcheck(v) }, "spellcheck", "spellcheck"},
		{"contenteditable", func(v bool) HTML { return Div().ContentEditable(v) }, "", "contenteditable"},
	} {
		t.Run(test.name, func(t *testing.T) {
			mounted, err := m.Mount(ctx, 0, test.makeElem(false))
			require.NoError(t, err)
			defer func() { m.Dismount(mounted) }()
			for _, value := range []bool{false, true, false, true} {
				elem := test.makeElem(value)
				var b bytes.Buffer
				m.Encode(ctx, &b, elem)
				container, _ := Window().createElement("div", "")
				container.Set("innerHTML", b.String())
				parsed := container.Get("firstElementChild")
				mounted, err = m.Update(ctx, mounted, elem)
				require.NoError(t, err)
				if test.property != "" {
					require.Equal(t, value, parsed.Get(test.property).Bool())
					require.Equal(t, value, mounted.JSValue().Get(test.property).Bool())
				}
				if test.attr != "" {
					want := "false"
					if value {
						want = "true"
					}
					require.Equal(t, want, parsed.Call("getAttribute", test.attr).String())
					require.Equal(t, want, mounted.JSValue().Call("getAttribute", test.attr).String())
				}
			}
		})
	}

	t.Run("hidden until found", func(t *testing.T) {
		mounted, err := m.Mount(ctx, 0, Div().Attr("hidden", "until-found"))
		require.NoError(t, err)
		defer m.Dismount(mounted)
		require.Equal(t, "until-found", mounted.JSValue().Call("getAttribute", "hidden").String())
	})

	t.Run("generic boolean values", func(t *testing.T) {
		for _, name := range []string{
			"disabled",
			"checked",
			"hidden",
		} {
			mounted, err := m.Mount(ctx, 0, Input().Type("checkbox"))
			require.NoError(t, err)
			for _, value := range []string{
				"",
				"false",
				name,
				"true",
				"false",
			} {
				elem := Input().Type("checkbox").Attr(name, value)
				var b bytes.Buffer
				m.Encode(ctx, &b, elem)
				container, _ := Window().createElement("div", "")
				container.Set("innerHTML", b.String())
				mounted, err = m.Update(ctx, mounted, elem)
				require.NoError(t, err)
				require.Equal(t, value != "false", mounted.JSValue().Get(name).Bool())
				require.Equal(t, value != "false", container.Get("firstElementChild").Get(name).Bool())
			}
			m.Dismount(mounted)
		}
	})

	t.Run("escaped and resolved values", func(t *testing.T) {
		ctx.resolveURL = func(v string) string { return "https://example.com" + v }
		elem := A().Title("quotes \" ' & <tags>\n世界").Href(`/web/file?a=1&b="quoted"`)
		var b bytes.Buffer
		m.Encode(ctx, &b, elem)
		container, _ := Window().createElement("div", "")
		container.Set("innerHTML", b.String())
		parsed := container.Get("firstElementChild")
		mounted, err := m.Mount(ctx, 0, elem)
		require.NoError(t, err)
		defer m.Dismount(mounted)
		require.Equal(t, 1, container.Get("childNodes").Get("length").Int())
		require.Equal(t, 0, parsed.Get("childNodes").Get("length").Int())
		require.Equal(t, 2, parsed.Get("attributes").Get("length").Int())
		for name, value := range elem.attrs() {
			want := resolveAttributeURLValue(name, value, ctx.ResolveStaticResource)
			require.Equal(t, want, parsed.Call("getAttribute", name).String())
			require.Equal(t, want, mounted.JSValue().Call("getAttribute", name).String())
		}
	})
}
