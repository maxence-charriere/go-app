package app

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompoMountDismount(t *testing.T) {
	testMountDismount(t, []mountTest{
		{
			scenario: "component",
			node:     &hello{},
		},
	})
}

func TestCompoUpdate(t *testing.T) {
	testUpdate(t, []updateTest{
		{
			scenario: "component is updated",
			a:        &bar{Value: "rab"},
			b:        &bar{Value: "bar"},
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: &bar{Value: "bar"},
				},
				{
					Path:     TestPath(0),
					Expected: Text("bar"),
				},
			},
		},
		{
			scenario: "component is updated",
			a:        &hello{},
			b:        &hello{Greeting: "world"},
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: &hello{Greeting: "world"},
				},
				{
					Path:     TestPath(0),
					Expected: Div(),
				},
				{
					Path:     TestPath(0, 0),
					Expected: H1(),
				},
				{
					Path:     TestPath(0, 0, 0),
					Expected: Text("hello, "),
				},
				{
					Path:     TestPath(0, 0, 1),
					Expected: Text("world"),
				},
			},
		},
		{
			scenario: "component is replaced by a text",
			a: Div().Body(
				&hello{},
			),
			b: Div().Body(
				Text("hello"),
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},
				{
					Path:     TestPath(0),
					Expected: Text("hello"),
				},
			},
		},
		{
			scenario: "component is replaced by an html element",
			a: Div().Body(
				&hello{},
			),
			b: Div().Body(
				H2().Text("hello"),
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},
				{
					Path:     TestPath(0),
					Expected: H2(),
				},
				{
					Path:     TestPath(0, 0),
					Expected: Text("hello"),
				},
			},
		},
		{
			scenario: "component is replaced by a raw html element",
			a: Div().Body(
				&hello{},
			),
			b: Div().Body(
				Raw("<svg></svg>"),
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},
				{
					Path:     TestPath(0),
					Expected: Raw("<svg></svg>"),
				},
			},
		},
		{
			scenario: "component is replaced by another component",
			a: Div().Body(
				&hello{},
			),
			b: Div().Body(
				&bar{},
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},
				{
					Path:     TestPath(0),
					Expected: &bar{},
				},
				{
					Path:     TestPath(0, 0),
					Expected: Text(""),
				},
			},
		},
		{
			scenario: "component root is updated",
			a: Div().Body(
				&foo{Bar: "hello"},
			),
			b: Div().Body(
				&foo{Bar: "goodbye"},
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},
				{
					Path:     TestPath(0),
					Expected: &foo{Bar: "goodbye"},
				},
				{
					Path:     TestPath(0, 0),
					Expected: &bar{Value: "goodbye"},
				},
				{
					Path:     TestPath(0, 0, 0),
					Expected: Text("goodbye"),
				},
			},
		},
		{
			scenario: "component root is replaced by a component",
			a: Div().Body(
				&foo{},
			),
			b: Div().Body(
				&foo{Bar: "test"},
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},
				{
					Path:     TestPath(0),
					Expected: &foo{Bar: "test"},
				},
				{
					Path:     TestPath(0, 0),
					Expected: &bar{Value: "test"},
				},
				{
					Path:     TestPath(0, 0, 0),
					Expected: Text("test"),
				},
			},
		},
		{
			scenario: "component root is replaced by a non-component",
			a: Div().Body(
				&foo{Bar: "test"},
			),
			b: Div().Body(
				&foo{},
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},
				{
					Path:     TestPath(0),
					Expected: &foo{},
				},
				{
					Path:     TestPath(0, 0),
					Expected: Text("bar"),
				},
			},
		},
	})
}

func TestNavigator(t *testing.T) {
	u, _ := url.Parse("https://murlok.io")
	h := &hello{}

	d := NewClientTester(h)
	defer d.Close()

	d.Nav(u)
	d.Consume()
	require.Equal(t, "https://murlok.io", h.onNavURL)
}

func TestNestedtNavigator(t *testing.T) {
	u, _ := url.Parse("https://murlok.io")

	h := &hello{}
	div := Div().Body(h)
	d := NewClientTester(div)
	defer d.Close()

	d.Nav(u)
	d.Consume()
	require.Equal(t, "https://murlok.io", h.onNavURL)
}

func TestNestedInComponentNavigator(t *testing.T) {
	u, _ := url.Parse("https://murlok.io")

	foo := &foo{Bar: "Bar"}
	d := NewClientTester(foo)
	defer d.Close()

	d.Nav(u)
	d.Consume()
	b := foo.getChildren()[0].(*bar)
	require.Equal(t, "https://murlok.io", b.onNavURL)
}

func TestAppUpdater(t *testing.T) {
	appUpdateAvailable = true
	defer func() {
		appUpdateAvailable = false
	}()

	h := &hello{}
	d := NewClientTester(h)
	defer d.Close()

	d.AppUpdate()
	d.Consume()
	require.True(t, h.appUpdated)
}

func TestNestedAppUpdater(t *testing.T) {
	appUpdateAvailable = true
	defer func() {
		appUpdateAvailable = false
	}()

	h := &hello{}
	div := Div().Body(h)
	d := NewClientTester(div)
	defer d.Close()

	d.AppUpdate()
	d.Consume()
	require.True(t, h.appUpdated)
}

func TestNestedInComponentAppUpdater(t *testing.T) {
	appUpdateAvailable = true
	defer func() {
		appUpdateAvailable = false
	}()

	foo := &foo{Bar: "Bar"}
	d := NewClientTester(foo)
	defer d.Close()

	d.AppUpdate()
	d.Consume()
	b := foo.getChildren()[0].(*bar)
	require.True(t, b.appUpdated)
}

func TestAppInstaller(t *testing.T) {
	h := &hello{}
	d := NewClientTester(h)
	defer d.Close()

	d.AppInstallChange()
	d.Consume()
	require.True(t, h.appInstalled)
}

func TestNestedAppInstaller(t *testing.T) {
	h := &hello{}
	div := Div().Body(h)
	d := NewClientTester(div)
	defer d.Close()

	d.AppInstallChange()
	d.Consume()
	require.True(t, h.appInstalled)
}

func TestNestedInComponentAppInstaller(t *testing.T) {
	foo := &foo{Bar: "Bar"}
	d := NewClientTester(foo)
	defer d.Close()

	d.AppInstallChange()
	d.Consume()
	b := foo.getChildren()[0].(*bar)
	require.True(t, b.appInstalled)
}

func TestResizer(t *testing.T) {
	h := &hello{}
	d := NewClientTester(h)
	defer d.Close()

	d.AppResize()
	d.Consume()
	require.True(t, h.appResized)
}

func TestNestedResizer(t *testing.T) {
	h := &hello{}
	div := Div().Body(h)
	d := NewClientTester(div)
	defer d.Close()

	d.AppResize()
	d.Consume()
	require.True(t, h.appResized)
}

func TestNestedInComponentResizer(t *testing.T) {
	foo := &foo{Bar: "Bar"}
	d := NewClientTester(foo)
	defer d.Close()

	d.AppResize()
	d.Consume()
	b := foo.getChildren()[0].(*bar)
	require.True(t, b.appRezized)
}

func TestPreRenderer(t *testing.T) {
	h := &hello{}
	d := NewServerTester(h)
	defer d.Close()

	d.Consume()
	require.True(t, h.preRenderer)
	require.Equal(t, "world", d.getCurrentPage().Title())
}

func TestNestedPreRenderer(t *testing.T) {
	h := &hello{}
	div := Div().Body(h)
	d := NewServerTester(div)
	defer d.Close()

	d.Consume()
	require.True(t, h.preRenderer)
	require.Equal(t, "world", d.getCurrentPage().Title())
}

func TestNestedInComponentPreRenderer(t *testing.T) {
	foo := &foo{Bar: "Bar"}
	d := NewServerTester(foo)
	defer d.Close()

	d.Consume()
	require.Equal(t, "bar", d.getCurrentPage().Title())
}

func TestUpdater(t *testing.T) {
	b := &bar{Value: "BAR"}
	d := NewClientTester(b)
	defer d.Close()

	require.False(t, b.updated)

	d.Mount(&bar{Value: "BAR"})
	d.Consume()
	require.False(t, b.updated)

	d.Mount(&bar{Value: "Bar"})
	d.Consume()
	require.Equal(t, "Bar", b.Value)
	require.True(t, b.updated)
}

func TestInitializerServer(t *testing.T) {
	b := &bar{}
	d := NewServerTester(b)
	defer d.Close()
	require.True(t, b.initialized)
}

func TestInitializerClient(t *testing.T) {
	b := &bar{}
	d := NewClientTester(b)
	defer d.Close()
	require.True(t, b.initialized) // b is mounted in NewClientTester
}

type hello struct {
	Compo

	Greeting     string
	onNavURL     string
	appUpdated   bool
	appInstalled bool
	appResized   bool
	preRenderer  bool
}

func (h *hello) OnMount(Context) {
}

func (h *hello) OnNav(ctx Context) {
	h.onNavURL = ctx.Page().URL().String()
}

func (h *hello) OnAppUpdate(ctx Context) {
	h.appUpdated = true
}

func (h *hello) OnAppInstallChange(ctx Context) {
	h.appInstalled = true
}

func (h *hello) OnResize(ctx Context) {
	h.appResized = true
}

func (h *hello) OnPreRender(ctx Context) {
	h.preRenderer = true
	ctx.Page().SetTitle("world")
}

func (h *hello) OnDismount(Context) {
}

func (h *hello) Render() UI {
	return Div().Body(
		H1().Body(
			Text("hello, "),
			Text(h.Greeting),
		),
	)
}

type foo struct {
	Compo
	Bar string
}

func (f *foo) Render() UI {
	return If(f.Bar != "",
		&bar{Value: f.Bar},
	).Else(
		Text("bar"),
	)
}

type bar struct {
	Compo
	Value string

	onNavURL     string
	appUpdated   bool
	appInstalled bool
	appRezized   bool
	updated      bool
	initialized  bool
}

func (b *bar) OnInit() {
	b.initialized = true
}

func (b *bar) OnPreRender(ctx Context) {
	ctx.Page().SetTitle("bar")
}

func (b *bar) OnNav(ctx Context) {
	b.onNavURL = ctx.Page().URL().String()
}

func (b *bar) OnAppUpdate(ctx Context) {
	b.appUpdated = true
}

func (b *bar) OnAppInstallChange(ctx Context) {
	b.appInstalled = true
}

func (b *bar) OnResize(ctx Context) {
	b.appRezized = true
}

func (b *bar) OnUpdate(ctx Context) {
	b.updated = true
}

func (b *bar) Render() UI {
	return Text(b.Value)
}
