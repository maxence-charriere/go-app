// +build darwin,amd64

package mac

import (
	"net/url"

	"github.com/murlokswarm/app"

	"github.com/google/uuid"
	"github.com/murlokswarm/app/bridge"
)

// Notification implements the app.Element interface.
type Notification struct {
	id      uuid.UUID
	onReply func(reply string)
}

func newNotification(c app.NotificationConfig) error {
	n := &Notification{
		id:      uuid.New(),
		onReply: c.OnReply,
	}

	if n.onReply != nil {
		driver.elements.Add(n)
	}

	return driver.macRPC.Call("notifications.New", nil, struct {
		ID        string
		Title     string
		Subtitle  string
		Text      string
		ImageName string
		Sound     bool
		Reply     bool
	}{
		ID:        n.ID().String(),
		Title:     c.Title,
		Subtitle:  c.Subtitle,
		Text:      c.Text,
		ImageName: c.ImageName,
		Sound:     c.Sound,
		Reply:     c.OnReply != nil,
	})
}

// ID satisfies the app.Element interface.
func (n *Notification) ID() uuid.UUID {
	return n.id
}

func onNotificationReply(n *Notification, u *url.URL, p bridge.Payload) (res bridge.Payload) {
	var reply string
	p.Unmarshal(&reply)

	if reply != "(null)" {
		n.onReply(reply)
	}

	driver.elements.Remove(n)
	return nil
}
