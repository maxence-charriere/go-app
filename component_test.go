package app

import "testing"

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
			scenario: "should register a component",
			function: testFactoryRegisterComponent,
		},
		{
			scenario: "register a component not implemented on pointer should fail",
			function: testFactoryRegisterComponentNoPtr,
		},
		{
			scenario: "register a component not implemented on a struct pointer should fail",
			function: testFactoryRegisterComponentNoStructPtr,
		},
		{
			scenario: "register a component implemented on an empty struct pointer should fail",
			function: testFactoryRegisterComponentEmptyStructPtr,
		},
		{
			scenario: "should create a component",
			function: testFactoryCreateComponent,
		},
		{
			scenario: "create a not registered component should fail",
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
		t.Error("name should be app.validcompo:", name)
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
		t.Errorf("name should be lib.foobar: %s", name)
	}

	if name := "main.FooBar"; normalizeComponentName(name) != "foobar" {
		t.Errorf("name should be foobar: %s", name)
	}
}
