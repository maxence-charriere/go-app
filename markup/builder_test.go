package markup

import (
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestNewCompoBuilder(t *testing.T) {
	NewCompoBuilder()
}

func TestCompoBuilderRegister(t *testing.T) {
	c := &ValidCompo{}
	ct := reflect.TypeOf(*c)
	cname := strings.ToLower(ct.String())

	b := make(compoBuilder)
	err := b.Register(c)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := b[cname]; !ok {
		t.Fatalf("%s should have been registered", cname)
	}

	empc := &EmptyCompo{}
	if err = b.Register(empc); err == nil {
		t.Fatal("register empty component should returns an error")
	}
}

func TestCompoBuilderNew(t *testing.T) {
	c := &ValidCompo{}
	cname := "markup.validcompo"
	b := make(compoBuilder)
	b.Register(c)

	n, err := b.New(cname)
	if err != nil {
		t.Fatal(err)
	}
	if n == nil {
		t.Fatalf("%s should have been created: %v", cname, n)
	}

	if _, err = b.New("unknown"); err == nil {
		t.Fatal("unknown should not have been created")
	}
}

func TestEnsureValidCompo(t *testing.T) {
	valc := &ValidCompo{}
	if err := ensureValidComponent(valc); err != nil {
		t.Error(err)
	}

	noptrc := NonPtrCompo{}
	if err := ensureValidComponent(noptrc); err == nil {
		t.Error("err should not be nil")
	}

	empc := &EmptyCompo{}
	if err := ensureValidComponent(empc); err == nil {
		t.Error("err should not be nil")
	}

	intc := IntCompo(42)
	if err := ensureValidComponent(&intc); err == nil {
		t.Error("err should not be nil")
	}
}

func TestNormalizeCompoName(t *testing.T) {
	if name := "lib.FooBar"; normalizeCompoName(name) != "lib.foobar" {
		t.Errorf(`name should be "lib.foobar": "%s"`, name)
	}

	if name := "main.FooBar"; normalizeCompoName(name) != "foobar" {
		t.Errorf(`name should be "foobar": "%s"`, name)
	}
}

func TestComponentNameFromURL(t *testing.T) {
	// Default.
	rawurl := "/markup.hello"
	u, err := url.Parse(rawurl)
	if err != nil {
		t.Fatal(err)
	}

	name, ok := ComponentNameFromURL(u)
	if !ok {
		t.Fatalf("%s should point to a component", rawurl)
	}
	if name != "markup.hello" {
		t.Fatal("component name should be markup.hello:", name)
	}

	// Component name as scheme.
	rawurl = "component://markup.hello"
	if u, err = url.Parse(rawurl); err != nil {
		t.Fatal(err)
	}

	if name, ok = ComponentNameFromURL(u); ok {
		t.Fatalf("%s should not point to a component", rawurl)
	}

	// Bad scheme.
	rawurl = "https://github.com"
	if u, err = url.Parse(rawurl); err != nil {
		t.Fatal(err)
	}

	if name, ok = ComponentNameFromURL(u); ok {
		t.Fatalf("%s should not point to a component", rawurl)
	}
}
