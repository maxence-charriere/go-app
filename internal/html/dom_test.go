package html

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/murlokswarm/app"
	"github.com/stretchr/testify/require"
)

type Foo struct {
	Value    string
	Disabled bool
}

func (f *Foo) Render() string {
	return `
	<div class="test" {{if .Disabled}}disabled{{end}}>
		{{.Value}}
	</div>
	`
}

type Bar struct {
	ReplaceTextByElem  bool
	ReplaceElemByElem  bool
	ReplaceCompoByElem bool
}

func (b *Bar) Render() string {
	return `
	<div>
		{{if .ReplaceTextByElem}}
			<span>hello</span>
		{{else}}
			hello
		{{end}}

		{{if .ReplaceElemByElem}}
			<h2>world</h2>
		{{else}}
			<h1>world</h1>
		{{end}}
	</div>
	`
}

type Boo struct {
	ReplaceCompoByElem bool
}

func (b *Boo) Render() string {
	return `
	<div>
		{{if .ReplaceCompoByElem}}
			<p>bar</p>
		{{else}}
			<html.Bar>
		{{end}}
	</div>
	`
}

func TestDOM(t *testing.T) {
	f := app.NewFactory()
	f.Register(&Foo{})

	tests := []struct {
		scenario string
		compo    app.Component
		modifier func(c app.Component)
		changes  []Change
		err      bool
	}{
		// Foo:
		{
			scenario: "create simple compo",
			compo:    &Foo{Value: "hello"},
			changes: []Change{
				createTextChange(""),
				setTextChange("", "hello"),
				createElemChange(newElemNode("div")),
				setAttrsChange("", map[string]string{"class": "test"}),
				appendChildChange("", ""),
				appendChildChange("", ""), // div -> root
			},
		},
		{
			scenario: "update simple compo",
			compo:    &Foo{Value: "hello"},
			modifier: func(c app.Component) {
				c.(*Foo).Value = "world"
			},
			changes: []Change{
				setTextChange("", "world"),
			},
		},
		{
			scenario: "append simple compo child",
			compo:    &Foo{},
			modifier: func(c app.Component) {
				c.(*Foo).Value = "hello"
			},
			changes: []Change{
				createTextChange(""),
				setTextChange("", "hello"),
				appendChildChange("", ""),
			},
		},
		{
			scenario: "remove simple compo child",
			compo:    &Foo{Value: "hello"},
			modifier: func(c app.Component) {
				c.(*Foo).Value = ""
			},
			changes: []Change{
				removeChildChange("", ""),
				deleteNodeChange(""),
			},
		},
		{
			scenario: "change simple compo root attrs",
			compo:    &Foo{},
			modifier: func(c app.Component) {
				c.(*Foo).Disabled = true
			},
			changes: []Change{
				setAttrsChange("", map[string]string{
					"class":    "test",
					"disabled": "",
				}),
			},
		},

		// Bar:
		{
			scenario: "create compo",
			compo:    &Bar{},
			changes: []Change{
				createTextChange(""),
				setTextChange("", "hello"),

				createTextChange(""),
				setTextChange("", "world"),
				createElemChange(newElemNode("h1")),
				setAttrsChange("", nil),
				appendChildChange("", ""), // world -> h1

				createElemChange(newElemNode("div")),
				setAttrsChange("", nil),
				appendChildChange("", ""), // hello -> div
				appendChildChange("", ""), // h1 -> div
				appendChildChange("", ""), // div -> root
			},
		},
		{
			scenario: "replace compo text by elem",
			compo:    &Bar{},
			modifier: func(c app.Component) {
				c.(*Bar).ReplaceTextByElem = true
			},
			changes: []Change{
				createTextChange(""),
				setTextChange("", "hello"),
				createElemChange(newElemNode("span")),
				setAttrsChange("", nil),
				appendChildChange("", ""), // hello -> span

				replaceChildChange("", "", ""),
				deleteNodeChange(""),
			},
		},
		{
			scenario: "replace compo elem by text",
			compo:    &Bar{ReplaceTextByElem: true},
			modifier: func(c app.Component) {
				c.(*Bar).ReplaceTextByElem = false
			},
			changes: []Change{
				createTextChange(""),
				setTextChange("", "hello"),
				replaceChildChange("", "", ""), // hello -> span
				deleteNodeChange(""),           // delete span.hello
				deleteNodeChange(""),           // delete span
			},
		},
		{
			scenario: "replace compo elem by elem",
			compo:    &Bar{},
			modifier: func(c app.Component) {
				c.(*Bar).ReplaceElemByElem = true
			},
			changes: []Change{
				createTextChange(""),
				setTextChange("", "world"),
				createElemChange(newElemNode("h2")),
				setAttrsChange("", nil),
				appendChildChange("", ""), // world -> h2

				replaceChildChange("", "", ""),
				deleteNodeChange(""), // delete h1.world
				deleteNodeChange(""), // delete h1
			},
		},

		// Boo:
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			dom := NewDOM(f, "test")
			changes, err := dom.Render(test.compo)

			if test.err {
				assert.Error(t, err)
				return
			}

			if test.modifier != nil {
				test.modifier(test.compo)
				changes, err = dom.Render(test.compo)
			}

			require.NoError(t, err)

			jsonChanges, _ := json.MarshalIndent(changes, "", "  ")
			t.Log("changes:", string(jsonChanges))

			require.Len(t, changes, len(test.changes))

			for i := range changes {
				requireEqualChange(t, test.changes[i], changes[i])
			}
		})
	}
}

func requireEqualChange(t require.TestingT, expected, actual Change) {
	require.Equal(t, expected.Type, actual.Type)

	switch expected.Type {
	case setText:
		require.Equal(t, expected.Value.(textValue).Text, actual.Value.(textValue).Text)

	case createElem:
		require.Equal(t, expected.Value.(elemValue).TagName, actual.Value.(elemValue).TagName)

	case setAttrs:
		attrs := actual.Value.(elemValue).Attrs
		delete(attrs, "data-goapp-id")
		require.Equal(t, expected.Value.(elemValue).Attrs, actual.Value.(elemValue).Attrs)

	default:
	}
}
