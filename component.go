package app

import (
	"net/url"
	"strings"
)

// Component is the interface that describes a component.
// Should be implemented on a non empty struct pointer.
type Component interface {
	// Render should return a string describing the component with HTML5
	// standard.
	// It supports standard Go html/template API.
	// Pipeline is based on the component struct.
	// See https://golang.org/pkg/text/template and
	// https://golang.org/pkg/html/template for template usage.
	Render() string
}

// Mounter is the interface that wraps OnMount method.
type Mounter interface {
	Component

	// OnMount is called when a component is mounted.
	// App.Render should not be called inside.
	OnMount()
}

// Dismounter is the interface that wraps OnDismount method.
type Dismounter interface {
	Component

	// OnDismount is called when a component is dismounted.
	// App.Render should not be called inside.
	OnDismount()
}

// Navigable is the interface that wraps OnNavigate method.
type Navigable interface {
	Component

	// OnNavigate is called when a component is loaded or navigated to.
	// It is called just after the component is mounted.
	OnNavigate(u *url.URL)
}

// Subscriber is the interface that describes a component that subscribes to
// events generated from actions.
type Subscriber interface {
	// Subscribe is called after a component is mounted.
	// The returned event subscriber is used to subscribe to events generated
	// from actions.
	// All the event subscribed are automatically unsuscribed when the component
	// is dismounted.
	Subscribe() EventSubscriber
}

// ComponentWithExtendedRender is the interface that wraps Funcs method.
type ComponentWithExtendedRender interface {
	Component

	// Funcs returns a map of funcs to use when rendering a component.
	// Funcs named raw, json and time are reserved.
	// They handle raw html code, json conversions and time format.
	// They can't be overloaded.
	// See https://golang.org/pkg/text/template/#Template.Funcs for more details.
	Funcs() map[string]interface{}
}

// ZeroCompo is the type to use as base for an empty compone.
// Every instances of an empty struct is given the same memory address, which
// causes problem for indexing components.
// ZeroCompo have a placeholder field to avoid that.
type ZeroCompo struct {
	placeholder byte
}

// ComponentNameFromURL is a helper function that returns the component name
// targeted by the given URL.
func ComponentNameFromURL(u *url.URL) string {
	if len(u.Scheme) != 0 && u.Scheme != "compo" {
		return ""
	}

	path := u.Path
	path = strings.TrimPrefix(path, "/")

	paths := strings.SplitN(path, "/", 2)
	if len(paths[0]) == 0 {
		return ""
	}
	return normalizeComponentName(paths[0])
}

// ComponentNameFromURLString is a helper function that returns the component
// name targeted by the given URL.
func ComponentNameFromURLString(rawurl string) string {
	u, _ := url.Parse(rawurl)
	return ComponentNameFromURL(u)
}

func normalizeComponentName(name string) string {
	name = strings.ToLower(name)
	if pkgsep := strings.IndexByte(name, '.'); pkgsep != -1 {
		if name[:pkgsep] == "main" {
			name = name[pkgsep+1:]
		}
	}
	return name
}
