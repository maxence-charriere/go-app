package app

import (
	"io"
	"reflect"
	"strings"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

// Composer is the interface that describes a customized, independent and
// reusable UI element.
//
// Satisfying this interface is done by embedding app.Compo into a struct and
// implementing the Render function.
//
// Example:
//
//	type Hello struct {
//	    app.Compo
//	}
//
//	func (c *Hello) Render() app.UI {
//	    return app.Text("hello")
//	}
type Composer interface {
	UI

	// Render returns the node tree that define how the component is desplayed.
	Render() UI

	// ValueTo stores the value of the DOM element (if exists) that emitted an
	// event into the given value.
	//
	// The given value must be a pointer to a signed integer, unsigned integer,
	// or a float.
	//
	// It panics if the given value is not a pointer.
	ValueTo(any) EventHandler

	setRef(Composer) Composer
	depth() uint
	setDepth(uint) Composer
	parent() UI
	root() UI
	setRoot(UI) Composer
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

// Initializer is the interface that describes a component that performs
// initialization instruction before being pre-rendered or mounted.
type Initializer interface {
	Composer

	// The function called before the component is pre-rendered or mounted.
	OnInit()
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

// Updater is the interface that describes a component that can do additional
// instructions when one of its exported fields is modified by its nearest
// parent component.
type Updater interface {
	// The function called when one of the component exported fields is modified
	// by its nearest parent component. It is always called on the UI goroutine.
	OnUpdate(Context)
}

// AppUpdater is the interface that describes a component that is notified when
// the application is updated.
type AppUpdater interface {
	// The function called when the application is updated. It is always called
	// on the UI goroutine.
	OnAppUpdate(Context)
}

// AppInstaller is the interface that describes a component that is notified
// when the application installation state changes.
type AppInstaller interface {
	// The function called when the application becomes installable or
	// installed. Use Context.IsAppInstallable() or Context.IsAppInstalled to
	// check the install state. OnAppInstallChange is always called on the UI
	// goroutine.
	OnAppInstallChange(Context)
}

// Resizer is the interface that describes a component that is notified when the
// app has been resized or a parent component calls the ResizeContent() method.
type Resizer interface {
	// The function called when the application is resized or a parent component
	// called its ResizeContent() method. It is always called on the UI
	// goroutine.
	OnResize(Context)
}

// UpdateNotifier defines a component that signals its parent component
// regarding the requirement for an update in response to an HTML event.
type UpdateNotifier interface {
	Composer

	// NotifyUpdate indicates whether the nearest parent component should be
	// queued for an update. It returns true to signal that the parent component
	// should be updated in the subsequent cycle,  and false to inhibit the
	// update.
	NotifyUpdate() bool
}

// Component events.
type nav struct{}
type appUpdate struct{}
type appInstallChange struct{}
type resize struct{}

// Compo represents the base struct to use in order to build a component.
type Compo struct {
	treeDepth     uint
	ref           Composer
	parentElement UI
	rootElement   UI
}

// JSValue returns the javascript value of the component root.
func (c *Compo) JSValue() Value {
	if c.rootElement == nil {
		return ValueOf(nil)
	}
	return c.rootElement.JSValue()
}

// Mounted reports whether the component is mounted.
func (c *Compo) Mounted() bool {
	return c.ref != nil
}

// Render describes the component content. This is a default implementation to
// satisfy the app.Composer interface. It should be redefined when app.Compo is
// embedded.
func (c *Compo) Render() UI {
	componentName := reflect.TypeOf(c.ref).Name()

	return Div().
		DataSet("compo-type", componentName).
		Style("border", "1px solid currentColor").
		Style("padding", "12px 0").
		Body(
			H1().Text("Component "+strings.TrimPrefix(componentName, "*")),
			P().Body(
				Text("Change appearance by implementing: "),
				Code().
					Style("color", "deepskyblue").
					Style("margin", "0 6px").
					Text("func (c "+componentName+") Render() app.UI"),
			),
		)
}

// ValueTo stores the value of the DOM element (if exists) that emitted an event
// into the given value.
//
// The given value must be a pointer to a signed integer, unsigned integer, or a
// float.
//
// It panics if the given value is not a pointer.
func (c *Compo) ValueTo(v any) EventHandler {
	return func(ctx Context, e Event) {
		value := ctx.JSSrc().Get("value")
		if err := stringTo(value.String(), v); err != nil {
			Log(errors.New("storing dom element value failed").Wrap(err))
			return
		}
	}
}

func (c *Compo) setRef(v Composer) Composer {
	c.ref = v
	return v
}

func (c *Compo) depth() uint {
	return c.treeDepth
}

func (c *Compo) setDepth(v uint) Composer {
	c.treeDepth = v
	return c.ref
}

func (c *Compo) parent() UI {
	return c.parentElement
}

func (c *Compo) setParent(p UI) UI {
	c.parentElement = p
	return c.ref
}

func (c *Compo) root() UI {
	return c.rootElement
}

func (c *Compo) setRoot(v UI) Composer {
	c.rootElement = v
	return c.ref
}

func (c *Compo) html(w io.Writer) {
	// if c.rootElement == nil {
	// 	c.rootElement = c.render()
	// 	c.rootElement.setSelf(c.rootElement)
	// }
	// c.rootElement.html(w)

	panic("not implemented")
}

func (c *Compo) htmlWithIndent(w io.Writer, indent int) {
	// if c.rootElement == nil {
	// 	c.rootElement = c.render()
	// 	c.rootElement.setSelf(c.rootElement)
	// }

	// c.rootElement.htmlWithIndent(w, indent)
	panic("not implemented")
}
