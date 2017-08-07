package markup

import (
	"strconv"
	"testing"
	"text/template"

	"github.com/google/uuid"
)

type Foo struct {
	Boo bool
}

func (c *Foo) OnMount() {}

func (c *Foo) OnDismount() {}

func (c *Foo) FuncMaps() template.FuncMap {
	return nil
}

func (c *Foo) Render() string {
	return `
<div>
	<h1>Foo</h1>
	<markup.bar>
</div>
	`
}

type Bar ZeroCompo

func (c *Bar) Render() string {
	return `<h2>Bar</h2>`
}

type CompoBadTmpl ZeroCompo

func (c *CompoBadTmpl) Render() string {
	return `<h2>{{.Hello}}</h2>`
}

type CompoBadTag ZeroCompo

func (c *CompoBadTag) Render() string {
	return `<h1><div/></h1>`
}

type CompoNotRegistered ZeroCompo

func (c *CompoNotRegistered) Render() string {
	return `
<div>
	<markup.unknown>
</div>
	`
}

type CompoBadChild ZeroCompo

func (c *CompoBadChild) Render() string {
	return `
<div>
	<markup.compobadtmpl>
</div>
	`
}

type CompoBadAttrs ZeroCompo

func (c *CompoBadAttrs) Render() string {
	return `
<div>
	<markup.foo boo="Holy Shit">
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
			<markup.world name="{{html .Name}}" err="{{.ChildErr}}" {{if .CompoFieldErr}}fielderr="-42"{{end}}>
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
		<markup.componotregistered>
	{{end}}
</div>
	`
}

func TestNewEnv(t *testing.T) {
	b := NewCompoBuilder()
	NewEnv(b)
}

func TestEnvComponent(t *testing.T) {
	compoID := uuid.New()
	foo := &Foo{}

	b := NewCompoBuilder()
	env := newEnv(b)
	env.components[compoID] = foo

	c, err := env.Component(compoID)
	if err != nil {
		t.Fatal(err)
	}
	if c != foo {
		t.Fatal("c and foo should point to the same component")
	}

	if _, err = env.Component(uuid.New()); err == nil {
		t.Fatal("err should not be nil")
	}
}

func TestEnvRoot(t *testing.T) {
	b := NewCompoBuilder()
	b.Register(&Foo{})
	b.Register(&Bar{})
	env := newEnv(b)

	foo := &Foo{}
	compoID := uuid.New()
	rootID := uuid.New()
	if _, err := env.mount(foo, rootID, compoID); err != nil {
		t.Fatal(err)
	}

	root, err := env.Root(foo)
	if err != nil {
		t.Fatal(err)
	}
	if root.ID != rootID {
		t.Fatal("rootID and root.ID should be equals")
	}

	if _, err = env.Root(&Bar{}); err == nil {
		t.Error("err should not be nil")
	}
	t.Log(err)
}

func TestEnv(t *testing.T) {
	b := NewCompoBuilder()
	b.Register(&Foo{})
	b.Register(&Bar{})
	b.Register(&CompoBadTmpl{})
	b.Register(&CompoBadTag{})
	b.Register(&CompoNotRegistered{})
	b.Register(&CompoBadChild{})
	b.Register(&Hello{})
	b.Register(&World{})

	env := newEnv(b)

	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "mount and dismount",
			test: func(t *testing.T) { testMountDismount(t, env, &Foo{}) },
		},
		{
			name: "mount mounted",
			test: func(t *testing.T) { testMountMounted(t, env, &Foo{}) },
		},
		{
			name: "mount component with bad template",
			test: func(t *testing.T) { testMountInvalid(t, env, &CompoBadTmpl{}) },
		},
		{
			name: "mount component with bad tag",
			test: func(t *testing.T) { testMountInvalid(t, env, &CompoBadTag{}) },
		},
		{
			name: "mount component with not registered child",
			test: func(t *testing.T) { testMountInvalid(t, env, &CompoNotRegistered{}) },
		},
		{
			name: "mount component with bad child",
			test: func(t *testing.T) { testMountInvalid(t, env, &CompoBadChild{}) },
		},
		{
			name: "mount component with bad attr",
			test: func(t *testing.T) { testMountInvalid(t, env, &CompoBadAttrs{}) },
		},
		{
			name: "dismount dismounted",
			test: func(t *testing.T) { testDismountDismounted(t, env, &Foo{}) },
		},
		{
			name: "dismount dismounted child",
			test: func(t *testing.T) { testDismountDismountedChild(t, env, &Foo{}) },
		},
		{
			name: "update should not do modifications",
			test: func(t *testing.T) { testEnvNoUpdate(t, env, &Hello{}) },
		},
		{
			name: "update should sync text",
			test: func(t *testing.T) { testEnvUpdateText(t, env, &Hello{Greeting: "Hi"}) },
		},
		{
			name: "update should merge html tag and component",
			test: func(t *testing.T) { testEnvUpdateMergeHTMLCompo(t, env, &Hello{}) },
		},
		{
			name: "update should merge html tag and text",
			test: func(t *testing.T) { testEnvUpdateMergeHTMLText(t, env, &Hello{}) },
		},
		{
			name: "update should merge text and html tag",
			test: func(t *testing.T) { testEnvUpdateMergeTextHTML(t, env, &Hello{TextBye: true}) },
		},
		{
			name: "update should sync components",
			test: func(t *testing.T) { testEnvUpdateComponent(t, env, &Hello{Name: "Jonhy"}) },
		},
		{
			name: "update should not sync component",
			test: func(t *testing.T) { testEnvUpdateComponentNoChange(t, env, &Hello{Name: "Jonhy"}) },
		},
		{
			name: "update should sync attributes",
			test: func(t *testing.T) { testEnvUpdateAttr(t, env, &Hello{}) },
		},
		{
			name: "update a not mounted component should fail",
			test: func(t *testing.T) { testEnvUpdateNotMounted(t, env, &Hello{}) },
		},
		{
			name: "update a component with bad template should fail",
			test: func(t *testing.T) { testEnvUpdateTemplateErr(t, env, &Hello{}) },
		},
		{
			name: "update a component with an error in its child should fail",
			test: func(t *testing.T) { testEnvUpdateTemplateChildErr(t, env, &Hello{Name: "Max"}) },
		},
		{
			name: "update should error when merge with error",
			test: func(t *testing.T) { testEnvUpdateMergeErr(t, env, &Hello{}) },
		},
		{
			name: "update should error when merge with component field error",
			test: func(t *testing.T) { testEnvUpdateCompoFieldErr(t, env, &Hello{Name: "Maxoo"}) },
		},
		{
			name: "update a component with dismounted child should error",
			test: func(t *testing.T) { testEnvUpdateSyncNotMountedComponent(t, env, &Hello{Name: "Maxoo"}) },
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.test)
	}
}

func testMountDismount(t *testing.T, env *env, c Componer) {
	// Mount.
	root, err := env.Mount(c)
	if err != nil {
		t.Fatal(err)
	}
	if count := len(env.components); count != 2 {
		t.Fatal("env shoud have 2 components:", count)
	}
	if count := len(env.compoRoots); count != 2 {
		t.Fatal("env shoud have 2 component roots:", count)
	}

	barTag := root.Children[1]
	if name := barTag.Name; name != "markup.bar" {
		t.Fatalf(`barTag.Name should be "markup.bar": "%s"`, name)
	}
	if _, err = env.Component(barTag.ID); err != nil {
		t.Fatal(err)
	}

	// Dismount
	env.Dismount(c)
	if count := len(env.components); count != 0 {
		t.Fatal("env shoud have 0 component:", count)
	}
	if count := len(env.compoRoots); count != 0 {
		t.Fatal("env shoud have 0 component root:", count)
	}
}

func testMountMounted(t *testing.T, env *env, c Componer) {
	if _, err := env.Mount(c); err != nil {
		t.Fatal(err)
	}
	defer env.Dismount(c)

	_, err := env.Mount(c)
	if err == nil {
		t.Fatal("err should not be nil")
	}
	t.Log(err)
}

func testMountInvalid(t *testing.T, env *env, c Componer) {
	_, err := env.Mount(c)
	if err == nil {
		t.Fatal("err should not be nil")
	}
	t.Log(err)
}

func testDismountDismounted(t *testing.T, env *env, c Componer) {
	if _, err := env.Mount(c); err != nil {
		t.Fatal(err)
	}
	env.Dismount(c)
	env.Dismount(c)
}

func testDismountDismountedChild(t *testing.T, env *env, c Componer) {
	root, err := env.Mount(c)
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range env.components {
		if k != root.CompoID {
			env.Dismount(v)
		}
	}
	env.Dismount(c)
}

func testEnvNoUpdate(t *testing.T, env *env, c *Hello) {
	if _, err := env.Mount(c); err != nil {
		t.Fatal(err)
	}

	syncs, err := env.Update(c)
	if err != nil {
		t.Fatal(err)
	}
	if len(syncs) != 0 {
		t.Error("syncs should be empty:", len(syncs))
	}
}

func testEnvUpdateText(t *testing.T, env *env, c *Hello) {
	if _, err := env.Mount(c); err != nil {
		t.Fatal(err)
	}

	c.Greeting = "Hello"

	syncs, err := env.Update(c)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(syncs); l != 1 {
		t.Fatal("syncs should have 1 element:", l)
	}

	s := syncs[0]
	if !s.Full {
		t.Fatal("s should be a full synchronization")
	}

	h := s.Tag
	if h.Name != "h1" {
		t.Fatal("tag to sync should be a h1:", h.Name)
	}

	if text := h.Children[0]; text.Text != c.Greeting {
		t.Fatalf(`text.Text should be "%s": "%s"`, c.Greeting, text.Text)
	}
}

func testEnvUpdateMergeHTMLCompo(t *testing.T, env *env, c *Hello) {
	if _, err := env.Mount(c); err != nil {
		t.Fatal(err)
	}

	c.Name = "Maxence"

	syncs, err := env.Update(c)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(syncs); l != 1 {
		t.Fatal("syncs should have 1 element:", l)
	}

	s := syncs[0]
	if !s.Full {
		t.Fatal("s should be a full synchronization")
	}

	world := s.Tag
	if world.Name != "markup.world" {
		t.Fatal("tag to sync should be a markup.world:", world.Name)
	}
	if name := world.Attrs["name"]; name != c.Name {
		t.Fatalf(`name should be "%s": "%s"`, c.Name, name)
	}
}

func testEnvUpdateMergeHTMLText(t *testing.T, env *env, c *Hello) {
	if _, err := env.Mount(c); err != nil {
		t.Fatal(err)
	}

	c.TextBye = true

	syncs, err := env.Update(c)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(syncs); l != 1 {
		t.Fatal("syncs should have 1 element:", l)
	}

	s := syncs[0]
	if !s.Full {
		t.Fatal("s should be a full synchronization")
	}

	root := s.Tag
	if root.Name != "div" {
		t.Fatal("root should be a div:", root.Name)
	}
	if l := len(root.Children); l != 4 {
		t.Fatal("root should have 4 children:", l)
	}
	if text := root.Children[3]; text.Text != "Goodbye" {
		t.Fatalf(`text should be "Goodbye": "%s"`, text.Text)
	}
}

func testEnvUpdateMergeTextHTML(t *testing.T, env *env, c *Hello) {
	if _, err := env.Mount(c); err != nil {
		t.Fatal(err)
	}

	c.TextBye = false

	syncs, err := env.Update(c)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(syncs); l != 1 {
		t.Fatal("syncs should have 1 element:", l)
	}

	s := syncs[0]
	if !s.Full {
		t.Fatal("s should be a full synchronization")
	}

	root := s.Tag
	if l := len(root.Children); l != 5 {
		t.Fatal("root should have 5 children:", l)
	}
	if span := root.Children[3]; span.Name != "span" {
		t.Fatalf(`span should be a span tag: %s`, span.Name)
	}
	if p := root.Children[4]; p.Name != "p" {
		t.Fatalf(`p should be a p tag: %s`, p.Name)
	}
}

func testEnvUpdateComponent(t *testing.T, env *env, c *Hello) {
	if _, err := env.Mount(c); err != nil {
		t.Fatal(err)
	}

	c.Name = "Maxence"

	syncs, err := env.Update(c)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(syncs); l != 1 {
		t.Fatal("syncs should have 1 element:", l)
	}

	s := syncs[0]
	if !s.Full {
		t.Fatal("s should be a full synchronization")
	}

	worldRoot := s.Tag
	if worldRoot.Name != "div" {
		t.Fatal("worldRoot should be a div:", worldRoot.Name)
	}
	if l := len(worldRoot.Children); l != 1 {
		t.Fatal("worldRoot should have 1 child:", l)
	}
	if text := worldRoot.Children[0]; text.Text != c.Name {
		t.Fatalf(`text.Text should be "%s": "%s"`, c.Name, text.Text)
	}
}

func testEnvUpdateComponentNoChange(t *testing.T, env *env, c *Hello) {
	if _, err := env.Mount(c); err != nil {
		t.Fatal(err)
	}

	syncs, err := env.Update(c)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(syncs); l != 0 {
		t.Fatal("syncs should be empty")
	}
}

func testEnvUpdateAttr(t *testing.T, env *env, c *Hello) {
	if _, err := env.Mount(c); err != nil {
		t.Fatal(err)
	}

	c.Placeholder = "Enter your name"

	syncs, err := env.Update(c)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(syncs); l != 1 {
		t.Fatal("syncs should have 1 element:", l)
	}

	s := syncs[0]
	if s.Full {
		t.Fatal("s should not be a full synchronization")
	}

	input := s.Tag
	if input.Name != "input" {
		t.Fatal("input should be an input:", input.Name)
	}
	if l := len(input.Children); l != 0 {
		t.Fatal("worldRoot should not have child")
	}
}

func testEnvUpdateNotMounted(t *testing.T, env *env, c *Hello) {
	_, err := env.Update(c)
	if err == nil {
		t.Error("err should not be nil")
	}
	t.Log(err)
}

func testEnvUpdateTemplateErr(t *testing.T, env *env, c *Hello) {
	if _, err := env.Mount(c); err != nil {
		t.Fatal(err)
	}

	c.TmplErr = true

	_, err := env.Update(c)
	if err == nil {
		t.Error("err should not be nil")
	}
	t.Log(err)
}

func testEnvUpdateTemplateChildErr(t *testing.T, env *env, c *Hello) {
	if _, err := env.Mount(c); err != nil {
		t.Fatal(err)
	}

	c.ChildErr = true

	_, err := env.Update(c)
	if err == nil {
		t.Error("err should not be nil")
	}
	t.Log(err)
}

func testEnvUpdateMergeErr(t *testing.T, env *env, c *Hello) {
	if _, err := env.Mount(c); err != nil {
		t.Fatal(err)
	}

	c.Name = "Maxence"
	c.ChildErr = true

	_, err := env.Update(c)
	if err == nil {
		t.Error("err should not be nil")
	}
	t.Log(err)
}

func testEnvUpdateCompoFieldErr(t *testing.T, env *env, c *Hello) {
	if _, err := env.Mount(c); err != nil {
		t.Fatal(err)
	}

	c.CompoFieldErr = true

	_, err := env.Update(c)
	if err == nil {
		t.Error("err should not be nil")
	}
	t.Log(err)
}

func testEnvUpdateSyncNotMountedComponent(t *testing.T, env *env, c *Hello) {
	root, err := env.Mount(c)
	if err != nil {
		t.Fatal(err)
	}

	world := root.Children[2].Children[0]
	env.dismountTag(world)

	c.Name = "Jonhy"

	if _, err = env.Update(c); err == nil {
		t.Error("err should not be nil")
	}
	t.Log(err)
}

func BenchmarkMount(b *testing.B) {
	bui := NewCompoBuilder()
	bui.Register(&Hello{})
	bui.Register(&World{})

	env := newEnv(bui)

	for i := 0; i < b.N; i++ {
		hello := &Hello{
			Name: "JonhyMaxoo",
		}
		env.Mount(hello)
		env.Dismount(hello)
	}
}

func BenchmarkSync(b *testing.B) {
	bui := NewCompoBuilder()
	bui.Register(&Hello{})
	bui.Register(&World{})

	env := newEnv(bui)

	hello := &Hello{
		Name: "JonhyMaxoo",
	}
	env.Mount(hello)

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

		env.Update(hello)

		alt = !alt
	}
}
