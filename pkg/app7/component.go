package app

import (
	"context"
	"net/url"
	"reflect"
	"strings"

	"github.com/maxence-charriere/go-app/v6/pkg/errors"
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

	// The function that is called when the component is mounted.
	OnMount(Context)
}

// Dismounter is the interface that describes a component that can perform
// additional actions when dismounted.
type Dismounter interface {
	Composer

	// The function that is called when the component is dismounted.
	OnDismount()
}

// Navigator is the interface that describes a component that can perform
// additional actions when navigated on.
type Navigator interface {
	Composer

	// The function that is called when the component is navigated on.
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

	root := c.this.Render()
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

func (c *Compo) update(UI) error {
	panic("not implemented")
}

// Update triggers a component appearance update. It should be called when a
// field used to render the component has been modified.
func (c *Compo) Update() {
	panic("not implemented")
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
