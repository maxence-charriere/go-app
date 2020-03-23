package app

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/maxence-charriere/go-app/v6/pkg/log"
)

// Composer is the interface that describes a component that embeds other nodes.
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
	nodeWithChildren

	// Render returns the node tree that define how the component is desplayed.
	Render() UI

	// Update update the component appearance. It should be called when a field
	// used to render the component has been modified.
	Update()

	setCompo(n Composer)
	mount(c Composer) error
	update(n Composer)
}

// Mounter is the interface that describes a component that can perform
// additional actions when mounted.
type Mounter interface {
	Composer

	// The function that is called when the component is mounted.
	OnMount()
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
	OnNav(u *url.URL)
}

// Compo represents the base struct to use in order to build a component.
type Compo struct {
	compo      Composer
	parentNode UI
	root       UI
}

func (c *Compo) nodeType() reflect.Type {
	return reflect.TypeOf(c.compo)
}

// JSValue returns the component root javascript value.
func (c *Compo) JSValue() Value {
	return c.root.JSValue()
}

func (c *Compo) parent() UI {
	return c.parentNode
}

func (c *Compo) setParent(p UI) {
	c.parentNode = p
}

func (c *Compo) setCompo(n Composer) {
	c.compo = n
}

func (c *Compo) dismount() {
	c.root.dismount()

	if dismounter, ok := c.root.(Dismounter); ok {
		dismounter.OnDismount()
	}
}

func (c *Compo) replaceChild(old, new UI) {
	if old == c.root {
		c.root = new
	}
}

// Update update the component appearance. It should be called when a field
// used to render the component has been modified.
func (c *Compo) Update() {
	dispatcher(func() {
		if c.compo == nil {
			return
		}

		current := c.root
		incoming := c.compo.Render().(UI)

		if err := update(current, incoming); err != nil {
			log.Error("updating component failed").
				T("component-type", reflect.TypeOf(c.compo)).
				T("error", err).
				Panic()
		}
	})
}

func (c *Compo) mount(compo Composer) error {
	c.setCompo(compo)

	root := compo.Render()
	if err := mount(root); err != nil {
		return fmt.Errorf("%T: invalid root: %w", compo, err)
	}
	c.root = root

	if mounter, ok := compo.(Mounter); ok {
		mounter.OnMount()
	}

	return nil
}

func (c *Compo) update(n Composer) {
	aval := reflect.Indirect(reflect.ValueOf(c.compo))
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
}

// Render describes the component content. This is a default implementation to
// satisfy the app.Composer interface. It should be redefined when app.Compo is
// embedded.
func (c *Compo) Render() UI {
	compoType := reflect.TypeOf(c)
	if c.compo != nil {
		compoType = reflect.TypeOf(c.compo)
	}
	compoName := compoType.String()
	compoName = strings.ReplaceAll(compoName, "main.", "")

	return Div().
		DataSet("compo-type", compoType).
		Style("border", "1px solid currentColor").
		Style("padding", "12px 0").
		Body(
			H1().Body(
				Text("Component "+strings.TrimPrefix(compoName, "*")),
			),
			P().Body(
				Text("Change appearance by implementing: "),
				Code().
					Style("color", "deepskyblue").
					Style("margin", "0 6px").
					Body(
						Text("func (c "+compoName+") Render() app.UI"),
					),
			),
		)
}
