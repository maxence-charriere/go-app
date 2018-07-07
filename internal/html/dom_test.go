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
	ChildErr           bool
	ChildNoImport      bool
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


		{{if .ChildErr}}
			<html.ErrCompo>
		{{end}}

		{{if .ChildNoImport}}
			<unknown>
		{{end}}
	</div>
	`
}

type Oob struct {
	Int int
}

func (o *Oob) Render() string {
	return `<p>{{if .Int}}{{.Int}}{{end}}</p>`
}

type Nested struct {
	Foo bool
}

func (n *Nested) Render() string {
	return `
		{{if .Foo}}
			<html.Foo>
		{{else}}
			<html.Oob>
		{{end}}
	`
}

type NestedNested struct {
	Foo bool
}

func (n *NestedNested) Render() string {
	return `
		{{if .Foo}}
			<html.Nested foo>
		{{else}}
			<html.Nested>
		{{end}}
	`
}

type CompoErr struct {
	DecodeErr       bool
	NoImport        bool
	ReplaceCompoErr bool
	AddChildErr     bool
	Int             interface{}
}

func (c *CompoErr) Render() string {
	return `
	<div>
		{{if .DecodeErr}}
			<html.DecodeErr>
		{{else}}
			<html.DecodeErr noerr>
		{{end}}

		{{if .NoImport}}
			<html.unknown>
		{{end}}

		{{if .ReplaceCompoErr}}
			<html.DecodeErr>
		{{else}}
			<html.Oob int="{{.Int}}">
		{{end}}

		{{if .AddChildErr}}
			<html.DecodeErr>
		{{end}}
	</div>
	`
}

type DecodeErr struct {
	NoErr bool
}

func (d *DecodeErr) Render() string {
	return `
	{{if .NoErr}}
		<div></div>
	{{else}}
		<div><div %error></div>
	{{end}}
	`
}

type NoPtrErr int

func (e NoPtrErr) Render() string {
	return `<p>42</p>`
}

type EmptyStructErr struct{}

func (e *EmptyStructErr) Render() string {
	return `<p></p>`
}

func TestDOM(t *testing.T) {
	f := app.NewFactory()
	f.Register(&Foo{})
	f.Register(&Bar{})
	f.Register(&Boo{})
	f.Register(&Oob{})
	f.Register(&Nested{})
	f.Register(&NestedNested{})
	f.Register(&CompoErr{})
	f.Register(&DecodeErr{})
	f.Register(NoPtrErr(0))
	f.Register(&EmptyStructErr{})

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
				createElemChange("", "div"),
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
				createElemChange("", "h1"),
				setAttrsChange("", nil),
				appendChildChange("", ""), // world -> h1

				createElemChange("", "div"),
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
				createElemChange("", "span"),
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
				createElemChange("", "h2"),
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
				createElemChange("", "div"),
				setAttrsChange("", map[string]string{"class": "test"}),
				createCompoChange("", "html.foo"),
				setCompoRootChange("", ""),

				createElemChange("", "div"),
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
				createElemChange("", "div"),
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
				createElemChange("", "p"), // oob.p
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
				createElemChange("", "p"),
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
				createElemChange("", "div"),
				setAttrsChange("", map[string]string{"class": "test"}),
				createCompoChange("", "html.foo"),
				setCompoRootChange("", ""),

				replaceChildChange("", "", ""), // p <-> foo
				deleteNodeChange(""),           // "foo"
				deleteNodeChange(""),           // p
			},
			compoCount: 2,
		},

		// Nested:
		{
			scenario: "create nested",
			compo:    &Nested{Foo: true},
			changes: []Change{
				createElemChange("", "div"), // foo.div
				setAttrsChange("", map[string]string{"class": "test"}),
				createCompoChange("", "html.foo"),
				setCompoRootChange("", ""),
				appendChildChange("", ""), // foo -> root
			},
			compoCount: 2,
		},
		{
			scenario: "replace nested",
			compo:    &Nested{},
			modifier: func(c app.Component) {
				c.(*Nested).Foo = true
			},
			changes: []Change{
				createElemChange("", "div"), // foo.div
				setAttrsChange("", map[string]string{"class": "test"}),
				createCompoChange("", "html.foo"),
				setCompoRootChange("", ""),
				replaceChildChange("", "", ""), // oob <=> foo

				deleteNodeChange(""), // oob.p
				deleteNodeChange(""), // oob
			},
			compoCount: 2,
		},
		{
			scenario: "create nested nested",
			compo:    &NestedNested{},
			changes: []Change{
				createElemChange("", "p"), // oob.p
				setAttrsChange("", nil),
				createCompoChange("", "html.oob"),
				setCompoRootChange("", ""),

				createCompoChange("", "html.nested"),
				setCompoRootChange("", ""), // nested.oob => nested

				appendChildChange("", ""), // nestednested -> root
			},
			compoCount: 3,
		},
		{
			scenario: "replace nested nested",
			compo:    &NestedNested{},
			modifier: func(c app.Component) {
				c.(*NestedNested).Foo = true
			},
			changes: []Change{
				createElemChange("", "div"), // foo.div
				setAttrsChange("", map[string]string{"class": "test"}),
				createCompoChange("", "html.foo"),
				setCompoRootChange("", ""),

				deleteNodeChange(""), // oob.p
				deleteNodeChange(""), // oob

				setCompoRootChange("", ""), // foo => nested
			},
			compoCount: 3,
		},

		// Err:
		{
			scenario: "fail decode",
			compo:    &DecodeErr{},
			err:      true,
		},
		{
			scenario: "fail decode update",
			compo:    &DecodeErr{NoErr: true},
			modifier: func(c app.Component) {
				c.(*DecodeErr).NoErr = false
			},
			err: true,
		},
		{
			scenario: "fail decode child",
			compo:    &CompoErr{DecodeErr: true},
			err:      true,
		},
		{
			scenario: "fail decode child update",
			compo:    &CompoErr{Int: 0},
			modifier: func(c app.Component) {
				c.(*CompoErr).DecodeErr = true
			},
			err: true,
		},
		{
			scenario: "fail map child fields",
			compo:    &CompoErr{Int: 42.42},
			err:      true,
		},
		{
			scenario: "fail update child fields",
			compo:    &CompoErr{Int: 42},
			modifier: func(c app.Component) {
				c.(*CompoErr).Int = 42.42
			},
			err: true,
		},
		{
			scenario: "fail child no import",
			compo:    &CompoErr{NoImport: true},
			err:      true,
		},
		{
			scenario: "replace compo err",
			compo:    &CompoErr{Int: 0},
			modifier: func(c app.Component) {
				c.(*CompoErr).ReplaceCompoErr = true
			},
			err: true,
		},
		{
			scenario: "fail add child",
			compo:    &CompoErr{Int: 42},
			modifier: func(c app.Component) {
				c.(*CompoErr).AddChildErr = true
			},
			err: true,
		},
		{
			scenario: "no ptr compo",
			compo:    NoPtrErr(42),
			err:      true,
		},
		{
			scenario: "empty compo",
			compo:    &EmptyStructErr{},
			err:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			dom := NewDOM(f, "test")
			changes, err := dom.Render(test.compo)

			if test.modifier != nil {
				test.modifier(test.compo)
				changes, err = dom.Render(test.compo)
			}

			if test.err {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.compoCount, dom.Len())

			jsonChanges, _ := json.MarshalIndent(changes, "", "  ")
			t.Log("changes:", string(jsonChanges))

			require.Len(t, changes, len(test.changes))
			require.True(t, dom.Contains(test.compo))

			for i := range changes {
				requireEqualChange(t, test.changes[i], changes[i])
			}
		})
	}
}

func TestRenderNewRoot(t *testing.T) {
	f := app.NewFactory()
	f.Register(&Oob{})

	dom := newDOM(f, "test")
	_, err := dom.Render(&Oob{})
	require.NoError(t, err)

	var changes []Change
	changes, err = dom.Render(&Oob{})
	require.NoError(t, err)

	jsonChanges, _ := json.MarshalIndent(changes, "", "  ")
	t.Log("changes:", string(jsonChanges))

	expected := []Change{
		createElemChange("", "p"),
		setAttrsChange("", nil),
		removeChildChange("", ""), // old oob.p
		deleteNodeChange(""),      // old oob.p
		appendChildChange("", ""), // new oob.p
	}
	require.Len(t, changes, len(expected))

	for i := range changes {
		requireEqualChange(t, expected[i], changes[i])
	}
}

func TestDOMComponentByID(t *testing.T) {
	f := app.NewFactory()
	f.Register(&Foo{})

	dom := newDOM(f, "test")
	foo := &Foo{}
	_, err := dom.Render(foo)
	require.NoError(t, err)

	var row compoRow
	for _, r := range dom.compoRowByCompo {
		row = r
		break
	}
	require.Equal(t, foo, row.component)

	var c app.Component
	c, err = dom.ComponentByID(row.id)
	require.NoError(t, err)
	require.Equal(t, foo, c)

	_, err = dom.ComponentByID("hello")
	require.Error(t, err)
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
