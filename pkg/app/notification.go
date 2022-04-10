package app

// A user notification.
type Notification struct {
	// The title shown at the top of the notification window.
	Title string

	// The notification's language, as specified in
	// https://www.sitepoint.com/iso-2-letter-language-codes.
	Lang string

	// The URL of the image used to represent the notification when there isn't
	// enough space to display the notification itself.
	Badge string

	// The body text of the notification, which is displayed below the title.
	Body string

	// An identifying tag for the notification.
	Tag string

	// The URL of an icon to be displayed in the notification.
	Icon string

	// The URL of an image to be displayed in the notification.
	Image string

	// Arbitrary data that to be associated with the notification.
	Data Value

	// Specifies whether the user should be notified after a new notification
	// replaces an old one.
	Renotify bool

	// Indicates whether a notification should remain active until the user
	// clicks or dismisses it, rather than closing automatically.
	RequireInteraction bool

	// specifies whether the notification is silent (no sounds or vibrations
	// issued), regardless of the device settings.
	Silent bool

	// A vibration pattern for the device's vibration hardware to emit with the
	// notification.
	//
	// See https://developer.mozilla.org/en-US/docs/Web/API/Vibration_API#vibration_patterns.
	Vibrate []int

	// The actions to display in the notification.
	Actions []NotificationAction

	// The function called when a notification is clicked by the user.
	OnClick EventHandler

	// The function called when a notification is closed.
	OnClose EventHandler

	// The function called when something goes wrong with a notification (in
	// many cases an error preventing the notification from being displayed).
	OnError EventHandler

	// The function called when a notification is displayed.
	OnShow EventHandler
}

// A notification action.
type NotificationAction struct {
	// The user action id to be displayed on the notification.
	Action string

	// The action text to be shown to the user.
	Title string

	// The URL of an icon to display with the action.
	Icon string
}
