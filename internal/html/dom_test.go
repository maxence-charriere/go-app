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

func (f *Foo) OnMount() {
}

func (f *Foo) OnDismount() {
}

func (f *Foo) Subscribe() app.EventSubscriber {
	return app.NewEventSubscriber()
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
	ReplaceCompoType   bool
	AddCompo           bool
	Value              string
}

func (b *Boo) Render() string {
	return `
	<div>
		{{if .ReplaceCompoByElem}}
			<p>foo</p>
		{{else if .ReplaceCompoType}}
			<html.Oob>
		{{else}}
			<html.Foo value="{{.Value}}">
		{{end}}

		{{if .AddCompo}}
			<html.Foo>
		{{end}}
	</div>
	`
}

type Oob app.ZeroCompo

func (o *Oob) Render() string {
	return `<p></p>`
}

func TestDOM(t *testing.T) {
	f := app.NewFactory()
	f.Register(&Foo{})
	f.Register(&Bar{})
	f.Register(&Boo{})
	f.Register(&Oob{})

	tests := []struct {
		scenario   string
		compo      app.Component
		modifier   func(c app.Component)
		changes    []Change
		compoCount int
		err        bool
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
			compoCount: 1,
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
			compoCount: 1,
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
			compoCount: 1,
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
			compoCount: 1,
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
			compoCount: 1,
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
			compoCount: 1,
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
			compoCount: 1,
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
			compoCount: 1,
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
			compoCount: 1,
		},

		// Boo:
		{
			scenario: "create compo with nested compo",
			compo:    &Boo{},
			changes: []Change{
				createElemChange(newElemNode("div")),
				setAttrsChange("", map[string]string{"class": "test"}),
				createCompoChange("", "html.foo"),
				setCompoRootChange("", ""),

				createElemChange(newElemNode("div")),
				setAttrsChange("", nil),
				appendChildChange("", ""), // foo.div -> div
				appendChildChange("", ""), // div -> root
			},
			compoCount: 2,
		},
		{
			scenario: "add compo to elem",
			compo:    &Boo{},
			modifier: func(c app.Component) {
				c.(*Boo).AddCompo = true
			},
			changes: []Change{
				createElemChange(newElemNode("div")),
				setAttrsChange("", map[string]string{"class": "test"}),
				createCompoChange("", "html.foo"),
				setCompoRootChange("", ""),
				appendChildChange("", ""), // foo2 -> div
			},
			compoCount: 3,
		},
		{
			scenario: "remove compo from elem",
			compo:    &Boo{AddCompo: true},
			modifier: func(c app.Component) {
				c.(*Boo).AddCompo = false
			},
			changes: []Change{
				removeChildChange("", ""), // foo2
				deleteNodeChange(""),      // foo2.div
				deleteNodeChange(""),      // foo2
			},
			compoCount: 2,
		},
		{
			scenario: "replace compo type",
			compo:    &Boo{},
			modifier: func(c app.Component) {
				c.(*Boo).ReplaceCompoType = true
			},
			changes: []Change{
				createElemChange(newElemNode("p")), // oob.p
				setAttrsChange("", nil),
				createCompoChange("", "html.oob"),
				setCompoRootChange("", ""), // oob.p -> oob

				replaceChildChange("", "", ""), // foo <-> oob
				deleteNodeChange(""),           // foo.text
				deleteNodeChange(""),           // foo
			},
			compoCount: 2,
		},
		{
			scenario: "change compo attrs",
			compo:    &Boo{Value: "hello"},
			modifier: func(c app.Component) {
				c.(*Boo).Value = "world"
			},
			changes: []Change{
				setTextChange("", "world"),
			},
			compoCount: 2,
		},
		{
			scenario: "replace compo by elem",
			compo:    &Boo{},
			modifier: func(c app.Component) {
				c.(*Boo).ReplaceCompoByElem = true
			},
			changes: []Change{
				createTextChange(""),
				setTextChange("", "foo"),
				createElemChange(newElemNode("p")),
				setAttrsChange("", nil),
				appendChildChange("", ""), // "foo" -> p

				replaceChildChange("", "", ""), // foo <-> p
				deleteNodeChange(""),           // foo.div
				deleteNodeChange(""),           // foo
			},
			compoCount: 1,
		},
		{
			scenario: "replace elem by compo",
			compo:    &Boo{ReplaceCompoByElem: true},
			modifier: func(c app.Component) {
				c.(*Boo).ReplaceCompoByElem = false
			},
			changes: []Change{
				createElemChange(newElemNode("div")),
				setAttrsChange("", map[string]string{"class": "test"}),
				createCompoChange("", "html.foo"),
				setCompoRootChange("", ""),

				replaceChildChange("", "", ""), // p <-> foo
				deleteNodeChange(""),           // "foo"
				deleteNodeChange(""),           // p
			},
			compoCount: 2,
		},
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
			require.Equal(t, test.compoCount, dom.Len())

			jsonChanges, _ := json.MarshalIndent(changes, "", "  ")
			t.Log("changes:", string(jsonChanges))

			require.Len(t, changes, len(test.changes))

			for i := range changes {
				requireEqualChange(t, test.changes[i], changes[i])
			}
		})
	}
}

func TestDOMRenderSubComponent(t *testing.T) {

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

	case createCompo:
		require.Equal(t, expected.Value.(compoValue).Name, actual.Value.(compoValue).Name)

	default:
	}
}
