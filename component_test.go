package app

import (
	"net/url"
	"testing"
)

type ValidCompo ZeroCompo

func (c *ValidCompo) Render() string {
	return `<p>Hello World</p>`
}

type NonPtrCompo ZeroCompo

func (c NonPtrCompo) Render() string {
	return `<p>Bye World</p>`
}

type IntCompo int

func (i *IntCompo) Render() string {
	return `<p>Aurevoir World</p>`
}

type EmptyCompo struct{}

func (c *EmptyCompo) Render() string {
	return `<p>Goodbye World</p>`
}

func TestFactory(t *testing.T) {
	tests := []struct {
		scenario string
		function func(t *testing.T, factory Factory)
	}{
		{
			scenario: "registers a component",
			function: testFactoryRegisterComponent,
		},
		{
			scenario: "registering a component not implemented on pointer returns an error",
			function: testFactoryRegisterComponentNoPtr,
		},
		{
			scenario: "registering a component not implemented on a struct pointer returns an error",
			function: testFactoryRegisterComponentNoStructPtr,
		},
		{
			scenario: "registering a component implemented on an empty struct pointer returns an error",
			function: testFactoryRegisterComponentEmptyStructPtr,
		},
		{
			scenario: "creates a component",
			function: testFactoryCreateComponent,
		},
		{
			scenario: "creating a not registered component returns an error",
			function: testFactoryCreateNotRegisteredComponent,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			test.function(t, NewFactory())
		})
	}
}

func testFactoryRegisterComponent(t *testing.T, factory Factory) {
	name, err := factory.RegisterComponent(&ValidCompo{})
	if err != nil {
		t.Fatal(err)
	}
	if name != "app.validcompo" {
		t.Error("name is not app.validcompo:", name)
	}
}

func testFactoryRegisterComponentNoPtr(t *testing.T, factory Factory) {
	_, err := factory.RegisterComponent(NonPtrCompo{})
	if err == nil {
		t.Fatal("err is nil")
	}
	t.Log(err)
}

func testFactoryRegisterComponentNoStructPtr(t *testing.T, factory Factory) {
	intc := IntCompo(42)
	_, err := factory.RegisterComponent(&intc)
	if err == nil {
		t.Fatal("err is nil")
	}
	t.Log(err)
}

func testFactoryRegisterComponentEmptyStructPtr(t *testing.T, factory Factory) {
	_, err := factory.RegisterComponent(&EmptyCompo{})
	if err == nil {
		t.Fatal("err is nil")
	}
	t.Log(err)
}

func testFactoryCreateComponent(t *testing.T, factory Factory) {
	if _, err := factory.RegisterComponent(&ValidCompo{}); err != nil {
		t.Fatal(err)
	}

	compo, err := factory.NewComponent("app.validcompo")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(compo)
}

func testFactoryCreateNotRegisteredComponent(t *testing.T, factory Factory) {
	_, err := factory.NewComponent("app.validcompo")
	if err == nil {
		t.Fatal("err is nil")
	}
	t.Log(err)
}

func TestNormalizeComponentName(t *testing.T) {
	if name := "lib.FooBar"; normalizeComponentName(name) != "lib.foobar" {
		t.Errorf("name is not lib.foobar: %s", name)
	}

	if name := "main.FooBar"; normalizeComponentName(name) != "foobar" {
		t.Errorf("name is not foobar: %s", name)
	}
}

func TestComponentNameFromURL(t *testing.T) {
	u1, _ := url.Parse("/hello")
	u2, _ := url.Parse("/hello?int=42")
	u3, _ := url.Parse("/hello/world")
	urls := []*url.URL{u1, u2, u3}

	for _, u := range urls {
		if name := ComponentNameFromURL(u); name != "hello" {
			t.Error("name is not hello:", name)
		}
	}

	u0 := &url.URL{
		Host: "test",
	}
	if name := ComponentNameFromURL(u0); len(name) != 0 {
		t.Error("name is not empty")
	}
}
