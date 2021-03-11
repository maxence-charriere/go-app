package app

import (
	"context"
	"io"
	"net/url"
	"reflect"
	"strings"

	"github.com/maxence-charriere/go-app/v8/pkg/errors"
)

// Composer is the interface that describes a customized, independent and
// reusable UI element.
//
// Satisfying this interface is done by embedding app.Compo into a struct and
// implementing the Render function.
//
// Example:
//  type Hello struct {
//      app.Compo
//  }
//
//  func (c *Hello) Render() app.UI {
//      return app.Text("hello")
//  }
type Composer interface {
	UI

	// Render returns the node tree that define how the component is desplayed.
	Render() UI

	// Update update the component appearance. It should be called when a field
	// used to render the component has been modified.
	Update()

	// Dispatch executes the given function on the goroutine dedicated to
	// updating the UI and ensures that the component is mounted when the
	// function is called.
	Defer(func(Context))

	// ResizeContent triggers OnResize() on all the component children that
	// implement the Resizer interface.
	ResizeContent()

	// ValueTo stores the value of the DOM element (if exists) that emitted an
	// event into the given value.
	//
	// The given value must be a pointer to a signed integer, unsigned integer,
	// or a float.
	//
	// It panics if the given value is not a pointer.
	ValueTo(interface{}) EventHandler
}

// PreRenderer is the interface that describes a component that performs
// instruction when it is server-side pre-rendered.
//
// A pre-rendered component helps in achieving SEO friendly content.
type PreRenderer interface {
	// The function called when the component is server-side pre-rendered.
	//
	// If pre-rendering requires blocking operations such as performing an HTTP
	// request, ensure that they are done synchronously. A good practice is to
	// avoid using goroutines during pre-rendering.
	OnPreRender(Context)
}

// Mounter is the interface that describes a component that can perform
// additional actions when mounted.
type Mounter interface {
	Composer

	// The function called when the component is mounted. It is always called on
	// the UI goroutine.
	OnMount(Context)
}

// Dismounter is the interface that describes a component that can perform
// additional actions when dismounted.
type Dismounter interface {
	Composer

	// The function called when the component is dismounted. It is always called
	// on the UI goroutine.
	OnDismount()
}

// Navigator is the interface that describes a component that can perform
// additional actions when navigated on.
type Navigator interface {
	Composer

	// The function that called when the component is navigated on. It is always
	// called on the UI goroutine.
	OnNav(Context)
}

type deprecatedNavigator interface {
	OnNav(Context, *url.URL)
}

// Updater is the interface that describes a component that is notified when the
// application is updated.
type Updater interface {
	// The function called when the application is updated. It is always called
	// on the UI goroutine.
	OnAppUpdate(Context)
}

// Resizer is the interface that describes a component that is notified when the
// app has been resized or a parent component calls the ResizeContent() method.
type Resizer interface {
	// The function called when the application is resized or a parent component
	// called its ResizeContent() method. It is always called on the UI
	// goroutine.
	OnResize(Context)
}

type deprecatedResizer interface {
	OnAppResize(Context)
}

// Compo represents the base struct to use in order to build a component.
type Compo struct {
	disp       Dispatcher
	ctx        context.Context
	ctxCancel  func()
	parentElem UI
	root       UI
	this       Composer
}

// Kind returns the ui element kind.
func (c *Compo) Kind() Kind {
	return Component
}

// JSValue returns the javascript value of the component root.
func (c *Compo) JSValue() Value {
	return c.root.JSValue()
}

// Mounted reports whether the component is mounted.
func (c *Compo) Mounted() bool {
	return c.dispatcher() != nil &&
		c.ctx != nil &&
		c.ctx.Err() == nil &&
		c.root != nil && c.root.Mounted() &&
		c.self() != nil
}

// Render describes the component content. This is a default implementation to
// satisfy the app.Composer interface. It should be redefined when app.Compo is
// embedded.
func (c *Compo) Render() UI {
	return Div().
		DataSet("compo-type", c.name()).
		Style("border", "1px solid currentColor").
		Style("padding", "12px 0").
		Body(
			H1().Text("Component "+strings.TrimPrefix(c.name(), "*")),
			P().Body(
				Text("Change appearance by implementing: "),
				Code().
					Style("color", "deepskyblue").
					Style("margin", "0 6px").
					Text("func (c "+c.name()+") Render() app.UI"),
			),
		)
}

// Defer executes the given function on the goroutine dedicated to
// updating the UI and ensures that the component is mounted when the
// function is called.
func (c *Compo) Defer(fn func(Context)) {
	c.dispatcher().Dispatch(func() {
		if !c.Mounted() {
			return
		}
		fn(makeContext(c.self()))
	})
}

// Update triggers a component appearance update. It should be called when a
// field used to render the component has been modified. Updates are always
// performed on the UI goroutine.
func (c *Compo) Update() {
	c.Defer(func(Context) {
		if err := c.updateRoot(); err != nil {
			panic(err)
		}
	})
}

// ResizeContent triggers OnResize() on all the component children that
// implement the Resizer interface.
func (c *Compo) ResizeContent() {
	c.Defer(func(Context) {
		c.root.onResize()
	})
}

// ValueTo stores the value of the DOM element (if exists) that emitted an event
// into the given value.
//
// The given value must be a pointer to a signed integer, unsigned integer, or a
// float.
//
// It panics if the given value is not a pointer.
func (c *Compo) ValueTo(v interface{}) EventHandler {
	return func(ctx Context, e Event) {
		value := ctx.JSSrc.Get("value")
		if !value.Truthy() {
			return
		}
		if err := stringTo(value.String(), v); err != nil {
			Log(errors.New("storing dom element value failed").Wrap(err))
			return
		}
		c.Update()
	}
}

func (c *Compo) name() string {
	name := reflect.TypeOf(c.self()).String()
	name = strings.ReplaceAll(name, "main.", "")
	return name
}

func (c *Compo) self() UI {
	return c.this
}

func (c *Compo) setSelf(n UI) {
	if n != nil {
		c.this = n.(Composer)
		return
	}

	c.this = nil
}

func (c *Compo) context() context.Context {
	return c.ctx
}

func (c *Compo) dispatcher() Dispatcher {
	return c.disp
}

func (c *Compo) attributes() map[string]string {
	return nil
}

func (c *Compo) eventHandlers() map[string]eventHandler {
	return nil
}

func (c *Compo) parent() UI {
	return c.parentElem
}

func (c *Compo) setParent(p UI) {
	c.parentElem = p
}

func (c *Compo) children() []UI {
	return []UI{c.root}
}

func (c *Compo) mount(d Dispatcher) error {
	if c.Mounted() {
		return errors.New("mounting component failed").
			Tag("reason", "already mounted").
			Tag("name", c.name()).
			Tag("kind", c.Kind())
	}

	c.disp = d
	c.ctx, c.ctxCancel = context.WithCancel(context.Background())

	root := c.render()
	if err := mount(d, root); err != nil {
		return errors.New("mounting component failed").
			Tag("name", c.name()).
			Tag("kind", c.Kind()).
			Wrap(err)
	}
	root.setParent(c.this)
	c.root = root

	if mounter, ok := c.self().(Mounter); ok && !c.dispatcher().isServerSideMode() {
		mounter.OnMount(makeContext(c.self()))
	}

	return nil
}

func (c *Compo) dismount() {
	dismount(c.root)
	c.ctxCancel()

	if dismounter, ok := c.this.(Dismounter); ok {
		dismounter.OnDismount()
	}
}

func (c *Compo) update(n UI) error {
	if c.self() == n || !c.Mounted() {
		return nil
	}

	if n.Kind() != c.Kind() || n.name() != c.name() {
		return errors.New("updating ui element failed").
			Tag("replace", true).
			Tag("reason", "different element types").
			Tag("current-kind", c.Kind()).
			Tag("current-name", c.name()).
			Tag("updated-kind", n.Kind()).
			Tag("updated-name", n.name())
	}

	aval := reflect.Indirect(reflect.ValueOf(c.self()))
	bval := reflect.Indirect(reflect.ValueOf(n))
	compotype := reflect.ValueOf(c).Elem().Type()

	for i := 0; i < aval.NumField(); i++ {
		a := aval.Field(i)
		b := bval.Field(i)

		if a.Type() == compotype {
			continue
		}

		if !a.CanSet() {
			continue
		}

		if !reflect.DeepEqual(a.Interface(), b.Interface()) {
			a.Set(b)
		}
	}

	return c.updateRoot()
}

func (c *Compo) updateRoot() error {
	a := c.root
	b := c.render()

	err := update(a, b)
	if isErrReplace(err) {
		err = c.replaceRoot(b)
	}

	if err != nil {
		return errors.New("updating component failed").
			Tag("kind", c.Kind()).
			Tag("name", c.name()).
			Wrap(err)
	}

	return nil
}

func (c *Compo) replaceRoot(n UI) error {
	old := c.root
	new := n

	if err := mount(c.dispatcher(), new); err != nil {
		return errors.New("replacing component root failed").
			Tag("kind", c.Kind()).
			Tag("name", c.name()).
			Tag("root-kind", old.Kind()).
			Tag("root-name", old.name()).
			Tag("new-root-kind", new.Kind()).
			Tag("new-root-name", new.name()).
			Wrap(err)
	}

	var parent UI
	for {
		parent = c.parent()
		if parent == nil || parent.Kind() == HTML {
			break
		}
	}

	if parent == nil {
		return errors.New("replacing component root failed").
			Tag("kind", c.Kind()).
			Tag("name", c.name()).
			Tag("reason", "coponent does not have html element parents")
	}

	c.root = new
	new.setParent(c.self())

	oldjs := old.JSValue()
	newjs := n.JSValue()
	parent.JSValue().replaceChild(newjs, oldjs)

	dismount(old)
	return nil
}

func (c *Compo) render() UI {
	elems := FilterUIElems(c.this.Render())
	return elems[0]
}

func (c *Compo) onNav(u *url.URL) {
	c.root.onNav(u)

	if nav, ok := c.self().(Navigator); ok {
		ctx := makeContext(c.self())
		nav.OnNav(ctx)
		return
	}

	if nav, ok := c.self().(deprecatedNavigator); ok {
		ctx := makeContext(c.self())
		Log(errors.New("a deprecated component interface is in use").
			Tag("component", reflect.TypeOf(c.self())).
			Tag("interface", "app.Navigator").
			Tag("method-current", "OnNav(app.Context, *url.URL)").
			Tag("method-fix", "OnNav(app.Context)").
			Tag("how-to-fix", "refactor component to use the right method"))
		nav.OnNav(ctx, u)
	}
}

func (c *Compo) onAppUpdate() {
	c.root.onAppUpdate()

	if updater, ok := c.self().(Updater); ok {
		updater.OnAppUpdate(makeContext(c.self()))
	}
}

func (c *Compo) onResize() {
	c.root.onResize()

	if resizer, ok := c.self().(Resizer); ok {
		resizer.OnResize(makeContext(c.self()))
		return
	}

	if resizer, ok := c.self().(deprecatedResizer); ok {
		Log(errors.New("a deprecated component interface is in use").
			Tag("component", reflect.TypeOf(c.self())).
			Tag("interface", "app.Resizer").
			Tag("method-current", "OnAppResize(app.Context)").
			Tag("method-fix", "OnResize(app.Context)").
			Tag("how-to-fix", "refactor component to use the right method"))
		resizer.OnAppResize(makeContext(c.self()))
	}
}

func (c *Compo) preRender(p Page) {
	c.root.preRender(p)

	if preRenderer, ok := c.self().(PreRenderer); ok {
		preRenderer.OnPreRender(makeContext(c.self()))
	}
}

func (c *Compo) html(w io.Writer) {
	if c.root == nil {
		c.root = c.render()
		c.root.setSelf(c.root)
	}
	c.root.html(w)
}

func (c *Compo) htmlWithIndent(w io.Writer, indent int) {
	if c.root == nil {
		c.root = c.render()
		c.root.setSelf(c.root)
	}

	c.root.htmlWithIndent(w, indent)
}
