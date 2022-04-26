## Push Notifications

Push notifications are notifications that are sent by a remote server and that can be displayed whether the app is running or closed.

Setting them up requires sending a subscription to a push notification server.

### Getting Notification Subscription

The push notification subscription is obtained with [Context.Notifications().Subscribe()](/reference#NotificationService.Subscribe):

```go
func (f *foo) enableNotifications(ctx app.Context, e app.Event) {
	f.notificationPermission = ctx.Notifications().RequestPermission()

	if f.notificationPermission == app.NotificationGranted {
		f.registerNotificationSubscription(ctx)
	}
}

func (f *foo) registerNotificationSubscription(ctx app.Context) {
	sub, err := ctx.Notifications().Subscribe("MY_VAPID_PUBLIC_KEY")
	if err != nil {
		log.Println("subscribing to push notifications failed:", err)
		return
	}
}
```

### Registering Notification Subscription

Once the subscription is obtained, it has to be registered on a push notification server. This is done with a classic HTTP request:

```go
func (f *foo) registerNotificationSubscription(ctx app.Context) {
	sub, err := ctx.Notifications().Subscribe("MY_VAPID_PUBLIC_KEY")
	if err != nil {
		log.Println("subscribing to push notifications failed:", err)
		return
	}

	ctx.Async(func() {
		var body bytes.Buffer
		if err := json.NewEncoder(&body).Encode(sub); err != nil {
			log.Println("encoding notification subscription failed:", err)
			return
		}

		res, err := http.Post("/PUSH_SERVER_ENDPOINT", "application/json", &body)
		if err != nil {
			log.Println("registering notification subscription failed:", err)
			return
		}
		defer res.Body.Close()
	})
}
```
