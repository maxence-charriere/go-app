package app

import (
	"encoding/json"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

// A user notification.
type Notification struct {
	// The title shown at the top of the notification window.
	Title string `json:"title"`

	// The URL path to navigate to when the notification is clicked.
	Path string `json:"path"`

	// The notification's language, as specified in
	// https://www.sitepoint.com/iso-2-letter-language-codes.
	Lang string `json:"lang,omitempty"`

	// The URL of the image used to represent the notification when there isn't
	// enough space to display the notification itself.
	Badge string `json:"badge,omitempty"`

	// The body text of the notification, which is displayed below the title.
	Body string `json:"body,omitempty"`

	// An identifying tag for the notification.
	Tag string `json:"tag,omitempty"`

	// The URL of an icon to be displayed in the notification.
	Icon string `json:"icon,omitempty"`

	// The URL of an image to be displayed in the notification.
	Image string `json:"image,omitempty"`

	// Arbitrary data that to be associated with the notification.
	//
	// The "goapp" key is reserved to go-app data.
	Data map[string]any `json:"data"`

	// Specifies whether the user should be notified after a new notification
	// replaces an old one.
	Renotify bool `json:"renotify,omitempty"`

	// Indicates whether a notification should remain active until the user
	// clicks or dismisses it, rather than closing automatically.
	RequireInteraction bool `json:"requireInteraction,omitempty"`

	// specifies whether the notification is silent (no sounds or vibrations
	// issued), regardless of the device settings.
	Silent bool `json:"silent,omitempty"`

	// A vibration pattern for the device's vibration hardware to emit with the
	// notification.
	//
	// See https://developer.mozilla.org/en-US/docs/Web/API/Vibration_API#vibration_patterns.
	Vibrate []int `json:"vibrate,omitempty"`

	// The actions to display in the notification.
	Actions []NotificationAction `json:"actions,omitempty"`
}

// A notification action.
type NotificationAction struct {
	// The user action id to be displayed on the notification.
	Action string `json:"action"`

	// The action text to be shown to the user.
	Title string `json:"title"`

	// The URL of an icon to display with the action.
	Icon string `json:"icon,omitempty"`

	// The URL path to navigate to when the action is clicked.
	Path string `json:"path"`
}

// NotificationSubscription represents a PushSubscription object from the Push
// API.
type NotificationSubscription struct {
	Endpoint string `json:"endpoint"`
	Keys     struct {
		Auth   string `json:"auth"`
		P256dh string `json:"p256dh"`
	} `json:"keys"`
}

// NotificationPermission a permission to display notifications.
type NotificationPermission string

const (
	// The user notifications choice is unknown and therefore the browser acts
	// as if the value were denied.
	NotificationDefault NotificationPermission = "default"

	// The user accepts having notifications displayed.
	NotificationGranted NotificationPermission = "granted"

	// The user refuses to have notifications displayed.
	NotificationDenied NotificationPermission = "denied"

	// Notifications are not supported by the browser.
	NotificationNotSupported NotificationPermission = "unsupported"
)

type NotificationService struct {
	dispatcher Dispatcher
}

// Returns the current notification permission.
func (s NotificationService) Permission() NotificationPermission {
	notification := Window().Get("Notification")
	if !notification.Truthy() {
		return NotificationNotSupported
	}

	return NotificationPermission(notification.Get("permission").String())
}

// Requests the user whether the app can use notifications.
func (s NotificationService) RequestPermission() NotificationPermission {
	notification := Window().Get("Notification")
	if !notification.Truthy() {
		return NotificationNotSupported
	}

	permission := make(chan string, 1)
	defer close(permission)

	notification.Call("requestPermission").Then(func(v Value) {
		permission <- v.String()
	})

	return NotificationPermission(<-permission)
}

// Creates and display a user notification.
func (s NotificationService) New(n Notification) {
	notification, _ := json.Marshal(n)
	Window().Call("goappNewNotification", string(notification))
}

// Returns a notification subscription with the given vap id.
func (s NotificationService) Subscribe(vapIDPublicKey string) (NotificationSubscription, error) {
	if vapIDPublicKey == "" {
		return NotificationSubscription{}, errors.New("vapid public key is empty")
	}

	subc := make(chan string, 1)
	defer close(subc)

	Window().Call("goappSubscribePushNotifications", vapIDPublicKey).Then(func(v Value) {
		subc <- v.String()
	})

	jsSub := <-subc
	if jsSub == "" {
		return NotificationSubscription{}, errors.
			New("push notifications are not supported by the browser")
	}

	var sub NotificationSubscription
	err := json.Unmarshal([]byte(jsSub), &sub)
	if err != nil {
		return NotificationSubscription{}, errors.
			New("decoding push notification subscription failed").
			Wrap(err)
	}
	return sub, nil
}
