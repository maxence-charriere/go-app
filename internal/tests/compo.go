package tests

import (
	"html/template"

	"github.com/murlokswarm/app"
)

// Foo is a test component.
type Foo struct {
	Boo bool
}

// OnMount satisfies the app.Mounter interface.
func (c *Foo) OnMount() {}

// OnDismount satisfies the app.Dismounter interface.
func (c *Foo) OnDismount() {}

// Subscribe satisfies the app.Subscriber interface.
func (c *Foo) Subscribe() *app.EventSubscriber {
	return app.NewEventSubscriber()
}

// Funcs satisfies the app.CompoWithExtendedRender interface.
func (c *Foo) Funcs() template.FuncMap {
	return nil
}

// Render satisfies the app.Compo interface.
func (c *Foo) Render() string {
	return `
<div>
	<h1>Foo</h1>
	<tests.bar>
</div>
	`
}

// Bar is a test component.
type Bar app.ZeroCompo

// Render satisfies the app.Compo interface.
func (c *Bar) Render() string {
	return `<h2>Bar</h2>`
}

// CompoWithBadTmpl is a test component that have a bad template.
type CompoWithBadTmpl app.ZeroCompo

// Render satisfies the app.Compo interface.
func (c *CompoWithBadTmpl) Render() string {
	return `<h2>{{.Hello}}</h2>`
}

// CompoWithBadTag is a test component that contains a bad tag.
type CompoWithBadTag app.ZeroCompo

// Render satisfies the app.Compo interface.
func (c *CompoWithBadTag) Render() string {
	return `<h1><div/></h1>`
}

// CompoWithNotRegisteredChild is a test component that constains a not
// registered child component.
type CompoWithNotRegisteredChild app.ZeroCompo

// Render satisfies the app.Compo interface.
func (c *CompoWithNotRegisteredChild) Render() string {
	return `
<div>
	<tests.unknown>
</div>
	`
}

// CompoWithBadChild is a test component that contains a bad child.
type CompoWithBadChild app.ZeroCompo

// Render satisfies the app.Compo interface.
func (c *CompoWithBadChild) Render() string {
	return `
<div>
	<tests.compobadtmpl>
</div>
	`
}

// CompoWithBadAttrs is a test component that contains a child set with bad
// attributes.
type CompoWithBadAttrs app.ZeroCompo

// Render satisfies the app.Compo interface.
func (c *CompoWithBadAttrs) Render() string {
	return `
<div>
	<tests.foo boo="Holy Shit">
</div>
	`
}

// NoPointerCompo is a test component that is not implemented as a pointer.
type NoPointerCompo app.ZeroCompo

// Render satisfies the app.Compo interface.
func (c NoPointerCompo) Render() string {
	return `<div>goodbye</div>`
}

// IntCompo is a test component that is not implemented as a pointer to struct.
type IntCompo int

// Render satisfies the app.Compo interface.
func (i *IntCompo) Render() string {
	return `<p>Aurevoir World</p>`
}

// EmptyCompo is a test component that has no fields.
type EmptyCompo struct{}

// Render satisfies the app.Compo interface.
func (c *EmptyCompo) Render() string {
	return `<p>Goodbye World</p>`
}

// Hello is a test component that allows to test multiple behaviors.
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

// Render satisfies the app.Compo interface.
func (h *Hello) Render() string {
	return `
<div>
	<h1>{{html .Greeting}}</h1>
	<input type="text" placeholder="{{.Placeholder}}" onchange="Name">
	<p>
		{{if .Name}}
			<tests.world name="{{html .Name}}" err="{{.ChildErr}}" {{if .CompoFieldErr}}fielderr="-42"{{end}}>
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

	<a href="tests.hello">hyperlink to a component</a>
	<a href="http://github.com">common hyperlink</a>

	<button onclick="js:console.console.log('hello')">js call</button>
</div>
	`
}

// World is a test component.
type World struct {
	Name     string
	Err      bool
	FieldErr uint
}

// Render satisfies the app.Compo interface.
func (w *World) Render() string {
	return `
<div>
	{{html .Name}}

	{{if .Err}}
		<tests.componotregistered>
	{{end}}
</div>
	`
}

// Mapping is a test component to verify mapping behaviors.
type Mapping struct {
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

// Method is a method to be mapped.
func (m *Mapping) Method() {
	m.method()
}

// Render satisfies the app.Compo interface.
func (m *Mapping) Render() string {
	return `<div>Some mappings</div>`
}

// MappingStruct is struct to test mapping struct behavior.
type MappingStruct struct {
	Exported   int
	unexported int
	method     func()
}

// Method is a method to be mapped.
func (s MappingStruct) Method() {
	s.method()
}

// MappingMap is map to test mapping map behavior.
type MappingMap map[string]func()

// Method is a method to be mapped.
func (m MappingMap) Method() {
	m["method"]()
}

// MappingSlice is slice to test mapping slice behavior.
type MappingSlice []func()

// Method is a method to be mapped.
func (s MappingSlice) Method() {
	s[0]()
}

// MappingInt is an int to test mapping value behavior.
type MappingInt int

// Method is a method to be mapped.
func (i MappingInt) Method(nb int) {
	mappedInt = nb
}

var mappedInt int

// RussianDoll is a component to test nested components behaviors.
type RussianDoll struct {
	Remaining int
}

// Render satisfies the app.Compo interface.
func (r *RussianDoll) Render() string {
	return `
<div>
	{{if gt .Remaining 0}}
		<tests.russiandoll remaining="{{sub .Remaining 1}}">
	{{end}}
</div>
	`
}

// Funcs satisfies the app.CompoWithExtendedRender interface.
func (r *RussianDoll) Funcs() map[string]interface{} {
	return map[string]interface{}{
		"sub": func(a, b int) int {
			return a - b
		},
	}
}

// Menu is a component to test a menu element.
type Menu app.ZeroCompo

// Render satisfies the app.Compo interface.
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
