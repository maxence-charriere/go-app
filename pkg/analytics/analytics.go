// Package analytics provides functions to send analytics to analytics stores
// such as Google Analytics.
package analytics

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// Backend is the interface that describes an analytics backend that sends
// analytics for a defined provider.
type Backend interface {
	// Links your users, and their actions, to a recognizable userID and traits.
	Identify(userID string, traits map[string]interface{})

	// Record actions your users perform.
	Track(event string, properties map[string]interface{})

	// Records page views on your website, along with optional extra information
	// about the page viewed by the user.
	Page(name string, properties map[string]interface{})
}

// Identify links your users, and their actions, to a recognizable userID and
// traits.
func Identify(userID string, traits map[string]interface{}) {
	if app.IsServer || userID == "" {
		return
	}

	traits = sanitizeValue(traits).(map[string]interface{})

	for _, b := range backends {
		b.Identify(userID, traits)
	}
}

// Track record actions your users perform.
func Track(event string, properties map[string]interface{}) {
	if app.IsServer || event == "" {
		return
	}

	properties = sanitizeValue(properties).(map[string]interface{})

	for _, b := range backends {
		b.Track(event, properties)
	}
}

// Page records page views on your website, along with optional extra
// information about the page viewed by the user.
//
// The following properties are automatically set: path, referrer, search, title
// and url.
func Page(name string, properties map[string]interface{}) {
	if app.IsServer {
		return
	}

	if properties == nil {
		properties = make(map[string]interface{}, 5)
	}

	properties["path"] = app.Window().Get("location").Get("pathname").String()
	properties["referrer"] = app.Window().Get("document").Get("referrer").String()
	properties["search"] = app.Window().Get("location").Get("search").String()
	properties["title"] = app.Window().Get("document").Get("title").String()
	properties["url"] = app.Window().Get("location").Get("href").String()

	if name == "" {
		name = app.Window().Get("document").Get("title").String()
	}

	properties = sanitizeValue(properties).(map[string]interface{})

	for _, b := range backends {
		b.Page(name, properties)
	}
}

// Add adds the given backend to the backends used to send analytics.
func Add(b Backend) {
	backends = append(backends, b)
}

var (
	backends []Backend
)

func sanitizeValue(v interface{}) interface{} {
	switch v := v.(type) {
	case app.Value,
		app.Func,
		nil,
		string,
		bool,
		int,
		int8,
		int16,
		int32,
		int64,
		uint,
		uint8,
		uint16,
		uint32,
		uint64,
		float32,
		float64:
		return v

	case []interface{}:
		for i, item := range v {
			v[i] = sanitizeValue(item)
		}
		return v

	case map[string]interface{}:
		for k, val := range v {
			v[k] = sanitizeValue(val)
		}
		return v

	default:
		return fmt.Sprintf("%v", v)
	}
}
