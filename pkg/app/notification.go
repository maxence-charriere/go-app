package app

import (
	"encoding/json"

	"github.com/maxence-charriere/go-app/v10/pkg/errors"
)

// Notification represents a user notification.
type Notification struct {
	// The title displayed prominently at the top of the notification.
	Title string `json:"title"`

	// Path is the URL to navigate to upon clicking the notification.
	Path string `json:"path"`

	// Lang specifies the notification's language, using ISO 2-letter codes.
	// See: https://www.sitepoint.com/iso-2-letter-language-codes
	Lang string `json:"lang,omitempty"`

	// Badge is the URL of an image to represent the notification when space
	// is limited.
	Badge string `json:"badge,omitempty"`

	// Body is the main content, displayed below the title.
	Body string `json:"body,omitempty"`

	// Tag provides a unique identifier for the notification.
	Tag string `json:"tag,omitempty"`

	// Icon is the URL of an icon shown in the notification.
	Icon string `json:"icon,omitempty"`

	// Image is the URL of an image displayed in the notification.
	Image string `json:"image,omitempty"`

	// Data contains arbitrary data associated with the notification.
	// Note: "goapp" key is reserved for go-app data.
	Data map[string]any `json:"data"`

	// Renotify determines if the user should be notified when a new
	// notification replaces an existing one.
	Renotify bool `json:"renotify,omitempty"`

	// RequireInteraction ensures the notification remains active until
	// the user acts on it, not closing automatically.
	RequireInteraction bool `json:"requireInteraction,omitempty"`

	// Silent specifies if the notification should be silent, suppressing
	// all sounds and vibrations, regardless of device settings.
	Silent bool `json:"silent,omitempty"`

	// Vibrate defines a vibration pattern for the device upon notification.
	// See: https://developer.mozilla.org/en-US/docs/Web/API/Vibration_API
	Vibrate []int `json:"vibrate,omitempty"`

	// Actions lists the available actions displayed within the notification.
	Actions []NotificationAction `json:"actions,omitempty"`
}

// NotificationAction represents an actionable item within a notification.
type NotificationAction struct {
	// Action is the unique ID associated with the user's action on
	// the notification.
	Action string `json:"action"`

	// Title is the descriptive text shown alongside the action.
	Title string `json:"title"`

	// Icon is the URL of an image to represent the action.
	Icon string `json:"icon,omitempty"`

	// Path is the URL to navigate to upon clicking the action.
	Path string `json:"path"`
}

// NotificationSubscription encapsulates a PushSubscription from the Push API.
type NotificationSubscription struct {
	// Endpoint is the push service endpoint URL.
	Endpoint string `json:"endpoint"`

	// Keys contains cryptographic keys used for the subscription.
	Keys struct {
		// Auth is the authentication secret.
		Auth string `json:"auth"`

		// P256dh is the user's public key, associated with the push
		// subscription.
		P256dh string `json:"p256dh"`
	} `json:"keys"`
}

// NotificationPermission represents permission levels for displaying
// notifications to users.
type NotificationPermission string

const (
	// NotificationDefault indicates the user's choice for notifications is
	// unknown, prompting the browser to treat it as "denied".
	NotificationDefault NotificationPermission = "default"

	// NotificationGranted means the user has allowed notifications.
	NotificationGranted NotificationPermission = "granted"

	// NotificationDenied means the user has declined notifications.
	NotificationDenied NotificationPermission = "denied"

	// NotificationNotSupported indicates that the browser doesn't support
	// notifications.
	NotificationNotSupported NotificationPermission = "unsupported"
)

// NotificationService provides functionalities related to user notifications.
type NotificationService struct{}

// Permission retrieves the current notification permission status.
func (s NotificationService) Permission() NotificationPermission {
	notification := Window().Get("Notification")
	if !notification.Truthy() {
		return NotificationNotSupported
	}

	return NotificationPermission(notification.Get("permission").String())
}

// RequestPermission prompts the user for permission to display notifications.
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

// New creates and displays a notification to the user.
func (s NotificationService) New(n Notification) {
	notification, _ := json.Marshal(n)
	Window().Call("goappNewNotification", string(notification))
}

// Subscribe retrieves a notification subscription using the provided VAPID
// public key. If the key is empty or push notifications aren't supported, an
// error is returned.
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
			New("failed to decode push notification subscription").Wrap(err)
	}
	return sub, nil
}
