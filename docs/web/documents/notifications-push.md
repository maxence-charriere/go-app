## Push Notifications

Push notifications are notifications that are sent by a remote server and that can be displayed whether the app is running or closed. Setting them up requires sending a subscription to a push notification server.

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

### Sending Push Notification

Sending a push notification is done on the server side by using the subscription previously created. Here is an [http.Handler](https://pkg.go.dev/net/http#Handler) implementation based on the [webpush-go](https://pkg.go.dev/github.com/SherClockHolmes/webpush-go@v1.2.0) package.

```go
func main() {
	// ...
	http.Handle("/test/notifications/", &notificationHandler{
		VAPIDPrivateKey: "MY_VAPID_PRIVATE_KEY",
		VAPIDPublicKey:  "MY_VAPID_PUBLIC_KEY",
	})
	// ...
}


type notificationHandler struct {
	VAPIDPrivateKey string
	VAPIDPublicKey  string

	mutex         sync.Mutex
	subscriptions map[string]webpush.Subscription
}

func (h *notificationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch path := path.Base(r.URL.Path); path {
	case "register":
		h.handleRegistrations(w, r)

	case "test":
		h.handleTests(w, r)
	}
}
```

The first step is to receive and store the previously created subscription:

```go
func (h *notificationHandler) handleRegistrations(w http.ResponseWriter, r *http.Request) {
	var sub webpush.Subscription
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	if h.subscriptions == nil {
		h.subscriptions = make(map[string]webpush.Subscription)
	}
	h.subscriptions[sub.Endpoint] = sub
}
```

Then create and send a JSON encoded [notification](/reference#Notification):

```go
// handleTests creates and sends a push notification for all the registered
// subscriptions.
func (h *notificationHandler) handleTests(w http.ResponseWriter, r *http.Request) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for _, sub := range h.subscriptions {
		go func(sub webpush.Subscription) {
			body, _ := json.Marshal(app.Notification{
				Title: "Push test from server",
				Body:  "go-app push notification number",
				Path: "/mypage"
			})

			res, err := webpush.SendNotification(body, &sub, &webpush.Options{
				VAPIDPrivateKey: h.VAPIDPrivateKey,
				VAPIDPublicKey:  h.VAPIDPublicKey,
				TTL:             30,
			})
			if err != nil {
				app.Log(err)
				return
			}
			defer res.Body.Close()
		}(sub)
	}
}
```

**Push servers can be implemented in various programming languages. The requirement to receive a push notification with go-app is that the notification message is a JSON encoded [Notification struct](/reference#Notification).**
