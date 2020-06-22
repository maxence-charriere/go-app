package app

import (
	"context"
	"net/url"
	"reflect"
	"strings"

	"github.com/maxence-charriere/go-app/v7/pkg/errors"
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
}

// Mounter is the interface that describes a component that can perform
// additional actions when mounted.
type Mounter interface {
	Composer

	// The function that is called when the component is mounted. It is always
	// called on the UI goroutine.
	OnMount(Context)
}

// Dismounter is the interface that describes a component that can perform
// additional actions when dismounted.
type Dismounter interface {
	Composer

	// The function that is called when the component is dismounted. It is
	// always called on the UI goroutine.
	OnDismount()
}

// Navigator is the interface that describes a component that can perform
// additional actions when navigated on.
type Navigator interface {
	Composer

	// The function that is called when the component is navigated on. It is
	// always called on the UI goroutine.
	OnNav(Context, *url.URL)
}

// Compo represents the base struct to use in order to build a component.
type Compo struct {
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
	return c.ctx != nil && c.ctx.Err() == nil &&
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

// Update triggers a component appearance update. It should be called when a
// field used to render the component has been modified. Updates are always
// performed on the UI goroutine.
func (c *Compo) Update() {
	dispatch(func() {
		if !c.Mounted() {
			return
		}

		if err := c.updateRoot(); err != nil {
			panic(err)
		}
	})
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

func (c *Compo) mount() error {
	if c.Mounted() {
		return errors.New("mounting component failed").
			Tag("reason", "already mounted").
			Tag("name", c.name()).
			Tag("kind", c.Kind())
	}

	c.ctx, c.ctxCancel = context.WithCancel(context.Background())

	root := c.render()
	if err := mount(root); err != nil {
		return errors.New("mounting component failed").
			Tag("name", c.name()).
			Tag("kind", c.Kind()).
			Wrap(err)
	}
	root.setParent(c.this)
	c.root = root

	if mounter, ok := c.this.(Mounter); ok {
		mounter.OnMount(Context{
			Context: c.ctx,
			Src:     c.this,
			JSSrc:   c.this.JSValue(),
		})
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

	if err := mount(new); err != nil {
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
	parent.JSValue().Call("replaceChild", newjs, oldjs)

	dismount(old)
	return nil
}

func (c *Compo) onNav(u *url.URL) {
	c.root.onNav(u)

	if nav, ok := c.self().(Navigator); ok {
		ctx := Context{
			Context: c.context(),
			Src:     c.self(),
			JSSrc:   c.JSValue(),
		}

		nav.OnNav(ctx, u)
	}
}

func (c *Compo) render() UI {
	elems := FilterUIElems(c.this.Render())
	return elems[0]
}
