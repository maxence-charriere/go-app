package dom

import (
	"testing"

	"github.com/murlokswarm/app"
	"github.com/stretchr/testify/assert"
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

func (f *Foo) Subscribe() *app.EventSubscriber {
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
	ReplaceTextByElem bool
	ReplaceElemByElem bool
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
			<dom.Oob>
		{{else}}
			<dom.Foo value="{{.Value}}">
		{{end}}

		{{if .AddCompo}}
			<dom.Foo>
		{{end}}


		{{if .ChildErr}}
			<dom.ErrCompo>
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
			<dom.Foo>
		{{else}}
			<dom.Oob>
		{{end}}
	`
}

type NestedNested struct {
	Foo bool
}

func (n *NestedNested) Render() string {
	return `
		{{if .Foo}}
			<dom.Nested foo>
		{{else}}
			<dom.Nested>
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
			<dom.DecodeErr>
		{{else}}
			<dom.DecodeErr noerr>
		{{end}}

		{{if .NoImport}}
			<dom.unknown>
		{{end}}

		{{if .ReplaceCompoErr}}
			<dom.DecodeErr>
		{{else}}
			<dom.Oob int="{{.Int}}">
		{{end}}

		{{if .AddChildErr}}
			<dom.DecodeErr>
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
	f.RegisterCompo(&Foo{})
	f.RegisterCompo(&Bar{})
	f.RegisterCompo(&Boo{})
	f.RegisterCompo(&Oob{})
	f.RegisterCompo(&Nested{})
	f.RegisterCompo(&NestedNested{})
	f.RegisterCompo(&CompoErr{})
	f.RegisterCompo(&DecodeErr{})
	f.RegisterCompo(NoPtrErr(0))
	f.RegisterCompo(&EmptyStructErr{})

	tests := []struct {
		scenario   string
		compo      app.Compo
		modifier   func(c app.Compo)
		changes    []Change
		compoCount int
		err        bool
	}{
		// Foo:
		{
			scenario: "create simple compo",
			compo:    &Foo{Value: "hello"},
			changes: []Change{
				createTextChange("text:"),
				setTextChange("text:", "hello"),

				createElemChange("div:", "div", ""),
				setAttrsChange("div:", map[string]string{"class": "test"}),
				appendChildChange("div:", "text:"),
				mountElemChange("div:", ""),

				appendChildChange("root:", "div:"), // div -> root
			},
			compoCount: 1,
		},
		{
			scenario: "update simple compo",
			compo:    &Foo{Value: "hello"},
			modifier: func(c app.Compo) {
				c.(*Foo).Value = "world"
			},
			changes: []Change{
				setTextChange("text:", "world"),
			},
			compoCount: 1,
		},
		{
			scenario: "append simple compo child",
			compo:    &Foo{},
			modifier: func(c app.Compo) {
				c.(*Foo).Value = "hello"
			},
			changes: []Change{
				createTextChange("text:"),
				setTextChange("text:", "hello"),
				appendChildChange("div:", "text:"),
			},
			compoCount: 1,
		},
		{
			scenario: "remove simple compo child",
			compo:    &Foo{Value: "hello"},
			modifier: func(c app.Compo) {
				c.(*Foo).Value = ""
			},
			changes: []Change{
				removeChildChange("div:", "text:"),
				deleteNodeChange("text:"),
			},
			compoCount: 1,
		},
		{
			scenario: "change simple compo root attrs",
			compo:    &Foo{},
			modifier: func(c app.Compo) {
				c.(*Foo).Disabled = true
			},
			changes: []Change{
				setAttrsChange("div:", map[string]string{
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
				createTextChange("text:"),
				setTextChange("text:", "hello"),

				createTextChange("text:"),
				setTextChange("text:", "world"),

				createElemChange("h1:", "h1", ""),
				setAttrsChange("h1:", nil),
				appendChildChange("h1:", "text:"),
				mountElemChange("h1:", ""),

				createElemChange("div:", "div", ""),
				setAttrsChange("div:", nil),
				appendChildChange("div:", "text:"),
				appendChildChange("div:", "h1:"),
				mountElemChange("div:", ""),

				appendChildChange("root:", "div:"),
			},
			compoCount: 1,
		},
		{
			scenario: "replace text by elem",
			compo:    &Bar{},
			modifier: func(c app.Compo) {
				c.(*Bar).ReplaceTextByElem = true
			},
			changes: []Change{
				createTextChange("text:"),
				setTextChange("text:", "hello"),

				createElemChange("span:", "span", ""),
				setAttrsChange("span:", nil),
				appendChildChange("span:", "text:"),
				mountElemChange("span:", ""),

				replaceChildChange("div:", "text:", "span:"),
				deleteNodeChange("text:"),
			},
			compoCount: 1,
		},
		{
			scenario: "replace elem by text",
			compo:    &Bar{ReplaceTextByElem: true},
			modifier: func(c app.Compo) {
				c.(*Bar).ReplaceTextByElem = false
			},
			changes: []Change{
				createTextChange("text:"),
				setTextChange("text:", "hello"),
				replaceChildChange("div:", "span:", "text:"),

				deleteNodeChange("text:"),
				deleteNodeChange("span:"),
			},
			compoCount: 1,
		},
		{
			scenario: "replace elem by elem",
			compo:    &Bar{},
			modifier: func(c app.Compo) {
				c.(*Bar).ReplaceElemByElem = true
			},
			changes: []Change{
				createTextChange("text:"),
				setTextChange("text:", "world"),

				createElemChange("h2:", "h2", ""),
				setAttrsChange("h2:", nil),
				appendChildChange("h2:", "text:"),
				mountElemChange("h2:", ""),

				replaceChildChange("div:", "h1:", "h2:"),
				deleteNodeChange("text:"),
				deleteNodeChange("h1:"),
			},
			compoCount: 1,
		},

		// Boo:
		{
			scenario: "create compo with nested compo",
			compo:    &Boo{},
			changes: []Change{
				createElemChange("div:", "div", ""),
				setAttrsChange("div:", map[string]string{"class": "test"}),
				mountElemChange("div:", ""),

				createCompoChange("dom.foo:", "dom.foo"),
				setCompoRootChange("dom.foo:", "div:"),

				createElemChange("div:", "div", ""),
				setAttrsChange("div:", nil),
				appendChildChange("div:", "dom.foo:"),
				mountElemChange("div:", ""),

				appendChildChange("root:", "div:"),
			},
			compoCount: 2,
		},
		{
			scenario: "add compo to elem",
			compo:    &Boo{},
			modifier: func(c app.Compo) {
				c.(*Boo).AddCompo = true
			},
			changes: []Change{
				createElemChange("div:", "div", ""),
				setAttrsChange("div:", map[string]string{"class": "test"}),
				mountElemChange("div:", ""),

				createCompoChange("dom.foo:", "dom.foo"),
				setCompoRootChange("dom.foo:", "div:"),

				appendChildChange("div:", "dom.foo:"),
			},
			compoCount: 3,
		},
		{
			scenario: "remove compo from elem",
			compo:    &Boo{AddCompo: true},
			modifier: func(c app.Compo) {
				c.(*Boo).AddCompo = false
			},
			changes: []Change{
				removeChildChange("div:", "dom.foo:"),
				deleteNodeChange("div:"),
				deleteNodeChange("dom.foo:"),
			},
			compoCount: 2,
		},
		{
			scenario: "replace compo type",
			compo:    &Boo{},
			modifier: func(c app.Compo) {
				c.(*Boo).ReplaceCompoType = true
			},
			changes: []Change{
				createElemChange("p:", "p", ""),
				setAttrsChange("p:", nil),
				mountElemChange("p:", ""),

				createCompoChange("dom.oob:", "dom.oob"),
				setCompoRootChange("dom.oob:", "p:"),

				replaceChildChange("div:", "dom.foo:", "dom.oob:"),
				deleteNodeChange("div:"),
				deleteNodeChange("dom.foo:"),
			},
			compoCount: 2,
		},
		{
			scenario: "change compo attrs",
			compo:    &Boo{Value: "hello"},
			modifier: func(c app.Compo) {
				c.(*Boo).Value = "world"
			},
			changes: []Change{
				setTextChange("text:", "world"),
			},
			compoCount: 2,
		},
		{
			scenario: "replace compo by elem",
			compo:    &Boo{},
			modifier: func(c app.Compo) {
				c.(*Boo).ReplaceCompoByElem = true
			},
			changes: []Change{
				createTextChange("text:"),
				setTextChange("text:", "foo"),

				createElemChange("p:", "p", ""),
				setAttrsChange("p:", nil),
				appendChildChange("p:", "text:"),
				mountElemChange("p:", ""),

				replaceChildChange("div:", "dom.foo:", "p:"),
				deleteNodeChange("div:"),
				deleteNodeChange("dom.foo:"),
			},
			compoCount: 1,
		},
		{
			scenario: "replace elem by compo",
			compo:    &Boo{ReplaceCompoByElem: true},
			modifier: func(c app.Compo) {
				c.(*Boo).ReplaceCompoByElem = false
			},
			changes: []Change{
				createElemChange("div:", "div", ""),
				setAttrsChange("div:", map[string]string{"class": "test"}),
				mountElemChange("div:", ""),

				createCompoChange("dom.foo:", "dom.foo"),
				setCompoRootChange("dom.foo:", "div:"),

				replaceChildChange("div:", "p:", "dom.foo:"),
				deleteNodeChange("text:"),
				deleteNodeChange("p:"),
			},
			compoCount: 2,
		},

		// Nested:
		{
			scenario: "create nested",
			compo:    &Nested{Foo: true},
			changes: []Change{
				createElemChange("div:", "div", ""),
				setAttrsChange("div:", map[string]string{"class": "test"}),
				mountElemChange("div:", ""),

				createCompoChange("dom.foo:", "dom.foo"),
				setCompoRootChange("dom.foo:", "div:"),

				appendChildChange("root:", "dom.foo:"),
			},
			compoCount: 2,
		},
		{
			scenario: "replace nested",
			compo:    &Nested{},
			modifier: func(c app.Compo) {
				c.(*Nested).Foo = true
			},
			changes: []Change{
				createElemChange("div:", "div", ""), // foo.div
				setAttrsChange("div:", map[string]string{"class": "test"}),
				mountElemChange("div:", ""),

				createCompoChange("dom.foo:", "dom.foo"),
				setCompoRootChange("dom.foo:", "div:"),

				replaceChildChange("root:", "dom.oob:", "dom.foo:"), // oob <=> foo

				deleteNodeChange("p:"),
				deleteNodeChange("dom.oob:"),
			},
			compoCount: 2,
		},
		{
			scenario: "create nested nested",
			compo:    &NestedNested{},
			changes: []Change{
				createElemChange("p:", "p", ""),
				setAttrsChange("p:", nil),
				mountElemChange("p:", ""),

				createCompoChange("dom:oob:", "dom.oob"),
				setCompoRootChange("dom:oob:", "p:"),

				createCompoChange("dom.nested:", "dom.nested"),
				setCompoRootChange("dom.nested:", "dom.oob:"),

				appendChildChange("root:", "dom.nested:"),
			},
			compoCount: 3,
		},
		{
			scenario: "replace nested nested",
			compo:    &NestedNested{},
			modifier: func(c app.Compo) {
				c.(*NestedNested).Foo = true
			},
			changes: []Change{
				createElemChange("div:", "div", ""),
				setAttrsChange("div:", map[string]string{"class": "test"}),
				mountElemChange("div:", ""),

				createCompoChange("dom.foo:", "dom.foo"),
				setCompoRootChange("dom.foo:", "div:"),

				deleteNodeChange("p:"),
				deleteNodeChange("dom.oob:"),

				setCompoRootChange("dom.nested:", "dom.foo:"),
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
			modifier: func(c app.Compo) {
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
			modifier: func(c app.Compo) {
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
			modifier: func(c app.Compo) {
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
			modifier: func(c app.Compo) {
				c.(*CompoErr).ReplaceCompoErr = true
			},
			err: true,
		},
		{
			scenario: "fail add child",
			compo:    &CompoErr{Int: 42},
			modifier: func(c app.Compo) {
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
			dom := NewDOM(f, JsToGoHandler, HrefCompoFmt)
			defer func() {
				dom.Clean()
				assert.Empty(t, dom.compoByID)
				assert.Empty(t, dom.compoByCompo)
			}()

			changes, err := dom.New(test.compo)

			if test.modifier != nil {
				test.modifier(test.compo)
				changes, err = dom.Update(test.compo)
			}

			if test.err {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.compoCount, dom.Len())
			require.True(t, dom.Contains(test.compo))

			t.Log(prettyChanges(changes))
			assertChangesEqual(t, test.changes, changes)
		})
	}
}

func TestDOMUpdateErrors(t *testing.T) {
	dom := NewDOM(app.NewFactory(), JsToGoHandler, HrefCompoFmt)
	defer func() {
		dom.Clean()
		assert.Empty(t, dom.compoByID)
		assert.Empty(t, dom.compoByCompo)
	}()

	_, err := dom.Update(&EmptyStructErr{})
	assert.Error(t, err)

	_, err = dom.Update(&Foo{})
	assert.Error(t, err)
}

func TestDOMCompoByID(t *testing.T) {
	f := app.NewFactory()
	f.RegisterCompo(&Foo{})

	dom := NewDOM(f, JsToGoHandler, HrefCompoFmt)
	foo := &Foo{}

	_, err := dom.New(foo)
	require.NoError(t, err)

	var c *component
	for _, c = range dom.compoByCompo {
		break
	}
	require.Equal(t, foo, c.compo)

	var foo2 app.Compo
	foo2, err = dom.CompoByID(c.id)
	require.NoError(t, err)
	require.Equal(t, foo, foo2)

	_, err = dom.CompoByID("hello")
	require.Error(t, err)
}
