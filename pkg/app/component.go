package app

import (
	"reflect"
	"strings"

	"github.com/maxence-charriere/go-app/v10/pkg/errors"
)

// Composer defines a contract for creating custom, independent, and reusable
// UI elements. Components that satisfy the Composer interface serve as building
// blocks for richer user interfaces in a structured manner.
//
// Implementing the Composer interface typically involves embedding app.Compo
// into a struct and defining the Render method, which dictates the component's
// visual representation.
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

	// Render constructs and returns the visual representation of the component
	// as a node tree.
	Render() UI

	setRef(Composer) Composer
	depth() uint
	setDepth(uint) Composer
	parent() UI
	root() UI
	setRoot(UI) Composer
}

// Initializer describes a component that requires initialization
// instructions to be executed before it is mounted.
type Initializer interface {
	// OnInit is invoked before the component is mounted.
	OnInit()
}

// PreRenderer is the interface that describes a component that performs
// additional instructions during server-side rendering.
//
// Implementing OnPreRender within a component can enhance SEO by allowing
// server-side preparations for rendering.
type PreRenderer interface {
	// OnPreRender is called during the server-side rendering process of the
	// component.
	OnPreRender(Context)
}

// Mounter represents components that require initialization or setup actions
// when they are integrated into the DOM. By implementing the Mounter interface,
// components gain the ability to define specific behaviors that occur right
// after they are visually rendered or integrated into the DOM hierarchy.
type Mounter interface {
	// OnMount is triggered right after the component is embedded into the DOM.
	// Use this hook to perform any post-render configurations or
	// initializations.
	// This method operates within the UI goroutine.
	OnMount(Context)
}

// Dismounter outlines the behavior for components that require specific
// tasks or cleanup operations when they're detached from the DOM. Components
// adhering to this interface can designate procedures to run immediately
// after their removal from the DOM structure.
type Dismounter interface {
	// OnDismount is invoked immediately after the component is detached from
	// the DOM. This method offers a stage for executing cleanup or other
	// concluding operations.
	// This function is executed within the UI goroutine.
	OnDismount()
}

// Navigator characterizes components that need to perform specific
// actions or initializations when they become the target of navigation.
// By adopting the Navigator interface, components can specify behaviors
// to be executed when they are navigated to within the application.
type Navigator interface {
	// OnNav is invoked when the component becomes the navigation target.
	// Use this method to handle actions or setups related to navigation events.
	// This function is always executed within the UI goroutine.
	OnNav(Context)
}

// DismountEnforcer defines a contract for components that enforce
// a dismount operation when they undergo updates. Implementing this interface
// allows components to specify whether they should be removed and re-added
// upon updates instead of being modified in place.
type DismountEnforcer interface {
	// CompoID returns an identifier used to determine a component's identity.
	// When a component update occurs, a mismatch in IDs triggers the dismount
	// of the current component and the mounting of a new version.
	CompoID() string
}

// Updater encapsulates components that require specific behaviors or reactions
// when one of their exported fields is updated by the closest parent component.
// Implementing the Updater interface allows components to define responsive
// actions that should be executed whenever they are modified by a parent.
type Updater interface {
	// OnUpdate is triggered whenever an exported field of the component gets
	// modified by its immediate parent component. This method is an opportunity
	// to handle related reactions or recalculations.
	// This function always runs within the UI goroutine context.
	OnUpdate(Context)
}

// AppUpdater defines components that are alerted when a newer version of the
// application is downloaded in the background. Implementing this interface
// allows components to proactively adapt to app updates, ensuring coherence
// with the most up-to-date version of the application.
type AppUpdater interface {
	// OnAppUpdate is called once a new version of the application has been
	// fetched in the background. It offers a window for components to execute
	// actions, such as prompting a page reload, to transition to the updated
	// app version.
	// This function always operates within the UI goroutine context.
	OnAppUpdate(Context)
}

// AppInstaller outlines components that receive notifications about changes in
// the application's installation status. Through this interface, components can
// actively respond to installation state transitions, facilitating dynamic
// user experiences tailored to the app's current status.
type AppInstaller interface {
	// OnAppInstallChange is invoked when the application shifts between the
	// states of being installable and actually installed.
	//
	// To determine the current installation state, one can use
	// Context.IsAppInstallable() or Context.IsAppInstalled().
	//
	// By leveraging this method, components can maintain alignment with the
	// app's installation status, potentially influencing UI elements like an
	// "Install" button visibility or behavior.
	// This method is always executed in the UI goroutine context.
	OnAppInstallChange(Context)
}

// Resizer identifies components that respond to size alterations within the
// application. These components can dynamically adjust to diverse size
// scenarios, ensuring they maintain both a visually appealing and functional
// display.
type Resizer interface {
	// OnResize is called whenever the application experiences a change in size.
	// Components can use this method to make appropriate adjustments,
	// recalculations, or layout shifts in response to the modified dimensions.
	// Note: This method operates exclusively within the UI goroutine context.
	OnResize(Context)
}

// Compo serves as the foundational struct for constructing a component. It
// provides basic methods and fields needed for component management.
type Compo struct {
	treeDepth     uint
	ref           Composer
	parentElement UI
	rootElement   UI
}

// JSValue retrieves the JavaScript value associated with the component's root.
// If the root element isn't defined, it returns a nil JavaScript value.
func (c *Compo) JSValue() Value {
	if c.rootElement == nil {
		return ValueOf(nil)
	}
	return c.rootElement.JSValue()
}

// Mounted checks if the component is currently mounted within the UI.
func (c *Compo) Mounted() bool {
	return c.ref != nil
}

// Render produces a visual representation of the component's content. This
// default implementation ensures the app.Composer interface is satisfied
// when app.Compo is embedded. However, developers are encouraged to redefine
// this method to customize the component's appearance.
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

// ValueTo captures the value of the DOM element (if it exists) that triggered
// an event, and assigns it to the provided receiver. The receiver must be a
// pointer pointing to either a string, integer, unsigned integer, or a float.
// This method panics if the provided value isn't a pointer.
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
