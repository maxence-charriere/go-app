package html

import (
	"bytes"
	"testing"

	"github.com/murlokswarm/app"
)

func TestDecoder(t *testing.T) {
	buff := bytes.NewBufferString(`
<div>
	<!-- Comment -->	
	<h1>hello</h1>
	<br>
	<input type="text" required>
	<lib.FooComponent Bar="42">
	<svg>
		<path d="M 42.42 Z "></path>
		<path d="M 21.21 Z " />
	</svg>
</div>
	`)
	dec := NewDecoder(buff)

	var root app.Tag
	if err := dec.Decode(&root); err != nil {
		t.Fatal(err)
	}

	testDecodeCheckRoot(t, root)
	testDecodeCheckH1(t, root.Children[0])
	testDecodeCheckBr(t, root.Children[1])
	testDecodeCheckInput(t, root.Children[2])
	testDecodeCheckFooComponent(t, root.Children[3])
	testDecodeCheckSvg(t, root.Children[4])
}

func testDecodeCheckRoot(t *testing.T, tag app.Tag) {
	if name := tag.Name; name != "div" {
		t.Fatalf(`tag name should be "div": "%s"`, name)
	}
	if typ := tag.Type; typ != app.SimpleTag {
		t.Fatal("tag should be a simple tag")
	}
	if count := len(tag.Children); count != 5 {
		t.Fatal("tag should have 5 children:", count)
	}
}

func testDecodeCheckH1(t *testing.T, tag app.Tag) {
	if name := tag.Name; name != "h1" {
		t.Fatalf(`tag name should be "h1": "%s"`, name)
	}
	if typ := tag.Type; typ != app.SimpleTag {
		t.Fatal("tag should be a simple tag")
	}
	if count := len(tag.Children); count != 1 {
		t.Fatal("tag should have 1 children:", count)
	}
	if text := tag.Children[0]; text.Text != "hello" {
		t.Fatalf(`text.Text should be "hello": "%s"`, text.Text)
	}
}

func testDecodeCheckBr(t *testing.T, tag app.Tag) {
	if name := tag.Name; name != "br" {
		t.Fatalf(`tag name should be "br": "%s"`, name)
	}
	if typ := tag.Type; typ != app.SimpleTag {
		t.Fatal("tag should be a simple tag")
	}
	if count := len(tag.Children); count != 0 {
		t.Fatal("root should not have children:", count)
	}
}

func testDecodeCheckInput(t *testing.T, tag app.Tag) {
	if name := tag.Name; name != "input" {
		t.Fatalf(`tag name should be "input": "%s"`, name)
	}
	if typ := tag.Type; typ != app.SimpleTag {
		t.Fatal("tag should be a simple tag")
	}
	if count := len(tag.Children); count != 0 {
		t.Fatal("tag should not have children:", count)
	}
	if count := len(tag.Attributes); count != 2 {
		t.Fatal("tag should have 2 attributes:", count)
	}
	if val, _ := tag.Attributes["type"]; val != "text" {
		t.Fatalf(`tag should have an attr with value = "text": %s`, val)
	}
	if _, ok := tag.Attributes["required"]; !ok {
		t.Fatal(`tag should have an attr with key = "required"`)
	}
}

func testDecodeCheckFooComponent(t *testing.T, tag app.Tag) {
	if name := tag.Name; name != "lib.foocomponent" {
		t.Fatalf(`tag name should be "lib.foocomponent": "%s"`, name)
	}
	if typ := tag.Type; typ != app.CompoTag {
		t.Fatal("tag should be a component tag")
	}
	if count := len(tag.Children); count != 0 {
		t.Fatal("tag should not have children:", count)
	}
	if count := len(tag.Attributes); count != 1 {
		t.Fatal("tag should have 1 attribute:", count)
	}
	if val, _ := tag.Attributes["bar"]; val != "42" {
		t.Fatalf(`tag should have an attr with value = "42": %s`, val)
	}
}

func testDecodeCheckSvg(t *testing.T, tag app.Tag) {
	if name := tag.Name; name != "svg" {
		t.Fatalf(`tag name should be "svg": "%s"`, name)
	}
	if typ := tag.Type; typ != app.SimpleTag {
		t.Fatal("tag should be a simple tag")
	}
	if count := len(tag.Children); count != 2 {
		t.Fatal("tag should have 2 children:", count)
	}

	path1 := tag.Children[0]
	if name := path1.Name; name != "path" {
		t.Fatalf(`path1 name should be "path": "%s"`, name)
	}
	if typ := tag.Type; typ != app.SimpleTag {
		t.Fatal("tag should be a simple tag")
	}
	if !path1.Svg {
		t.Fatal("tag should have a svg context")
	}
	if count := len(path1.Children); count != 0 {
		t.Fatal("path1 should not have children:", count)
	}
	if count := len(path1.Attributes); count != 1 {
		t.Fatal("path1 should have 1 attribute:", count)
	}
	if d := path1.Attributes["d"]; d != "M 42.42 Z " {
		t.Fatalf(`path1 should have the attribute d="M 42.42 Z ": "%s"`, d)
	}

	path2 := tag.Children[1]
	if name := path2.Name; name != "path" {
		t.Fatalf(`path2 name should be "path": "%s"`, name)
	}
	if typ := tag.Type; typ != app.SimpleTag {
		t.Fatal("tag should be a simple tag")
	}
	if !path2.Svg {
		t.Fatal("tag should have a svg context")
	}
	if count := len(path2.Children); count != 0 {
		t.Fatal("path2 should not have children:", count)
	}
	if count := len(path2.Attributes); count != 1 {
		t.Fatal("path2 should have 1 attribute:", count)
	}
	if d := path2.Attributes["d"]; d != "M 21.21 Z " {
		t.Fatalf(`path2 should have the attribute d="M 21.21 Z ": "%s"`, d)
	}
}

func TestDecoderEmptyHTML(t *testing.T) {
	dec := NewDecoder(bytes.NewBufferString(""))

	var root app.Tag
	if err := dec.Decode(&root); err == nil {
		t.Fatal("err should not be nil")
	}
}

func TestDecoderSelfClosingTag(t *testing.T) {
	buff := bytes.NewBufferString(`
<p>
	<input/>
</p>
`)
	dec := NewDecoder(buff)
	var root app.Tag
	if err := dec.Decode(&root); err == nil {
		t.Fatal("err should not be nil")
	}
}

func TestDecoderNonClosingHTML(t *testing.T) {
	dec := NewDecoder(bytes.NewBufferString("<body><div>"))

	var root app.Tag
	if err := dec.Decode(&root); err != nil {
		t.Fatal(err)
	}
}

func TestIsVoidElement(t *testing.T) {
	if isVoidElement("path", true) {
		t.Error("path should not a void element")
	}

	if !isVoidElement("img", false) {
		t.Error("img should be a void element")
	}

	if isVoidElement("div", false) {
		t.Error("div should not be a void element")
	}
}

func TestIsComponent(t *testing.T) {
	if isComponent("", false) {
		t.Error("empty name should not be a component")
	}

	if isComponent("app.component", true) {
		t.Error("app.component should not be a component")
	}

	if !isComponent("app.component", false) {
		t.Error("app.component should be a component")
	}
}
