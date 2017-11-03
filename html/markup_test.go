package html

import (
	"encoding/json"
	"html/template"
	"strconv"
	"testing"
	"time"

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
			function: testMarkupMountDismount,
		},
		{
			scenario: "mount a mounted component should fail",
			function: testMarkupMountMounted,
		},
		{
			scenario: "mount a component with a bad template should fail",
			function: testMarkupMountComponentWithBadTemplate,
		},
		{
			scenario: "mount a component with a bad tag",
			function: testMarkuptMountComponentWithBadTag,
		},
		{
			scenario: "mount a component with a not registered child",
			function: testMarkuptMountComponentWithNotRegistedChild,
		},
		{
			scenario: "mount a component with bad attributes",
			function: testMarkuptMountComponentWithBadAttrs,
		},
		{
			scenario: "dismount a dismounted should do nothing",
			function: testMarkupDismountDismounted,
		},
		{
			scenario: "dismount a component with dismounted child should do nothing",
			function: testMarkupDismountDismountedChild,
		},
		{
			scenario: "update should not trigger changes",
			function: testMarkupUpdateNoChanges,
		},
		{
			scenario: "should update text",
			function: testMarkupUpdateText,
		},
		{
			scenario: "should update simple tag to component",
			function: testMarkupUpdateSimpleToCompo,
		},
		{
			scenario: "should update simple tag to text",
			function: testMarkupUpdateSimpleToText,
		},
		{
			scenario: "should update text to simple tag",
			function: testMarkupUpdateTextToSimple,
		},
		{
			scenario: "should update component",
			function: testMarkupUpdateComponent,
		},
		{
			scenario: "update a unchanged component should do nothing",
			function: testMarkupUpdateComponentNoChange,
		},
		{
			scenario: "should update attributes",
			function: testMarkupUpdateUpdateAttributes,
		},
		{
			scenario: "update not mounted component should fail",
			function: testMarkupUpdateUpdateNotMountedComponent,
		},
		{
			scenario: "update component with bad template should fail",
			function: testMarkupUpdateComponentWithBadTemplate,
		},
		{
			scenario: "update component with bad child should fail",
			function: testMarkupUpdateComponentWithBadChild,
		},
		{
			scenario: "update component with an error should fail",
			function: testMarkupUpdateComponentWithError,
		},
		{
			scenario: "update tag with bad attribute shoulld fail",
			function: testMarkupUpdateBadAttribute,
		},
		{
			scenario: "update component with dismounted child should fail",
			function: testMarkupUpdateComponentWithDismountedChild,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			markup := NewMarkup(factory)
			test.function(t, markup)
		})
	}
}

func testMarkupMountDismount(t *testing.T, markup *Markup) {
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

func testMarkupMountMounted(t *testing.T, markup *Markup) {
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

func testMarkupMountComponentWithBadTemplate(t *testing.T, markup *Markup) {
	testMarkuptMountInvalidComponent(t, markup, &CompoWithBadTmpl{})

}

func testMarkuptMountInvalidComponent(t *testing.T, markup *Markup, compo app.Component) {
	_, err := markup.Mount(compo)
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testMarkuptMountComponentWithBadTag(t *testing.T, markup *Markup) {
	testMarkuptMountInvalidComponent(t, markup, &CompoWithBadTag{})

}

func testMarkuptMountComponentWithNotRegistedChild(t *testing.T, markup *Markup) {
	testMarkuptMountInvalidComponent(t, markup, &CompoWithNotRegisteredChild{})

}

func testMarkuptMountComponentWithBadAttrs(t *testing.T, markup *Markup) {
	testMarkuptMountInvalidComponent(t, markup, &CompoWithBadAttrs{})
}

func testMarkupDismountDismounted(t *testing.T, markup *Markup) {
	compo := &Foo{}
	if _, err := markup.Mount(compo); err != nil {
		t.Fatal(err)
	}
	markup.Dismount(compo)
	markup.Dismount(compo)
}

func testMarkupDismountDismountedChild(t *testing.T, markup *Markup) {
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

func testMarkupUpdateNoChanges(t *testing.T, markup *Markup) {
	compo := &Hello{}
	if _, err := markup.Mount(compo); err != nil {
		t.Fatal(err)
	}

	syncs, err := markup.Update(compo)
	if err != nil {
		t.Fatal(err)
	}
	if len(syncs) != 0 {
		t.Error("syncs is not empty:", len(syncs))
	}
}

func testMarkupUpdateText(t *testing.T, markup *Markup) {
	compo := &Hello{Greeting: "Hi"}
	if _, err := markup.Mount(compo); err != nil {
		t.Fatal(err)
	}

	compo.Greeting = "Hello"

	syncs, err := markup.Update(compo)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(syncs); l != 1 {
		t.Fatal("syncs doesn't have 1 element:", l)
	}

	sync := syncs[0]
	if !sync.Replace {
		t.Error("sync is not a replace")
	}

	h1 := sync.Tag
	if h1.Name != "h1" {
		t.Fatal("tag updated is not a h1:", h1.Name)
	}

	if text := h1.Children[0]; text.Text != compo.Greeting {
		t.Errorf(`text is not "%s": "%s"`, compo.Greeting, text.Text)
	}
}

func testMarkupUpdateSimpleToCompo(t *testing.T, markup *Markup) {
	compo := &Hello{}
	if _, err := markup.Mount(compo); err != nil {
		t.Fatal(err)
	}

	compo.Name = "Maxence"

	syncs, err := markup.Update(compo)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(syncs); l != 1 {
		t.Fatal("syncs doesn't have 1 element:", l)
	}

	sync := syncs[0]
	if !sync.Replace {
		t.Error("sync is not a replace")
	}

	world := sync.Tag
	if world.Name != "html.world" {
		t.Fatal("tag updated is not a component html.world:", world.Name)
	}
	if name := world.Attributes["name"]; name != compo.Name {
		t.Fatalf(`name is not "%s": "%s"`, compo.Name, name)
	}
	if l := len(world.Children); l != 0 {
		t.Fatal("world has children", l)
	}
}

func testMarkupUpdateSimpleToText(t *testing.T, markup *Markup) {
	compo := &Hello{}
	if _, err := markup.Mount(compo); err != nil {
		t.Fatal(err)
	}

	compo.TextBye = true

	syncs, err := markup.Update(compo)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(syncs); l != 1 {
		t.Fatal("syncs doesn't have 1 element:", l)
	}

	sync := syncs[0]
	if !sync.Replace {
		t.Error("sync is not a replace")
	}

	root := sync.Tag
	if root.Name != "div" {
		t.Fatal("root is not a div:", root.Name)
	}
	if l := len(root.Children); l != 4 {
		t.Fatal("root doesn't have 4 children:", l)
	}
	if text := root.Children[3]; text.Text != "Goodbye" {
		t.Fatalf(`text is not "Goodbye": "%s"`, text.Text)
	}
}

func testMarkupUpdateTextToSimple(t *testing.T, markup *Markup) {
	compo := &Hello{TextBye: true}
	if _, err := markup.Mount(compo); err != nil {
		t.Fatal(err)
	}

	compo.TextBye = false

	syncs, err := markup.Update(compo)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(syncs); l != 1 {
		t.Fatal("syncs doesn't have 1 element:", l)
	}

	sync := syncs[0]
	if !sync.Replace {
		t.Error("sync is not a replace")
	}

	root := sync.Tag
	if l := len(root.Children); l != 5 {
		t.Fatal("root doesn't have 5 children:", l)
	}
	if span := root.Children[3]; span.Name != "span" {
		t.Fatalf(`span is not a span tag: %s`, span.Name)
	}
	if p := root.Children[4]; p.Name != "p" {
		t.Fatalf(`p is not a p tag: %s`, p.Name)
	}
}

func testMarkupUpdateComponent(t *testing.T, markup *Markup) {
	compo := &Hello{Name: "Jonhy"}
	if _, err := markup.Mount(compo); err != nil {
		t.Fatal(err)
	}

	compo.Name = "Maxence"

	syncs, err := markup.Update(compo)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(syncs); l != 1 {
		t.Fatal("syncs doesn't have 1 element:", l)
	}

	sync := syncs[0]
	if !sync.Replace {
		t.Error("sync is not a replace")
	}

	worldRoot := sync.Tag
	if worldRoot.Name != "div" {
		t.Fatal("root of world is not a div:", worldRoot.Name)
	}
	if l := len(worldRoot.Children); l != 1 {
		t.Fatal("root of world doesn't have 1 child:", l)
	}
	if text := worldRoot.Children[0]; text.Text != compo.Name {
		t.Fatalf(`text is not "%s": "%s"`, compo.Name, text.Text)
	}
}

func testMarkupUpdateComponentNoChange(t *testing.T, markup *Markup) {
	compo := &Hello{Name: "JonhyMaxoo"}
	if _, err := markup.Mount(compo); err != nil {
		t.Fatal(err)
	}

	syncs, err := markup.Update(compo)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(syncs); l != 0 {
		t.Error("syncs is not empty:", l)
	}
}

func testMarkupUpdateUpdateAttributes(t *testing.T, markup *Markup) {
	compo := &Hello{}
	if _, err := markup.Mount(compo); err != nil {
		t.Fatal(err)
	}

	compo.Placeholder = "Enter your name"

	syncs, err := markup.Update(compo)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(syncs); l != 1 {
		t.Fatal("syncs doesn't have 1 element:", l)
	}

	sync := syncs[0]
	if sync.Replace {
		t.Error("sync is a replace")
	}

	input := sync.Tag
	if input.Name != "input" {
		t.Error("input is not an input tag:", input.Name)
	}
	if placeholder := input.Attributes["placeholder"]; placeholder != compo.Placeholder {
		t.Errorf("input placeholder is not %s: %s", compo.Placeholder, placeholder)
	}
	if l := len(input.Children); l != 0 {
		t.Error("input has child")
	}
}

func testMarkupUpdateUpdateNotMountedComponent(t *testing.T, markup *Markup) {
	_, err := markup.Update(&Hello{})
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testMarkupUpdateComponentWithBadTemplate(t *testing.T, markup *Markup) {
	compo := &Hello{}
	if _, err := markup.Mount(compo); err != nil {
		t.Fatal(err)
	}

	compo.TmplErr = true

	_, err := markup.Update(compo)
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testMarkupUpdateComponentWithBadChild(t *testing.T, markup *Markup) {
	compo := &Hello{Name: "Max"}
	if _, err := markup.Mount(compo); err != nil {
		t.Fatal(err)
	}

	compo.ChildErr = true

	_, err := markup.Update(compo)
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testMarkupUpdateComponentWithError(t *testing.T, markup *Markup) {
	compo := &Hello{}
	if _, err := markup.Mount(compo); err != nil {
		t.Fatal(err)
	}

	compo.Name = "Jonhy"
	compo.ChildErr = true

	_, err := markup.Update(compo)
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testMarkupUpdateBadAttribute(t *testing.T, markup *Markup) {
	compo := &Hello{Name: "Maxoo"}
	if _, err := markup.Mount(compo); err != nil {
		t.Fatal(err)
	}

	compo.CompoFieldErr = true

	_, err := markup.Update(compo)
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testMarkupUpdateComponentWithDismountedChild(t *testing.T, markup *Markup) {
	compo := &Hello{Name: "Maxoo"}
	root, err := markup.Mount(compo)
	if err != nil {
		t.Fatal(err)
	}

	world := root.Children[2].Children[0]
	markup.dismountTag(world)

	compo.Name = "Jonhy"

	if _, err = markup.Update(compo); err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
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

type CompoWithFields struct {
	app.ZeroCompo
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
</div>
	`
}

func (c *CompoWithFields) Funcs() map[string]interface{} {
	return map[string]interface{}{
		"hello": func(string) string { return "hello" },
	}
}

func TestDecodeComponent(t *testing.T) {
	var tag app.Tag

	s := struct {
		A int
		B string
	}{
		A: 42,
		B: "foobar",
	}

	data, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	sjson := string(data)

	compo := &CompoWithFields{
		String: "Hi",
		Time:   time.Now(),
		Struct: s,
	}

	if err := decodeComponent(compo, &tag); err != nil {
		t.Fatal(err)
	}

	raw := tag.Children[1].Children[0]
	if raw.Text != "raw String: Hi" {
		t.Error(`raw is not "raw String: Hi":`, raw.Text)
	}

	component := tag.Children[7]
	if component.Attributes["obj"] != sjson {
		t.Errorf("component obj attribute should be %s: %s", sjson, component.Attributes["obj"])
	}

	year := strconv.Itoa(time.Now().Year())
	timet := tag.Children[8].Children[0]
	if timet.Text != "Time: "+year {
		t.Errorf(`time text should be "Time: %s": %s`, year, timet.Text)
	}

	hello := tag.Children[9].Children[0]
	if hello.Text != "hello" {
		t.Error("hello text should be hello:", hello.Text)
	}
}

func TestMapComponentFields(t *testing.T) {
	tests := []struct {
		scenario string
		function func(t *testing.T)
	}{
		{
			scenario: "map nil should do nothing",
			function: testMapComponentFieldsNil,
		},
		{
			scenario: "map anonymous field should do nothing",
			function: testMapComponentFieldsAnonymous,
		},
		{
			scenario: "map unexported field should do nothing",
			function: testMapComponentFieldsUnexported,
		},
		{
			scenario: "should map string",
			function: testMapComponentFieldsString,
		},
		{
			scenario: "should map bool",
			function: testMapComponentFieldsBool,
		},
		{
			scenario: "should map naked bool",
			function: testMapComponentFieldsBoolNaked,
		},
		{
			scenario: "map a non boolean value to bool should fail",
			function: testMapComponentFieldsBoolError,
		},
		{
			scenario: "should map int",
			function: testMapComponentFieldsInt,
		},
		{
			scenario: "map a non int value to int should fail",
			function: testMapComponentFieldsIntError,
		},
		{
			scenario: "should map uint",
			function: testMapComponentFieldsUint,
		},
		{
			scenario: "map a non uint value to uint should fail",
			function: testMapComponentFieldsUintError,
		},
		{
			scenario: "should map float",
			function: testMapComponentFieldsFloat,
		},
		{
			scenario: "map a non float value to float should fail",
			function: testMapComponentFieldsFloatError,
		},
		{
			scenario: "should map struct",
			function: testMapComponentFieldsStruct,
		},
		{
			scenario: "map a struct with invalid fields should fail",
			function: testMapComponentFieldsStructError,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, test.function)
	}
}

func testMapComponentFieldsNil(t *testing.T) {
	compo := &CompoWithFields{}
	if err := mapComponentFields(compo, nil); err != nil {
		t.Fatal(err)
	}
}

func testMapComponentFieldsAnonymous(t *testing.T) {
	compo := &CompoWithFields{}
	if err := mapComponentFields(compo, app.AttributeMap{"zerocompo": `{"placeholder": 42}`}); err != nil {
		t.Fatal(err)
	}
}

func testMapComponentFieldsUnexported(t *testing.T) {
	compo := &CompoWithFields{}
	if err := mapComponentFields(compo, app.AttributeMap{"secret": "pandore"}); err != nil {
		t.Fatal(err)
	}
	if len(compo.secret) != 0 {
		t.Error("secret is not empty:", compo.secret)
	}
}

func testMapComponentFieldsString(t *testing.T) {
	compo := &CompoWithFields{}
	s := "hello"
	if err := mapComponentFields(compo, app.AttributeMap{"string": s}); err != nil {
		t.Fatal(err)
	}
	if compo.String != s {
		t.Errorf("string is not %s: %s", s, compo.String)
	}
}

func testMapComponentFieldsBool(t *testing.T) {
	compo := &CompoWithFields{}
	if err := mapComponentFields(compo, app.AttributeMap{"bool": "true"}); err != nil {
		t.Fatal(err)
	}
	if !compo.Bool {
		t.Error("bool is not true")
	}
}

func testMapComponentFieldsBoolNaked(t *testing.T) {
	compo := &CompoWithFields{}
	if err := mapComponentFields(compo, app.AttributeMap{"bool": ""}); err != nil {
		t.Fatal(err)
	}
	if !compo.Bool {
		t.Error("bool is not true")
	}
}

func testMapComponentFieldsBoolError(t *testing.T) {
	compo := &CompoWithFields{}
	err := mapComponentFields(compo, app.AttributeMap{"bool": "lolilol"})
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testMapComponentFieldsInt(t *testing.T) {
	compo := &CompoWithFields{}
	if err := mapComponentFields(compo, app.AttributeMap{"int": "42"}); err != nil {
		t.Fatal(err)
	}
	if compo.Int != 42 {
		t.Error("int is not 42:", compo.Int)
	}
}

func testMapComponentFieldsIntError(t *testing.T) {
	compo := &CompoWithFields{}
	err := mapComponentFields(compo, app.AttributeMap{"int": "lolilol"})
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testMapComponentFieldsUint(t *testing.T) {
	compo := &CompoWithFields{}
	if err := mapComponentFields(compo, app.AttributeMap{"uint": "42"}); err != nil {
		t.Fatal(err)
	}
	if compo.Uint != 42 {
		t.Error("uint is not 42:", compo.Uint)
	}
}

func testMapComponentFieldsUintError(t *testing.T) {
	compo := &CompoWithFields{}
	err := mapComponentFields(compo, app.AttributeMap{"uint": "lolilol"})
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testMapComponentFieldsFloat(t *testing.T) {
	compo := &CompoWithFields{}
	if err := mapComponentFields(compo, app.AttributeMap{"float": "42.42"}); err != nil {
		t.Fatal(err)
	}
	if compo.Float != 42.42 {
		t.Error("float is not 42.42:", compo.Float)
	}
}

func testMapComponentFieldsFloatError(t *testing.T) {
	compo := &CompoWithFields{}
	err := mapComponentFields(compo, app.AttributeMap{"float": "42.world"})
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testMapComponentFieldsStruct(t *testing.T) {
	compo := &CompoWithFields{}
	if err := mapComponentFields(compo, app.AttributeMap{"struct": `{"A": 42, "B": "world"}`}); err != nil {
		t.Fatal(err)
	}
	if compo.Struct.A != 42 {
		t.Error("struct.A is not 42:", compo.Struct.A)
	}
	if compo.Struct.B != "world" {
		t.Error("struct.B is not world:", compo.Struct.B)
	}
}

func testMapComponentFieldsStructError(t *testing.T) {
	compo := &CompoWithFields{}
	err := mapComponentFields(compo, app.AttributeMap{"struct": `{"A": "world", "B": 42}`})
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func BenchmarkMarkupMount(b *testing.B) {
	factory := app.NewFactory()
	factory.RegisterComponent(&Hello{})
	factory.RegisterComponent(&World{})

	markup := NewMarkup(factory)

	for i := 0; i < b.N; i++ {
		hello := &Hello{
			Name: "JonhyMaxoo",
		}
		markup.Mount(hello)
		markup.Dismount(hello)
	}
}

func BenchmarkMarkupUpdate(b *testing.B) {
	factory := app.NewFactory()
	factory.RegisterComponent(&Hello{})
	factory.RegisterComponent(&World{})

	markup := NewMarkup(factory)

	hello := &Hello{
		Name: "JonhyMaxoo",
	}
	markup.Mount(hello)

	alt := false

	for i := 0; i < b.N; i++ {
		if alt {
			hello.Greeting = "Jon"
		} else {
			hello.Greeting = ""
		}
		hello.TextBye = alt
		hello.Placeholder = strconv.Itoa(i)
		hello.Greeting = strconv.Itoa(i)

		markup.Update(hello)

		alt = !alt
	}
}
