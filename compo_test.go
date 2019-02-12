package app

import (
	"html/template"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type CompoWithFields struct {
	ZeroCompo
	secret             string
	funcHandler        func()
	funcWithArgHandler func(int)

	String     string
	Bool       bool
	NotSetBool bool
	Int        int
	Uint       uint
	Float      float64
	Struct     struct {
		A int
		B string
	}
	Time time.Time
}

func (c *CompoWithFields) Render() string {
	return `
<div>
	<div>String: {{.String}}</div>
	<div>raw String: {{raw .String}}</div>
	<div>Bool: {{.Bool}}</div>
	<div>Int: {{.Int}}</div>
	<div>Uint: {{.Uint}}</div>
	<div>Float: {{.Float}}</div>
	<div>Struct: {{.Struct}}</div>
	<html.compo obj="{{json .Struct}}">	
	<div>Time: {{time .Time "2006"}}</div>
	<div>{{hello .String}}</div>
	<div>compo String: {{compo "html.compo"}}</div>	
</div>
	`
}

func TestMapComponentFields(t *testing.T) {
	tests := []struct {
		scenario string
		attrs    map[string]string
		expected CompoWithFields
		err      bool
	}{
		{
			scenario: "skip mapping nil",
			attrs:    nil,
		},
		{
			scenario: "skip mapping an anonymous field",
			attrs:    map[string]string{"zerocompo": `{"placeholder": 42}`},
		},
		{
			scenario: "skip mapping an unexported field",
			attrs:    map[string]string{"secret": "pandore"},
		},
		{
			scenario: "map a string",
			attrs:    map[string]string{"string": "hello"},
			expected: CompoWithFields{
				String: "hello",
			},
		},
		{
			scenario: "map a bool",
			attrs:    map[string]string{"bool": "true"},
			expected: CompoWithFields{
				Bool: true,
			},
		},
		{
			scenario: "map a naked bool",
			attrs:    map[string]string{"bool": ""},
			expected: CompoWithFields{
				Bool: true,
			},
		},
		{
			scenario: "map a non boolean value to bool returns an error",
			attrs:    map[string]string{"bool": "lolilol"},
			err:      true,
		},
		{
			scenario: "map an int",
			attrs:    map[string]string{"int": "-42"},
			expected: CompoWithFields{
				Int: -42,
			},
		},
		{
			scenario: "map a non int value to int returns an error",
			attrs:    map[string]string{"int": "lolilol"},
			err:      true,
		},
		{
			scenario: "map an uint",
			attrs:    map[string]string{"uint": "21"},
			expected: CompoWithFields{
				Uint: 21,
			},
		},
		{
			scenario: "map a non uint value to uint returns an error",
			attrs:    map[string]string{"uint": "lolilol"},
			err:      true,
		},
		{
			scenario: "map a float",
			attrs:    map[string]string{"float": "42.42"},
			expected: CompoWithFields{
				Float: 42.42,
			},
		},
		{
			scenario: "map a non float value to float returns an error",
			attrs:    map[string]string{"float": "42.world"},
			err:      true,
		},
		{
			scenario: "map a struct",
			attrs:    map[string]string{"struct": `{"A": 42, "B": "world"}`},
			expected: CompoWithFields{
				Struct: struct {
					A int
					B string
				}{
					A: 42,
					B: "world",
				},
			},
		},
		{
			scenario: "map a struct with invalid fields returns an error",
			attrs:    map[string]string{"struct": `{"A": "world", "B": 42}`},
			err:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			var c CompoWithFields

			err := mapCompoFields(&c, test.attrs)
			if test.err {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, test.expected, c)
		})
	}
}

type Foo struct {
	Boo bool
}

func (c *Foo) OnMount() {}

func (c *Foo) OnDismount() {}

func (c *Foo) Subscribe() *Subscriber {
	return NewSubscriber()
}

func (c *Foo) Funcs() template.FuncMap {
	return nil
}

func (c *Foo) Render() string {
	return `
<div>
	<h1>Foo</h1>
	<app.bar>
</div>
	`
}

type Bar ZeroCompo

func (c *Bar) Render() string {
	return `<h2>Bar</h2>`
}

type CompoWithBadTmpl ZeroCompo

func (c *CompoWithBadTmpl) Render() string {
	return `<h2>{{.Hello}}</h2>`
}

type CompoWithBadTag ZeroCompo

func (c *CompoWithBadTag) Render() string {
	return `<h1><div/></h1>`
}

type CompoWithNotRegisteredChild ZeroCompo

func (c *CompoWithNotRegisteredChild) Render() string {
	return `
<div>
	<app.unknown>
</div>
	`
}

type CompoWithBadChild ZeroCompo

func (c *CompoWithBadChild) Render() string {
	return `
<div>
	<app.compobadtmpl>
</div>
	`
}

type CompoWithBadAttrs ZeroCompo

func (c *CompoWithBadAttrs) Render() string {
	return `
<div>
	<app.foo boo="Holy Shit">
</div>
	`
}

type NoPointerCompo ZeroCompo

func (c NoPointerCompo) Render() string {
	return `<div>goodbye</div>`
}

type IntCompo int

// Render satisfies the app.Compo interface.
func (i *IntCompo) Render() string {
	return `<p>Aurevoir World</p>`
}

type EmptyCompo struct{}

// Render satisfies the app.Compo interface.
func (c *EmptyCompo) Render() string {
	return `<p>Goodbye World</p>`
}

type Hello struct {
	Greeting      string
	Name          string
	Placeholder   string
	TextBye       bool
	SizeDiff      bool
	TmplErr       bool
	ChildErr      bool
	CompoFieldErr bool
}

func (h *Hello) Render() string {
	return `
<div>
	<h1>{{html .Greeting}}</h1>
	<input type="text" placeholder="{{.Placeholder}}" onchange="Name">
	<p>
		{{if .Name}}
			<app.world name="{{html .Name}}" err="{{.ChildErr}}" {{if .CompoFieldErr}}fielderr="-42"{{end}}>
		{{else}}
			<span>World</span>
		{{end}}
	</p>

	{{if .TmplErr}}
		<div>{{.UnknownField}}</div>
	{{end}}

	{{if .TextBye}}
		Goodbye
	{{else}}
		<span>Goodbye</span>
		<p>world</p>
	{{end}}

	{{if .SizeDiff}}
		<a>another tag</a>
	{{end}}

	<a href="app.hello">hyperlink to a component</a>
	<a href="http://github.com">common hyperlink</a>

	<button onclick="js:console.console.log('hello')">js call</button>
</div>
	`
}

type World struct {
	Name     string
	Err      bool
	FieldErr uint
}

func (w *World) Render() string {
	return `
<div>
	{{html .Name}}

	{{if .Err}}
		<app.componotregistered>
	{{end}}
</div>
	`
}

type MapperComp struct {
	String              string
	Int                 int
	IntWithMethod       MappingInt
	IntPtr              *int
	Struct              MappingStruct
	Map                 map[string]string
	MapWithMethod       MappingMap
	Slice               []int
	SliceWithMethod     MappingSlice
	Array               [5]int
	Func                func()
	FuncWithArg         func(i int)
	FuncWithMultipleArg func(x, y int)
	method              func()
}

func (m *MapperComp) Method() {
	m.method()
}

func (m *MapperComp) Render() string {
	return `<div>Some mappings</div>`
}

type MappingStruct struct {
	Exported   int
	unexported int
	method     func()
}

func (s MappingStruct) Method() {
	s.method()
}

type MappingMap map[string]func()

func (m MappingMap) Method() {
	m["method"]()
}

type MappingSlice []func()

func (s MappingSlice) Method() {
	s[0]()
}

type MappingInt int

func (i MappingInt) Method(nb int) {
	mappedInt = nb
}

var mappedInt int

type RussianDoll struct {
	Remaining int
}

func (r *RussianDoll) Render() string {
	return `
<div>
	{{if gt .Remaining 0}}
		<app.russiandoll remaining="{{sub .Remaining 1}}">
	{{end}}
</div>
	`
}

func (r *RussianDoll) Funcs() map[string]interface{} {
	return map[string]interface{}{
		"sub": func(a, b int) int {
			return a - b
		},
	}
}

type Menu ZeroCompo

func (m *Menu) Render() string {
	return `
<menu>
	<menu label="app">
		<menuitem label="a menu"></menuitem>
	</menu>
	<menu label="edit">
		<menuitem label="a menu"></menuitem>
	</menu>
</menu>
	`
}

func TestCompoNameFromURL(t *testing.T) {
	tests := []struct {
		rawurl       string
		expectedName string
	}{
		{
			rawurl:       "/hello",
			expectedName: "hello",
		},
		{
			rawurl:       "/Hello",
			expectedName: "hello",
		},
		{
			rawurl:       "/hello?int=42",
			expectedName: "hello",
		},
		{
			rawurl:       "/hello/world",
			expectedName: "hello",
		},
		{
			rawurl:       "hello",
			expectedName: "hello",
		},
		{
			rawurl:       "main.hello",
			expectedName: "hello",
		},
		{
			rawurl:       "main.hello?foo=bar",
			expectedName: "hello",
		},
		{
			rawurl:       "hello?foo=bar",
			expectedName: "hello",
		},
		{
			rawurl: "test://hello",
		},
		{
			rawurl: "compo://",
		},
		{
			rawurl: "http://www.github.com",
		},
	}

	for _, test := range tests {
		name := compoNameFromURLString(test.rawurl)
		assert.Equal(t, test.expectedName, name)
	}
}
