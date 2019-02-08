package app

import (
	"encoding/json"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Mur struct {
	Value    string
	Disabled bool
}

func (f *Mur) OnMount() {
}

func (f *Mur) OnDismount() {
}

func (f *Mur) Subscribe() *Subscriber {
	return NewSubscriber()
}

func (f *Mur) Render() string {
	return `
	<div class="test" {{if .Disabled}}disabled{{end}}>
		{{.Value}}
	</div>
	`
}

type Lok struct {
	ReplaceTextByNode bool
	ReplaceNodeByNode bool
}

func (b *Lok) Render() string {
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
			<app.Oob />
		{{else}}
			<app.Mur value="{{.Value}}">
		{{end}}

		{{if .AddCompo}}
			<app.Mur>
		{{end}}


		{{if .ChildErr}}
			<app.ErrCompo>
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
	Mur bool
}

func (n *Nested) Render() string {
	return `
		{{if .Mur}}
			<app.Mur>
		{{else}}
			<app.Oob>
		{{end}}
	`
}

type NestedNested struct {
	Mur bool
}

func (n *NestedNested) Render() string {
	return `
		{{if .Mur}}
			<app.Nested mur>
		{{else}}
			<app.Nested>
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

type VoidElem ZeroCompo

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
			<app.BadTemplateRead err>
		{{end}}

		{{if .TemplateExecErr}}
			<app.BadTemplateExec err>
		{{end}}

		{{if .DecodeErr}}
			<app.DecodeErr>
		{{end}}

		{{if .NoImport}}
			<app.unknown>
		{{end}}

		{{if .ReplaceCompoErr}}
			<app.badtemplate TemplateExecErr>
		{{else}}
			<app.Oob int="{{.Int}}">
		{{end}}

		{{if .AddChildErr}}
			<app.DecodeErr>
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

type DecodeErr ZeroCompo

func (d *DecodeErr) Render() string {
	return `<div %error="42">`
}

type EmptyRender ZeroCompo

func (e *EmptyRender) Render() string {
	return ""
}

func TestEngine(t *testing.T) {
	f := newCompoBuilder()
	f.register(&Mur{})
	f.register(&Lok{})
	f.register(&Boo{})
	f.register(&Oob{})
	f.register(&Nested{})
	f.register(&NestedNested{})
	f.register(&Svg{})
	f.register(&SelfClosing{})
	f.register(&VoidElem{})
	f.register(&CompoErr{})
	f.register(&BadTemplateRead{})
	f.register(&BadTemplateExec{})
	f.register(&DecodeErr{})
	f.register(NoPtrErr(0))
	f.register(&EmptyStructErr{})
	f.register(&EmptyRender{})

	tests := []struct {
		scenario     string
		allowedNodes []string
		compo        Compo
		mutate       func(c Compo)
		changes      []change
		compoCount   int
		nodeCount    int
		err          bool
	}{
		// Mur:
		{
			scenario: "create compo nodes",
			compo:    &Mur{Value: "hello"},
			changes: []change{
				{Action: newNode, NodeID: "app.mur:", Type: "app.mur", IsCompo: true},
				{Action: newNode, NodeID: "div:", Type: "div", CompoID: "app.mur:"},
				{Action: newNode, NodeID: "text:", Type: "text", CompoID: "app.mur:"},

				{Action: setAttr, NodeID: "div:", Key: "class", Value: "test"},
				{Action: setText, NodeID: "text:", Value: "hello"},
				{Action: appendChild, NodeID: "div:", ChildID: "text:"},
				{Action: appendChild, NodeID: "app.mur:", ChildID: "div:"},
				{Action: setRoot, NodeID: "app.mur:"},
			},
			compoCount: 1,
			nodeCount:  3,
		},
		{
			scenario: "update node",
			compo:    &Mur{Value: "hello"},
			mutate: func(c Compo) {
				c.(*Mur).Value = "world"
			},
			changes: []change{
				{Action: setText, NodeID: "text:", Value: "world"},
			},
			compoCount: 1,
			nodeCount:  3,
		},
		{
			scenario: "append child",
			compo:    &Mur{},
			mutate: func(c Compo) {
				c.(*Mur).Value = "hello"
			},
			changes: []change{
				{Action: newNode, NodeID: "text:", Type: "text", CompoID: "app.mur:"},

				{Action: setText, NodeID: "text:", Value: "hello"},
				{Action: appendChild, NodeID: "div:", ChildID: "text:"},
			},
			compoCount: 1,
			nodeCount:  3,
		},
		{
			scenario: "remove child",
			compo:    &Mur{Value: "hello"},
			mutate: func(c Compo) {
				c.(*Mur).Value = ""
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
			compo:    &Mur{},
			mutate: func(c Compo) {
				c.(*Mur).Disabled = true
			},
			changes: []change{
				{Action: setAttr, NodeID: "div:", Key: "disabled"},
			},
			compoCount: 1,
			nodeCount:  2,
		},
		{
			scenario: "delete attr",
			compo:    &Mur{Disabled: true},
			mutate: func(c Compo) {
				c.(*Mur).Disabled = false
			},
			changes: []change{
				{Action: delAttr, NodeID: "div:", Key: "disabled"},
			},
			compoCount: 1,
			nodeCount:  2,
		},

		// Lok:
		{
			scenario: "replace text by node",
			compo:    &Lok{},
			mutate: func(c Compo) {
				c.(*Lok).ReplaceTextByNode = true
			},
			changes: []change{
				{Action: newNode, NodeID: "span:", Type: "span", CompoID: "app.lok:"},
				{Action: newNode, NodeID: "text:", Type: "text", CompoID: "app.lok:"},

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
			compo:    &Lok{ReplaceTextByNode: true},
			mutate: func(c Compo) {
				c.(*Lok).ReplaceTextByNode = false
			},
			changes: []change{
				{Action: newNode, NodeID: "text:", Type: "text", CompoID: "app.lok:"},

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
			compo:    &Lok{},
			mutate: func(c Compo) {
				c.(*Lok).ReplaceNodeByNode = true
			},
			changes: []change{
				{Action: newNode, NodeID: "h2:", Type: "h2", CompoID: "app.lok:"},
				{Action: newNode, NodeID: "text:", Type: "text", CompoID: "app.lok:"},

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
				{Action: newNode, NodeID: "app.boo:", Type: "app.boo", IsCompo: true},
				{Action: newNode, NodeID: "div:", Type: "div", CompoID: "app.boo:"},

				{Action: newNode, NodeID: "app.mur:", Type: "app.mur", IsCompo: true, CompoID: "app.boo:"},
				{Action: newNode, NodeID: "div:", Type: "div", CompoID: "app.mur:"},
				{Action: setAttr, NodeID: "div:", Key: "class", Value: "test"},
				{Action: appendChild, NodeID: "app.mur:", ChildID: "div:"},

				{Action: appendChild, NodeID: "div:", ChildID: "app.mur:"},
				{Action: appendChild, NodeID: "app.boo:", ChildID: "div:"},
				{Action: setRoot, NodeID: "app.boo:"},
			},
			compoCount: 2,
			nodeCount:  4,
		},
		{
			scenario: "add compo",
			compo:    &Boo{},
			mutate: func(c Compo) {
				c.(*Boo).AddCompo = true
			},
			changes: []change{
				{Action: newNode, NodeID: "app.mur:", Type: "app.mur", IsCompo: true, CompoID: "app.boo:"},
				{Action: newNode, NodeID: "div:", Type: "div", CompoID: "app.mur:"},

				{Action: setAttr, NodeID: "div:", Key: "class", Value: "test"},
				{Action: appendChild, NodeID: "app.mur:", ChildID: "div:"},
				{Action: appendChild, NodeID: "div:", ChildID: "app.mur:"},
			},
			compoCount: 3,
			nodeCount:  6,
		},
		{
			scenario: "remove compo",
			compo:    &Boo{AddCompo: true},
			mutate: func(c Compo) {
				c.(*Boo).AddCompo = false
			},
			changes: []change{
				{Action: removeChild, NodeID: "div:", ChildID: "app.mur:"},

				{Action: delNode, NodeID: "div:"},
				{Action: delNode, NodeID: "app.mur:"},
			},
			compoCount: 2,
			nodeCount:  4,
		},
		{
			scenario: "replace compo by compo",
			compo:    &Boo{},
			mutate: func(c Compo) {
				c.(*Boo).ReplaceCompoByCompo = true
			},
			changes: []change{
				{Action: newNode, NodeID: "app.oob:", Type: "app.oob", IsCompo: true, CompoID: "app.boo:"},
				{Action: newNode, NodeID: "p:", Type: "p", CompoID: "app.oob:"},

				{Action: appendChild, NodeID: "app.oob:", ChildID: "p:"},
				{Action: replaceChild, NodeID: "div:", ChildID: "app.mur:", NewChildID: "app.oob:"},

				{Action: delNode, NodeID: "div:"},
				{Action: delNode, NodeID: "app.mur:"},
			},
			compoCount: 2,
			nodeCount:  4,
		},
		{
			scenario: "set compo attr",
			compo:    &Boo{Value: "hello"},
			mutate: func(c Compo) {
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
			mutate: func(c Compo) {
				c.(*Boo).ReplaceCompoByNode = true
			},
			changes: []change{
				{Action: newNode, NodeID: "p:", Type: "p", CompoID: "app.boo:"},
				{Action: newNode, NodeID: "text:", Type: "text", CompoID: "app.boo:"},

				{Action: setText, NodeID: "text:", Value: "foo"},
				{Action: appendChild, NodeID: "p:", ChildID: "text:"},
				{Action: replaceChild, NodeID: "div:", ChildID: "app.mur:", NewChildID: "p:"},

				{Action: delNode, NodeID: "div:"},
				{Action: delNode, NodeID: "app.mur:"},
			},
			compoCount: 1,
			nodeCount:  4,
		},
		{
			scenario: "replace node by compo",
			compo:    &Boo{ReplaceCompoByNode: true},
			mutate: func(c Compo) {
				c.(*Boo).ReplaceCompoByNode = false
			},
			changes: []change{
				{Action: newNode, NodeID: "app.mur:", Type: "app.mur", IsCompo: true, CompoID: "app.boo:"},
				{Action: newNode, NodeID: "div:", Type: "div", CompoID: "app.mur:"},

				{Action: setAttr, NodeID: "div:", Key: "class", Value: "test"},
				{Action: appendChild, NodeID: "app.mur:", ChildID: "div:"},
				{Action: replaceChild, NodeID: "div:", ChildID: "p:", NewChildID: "app.mur:"},

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
			mutate: func(c Compo) {
				c.(*Nested).Mur = true
			},
			changes: []change{
				{Action: newNode, NodeID: "app.mur:", Type: "app.mur", IsCompo: true, CompoID: "app.nested:"},
				{Action: newNode, NodeID: "div:", Type: "div", CompoID: "app.mur:"},

				{Action: setAttr, NodeID: "div:", Key: "class", Value: "test"},
				{Action: appendChild, NodeID: "app.mur:", ChildID: "div:"},
				{Action: replaceChild, NodeID: "app.nested:", ChildID: "app.oob:", NewChildID: "app.mur:"},

				{Action: delNode, NodeID: "p:"},
				{Action: delNode, NodeID: "app.oob:"},
			},
			compoCount: 2,
			nodeCount:  3,
		},
		{
			scenario: "replace nested compo first child",
			compo:    &NestedNested{},
			mutate: func(c Compo) {
				c.(*NestedNested).Mur = true
			},
			changes: []change{
				{Action: newNode, NodeID: "app.mur:", Type: "app.mur", IsCompo: true, CompoID: "app.nested:"},
				{Action: newNode, NodeID: "div:", Type: "div", CompoID: "app.mur:"},

				{Action: setAttr, NodeID: "div:", Key: "class", Value: "test"},
				{Action: appendChild, NodeID: "app.mur:", ChildID: "div:"},
				{Action: replaceChild, NodeID: "app.nested:", ChildID: "app.oob:", NewChildID: "app.mur:"},

				{Action: delNode, NodeID: "p:"},
				{Action: delNode, NodeID: "app.oob:"},
			},
			compoCount: 3,
			nodeCount:  4,
		},

		// Svg:
		{
			scenario: "create node with namespace",
			compo:    &Svg{},
			changes: []change{
				{Action: newNode, NodeID: "app.svg:", Type: "app.svg", IsCompo: true},
				{Action: newNode, NodeID: "svg:", Type: "svg", Namespace: svg, CompoID: "app.svg:"},
				{Action: newNode, NodeID: "path:", Type: "path", Namespace: svg, CompoID: "app.svg:"},
				{Action: newNode, NodeID: "path:", Type: "path", Namespace: svg, CompoID: "app.svg:"},

				{Action: setAttr, NodeID: "path:", Key: "data"},
				{Action: appendChild, NodeID: "svg:", ChildID: "path:"},
				{Action: setAttr, NodeID: "path:", Key: "data"},
				{Action: appendChild, NodeID: "svg:", ChildID: "path:"},
				{Action: appendChild, NodeID: "app.svg:", ChildID: "svg:"},

				{Action: setRoot, NodeID: "app.svg:"},
			},
			compoCount: 1,
			nodeCount:  4,
		},
		{
			scenario: "update node with namespace",
			compo:    &Svg{},
			mutate: func(c Compo) {
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
			mutate: func(c Compo) {
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
				{Action: newNode, NodeID: "app.selfclosing:", Type: "app.selfclosing", IsCompo: true},
				{Action: newNode, NodeID: "div:", Type: "div", CompoID: "app.selfclosing:"},
				{Action: newNode, NodeID: "div:", Type: "div", CompoID: "app.selfclosing:"},
				{Action: newNode, NodeID: "svg:", Type: "svg", Namespace: svg, CompoID: "app.selfclosing:"},

				{Action: appendChild, NodeID: "div:", ChildID: "div:"},
				{Action: appendChild, NodeID: "div:", ChildID: "svg:"},
				{Action: appendChild, NodeID: "app.selfclosing:", ChildID: "div:"},

				{Action: setRoot, NodeID: "app.selfclosing:"},
			},
			compoCount: 1,
			nodeCount:  4,
		},

		// Void elem:
		{
			scenario: "void elem node",
			compo:    &VoidElem{},
			changes: []change{
				{Action: newNode, NodeID: "app.voidelem:", Type: "app.voidelem", IsCompo: true},
				{Action: newNode, NodeID: "div:", Type: "div", CompoID: "app.voidelem:"},
				{Action: newNode, NodeID: "img:", Type: "img", CompoID: "app.voidelem:"},
				{Action: newNode, NodeID: "p:", Type: "p", CompoID: "app.voidelem:"},

				{Action: appendChild, NodeID: "div:", ChildID: "img:"},
				{Action: appendChild, NodeID: "div:", ChildID: "p:"},
				{Action: appendChild, NodeID: "app.voidelem:", ChildID: "div:"},

				{Action: setRoot, NodeID: "app.voidelem:"},
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
			mutate: func(c Compo) {
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
			mutate: func(c Compo) {
				c.(*CompoErr).ReplaceCompoErr = true
			},
			err: true,
		},
		{
			scenario: "fail add child",
			compo:    &CompoErr{Int: 42},
			mutate: func(c Compo) {
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

			e := domEngine{
				CompoBuilder: f,
				AllowedNodes: test.allowedNodes,
				AttrTransforms: []attrTransform{
					JsToGoHandler,
					HrefCompoFmt,
				},
				Sync: func(v []change) error {
					changes = make([]change, len(v))
					copy(changes, v)
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
	e := domEngine{
		Sync: func([]change) error {
			return errors.New("simulated err")
		},
	}

	err := e.Render(&Mur{})
	assert.Error(t, err)
}

func TestEngineSyncError(t *testing.T) {
	f := newCompoBuilder()
	f.register(&Mur{})

	e := domEngine{
		CompoBuilder: f,
		Sync: func([]change) error {
			return errors.New("simulated err")
		},
	}

	err := e.New(&Mur{})
	assert.Error(t, err)
}

func TestEngineEmptySync(t *testing.T) {
	f := newCompoBuilder()
	f.register(&Mur{})

	e := domEngine{CompoBuilder: f}
	err := e.New(&Mur{})
	assert.NoError(t, err)
}

func TestDOMCompoByID(t *testing.T) {
	f := newCompoBuilder()
	f.register(&Mur{})

	e := domEngine{CompoBuilder: f}
	foo := &Mur{}

	err := e.New(foo)
	require.NoError(t, err)

	c, ok := e.compos[foo]
	require.True(t, ok)
	require.Equal(t, foo, c.Compo)

	var foo2 Compo
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
