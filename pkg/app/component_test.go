package app

type hello struct {
	Compo

	Greeting     string
	onNavURL     string
	appUpdated   bool
	appInstalled bool
	appResized   bool

	mounted     bool
	preRendered bool
}

func (h *hello) OnMount(Context) {
	h.mounted = true
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
	h.preRendered = true
	// ctx.Page().SetTitle("world")
}

func (h *hello) OnDismount() {
	h.mounted = false
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
	return If(f.Bar != "", func() UI {
		return &bar{Value: f.Bar}
	}).Else(func() UI {
		return Text("bar")
	})
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

func (b *bar) OnMount(ctx Context) {}

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

type compoWithNilRendering struct {
	Compo
	NilOverride UI
}

func (c *compoWithNilRendering) Render() UI {
	return nil
}

type compoWithNonMountableRoot struct {
	Compo
}

func (c *compoWithNonMountableRoot) Render() UI {
	return &compoWithNilRendering{}
}

type compoWithCustomRoot struct {
	Compo

	Root UI
}

func (c *compoWithCustomRoot) Render() UI {
	return c.Root
}

type updateNotifierCompo struct {
	Compo
	Root   UI
	notify bool
}

func (c *updateNotifierCompo) NotifyUpdate() bool {
	return c.notify
}

func (c *updateNotifierCompo) Render() UI {
	if c.Root != nil {
		return c.Root
	}
	return Span()
}

type navigatorComponent struct {
	Compo

	onNav func(Context)
}

func (c *navigatorComponent) OnNav(ctx Context) {
	if c.onNav != nil {
		c.onNav(ctx)
	}
}

type replacerComponent struct {
	Compo

	replace bool
}

func (c *replacerComponent) ReplaceOnUpdate() bool {
	return c.replace
}
