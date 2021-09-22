## Intro

Apps built with this package are [progressive web apps (PWA)](https://web.dev/progressive-web-apps). They are **out of the box installable**, which mean they can run in their own window and can be pinned on the dock, taskbar, or home screen.

## Desktop

![Desktop install](/web/images/desktop-install.png)

- Open the app in [Chrome](https://www.google.com/chrome) or [Edge](https://www.microsoft.com/edge)
- Click on the `Install` button on the right inside the search bar
- In the install popup, click on the `Install` button

## IOS

![IOS install](/web/images/ios-install.png)

- Open the app in [Safari](https://www.apple.com/safari)
- Tap on the `Share` button
- Scroll and tap on the `Add to Home Screen` button
- Tapp on the `Add` button

## Android

![Android install](/web/images/android-install.png)

- Open the app in [Chrome](https://www.google.com/chrome)
- Tap on the `dot menu` button
- Tap `Add to home screen` button
- Tapp on the `Add` button

## Programmatically

On supported web browsers (usually Chromium-based), it is possible to know whether an app is installable and manually show the Install prompt.

### Detect Install Support

Detecting whether the app is installable from a web browser is done the `IsAppInstallable` [Context](/reference#Context) method:

```go
type hello struct {
	app.Compo

	name             string
	isAppInstallable bool
}

func (h *hello) OnMount(ctx app.Context) {
	h.isAppInstallable = ctx.IsAppInstallable()
}
```

Since the installable state can change depending on whether the app is installed, the component that displays an install button should be notified of the change. This is done by implementing the [AppInstaller](/reference#AppInstaller) interface.

```go
func (h *hello) OnAppInstallChange(ctx app.Context) {
	h.isAppInstallable = ctx.IsAppInstallable()
}
```

### Display Install Popup

Displaying the web browser install prompt is done by calling the `ShowAppInstallPrompt` [Context](/reference#Context) method.

```go
func (h *hello) Render() app.UI {
	return app.Div().
		Body(
			app.H1().Text("Hello World!"),

			app.If(h.isAppInstallable,
				app.Button().
					Text("Install App").
					OnClick(h.onInstallButtonClicked),
			),
		)
}

func (h *hello) onInstallButtonClicked(ctx app.Context, e app.Event) {
	ctx.ShowAppInstallPrompt()
}
```

## Next

- [Testing](/testing)
- [Reference](/reference)
