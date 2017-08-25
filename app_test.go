package app

import (
	"testing"

	"github.com/murlokswarm/app/markup"
)

type Component markup.ZeroCompo

func (c *Component) Render() string {
	return `<div>Hello</div>`
}

type InvalidComponent markup.ZeroCompo

func (c InvalidComponent) Render() string {
	return ``
}

func TestApp(t *testing.T) {
	d := &driverTest{
		Test: t,
	}

	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "should import component",
			test: testImport,
		},
		{
			name: "import invalid component should panic",
			test: testImportPanic,
		},
		{
			name: "should run",
			test: func(t *testing.T) { testRun(t, d) },
		},
		{
			name: "second run should panic",
			test: testRunPanic,
		},
		{
			name: "should return the running driver",
			test: func(t *testing.T) { testRunningDriver(t, d) },
		},
		{
			name: "running driver when app is not running should panic",
			test: testRunningDriverPanic,
		},
		{
			name: "should render a component",
			test: testRender,
		},
		{
			name: "render should log an error",
			test: testRenderLogError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.test)
	}
}

func testImport(t *testing.T) {
	Import(&Component{})
}

func testImportPanic(t *testing.T) {
	defer func() { recover() }()

	Import(InvalidComponent{})
	t.Error("should panic")
}

func testRun(t *testing.T, d Driver) {
	if err := Run(d); err != nil {
		t.Fatal(err)
	}
}

func testRunPanic(t *testing.T) {
	defer func() { recover() }()

	Run(&driverTest{
		Test: t,
	})
	t.Error("should panic")
}

func testRunningDriver(t *testing.T, d Driver) {
	if RunningDriver() != d {
		t.Fatal("running driver should be d")
	}
}

func testRunningDriverPanic(t *testing.T) {
	d := driver
	driver = nil
	defer func() { driver = d }()
	defer func() { recover() }()

	RunningDriver()
	t.Error("should panic")
}

func testRender(t *testing.T) {
	Render(&Component{})
}

func testRenderLogError(t *testing.T) {
	Render(&Component{})
}
