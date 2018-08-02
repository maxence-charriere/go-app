package html

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/tests"
)

func TestMarkup(t *testing.T) {
	tests.TestMarkup(t, func(factory *app.Factory) app.Markup {
		return NewMarkup(factory)
	})
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
	<div>compo String: {{compo "html.compo"}}</div>	
</div>
	`
}

func (c *CompoWithFields) Funcs() map[string]interface{} {
	return map[string]interface{}{
		"hello": func(string) string { return "hello" },
	}
}

func TestDecodeCompo(t *testing.T) {
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

	if err := decodeCompo(compo, &tag); err != nil {
		t.Fatal(err)
	}

	raw := tag.Children[1].Children[0]
	if raw.Text != "raw String: Hi" {
		t.Error(`raw is not "raw String: Hi":`, raw.Text)
	}

	component := tag.Children[7]
	if component.Attributes["obj"] != sjson {
		t.Errorf("component obj attribute is not %s: %s", sjson, component.Attributes["obj"])
	}

	year := strconv.Itoa(time.Now().Year())
	timet := tag.Children[8].Children[0]
	if timet.Text != "Time: "+year {
		t.Errorf(`time text is not "Time: %s": %s`, year, timet.Text)
	}

	hello := tag.Children[9].Children[0]
	if hello.Text != "hello" {
		t.Error("hello text is not hello:", hello.Text)
	}
}

func TestMapCompoFields(t *testing.T) {
	tests := []struct {
		scenario string
		function func(t *testing.T)
	}{
		{
			scenario: "skip mapping nil",
			function: testMapCompoFieldsNil,
		},
		{
			scenario: "skip mapping an anonymous field",
			function: testMapCompoFieldsAnonymous,
		},
		{
			scenario: "skip mapping an unexported field",
			function: testMapCompoFieldsUnexported,
		},
		{
			scenario: "map a string",
			function: testMapCompoFieldsString,
		},
		{
			scenario: "map a bool",
			function: testMapCompoFieldsBool,
		},
		{
			scenario: "map a naked bool",
			function: testMapCompoFieldsBoolNaked,
		},
		{
			scenario: "map a non boolean value to bool returns an error",
			function: testMapCompoFieldsBoolError,
		},
		{
			scenario: "map an int",
			function: testMapCompoFieldsInt,
		},
		{
			scenario: "map a non int value to int returns an error",
			function: testMapCompoFieldsIntError,
		},
		{
			scenario: "map an uint",
			function: testMapCompoFieldsUint,
		},
		{
			scenario: "map a non uint value to uint returns an error",
			function: testMapCompoFieldsUintError,
		},
		{
			scenario: "map a float",
			function: testMapCompoFieldsFloat,
		},
		{
			scenario: "map a non float value to float returns an error",
			function: testMapCompoFieldsFloatError,
		},
		{
			scenario: "map a struct",
			function: testMapCompoFieldsStruct,
		},
		{
			scenario: "map a struct with invalid fields returns an error",
			function: testMapCompoFieldsStructError,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, test.function)
	}
}

func testMapCompoFieldsNil(t *testing.T) {
	compo := &CompoWithFields{}
	if err := mapCompoFields(compo, nil); err != nil {
		t.Fatal(err)
	}
}

func testMapCompoFieldsAnonymous(t *testing.T) {
	compo := &CompoWithFields{}
	if err := mapCompoFields(compo, app.AttributeMap{"zerocompo": `{"placeholder": 42}`}); err != nil {
		t.Fatal(err)
	}
}

func testMapCompoFieldsUnexported(t *testing.T) {
	compo := &CompoWithFields{}
	if err := mapCompoFields(compo, app.AttributeMap{"secret": "pandore"}); err != nil {
		t.Fatal(err)
	}
	if len(compo.secret) != 0 {
		t.Error("secret is not empty:", compo.secret)
	}
}

func testMapCompoFieldsString(t *testing.T) {
	compo := &CompoWithFields{}
	s := "hello"
	if err := mapCompoFields(compo, app.AttributeMap{"string": s}); err != nil {
		t.Fatal(err)
	}
	if compo.String != s {
		t.Errorf("string is not %s: %s", s, compo.String)
	}
}

func testMapCompoFieldsBool(t *testing.T) {
	compo := &CompoWithFields{}
	if err := mapCompoFields(compo, app.AttributeMap{"bool": "true"}); err != nil {
		t.Fatal(err)
	}
	if !compo.Bool {
		t.Error("bool is not true")
	}
}

func testMapCompoFieldsBoolNaked(t *testing.T) {
	compo := &CompoWithFields{}
	if err := mapCompoFields(compo, app.AttributeMap{"bool": ""}); err != nil {
		t.Fatal(err)
	}
	if !compo.Bool {
		t.Error("bool is not true")
	}
}

func testMapCompoFieldsBoolError(t *testing.T) {
	compo := &CompoWithFields{}
	err := mapCompoFields(compo, app.AttributeMap{"bool": "lolilol"})
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testMapCompoFieldsInt(t *testing.T) {
	compo := &CompoWithFields{}
	if err := mapCompoFields(compo, app.AttributeMap{"int": "42"}); err != nil {
		t.Fatal(err)
	}
	if compo.Int != 42 {
		t.Error("int is not 42:", compo.Int)
	}
}

func testMapCompoFieldsIntError(t *testing.T) {
	compo := &CompoWithFields{}
	err := mapCompoFields(compo, app.AttributeMap{"int": "lolilol"})
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testMapCompoFieldsUint(t *testing.T) {
	compo := &CompoWithFields{}
	if err := mapCompoFields(compo, app.AttributeMap{"uint": "42"}); err != nil {
		t.Fatal(err)
	}
	if compo.Uint != 42 {
		t.Error("uint is not 42:", compo.Uint)
	}
}

func testMapCompoFieldsUintError(t *testing.T) {
	compo := &CompoWithFields{}
	err := mapCompoFields(compo, app.AttributeMap{"uint": "lolilol"})
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testMapCompoFieldsFloat(t *testing.T) {
	compo := &CompoWithFields{}
	if err := mapCompoFields(compo, app.AttributeMap{"float": "42.42"}); err != nil {
		t.Fatal(err)
	}
	if compo.Float != 42.42 {
		t.Error("float is not 42.42:", compo.Float)
	}
}

func testMapCompoFieldsFloatError(t *testing.T) {
	compo := &CompoWithFields{}
	err := mapCompoFields(compo, app.AttributeMap{"float": "42.world"})
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testMapCompoFieldsStruct(t *testing.T) {
	compo := &CompoWithFields{}
	if err := mapCompoFields(compo, app.AttributeMap{"struct": `{"A": 42, "B": "world"}`}); err != nil {
		t.Fatal(err)
	}
	if compo.Struct.A != 42 {
		t.Error("struct.A is not 42:", compo.Struct.A)
	}
	if compo.Struct.B != "world" {
		t.Error("struct.B is not world:", compo.Struct.B)
	}
}

func testMapCompoFieldsStructError(t *testing.T) {
	compo := &CompoWithFields{}
	err := mapCompoFields(compo, app.AttributeMap{"struct": `{"A": "world", "B": 42}`})
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func BenchmarkMarkupMount(b *testing.B) {
	factory := app.NewFactory()
	factory.RegisterCompo(&tests.Hello{})
	factory.RegisterCompo(&tests.World{})

	markup := NewMarkup(factory)

	for i := 0; i < b.N; i++ {
		hello := &tests.Hello{
			Name: "JonhyMaxoo",
		}
		markup.Mount(hello)
		markup.Dismount(hello)
	}
}

func BenchmarkMarkupUpdate(b *testing.B) {
	factory := app.NewFactory()
	factory.RegisterCompo(&tests.Hello{})
	factory.RegisterCompo(&tests.World{})

	markup := NewMarkup(factory)

	hello := &tests.Hello{
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
