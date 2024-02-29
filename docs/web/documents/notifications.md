## Enable Notifications

Enabling notifications requires the user to give permission to display them.

### Current Permission

The current notifications permission is retrieved by calling [Context.Notifications().Permission()](/reference#NotificationService.Permission):

```go
type foo struct {
	app.Compo

	notificationPermission app.NotificationPermission
}

func (f *foo) OnMount(ctx app.Context) {
	f.notificationPermission = ctx.Notifications().Permission()
}
```

### Request Permission

The Notification permission is given by requesting the user permission with [Context.Notifications().RequestPermission()](/reference#NotificationService.RequestPermission):

```go
func (f *foo) Render() app.UI {
	return app.Div().Body(
		app.If(f.notificationPermission == app.NotificationDefault, func() app.UI {
			return app.Button().
				Text("Enable Notifications").
				OnClick(f.enableNotifications)
		}).ElseIf(f.notificationPermission == app.NotificationDenied, func() app.UI {
			return app.Text("Notification permission is denied")
		}).ElseIf(f.notificationPermission == app.NotificationGranted, func() app.UI {
			return app.Text("Notification permission is already granted")
		}).Else(func() app.UI {
			return app.Text("Notification are not supported")
		}),
	)
}

func (f *foo) enableNotifications(ctx app.Context, e app.Event) {
	// Triggers a browser popup that asks for user permission.
	f.notificationPermission = ctx.Notifications().RequestPermission()
}
```

### Display Local Notifications

A local notification is a notification created in the app with [Context.Notifications().New()](/reference#NotificationService.New):

```go
func (f *foo) Render() app.UI {
	return app.Div().Body(
		app.If(f.notificationPermission == app.NotificationDefault, func() app.UI {
			return app.Button().
				Text("Enable Notifications").
				OnClick(f.enableNotifications)
		}).ElseIf(f.notificationPermission == app.NotificationDenied, func() app.UI {
			return app.Text("Notification permission is denied")
		}).ElseIf(f.notificationPermission == app.NotificationGranted, func() app.UI {
			return app.Button().
				Text("Test Notification").
				OnClick(f.enableNotifications)
		}).Else(func() app.UI {
			return app.Text("Notification are not supported")
		}),
	)
}

func (f *foo) testNotification(ctx app.Context, e app.Event) {
	ctx.Notifications().New(app.Notification{
		Title:  "Test",
		Body:   "A test notification",
		Path: "/mypage",
	})
}
```

**[Notification.Path](/reference#Notification) is a URL path that targets a page in the app. When a notification is clicked, the app will be navigated on this URL path.**

### Example
