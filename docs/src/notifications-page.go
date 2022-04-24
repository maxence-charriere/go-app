package main

import (
	"fmt"
	"math/rand"

	"github.com/maxence-charriere/go-app/v9/pkg/analytics"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type notificationsPage struct {
	app.Compo

	notificationPermission app.NotificationPermission
}

func newNotificationsPage() *notificationsPage {
	return &notificationsPage{}
}

func (p *notificationsPage) OnPreRender(ctx app.Context) {
	p.initPage(ctx)
}

func (p *notificationsPage) OnMount(ctx app.Context) {
	p.notificationPermission = ctx.NotificationPermission()
}

func (p *notificationsPage) OnNav(ctx app.Context) {
	p.initPage(ctx)
}

func (p *notificationsPage) initPage(ctx app.Context) {
	ctx.Page().SetTitle("Receive And Display Notifications")
	ctx.Page().SetDescription("Documentation about how receive and display notifications.")
	analytics.Page("notifications", nil)
}

func (p *notificationsPage) Render() app.UI {
	requestEnabled := ""
	if p.notificationPermission != app.NotificationDefault {
		requestEnabled = "disabled"
	}

	testEnabled := "disabled"
	if p.notificationPermission == app.NotificationGranted {
		testEnabled = ""
	}

	return newPage().
		Title("Notifications").
		Icon(bellSVG).
		Index(
			newIndexLink().Title("Enable Notifications"),
			newIndexLink().Title("    Current Permission"),
			newIndexLink().Title("    Request Permission"),
			newIndexLink().Title("    Display Local Notifications"),
			newIndexLink().Title("    Example"),

			app.Div().Class("separator"),

			newIndexLink().Title("Push Notifications"),

			app.Div().Class("separator"),

			newIndexLink().Title("Next"),
		).
		Content(
			newRemoteMarkdownDoc().Src("/web/documents/notifications.md"),

			app.P().Body(
				app.Button().
					Class("button").
					Class(requestEnabled).
					Text("Enable Notifications").
					OnClick(p.enableNotifications),
				app.Button().
					Class("button").
					Class(testEnabled).
					Text("Test Notification").
					OnClick(p.testNotification),
			),

			app.Div().Class("separator"),

			newRemoteMarkdownDoc().Src("/web/documents/notifications-push.md"),
		)
}

func (p *notificationsPage) enableNotifications(ctx app.Context, e app.Event) {
	p.notificationPermission = ctx.RequestNotificationPermission()
}

func (p *notificationsPage) testNotification(ctx app.Context, e app.Event) {
	n := rand.Intn(43)

	ctx.NewNotification(app.Notification{
		Title:  fmt.Sprintln("go-app test", n),
		Body:   fmt.Sprintln("Test notification for go-app number", n),
		Target: "/notifications#example",
	})
}
