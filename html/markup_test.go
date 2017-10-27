package html

import (
	"html/template"
	"testing"

	"github.com/murlokswarm/app"
)

type Foo struct {
	Boo bool
}

func (c *Foo) OnMount() {}

func (c *Foo) OnDismount() {}

func (c *Foo) Funcs() template.FuncMap {
	return nil
}

func (c *Foo) Render() string {
	return `
<div>
	<h1>Foo</h1>
	<html.bar>
</div>
	`
}

type Bar app.ZeroCompo

func (c *Bar) Render() string {
	return `<h2>Bar</h2>`
}

type CompoWithBadTmpl app.ZeroCompo

func (c *CompoWithBadTmpl) Render() string {
	return `<h2>{{.Hello}}</h2>`
}

type CompoWithBadTag app.ZeroCompo

func (c *CompoWithBadTag) Render() string {
	return `<h1><div/></h1>`
}

type CompoWithNotRegisteredChild app.ZeroCompo

func (c *CompoWithNotRegisteredChild) Render() string {
	return `
<div>
	<html.unknown>
</div>
	`
}

type CompoWithBadChild app.ZeroCompo

func (c *CompoWithBadChild) Render() string {
	return `
<div>
	<html.compobadtmpl>
</div>
	`
}

type CompoWithBadAttrs app.ZeroCompo

func (c *CompoWithBadAttrs) Render() string {
	return `
<div>
	<html.foo boo="Holy Shit">
</div>
	`
}

type Hello struct {
	Greeting      string
	Name          string
	Placeholder   string
	TextBye       bool
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
			<html.world name="{{html .Name}}" err="{{.ChildErr}}" {{if .CompoFieldErr}}fielderr="-42"{{end}}>
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
		<html.componotregistered>
	{{end}}
</div>
	`
}

func TestMarkup(t *testing.T) {
	factory := app.NewFactory()
	factory.RegisterComponent(&Foo{})
	factory.RegisterComponent(&Bar{})
	factory.RegisterComponent(&CompoWithBadTmpl{})
	factory.RegisterComponent(&CompoWithBadTag{})
	factory.RegisterComponent(&CompoWithNotRegisteredChild{})
	factory.RegisterComponent(&CompoWithBadChild{})
	factory.RegisterComponent(&Hello{})
	factory.RegisterComponent(&World{})

	tests := []struct {
		scenario string
		function func(t *testing.T, markup *Markup)
	}{
		{
			scenario: "should mount and dismount a component",
			function: testMountDismount,
		},
		{
			scenario: "mount a mounted component should fail",
			function: testMountMounted,
		},
		{
			scenario: "mount a component with a bad template should fail",
			function: testMountComponentWithBadTemplate,
		},
		{
			scenario: "mount a component with a bad tag",
			function: testMountComponentWithBadTag,
		},
		{
			scenario: "mount a component with a not registered child",
			function: testMountComponentWithNotRegistedChild,
		},
		{
			scenario: "mount a component with bad attributes",
			function: testMountComponentWithBadAttrs,
		},
		{
			scenario: "dismount a dismounted should do nothing",
			function: testDismountDismounted,
		},
		{
			scenario: "dismount a component with dismounted child should do nothing",
			function: testDismountDismountedChild,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			markup := NewMarkup(factory)
			test.function(t, markup)
		})
	}
}

func testMountDismount(t *testing.T, markup *Markup) {
	compo := &Foo{}

	root, err := markup.Mount(compo)
	if err != nil {
		t.Fatal(err)
	}
	if count := len(markup.components); count != 2 {
		t.Fatal("markup doesn't have 2 components:", count)
	}
	if count := len(markup.roots); count != 2 {
		t.Fatal("markup doesn't have 2 roots:", count)
	}

	barTag := root.Children[1]
	if name := barTag.Name; name != "html.bar" {
		t.Fatalf("bar tag is not a html.bar: %s", name)
	}
	if _, err = markup.Component(barTag.ID); err != nil {
		t.Fatal(err)
	}

	markup.Dismount(compo)
	if count := len(markup.components); count != 0 {
		t.Fatal("markup should not have components")
	}
	if count := len(markup.roots); count != 0 {
		t.Fatal("markup should not have roots")
	}
}

func testMountMounted(t *testing.T, markup *Markup) {
	compo := &Foo{}

	if _, err := markup.Mount(compo); err != nil {
		t.Fatal(err)
	}

	_, err := markup.Mount(compo)
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testMountComponentWithBadTemplate(t *testing.T, markup *Markup) {
	testMountInvalidComponent(t, markup, &CompoWithBadTmpl{})

}

func testMountInvalidComponent(t *testing.T, markup *Markup, compo app.Component) {
	_, err := markup.Mount(compo)
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testMountComponentWithBadTag(t *testing.T, markup *Markup) {
	testMountInvalidComponent(t, markup, &CompoWithBadTag{})

}

func testMountComponentWithNotRegistedChild(t *testing.T, markup *Markup) {
	testMountInvalidComponent(t, markup, &CompoWithNotRegisteredChild{})

}

func testMountComponentWithBadAttrs(t *testing.T, markup *Markup) {
	testMountInvalidComponent(t, markup, &CompoWithBadAttrs{})
}

func testDismountDismounted(t *testing.T, markup *Markup) {
	compo := &Foo{}
	if _, err := markup.Mount(compo); err != nil {
		t.Fatal(err)
	}
	markup.Dismount(compo)
	markup.Dismount(compo)
}

func testDismountDismountedChild(t *testing.T, markup *Markup) {
	compo := &Foo{}
	root, err := markup.Mount(compo)
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range markup.components {
		if k != root.CompoID {
			markup.Dismount(v)
		}
	}
	markup.Dismount(compo)
}

func TestAttributesEquals(t *testing.T) {
	attrs := app.AttributeMap{
		"hello": "world",
		"foo":   "bar",
		"value": "",
	}

	attrs2 := app.AttributeMap{
		"foo":   "bar",
		"hello": "world",
		"value": "",
	}

	if !attributesEquals("div", attrs, attrs2) {
		t.Error("attrs and attrs2 are not equals")
	}

	if attributesEquals("div", attrs, nil) {
		t.Error("attrs and nil are equals")
	}

	attrs3 := app.AttributeMap{
		"foo":   "bar",
		"hello": "maxoo",
		"value": "",
	}

	if attributesEquals("div", attrs, attrs3) {
		t.Error("attrs and attrs3 are equals")
	}

	attrs4 := app.AttributeMap{
		"foo":   "bar",
		"bye":   "world",
		"value": "",
	}

	if attributesEquals("div", attrs, attrs4) {
		t.Error("attrs and attrs4 are equals")
	}

	attrs5 := app.AttributeMap{
		"hello": "world",
		"foo":   "bar",
		"value": "",
	}

	if attributesEquals("input", attrs, attrs5) {
		t.Error("attrs and attrs5 are equals")
	}
}
