package markup

import (
	"bytes"
	"testing"
)

func TestTagIsEmpty(t *testing.T) {
	tag := Tag{}
	if !tag.IsEmpty() {
		t.Error("tag should be empty")
	}
}

func TestTagIsText(t *testing.T) {
	tag := Tag{Text: "foo"}
	if !tag.IsText() {
		t.Error("tag should be a text")
	}

	tag = Tag{Name: "div", Text: "foo"}
	if tag.IsText() {
		t.Error("tag should not be a text")
	}

	tag = Tag{}
	if tag.IsText() {
		t.Error("tag should not be a text")
	}
}

func TestTagIsComponent(t *testing.T) {
	tag := Tag{Name: "foo"}
	if !tag.IsComponent() {
		t.Error("tag should be a component")
	}

	tag = Tag{Name: "foo", Svg: true}
	if tag.IsComponent() {
		t.Error("tag should not be a component")
	}

	tag = Tag{Name: "div"}
	if tag.IsComponent() {
		t.Error("tag should not be a component")
	}

	tag = Tag{}
	if tag.IsComponent() {
		t.Error("tag should not be a component")
	}
}

func TestTagIsVoidElement(t *testing.T) {
	tag := Tag{Name: "input"}
	if !tag.IsVoidElem() {
		t.Error("tag should be a void element")
	}

	tag = Tag{Name: "link", Svg: true}
	if tag.IsVoidElem() {
		t.Error("tag should be a void element")
	}

	tag = Tag{Name: "div"}
	if tag.IsVoidElem() {
		t.Error("tag should not be a void element")
	}
}

func TestAttrEquals(t *testing.T) {
	attr := AttrMap{
		"hello": "world",
		"foo":   "bar",
	}

	attr2 := AttrMap{
		"foo":   "bar",
		"hello": "world",
	}

	if !AttrEquals(attr, attr2) {
		t.Error("attr and attr2 should be equals")
	}

	if AttrEquals(attr, nil) {
		t.Error("attr and nil should not be equals")
	}

	attr3 := AttrMap{
		"foo":   "bar",
		"hello": "maxoo",
	}

	if AttrEquals(attr, attr3) {
		t.Error("attr and attr3 should not be equals")
	}

	attr4 := AttrMap{
		"foo": "bar",
		"bye": "world",
	}

	if AttrEquals(attr, attr4) {
		t.Error("attr and attr4 should not be equals")
	}
}

func TestTagEncoderEncode(t *testing.T) {
	b := NewCompoBuilder()
	b.Register(&Hello{})
	b.Register(&World{})

	env := newEnv(b)

	hello := &Hello{
		Name: "JonhyMaxoo",
	}

	root, err := env.Mount(hello)
	if err != nil {
		t.Fatal(err)
	}

	w := &bytes.Buffer{}
	enc := NewTagEncoder(w, env)
	if err := enc.Encode(root); err != nil {
		t.Fatal(err)
	}
	t.Log(w.String())

	errRoot := Tag{
		Name: "markup.hello",
	}
	if err = enc.Encode(errRoot); err == nil {
		t.Fatal("err should not be nil")
	}
	t.Log(err)
}

func BenchmarkTagEncoder(b *testing.B) {
	bui := NewCompoBuilder()
	bui.Register(&Hello{})
	bui.Register(&World{})

	env := newEnv(bui)

	hello := &Hello{
		Name: "JonhyMaxoo",
	}

	root, _ := env.Mount(hello)
	for i := 0; i < b.N; i++ {
		var v bytes.Buffer
		enc := NewTagEncoder(&v, env)
		enc.Encode(root)
	}
}

func TestDecode(t *testing.T) {
	h := `
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
	`

	b := bytes.NewBufferString(h)
	d := NewTagDecoder(b)
	root := Tag{}
	if err := d.Decode(&root); err != nil {
		t.Fatal(err)
	}

	testDecodeCheckRoot(t, root)
	testDecodeCheckH1(t, root.Children[0])
	testDecodeCheckBr(t, root.Children[1])
	testDecodeCheckInput(t, root.Children[2])
	testDecodeCheckFooComponent(t, root.Children[3])
	testDecodeCheckSvg(t, root.Children[4])
}

func testDecodeCheckRoot(t *testing.T, tag Tag) {
	if name := tag.Name; name != "div" {
		t.Fatalf(`tag name should be "div": "%s"`, name)
	}
	if count := len(tag.Children); count != 5 {
		t.Fatal("tag should have 5 children:", count)
	}
}

func testDecodeCheckH1(t *testing.T, tag Tag) {
	if name := tag.Name; name != "h1" {
		t.Fatalf(`tag name should be "h1": "%s"`, name)
	}
	if count := len(tag.Children); count != 1 {
		t.Fatal("tag should have 1 children:", count)
	}
	if text := tag.Children[0]; text.Text != "hello" {
		t.Fatalf(`text.Text should be "hello": "%s"`, text.Text)
	}
}

func testDecodeCheckBr(t *testing.T, tag Tag) {
	if name := tag.Name; name != "br" {
		t.Fatalf(`tag name should be "br": "%s"`, name)
	}
	if count := len(tag.Children); count != 0 {
		t.Fatal("root should not have children:", count)
	}
}

func testDecodeCheckInput(t *testing.T, tag Tag) {
	if name := tag.Name; name != "input" {
		t.Fatalf(`tag name should be "input": "%s"`, name)
	}
	if count := len(tag.Children); count != 0 {
		t.Fatal("tag should not have children:", count)
	}
	if count := len(tag.Attrs); count != 2 {
		t.Fatal("tag should have 2 attributes:", count)
	}
	if val, _ := tag.Attrs["type"]; val != "text" {
		t.Fatalf(`tag should have an attr with value = "text": %s`, val)
	}
	if _, ok := tag.Attrs["required"]; !ok {
		t.Fatal(`tag should have an attr with key = "required"`)
	}
}

func testDecodeCheckFooComponent(t *testing.T, tag Tag) {
	if name := tag.Name; name != "lib.foocomponent" {
		t.Fatalf(`tag name should be "lib.foocomponent": "%s"`, name)
	}
	if count := len(tag.Children); count != 0 {
		t.Fatal("tag should not have children:", count)
	}
	if count := len(tag.Attrs); count != 1 {
		t.Fatal("tag should have 1 attribute:", count)
	}
	if val, _ := tag.Attrs["bar"]; val != "42" {
		t.Fatalf(`tag should have an attr with value = "42": %s`, val)
	}
}

func testDecodeCheckSvg(t *testing.T, tag Tag) {
	if name := tag.Name; name != "svg" {
		t.Fatalf(`tag name should be "svg": "%s"`, name)
	}
	if count := len(tag.Children); count != 2 {
		t.Fatal("tag should have 2 children:", count)
	}

	path1 := tag.Children[0]
	if name := path1.Name; name != "path" {
		t.Fatalf(`path1 name should be "path": "%s"`, name)
	}
	if count := len(path1.Children); count != 0 {
		t.Fatal("path1 should not have children:", count)
	}
	if count := len(path1.Attrs); count != 1 {
		t.Fatal("path1 should have 1 attribute:", count)
	}
	if d := path1.Attrs["d"]; d != "M 42.42 Z " {
		t.Fatalf(`path1 should have the attribute d="M 42.42 Z ": "%s"`, d)
	}

	path2 := tag.Children[1]
	if name := path2.Name; name != "path" {
		t.Fatalf(`path2 name should be "path": "%s"`, name)
	}
	if count := len(path2.Children); count != 0 {
		t.Fatal("path2 should not have children:", count)
	}
	if count := len(path2.Attrs); count != 1 {
		t.Fatal("path2 should have 1 attribute:", count)
	}
	if d := path2.Attrs["d"]; d != "M 21.21 Z " {
		t.Fatalf(`path2 should have the attribute d="M 21.21 Z ": "%s"`, d)
	}
}

func TestDecodeSelfClosingTagError(t *testing.T) {
	h := `
<p>
	<input/>
</p>
`

	b := bytes.NewBufferString(h)
	d := NewTagDecoder(b)
	root := Tag{}
	if err := d.Decode(&root); err == nil {
		t.Fatal("err should not be nil")
	}
}

func TestDecodeEmptyHTML(t *testing.T) {
	h := ""

	b := bytes.NewBufferString(h)
	d := NewTagDecoder(b)

	root := Tag{}
	if err := d.Decode(&root); err == nil {
		t.Fatal("err should not be nil")
	}
}

func TestDecodeNonClosingHTML(t *testing.T) {
	h := "<body><div>"

	b := bytes.NewBufferString(h)
	d := NewTagDecoder(b)

	root := Tag{}
	if err := d.Decode(&root); err != nil {
		t.Fatal(err)
	}
}
