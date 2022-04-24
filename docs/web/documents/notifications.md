## Enable Notifications

Enabling notifications requires the user to give permission to display them.

### Current Permission

The current notifications permission is retrieved by calling [Context.NotificationPermission](/reference#Context):

```go
type foo struct {
	app.Compo

	notificationPermission app.NotificationPermission
}

func (f *foo) OnMount(ctx app.Context) {
	f.notificationPermission = ctx.NotificationPermission()
}
```

### Request Permission

The Notification permission is given by requesting the user permission with [Context.RequestNotificationPermission](/reference#Context):

```go
func (f *foo) Render() app.UI {
	return app.Div().Body(
		app.If(f.notificationPermission == app.NotificationGranted,
			app.Text("Notification permission is already granted"),
		).Else(
			app.Button().
				Text("Enable Notifications").
				OnClick(f.enableNotifications),
		),
	)
}

func (f *foo) enableNotifications(ctx app.Context, e app.Event) {
    // Triggers a browser popup that asks for user permission.
	f.notificationPermission = ctx.RequestNotificationPermission()
}
```

### Display Local Notifications

A local notification is a notification created in the app with [Context.NewNotification](/reference#Context):

```go
func (f *foo) Render() app.UI {
	return app.Div().Body(
		app.If(f.notificationPermission == app.NotificationGranted,
			app.Button().
				Text("Test Notification").
				OnClick(f.enableNotifications),
		).Else(
			app.Button().
				Text("Enable Notifications").
				OnClick(f.enableNotifications),
		),
	)
}

func (f *foo) testNotification(ctx app.Context, e app.Event) {
	ctx.NewNotification(app.Notification{
		Title: "Test",
		Body:  "A test notification",
        Target: "/mypage",
	})
}
```

### Example
