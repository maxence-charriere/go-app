package dom

import (
	"encoding/json"
	"testing"

	"github.com/murlokswarm/app"
	"github.com/pkg/errors"
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
	ReplaceTextByNode bool
	ReplaceNodeByNode bool
}

func (b *Bar) Render() string {
	return `
	<div>
		{{if .ReplaceTextByNode}}
			<span>hello</span>
		{{else}}
			hello
		{{end}}

		{{if .ReplaceNodeByNode}}
			<h2>world</h2>
		{{else}}
			<h1>world</h1>
		{{end}}
	</div>
	`
}

type Boo struct {
	ReplaceCompoByNode  bool
	ReplaceCompoByCompo bool
	AddCompo            bool
	ChildErr            bool
	ChildNoImport       bool
	Value               string
}

func (b *Boo) Render() string {
	return `
	<div>
		{{if .ReplaceCompoByNode}}
			<p>foo</p>
		{{else if .ReplaceCompoByCompo}}
			<dom.Oob />
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
	Int             int
	BadExtendedFunc bool
}

func (o *Oob) Funcs() map[string]interface{} {
	return map[string]interface{}{
		"hello": func(s string) string {
			return "hello " + s
		},
	}
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

type Svg struct {
	Path string
}

func (s *Svg) Render() string {
	return `
	<svg>
		<path data="{{.Path}}"></path>
		<path data="" />
	</svg>
	`
}

type SelfClosing struct {
	NoClose bool
	Svg     bool
}

func (c *SelfClosing) Render() string {
	return `
	<div>
		{{if .NoClose}}
			<div>
				<p></p>
			</div>
		{{else}}
			<div />
		{{end}}

		{{if .Svg}}
		<svg />
		{{end}}
	</div>
	`
}

type VoidElem app.ZeroCompo

func (v *VoidElem) Render() string {
	return `
	<div>
		<img>
		<p></p>
	</div>
	`
}

type CompoErr struct {
	TemplateReadErr bool
	TemplateExecErr bool
	DecodeErr       bool
	BadExtendedFunc bool
	NoImport        bool
	ReplaceCompoErr bool
	AddChildErr     bool
	Int             interface{}
}

func (c *CompoErr) Funcs() map[string]interface{} {
	if c.BadExtendedFunc {
		return map[string]interface{}{
			"raw": func(s string) string {
				panic("should not be overridden")
			},
		}
	}

	return nil
}

func (c *CompoErr) Render() string {
	return `
	<!DOCTYPE html>
	<div>
		{{if .TemplateReadErr}}
			<dom.BadTemplateRead err>
		{{end}}

		{{if .TemplateExecErr}}
			<dom.BadTemplateExec err>
		{{end}}

		{{if .DecodeErr}}
			<dom.DecodeErr>
		{{end}}

		{{if .NoImport}}
			<dom.unknown>
		{{end}}

		{{if .ReplaceCompoErr}}
			<dom.badtemplate TemplateExecErr>
		{{else}}
			<dom.Oob int="{{.Int}}">
		{{end}}

		{{if .AddChildErr}}
			<dom.DecodeErr>
		{{end}}
	</div>
	`
}

type BadTemplateRead struct {
	Err bool
}

func (b *BadTemplateRead) Render() string {
	return `
	<div>
		{{if .Err}}
			{{print :)}}
		{{else}}
			<div></div>
		{{end}}
	</div>
	`
}

type BadTemplateExec struct {
	Err bool
}

func (b *BadTemplateExec) Render() string {
	return `
	<div>
		{{if .Err}}
			{{.KDNDSLndslj}}
		{{else}}
			<div></div>
		{{end}}
	</div>
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

type DecodeErr app.ZeroCompo

func (d *DecodeErr) Render() string {
	return `<div %error="42">`
}

type EmptyRender app.ZeroCompo

func (e *EmptyRender) Render() string {
	return ""
}

func TestEngine(t *testing.T) {
	f := app.NewFactory()
	f.RegisterCompo(&Foo{})
	f.RegisterCompo(&Bar{})
	f.RegisterCompo(&Boo{})
	f.RegisterCompo(&Oob{})
	f.RegisterCompo(&Nested{})
	f.RegisterCompo(&NestedNested{})
	f.RegisterCompo(&Svg{})
	f.RegisterCompo(&SelfClosing{})
	f.RegisterCompo(&VoidElem{})
	f.RegisterCompo(&CompoErr{})
	f.RegisterCompo(&BadTemplateRead{})
	f.RegisterCompo(&BadTemplateExec{})
	f.RegisterCompo(&DecodeErr{})
	f.RegisterCompo(NoPtrErr(0))
	f.RegisterCompo(&EmptyStructErr{})
	f.RegisterCompo(&EmptyRender{})

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
			scenario: "create compo nodes",
			compo:    &Foo{Value: "hello"},
			changes: []change{
				{Action: newNode, NodeID: "dom.foo:", Type: "dom.foo", IsCompo: true},
				{Action: newNode, NodeID: "div:", Type: "div", CompoID: "dom.foo:"},
				{Action: newNode, NodeID: "text:", Type: "text", CompoID: "dom.foo:"},

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
			scenario: "update node",
			compo:    &Foo{Value: "hello"},
			mutate: func(c app.Compo) {
				c.(*Foo).Value = "world"
			},
			changes: []change{
				{Action: setText, NodeID: "text:", Value: "world"},
			},
			compoCount: 1,
			nodeCount:  3,
		},
		{
			scenario: "append child",
			compo:    &Foo{},
			mutate: func(c app.Compo) {
				c.(*Foo).Value = "hello"
			},
			changes: []change{
				{Action: newNode, NodeID: "text:", Type: "text", CompoID: "dom.foo:"},

				{Action: setText, NodeID: "text:", Value: "hello"},
				{Action: appendChild, NodeID: "div:", ChildID: "text:"},
			},
			compoCount: 1,
			nodeCount:  3,
		},
		{
			scenario: "remove child",
			compo:    &Foo{Value: "hello"},
			mutate: func(c app.Compo) {
				c.(*Foo).Value = ""
			},
			changes: []change{
				{Action: removeChild, NodeID: "div:", ChildID: "text:"},
				{Action: delNode, NodeID: "text:"},
			},
			compoCount: 1,
			nodeCount:  2,
		},
		{
			scenario: "set attr",
			compo:    &Foo{},
			mutate: func(c app.Compo) {
				c.(*Foo).Disabled = true
			},
			changes: []change{
				{Action: setAttr, NodeID: "div:", Key: "disabled"},
			},
			compoCount: 1,
			nodeCount:  2,
		},
		{
			scenario: "delete attr",
			compo:    &Foo{Disabled: true},
			mutate: func(c app.Compo) {
				c.(*Foo).Disabled = false
			},
			changes: []change{
				{Action: delAttr, NodeID: "div:", Key: "disabled"},
			},
			compoCount: 1,
			nodeCount:  2,
		},

		// Bar:
		{
			scenario: "replace text by node",
			compo:    &Bar{},
			mutate: func(c app.Compo) {
				c.(*Bar).ReplaceTextByNode = true
			},
			changes: []change{
				{Action: newNode, NodeID: "span:", Type: "span", CompoID: "dom.bar:"},
				{Action: newNode, NodeID: "text:", Type: "text", CompoID: "dom.bar:"},

				{Action: setText, NodeID: "text:", Value: "hello"},
				{Action: appendChild, NodeID: "span:", ChildID: "text:"},
				{Action: replaceChild, NodeID: "div:", ChildID: "text:", NewChildID: "span:"},

				{Action: delNode, NodeID: "text:"},
			},
			compoCount: 1,
			nodeCount:  6,
		},
		{
			scenario: "replace node by text",
			compo:    &Bar{ReplaceTextByNode: true},
			mutate: func(c app.Compo) {
				c.(*Bar).ReplaceTextByNode = false
			},
			changes: []change{
				{Action: newNode, NodeID: "text:", Type: "text", CompoID: "dom.bar:"},

				{Action: setText, NodeID: "text:", Value: "hello"},
				{Action: replaceChild, NodeID: "div:", ChildID: "span:", NewChildID: "text:"},

				{Action: delNode, NodeID: "text:"},
				{Action: delNode, NodeID: "span:"},
			},
			compoCount: 1,
			nodeCount:  5,
		},
		{
			scenario: "replace node by node",
			compo:    &Bar{},
			mutate: func(c app.Compo) {
				c.(*Bar).ReplaceNodeByNode = true
			},
			changes: []change{
				{Action: newNode, NodeID: "h2:", Type: "h2", CompoID: "dom.bar:"},
				{Action: newNode, NodeID: "text:", Type: "text", CompoID: "dom.bar:"},

				{Action: setText, NodeID: "text:", Value: "world"},
				{Action: appendChild, NodeID: "h2:", ChildID: "text:"},
				{Action: replaceChild, NodeID: "div:", ChildID: "h1:", NewChildID: "h2:"},

				{Action: delNode, NodeID: "text:"},
				{Action: delNode, NodeID: "h1:"},
			},
			compoCount: 1,
			nodeCount:  5,
		},

		// Boo:
		{
			scenario: "create nested compo",
			compo:    &Boo{},
			changes: []change{
				{Action: newNode, NodeID: "dom.boo:", Type: "dom.boo", IsCompo: true},
				{Action: newNode, NodeID: "div:", Type: "div", CompoID: "dom.boo:"},

				{Action: newNode, NodeID: "dom.foo:", Type: "dom.foo", IsCompo: true, CompoID: "dom.boo:"},
				{Action: newNode, NodeID: "div:", Type: "div", CompoID: "dom.foo:"},
				{Action: setAttr, NodeID: "div:", Key: "class", Value: "test"},
				{Action: appendChild, NodeID: "dom.foo:", ChildID: "div:"},

				{Action: appendChild, NodeID: "div:", ChildID: "dom.foo:"},
				{Action: appendChild, NodeID: "dom.boo:", ChildID: "div:"},
				{Action: setRoot, NodeID: "dom.boo:"},
			},
			compoCount: 2,
			nodeCount:  4,
		},
		{
			scenario: "add compo",
			compo:    &Boo{},
			mutate: func(c app.Compo) {
				c.(*Boo).AddCompo = true
			},
			changes: []change{
				{Action: newNode, NodeID: "dom.foo:", Type: "dom.foo", IsCompo: true, CompoID: "dom.boo:"},
				{Action: newNode, NodeID: "div:", Type: "div", CompoID: "dom.foo:"},

				{Action: setAttr, NodeID: "div:", Key: "class", Value: "test"},
				{Action: appendChild, NodeID: "dom.foo:", ChildID: "div:"},
				{Action: appendChild, NodeID: "div:", ChildID: "dom.foo:"},
			},
			compoCount: 3,
			nodeCount:  6,
		},
		{
			scenario: "remove compo",
			compo:    &Boo{AddCompo: true},
			mutate: func(c app.Compo) {
				c.(*Boo).AddCompo = false
			},
			changes: []change{
				{Action: removeChild, NodeID: "div:", ChildID: "dom.foo:"},

				{Action: delNode, NodeID: "div:"},
				{Action: delNode, NodeID: "dom.foo:"},
			},
			compoCount: 2,
			nodeCount:  4,
		},
		{
			scenario: "replace compo by compo",
			compo:    &Boo{},
			mutate: func(c app.Compo) {
				c.(*Boo).ReplaceCompoByCompo = true
			},
			changes: []change{
				{Action: newNode, NodeID: "dom.oob:", Type: "dom.oob", IsCompo: true, CompoID: "dom.boo:"},
				{Action: newNode, NodeID: "p:", Type: "p", CompoID: "dom.oob:"},

				{Action: appendChild, NodeID: "dom.oob:", ChildID: "p:"},
				{Action: replaceChild, NodeID: "div:", ChildID: "dom.foo:", NewChildID: "dom.oob:"},

				{Action: delNode, NodeID: "div:"},
				{Action: delNode, NodeID: "dom.foo:"},
			},
			compoCount: 2,
			nodeCount:  4,
		},
		{
			scenario: "set compo attr",
			compo:    &Boo{Value: "hello"},
			mutate: func(c app.Compo) {
				c.(*Boo).Value = "world"
			},
			changes: []change{
				{Action: setText, NodeID: "text:", Value: "world"},
			},
			compoCount: 2,
			nodeCount:  5,
		},
		{
			scenario: "replace compo by node",
			compo:    &Boo{},
			mutate: func(c app.Compo) {
				c.(*Boo).ReplaceCompoByNode = true
			},
			changes: []change{
				{Action: newNode, NodeID: "p:", Type: "p", CompoID: "dom.boo:"},
				{Action: newNode, NodeID: "text:", Type: "text", CompoID: "dom.boo:"},

				{Action: setText, NodeID: "text:", Value: "foo"},
				{Action: appendChild, NodeID: "p:", ChildID: "text:"},
				{Action: replaceChild, NodeID: "div:", ChildID: "dom.foo:", NewChildID: "p:"},

				{Action: delNode, NodeID: "div:"},
				{Action: delNode, NodeID: "dom.foo:"},
			},
			compoCount: 1,
			nodeCount:  4,
		},
		{
			scenario: "replace node by compo",
			compo:    &Boo{ReplaceCompoByNode: true},
			mutate: func(c app.Compo) {
				c.(*Boo).ReplaceCompoByNode = false
			},
			changes: []change{
				{Action: newNode, NodeID: "dom.foo:", Type: "dom.foo", IsCompo: true, CompoID: "dom.boo:"},
				{Action: newNode, NodeID: "div:", Type: "div", CompoID: "dom.foo:"},

				{Action: setAttr, NodeID: "div:", Key: "class", Value: "test"},
				{Action: appendChild, NodeID: "dom.foo:", ChildID: "div:"},
				{Action: replaceChild, NodeID: "div:", ChildID: "p:", NewChildID: "dom.foo:"},

				{Action: delNode, NodeID: "text:"},
				{Action: delNode, NodeID: "p:"},
			},
			compoCount: 2,
			nodeCount:  4,
		},

		// Nested:
		{
			scenario: "replace compo first child",
			compo:    &Nested{},
			mutate: func(c app.Compo) {
				c.(*Nested).Foo = true
			},
			changes: []change{
				{Action: newNode, NodeID: "dom.foo:", Type: "dom.foo", IsCompo: true, CompoID: "dom.nested:"},
				{Action: newNode, NodeID: "div:", Type: "div", CompoID: "dom.foo:"},

				{Action: setAttr, NodeID: "div:", Key: "class", Value: "test"},
				{Action: appendChild, NodeID: "dom.foo:", ChildID: "div:"},
				{Action: replaceChild, NodeID: "dom.nested:", ChildID: "dom.oob:", NewChildID: "dom.foo:"},

				{Action: delNode, NodeID: "p:"},
				{Action: delNode, NodeID: "dom.oob:"},
			},
			compoCount: 2,
			nodeCount:  3,
		},
		{
			scenario: "replace nested compo first child",
			compo:    &NestedNested{},
			mutate: func(c app.Compo) {
				c.(*NestedNested).Foo = true
			},
			changes: []change{
				{Action: newNode, NodeID: "dom.foo:", Type: "dom.foo", IsCompo: true, CompoID: "dom.nested:"},
				{Action: newNode, NodeID: "div:", Type: "div", CompoID: "dom.foo:"},

				{Action: setAttr, NodeID: "div:", Key: "class", Value: "test"},
				{Action: appendChild, NodeID: "dom.foo:", ChildID: "div:"},
				{Action: replaceChild, NodeID: "dom.nested:", ChildID: "dom.oob:", NewChildID: "dom.foo:"},

				{Action: delNode, NodeID: "p:"},
				{Action: delNode, NodeID: "dom.oob:"},
			},
			compoCount: 3,
			nodeCount:  4,
		},

		// Svg:
		{
			scenario: "create node with namespace",
			compo:    &Svg{},
			changes: []change{
				{Action: newNode, NodeID: "dom.svg:", Type: "dom.svg", IsCompo: true},
				{Action: newNode, NodeID: "svg:", Type: "svg", Namespace: svg, CompoID: "dom.svg:"},
				{Action: newNode, NodeID: "path:", Type: "path", Namespace: svg, CompoID: "dom.svg:"},
				{Action: newNode, NodeID: "path:", Type: "path", Namespace: svg, CompoID: "dom.svg:"},

				{Action: setAttr, NodeID: "path:", Key: "data"},
				{Action: appendChild, NodeID: "svg:", ChildID: "path:"},
				{Action: setAttr, NodeID: "path:", Key: "data"},
				{Action: appendChild, NodeID: "svg:", ChildID: "path:"},
				{Action: appendChild, NodeID: "dom.svg:", ChildID: "svg:"},

				{Action: setRoot, NodeID: "dom.svg:"},
			},
			compoCount: 1,
			nodeCount:  4,
		},
		{
			scenario: "update node with namespace",
			compo:    &Svg{},
			mutate: func(c app.Compo) {
				c.(*Svg).Path = "M42"
			},
			changes: []change{
				{Action: setAttr, NodeID: "path:", Key: "data", Value: "M42"},
			},
			compoCount: 1,
			nodeCount:  4,
		},

		// Self closing:
		{
			scenario: "replace node by self closing node",
			compo:    &SelfClosing{NoClose: true},
			mutate: func(c app.Compo) {
				c.(*SelfClosing).NoClose = false
			},
			changes: []change{
				{Action: removeChild, NodeID: "div:", ChildID: "p:"},
				{Action: delNode, NodeID: "p:"},
			},
			compoCount: 1,
			nodeCount:  3,
		},
		{
			scenario: "self closing svg",
			compo:    &SelfClosing{Svg: true},
			changes: []change{
				{Action: newNode, NodeID: "dom.selfclosing:", Type: "dom.selfclosing", IsCompo: true},
				{Action: newNode, NodeID: "div:", Type: "div", CompoID: "dom.selfclosing:"},
				{Action: newNode, NodeID: "div:", Type: "div", CompoID: "dom.selfclosing:"},
				{Action: newNode, NodeID: "svg:", Type: "svg", Namespace: svg, CompoID: "dom.selfclosing:"},

				{Action: appendChild, NodeID: "div:", ChildID: "div:"},
				{Action: appendChild, NodeID: "div:", ChildID: "svg:"},
				{Action: appendChild, NodeID: "dom.selfclosing:", ChildID: "div:"},

				{Action: setRoot, NodeID: "dom.selfclosing:"},
			},
			compoCount: 1,
			nodeCount:  4,
		},

		// Void elem:
		{
			scenario: "void elem node",
			compo:    &VoidElem{},
			changes: []change{
				{Action: newNode, NodeID: "dom.voidelem:", Type: "dom.voidelem", IsCompo: true},
				{Action: newNode, NodeID: "div:", Type: "div", CompoID: "dom.voidelem:"},
				{Action: newNode, NodeID: "img:", Type: "img", CompoID: "dom.voidelem:"},
				{Action: newNode, NodeID: "p:", Type: "p", CompoID: "dom.voidelem:"},

				{Action: appendChild, NodeID: "div:", ChildID: "img:"},
				{Action: appendChild, NodeID: "div:", ChildID: "p:"},
				{Action: appendChild, NodeID: "dom.voidelem:", ChildID: "div:"},

				{Action: setRoot, NodeID: "dom.voidelem:"},
			},
			compoCount: 1,
			nodeCount:  4,
		},

		// Errors:
		{
			scenario: "fail no import",
			compo:    &CompoErr{NoImport: true},
			err:      true,
		},
		{
			scenario: "fail read template",
			compo:    &CompoErr{TemplateReadErr: true},
			err:      true,
		},
		{
			scenario: "fail exec template",
			compo:    &CompoErr{},
			mutate: func(c app.Compo) {
				c.(*CompoErr).TemplateExecErr = true
			},
			err: true,
		},
		{
			scenario: "fail bad extended func",
			compo:    &CompoErr{BadExtendedFunc: true},
			err:      true,
		},
		{
			scenario:     "fail with not allowed node",
			allowedNodes: []string{"menu", "menuitem"},
			compo:        &CompoErr{},
			err:          true,
		},
		{
			scenario:     "fail with not allowed self closing node",
			allowedNodes: []string{"div"},
			compo:        &SelfClosing{Svg: true},
			err:          true,
		},
		{
			scenario: "fail map child fields",
			compo:    &CompoErr{Int: 42.42},
			err:      true,
		},
		{
			scenario: "replace compo err",
			compo:    &CompoErr{Int: 0},
			mutate: func(c app.Compo) {
				c.(*CompoErr).ReplaceCompoErr = true
			},
			err: true,
		},
		{
			scenario: "fail add child",
			compo:    &CompoErr{Int: 42},
			mutate: func(c app.Compo) {
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
		{
			scenario: "empty render",
			compo:    &EmptyRender{},
			err:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			changes := []change{}

			e := Engine{
				Factory:      f,
				AllowedNodes: test.allowedNodes,
				AttrTransforms: []Transform{
					JsToGoHandler,
					HrefCompoFmt,
				},
				Sync: func(v interface{}) error {
					c := v.([]change)
					changes = make([]change, len(c))
					copy(changes, c)
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

			require.True(t, e.Contains(test.compo))
		})
	}
}

func TestEngineRenderNotMounted(t *testing.T) {
	e := Engine{
		Sync: func(v interface{}) error {
			return errors.New("simulated err")
		},
	}

	err := e.Render(&Foo{})
	assert.Error(t, err)
}

func TestEngineSyncError(t *testing.T) {
	f := app.NewFactory()
	f.RegisterCompo(&Foo{})

	e := Engine{
		Factory: f,
		Sync: func(v interface{}) error {
			return errors.New("simulated err")
		},
	}

	err := e.New(&Foo{})
	assert.Error(t, err)
}

func TestEngineEmptySync(t *testing.T) {
	f := app.NewFactory()
	f.RegisterCompo(&Foo{})

	e := Engine{Factory: f}
	err := e.New(&Foo{})
	assert.NoError(t, err)
}

func TestDOMCompoByID(t *testing.T) {
	f := app.NewFactory()
	f.RegisterCompo(&Foo{})

	e := Engine{Factory: f}
	foo := &Foo{}

	err := e.New(foo)
	require.NoError(t, err)

	c, ok := e.compos[foo]
	require.True(t, ok)
	require.Equal(t, foo, c.Compo)

	var foo2 app.Compo
	foo2, err = e.CompoByID(c.ID)
	require.NoError(t, err)
	require.Equal(t, foo, foo2)

	_, err = e.CompoByID("unknownID")
	require.Error(t, err)
}

func pretty(v interface{}) string {
	s, _ := json.MarshalIndent(v, "", "    ")
	return string(s)
}
