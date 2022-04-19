package app

// A user notification.
type Notification struct {
	// The title shown at the top of the notification window.
	Title string `json:"title"`

	// The URL to navigate to when the notification is clicked.
	Target string `json:"target"`

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
	Data interface{} `json:"data,omitempty"`

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
	Actions []NotificationAction
}

// A notification action.
type NotificationAction struct {
	// The user action id to be displayed on the notification.
	Action string `json:"action"`

	// The action text to be shown to the user.
	Title string `json:"title"`

	// The URL of an icon to display with the action.
	Icon string `json:"icon,omitempty"`

	// The URL to navigate to when the action is clicked.
	Target string `json:"target"`
}

// The configuration to subscribe to and receive push notifications.
type PushNotificationsConfig struct {
	// The public VAPID key.
	VAPIDPublicKey string

	// The URL where push notification subscriptions are sent.
	RegistrationURL string
}
