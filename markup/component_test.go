package markup

import (
	"testing"
	"time"
)

type ValidCompo ZeroCompo

func (c *ValidCompo) Render() string {
	return `<p>Hello World</p>`
}

type EmptyCompo struct{}

func (c *EmptyCompo) Render() string {
	return `<p>Goodbye World</p>`
}

type NonPtrCompo ZeroCompo

func (c NonPtrCompo) Render() string {
	return `<p>Bye World</p>`
}

type IntCompo int

func (i *IntCompo) Render() string {
	return `<p>Aurevoir World</p>`
}

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
}

func (c *CompoWithFields) Render() string {
	return `<div></div>`
}

func (c *CompoWithFields) Func() {
	c.funcHandler()
}

func (c *CompoWithFields) FuncWithArg(nb int) {
	c.funcWithArgHandler(nb)
}

func (c *CompoWithFields) FuncWithMultipleArg(a, b int) {
	panic("should not be called")
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

func TestMapComponentFields(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "empty",
			test: func(t *testing.T) {
				testMapComponentFields(t, nil)
			},
		},
		{
			name: "anonymous",
			test: func(t *testing.T) {
				attrs := AttrMap{"zerocompo": `{"placeholder": 42}`}
				testMapComponentFields(t, attrs)
			},
		},
		{
			name: "unexported",
			test: func(t *testing.T) {
				attrs := AttrMap{"secret": "pandore"}
				testMapComponentFields(t, attrs)
			},
		},
		{
			name: "string",
			test: func(t *testing.T) {
				attrs := AttrMap{"string": "hello"}
				testMapComponentFields(t, attrs)
			},
		},
		{
			name: "bool",
			test: func(t *testing.T) {
				attrs := AttrMap{"bool": "true"}
				testMapComponentFields(t, attrs)
			},
		},
		{
			name: "bool error",
			test: func(t *testing.T) {
				attrs := AttrMap{"bool": "lkdsja"}
				testMapComponentFieldsErrors(t, attrs)
			},
		},
		{
			name: "int",
			test: func(t *testing.T) {
				attrs := AttrMap{"int": "42"}
				testMapComponentFields(t, attrs)
			},
		},
		{
			name: "int error",
			test: func(t *testing.T) {
				attrs := AttrMap{"int": "zzedgw"}
				testMapComponentFieldsErrors(t, attrs)
			},
		},
		{
			name: "uint",
			test: func(t *testing.T) {
				attrs := AttrMap{"uint": "42"}
				testMapComponentFields(t, attrs)
			},
		},
		{
			name: "uint error",
			test: func(t *testing.T) {
				attrs := AttrMap{"uint": "-42"}
				testMapComponentFieldsErrors(t, attrs)
			},
		},
		{
			name: "float",
			test: func(t *testing.T) {
				attrs := AttrMap{"float": "42.42"}
				testMapComponentFields(t, attrs)
			},
		},
		{
			name: "float error",
			test: func(t *testing.T) {
				attrs := AttrMap{"float": "-42.zdf"}
				testMapComponentFieldsErrors(t, attrs)
			},
		},
		{
			name: "struct",
			test: func(t *testing.T) {
				attrs := AttrMap{"struct": `{"A": 42, "B": "world"}`}
				testMapComponentFields(t, attrs)
			},
		},
		{
			name: "struct error",
			test: func(t *testing.T) {
				attrs := AttrMap{"struct": `{"A": "world", "B": 42}`}
				testMapComponentFieldsErrors(t, attrs)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.test)
	}
}

func testMapComponentFields(t *testing.T, attrs AttrMap) {
	c := &CompoWithFields{}
	if err := mapComponentFields(c, attrs); err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", c)
}

func testMapComponentFieldsErrors(t *testing.T, attrs AttrMap) {
	c := CompoWithFields{}
	err := mapComponentFields(&c, attrs)
	if err == nil {
		t.Fatal("err should not be nil")
	}
	t.Log(err)
}

func TestConvertToJSON(t *testing.T) {
	c := &CompoWithFields{}
	t.Log(convertToJSON(c))
}

func TestFormatTime(t *testing.T) {
	t.Log(formatTime(time.Now(), "2006"))
}

func TestCallOrAssignMethod(t *testing.T) {
	funcCalled := false
	funcWithArgValue := 0

	c := &CompoWithFields{
		funcHandler: func() {
			funcCalled = true
		},
		funcWithArgHandler: func(nb int) {
			funcWithArgValue = nb
		},
	}

	if err := CallOrAssign(c, "Func", ""); err != nil {
		t.Fatal(err)
	}
	if !funcCalled {
		t.Fatal("funcCalled should be true")
	}

	if err := CallOrAssign(c, "FuncWithArg", "42"); err != nil {
		t.Fatal(err)
	}
	if funcWithArgValue != 42 {
		t.Fatal("funcWithArgValue should be 42:", funcWithArgValue)
	}
}

func TestCallOrAssignField(t *testing.T) {
	c := &CompoWithFields{}

	if err := CallOrAssign(c, "String", `"hello world"`); err != nil {
		t.Fatal(err)
	}
	if c.String != "hello world" {
		t.Fatalf(`c.String should be "hello world": "%s"`, c.String)
	}

	if err := CallOrAssign(c, "Int", "42"); err != nil {
		t.Fatal(err)
	}
	if c.Int != 42 {
		t.Fatal("c.Int should be 42:", c.Int)
	}
}

func TestCallOrAssignErrors(t *testing.T) {
	c := &CompoWithFields{}

	err := CallOrAssign(c, "NonexitentMethodOrField", "")
	if err == nil {
		t.Fatal("err should not be nil")
	}
	t.Log(err)

	if err = CallOrAssign(c, "Render", ""); err == nil {
		t.Fatal("err should not be nil")
	}
	t.Log(err)

	if err = CallOrAssign(c, "FuncWithMultipleArg", ""); err == nil {
		t.Fatal("err should not be nil")
	}
	t.Log(err)

	if err = CallOrAssign(c, "FuncWithArg", "}{"); err == nil {
		t.Fatal("err should not be nil")
	}
	t.Log(err)
}
