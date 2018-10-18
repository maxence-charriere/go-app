package dom

import (
	"encoding/json"
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

func TestEngine(t *testing.T) {
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
		scenario     string
		allowedNodes []string
		compo        app.Compo
		mutate       func(c app.Compo)
		changes      []change
		compoCount   int
		nodeCount    int
		err          bool
	}{
		// Foo:
		{
			scenario: "create simple compo",
			compo:    &Foo{Value: "hello"},
			changes: []change{
				{Action: newNode, NodeID: "dom.foo:", Type: "dom.foo", IsCompo: true},
				{Action: newNode, NodeID: "div:", Type: "div"},
				{Action: newNode, NodeID: "text:", Type: "text"},

				{Action: setAttr, NodeID: "div:", Key: "class", Value: "test"},
				{Action: setText, NodeID: "text:", Value: "hello"},
				{Action: appendChild, NodeID: "div:", ChildID: "text:"},
				{Action: appendChild, NodeID: "dom.foo:", ChildID: "div:"},
				{Action: setRoot, NodeID: "dom.foo:"},
			},
			compoCount: 1,
			nodeCount:  3,
		},
		{
			scenario: "append simple compo child",
			compo:    &Foo{},
			mutate: func(c app.Compo) {
				c.(*Foo).Value = "hello"
			},
			changes: []change{
				{Action: newNode, NodeID: "text:", Type: "text"},

				{Action: setText, NodeID: "text:", Value: "hello"},
				{Action: appendChild, NodeID: "div:", ChildID: "text:"},
			},
			compoCount: 1,
			nodeCount:  3,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			changes := []change{}

			e := Engine{
				Factory:      f,
				AllowedNodes: test.allowedNodes,
				Sync: func(v interface{}) error {
					changes, _ = v.([]change)
					return nil
				},
			}

			defer func() {
				e.Close()

				require.Empty(t, e.compos)
				require.Empty(t, e.compoIDs)
				require.Empty(t, e.nodes)

				require.Empty(t, e.creates)
				require.Empty(t, e.changes)
				require.Empty(t, e.deletes)
				require.Empty(t, e.toSync)
			}()

			err := e.New(test.compo)

			if test.mutate != nil {
				test.mutate(test.compo)
				err = e.Render(test.compo)
			}

			if test.err {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			t.Log(pretty(changes))

			require.Len(t, e.compos, test.compoCount)
			require.Len(t, e.compoIDs, test.compoCount)
			require.Len(t, e.nodes, test.nodeCount)
			require.NotEmpty(t, e.rootID)
			requireChangesMatches(t, test.changes, changes)
		})
	}
}

func pretty(v interface{}) string {
	s, _ := json.MarshalIndent(v, "", "    ")
	return string(s)
}
